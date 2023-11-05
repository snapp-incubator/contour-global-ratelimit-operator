package parser

import (
	"fmt"
	"strconv"
	"strings"

	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
)

type HTTPProxyGlobalRateLimitPolicy struct {
	Name                  string
	Namespace             string
	IngressClass          string // IngressClass is equivalent to the domain
	RateLimitsDescriptors []Descriptor
}

type Descriptor struct {
	Key         string
	Value       string
	RateLimit   RateLimit
	ShadowMode  bool
	Descriptors []Descriptor
}

type RateLimit struct {
	Unit            string
	RequestsPerUnit string
}

// ExtractDescriptorsFromHTTPProxy extracts rate limit descriptors from an HTTPProxy.
func ExtractDescriptorsFromHTTPProxy(httpProxy *contourv1.HTTPProxy) (hasRateLimitConfig bool, policies HTTPProxyGlobalRateLimitPolicy, err error) {
	// Initialize the global rate limit policy
	globalRateLimitPolicy := HTTPProxyGlobalRateLimitPolicy{
		Name:         httpProxy.Name,
		Namespace:    httpProxy.Namespace,
		IngressClass: httpProxy.Spec.IngressClassName,
	}

	for _, route := range httpProxy.Spec.Routes {
		if route.RateLimitPolicy != nil && route.RateLimitPolicy.Global != nil {
			hasRateLimitConfig = true

			// Extract descriptors from the global rate limit policy
			descriptors, extractErr := extractDescriptorsFromGlobalRateLimitPolicy(route.RateLimitPolicy.Global, globalRateLimitPolicy.Name, globalRateLimitPolicy.Namespace)
			if extractErr != nil {
				return false, HTTPProxyGlobalRateLimitPolicy{}, extractErr
			}
			// Append the extracted descriptors to the global rate limit policy
			globalRateLimitPolicy.RateLimitsDescriptors = append(globalRateLimitPolicy.RateLimitsDescriptors, descriptors...)

		}
	}

	return hasRateLimitConfig, globalRateLimitPolicy, err

	//return false, HTTPProxyGlobalRateLimitPolicy{}, err
}

func extractDescriptorsFromGlobalRateLimitPolicy(policy *contourv1.GlobalRateLimitPolicy, name string, namespace string) ([]Descriptor, error) {
	var descriptors []Descriptor
	for _, contourDescriptor := range policy.Descriptors {
		var entryDescriptor Descriptor
		var limit RateLimit
		for i, entry := range contourDescriptor.Entries {
			descriptor, rateLimit, err := extractDescriptorFromEntry(entry, name, namespace)
			if err != nil {
				//log.Println("line 78#####", err, entry.GenericKey)
				//continue
				return descriptors, err
			}
			if i == 0 {
				limit = rateLimit
				entryDescriptor = descriptor
				//first entry must be genricKey and also maybe dosn't have second entry
				if len(contourDescriptor.Entries) == 1 && entry.GenericKey != nil {
					entryDescriptor.RateLimit = limit
					descriptors = append(descriptors, entryDescriptor)
				}
			}
			//Todo: Add name_namspace validation here
			if i == 1 && entryDescriptor.Key != "" {
				descriptor.RateLimit = limit
				entryDescriptor.Descriptors = append(entryDescriptor.Descriptors, descriptor)
				descriptors = append(descriptors, entryDescriptor)
			}

		}
	}
	return descriptors, nil
}

func extractDescriptorFromEntry(entry contourv1.RateLimitDescriptorEntry, name string, namespace string) (Descriptor, RateLimit, error) {
	switch {
	case entry.GenericKey != nil:
		genericDescriptor, limit, err := extractDescriptorFromGenericKey(entry.GenericKey, name, namespace)
		if err != nil {
			return Descriptor{}, RateLimit{}, err
		}
		return genericDescriptor, limit, nil

	case entry.RequestHeader != nil:
		requestHeaderDescriptor, err := extractDescriptorFromRequestHeader(entry.RequestHeader)
		if err != nil {
			return Descriptor{}, RateLimit{}, err
		}
		return requestHeaderDescriptor, RateLimit{}, nil

	case entry.RequestHeaderValueMatch != nil:
		headerValueMatchDescriptor, err := extractDescriptorFromRequestHeaderValueMatch(entry.RequestHeaderValueMatch)
		if err != nil {
			return Descriptor{}, RateLimit{}, err
		}
		return headerValueMatchDescriptor, RateLimit{}, nil

	case entry.RemoteAddress != nil:
		remoteAddr, err := extractDescriptorsFromRemoteAddress(*entry.RemoteAddress)
		return remoteAddr, RateLimit{}, err
	default:
		return Descriptor{}, RateLimit{}, fmt.Errorf("RateLimitDescriptorEntry not found")
	}
}

func extractDescriptorFromGenericKey(genericKey *contourv1.GenericKeyDescriptor, name string, namepsace string) (Descriptor, RateLimit, error) {
	key := genericKey.Key
	val := genericKey.Value
	if !isGenericKeyContaineName_Namespace(key, name, namepsace) {
		err := fmt.Errorf("%v is not valid", key)
		return Descriptor{}, RateLimit{}, err
	}
	limit, err := generateRateLimitFromGenericKeyValue(val, name, namepsace)
	if err != nil {
		return Descriptor{}, RateLimit{}, err
	}
	return Descriptor{
		Key:   key,
		Value: val,
	}, limit, nil
}
func isGenericKeyContaineName_Namespace(key string, name string, namespace string) bool {
	return strings.HasPrefix(key, fmt.Sprintf("%v.%v.", namespace, name))
}

func extractDescriptorFromRequestHeader(requestHeader *contourv1.RequestHeaderDescriptor) (Descriptor, error) {
	key := requestHeader.DescriptorKey
	return Descriptor{
		Key: key,
	}, nil
}
func extractDescriptorFromRequestHeaderValueMatch(headerValueMattch *contourv1.RequestHeaderValueMatchDescriptor) (Descriptor, error) {

	return Descriptor{
		Key:   "header_match",
		Value: headerValueMattch.Value,
	}, nil
}
func extractDescriptorsFromRemoteAddress(addresse contourv1.RemoteAddressDescriptor) (Descriptor, error) {
	return Descriptor{
		Key: "remote_address",
	}, nil
}

func generateRateLimitFromGenericKeyValue(value string, name string, namespace string) (RateLimit, error) {
	valParts := strings.Split(value, "/")
	if len(valParts) != 2 {
		return RateLimit{}, fmt.Errorf("value %v invalid format", value)
	}
	requestsPerUnit, err := strconv.ParseUint(valParts[0], 10, 32)
	if err != nil {
		return RateLimit{}, err
	}
	unit := valParts[1]
	validUnits := map[string]bool{"s": true, "m": true, "h": true, "d": true}
	if !validUnits[unit] {
		return RateLimit{}, fmt.Errorf("value must be s, m, h, or d: %v", value)
	}
	return RateLimit{
		RequestsPerUnit: fmt.Sprintf("%v", requestsPerUnit),
		Unit:            unit,
	}, nil
}
