package parser

import (
	"strconv"
	"strings"

	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
)

// convertToRateLimitDescriptors converts a slice of localDescriptors to a slice of RateLimitDescriptors.
func convertToRateLimitDescriptors(localDescriptors []Descriptor) []*rls_config.RateLimitDescriptor {
	descriptors := make([]*rls_config.RateLimitDescriptor, len(localDescriptors))

	for i, localDescriptor := range localDescriptors {
		descriptors[i] = convertToRateLimitDescriptor(localDescriptor)
	}

	return descriptors
}

// convertToRateLimitDescriptor converts a localDescriptor to a RateLimitDescriptor.
func convertToRateLimitDescriptor(localDescriptor Descriptor) *rls_config.RateLimitDescriptor {
	// Helper function to convert a RateLimit struct to a RateLimitPolicy struct
	convertRateLimit := func(r RateLimit) *rls_config.RateLimitPolicy {

		if r.Unit == "" || r.RequestsPerUnit == "" {
			return nil
		}

		policy := &rls_config.RateLimitPolicy{}
		switch strings.ToLower(r.Unit) {
		case "s":
			policy.Unit = rls_config.RateLimitUnit_SECOND
		case "m":
			policy.Unit = rls_config.RateLimitUnit_MINUTE
		case "h":
			policy.Unit = rls_config.RateLimitUnit_HOUR
		case "d":
			policy.Unit = rls_config.RateLimitUnit_DAY
		}
		requestPerUnit, _ := strconv.ParseUint(r.RequestsPerUnit, 10, 32)

		policy.RequestsPerUnit = uint32(requestPerUnit)

		return policy
	}

	// Helper function to convert a slice of Descriptors to a slice of RateLimitDescriptors
	convertDescriptors := func(descriptors []Descriptor) []*rls_config.RateLimitDescriptor {
		if len(descriptors) == 0 {
			return nil
		}

		rateLimitDescriptors := make([]*rls_config.RateLimitDescriptor, 0)

		for _, desc := range descriptors {
			rateLimitDescriptors = append(rateLimitDescriptors, convertToRateLimitDescriptor(desc))
		}

		return rateLimitDescriptors
	}

	// Create and return the RateLimitDescriptor
	return &rls_config.RateLimitDescriptor{
		Key:         localDescriptor.Key,
		Value:       localDescriptor.Value,
		RateLimit:   convertRateLimit(localDescriptor.RateLimit),
		Descriptors: convertDescriptors(localDescriptor.Descriptors),
	}
}
