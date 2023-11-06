package parser

import (
	"fmt"
	"sync"

	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
)

// ContourLimitConfigs is a global instance of LimitConfigs for the "contour" domain.
var ContourLimitConfigs = NewLimitConfigs("contour")

// LimitConfigs is a configuration for rate limiting.
type LimitConfigs struct {
	rateLimitConfig rls_config.RateLimitConfig
	descriptorMaps  map[string]*rls_config.RateLimitDescriptor
	mutex           sync.Mutex
}

// NewLimitConfigs creates a new instance of LimitConfigs with the specified domain.
func NewLimitConfigs(domain string) *LimitConfigs {
	limitConf := &LimitConfigs{}
	limitConf.setDomain(domain)
	limitConf.descriptorMaps = make(map[string]*rls_config.RateLimitDescriptor)
	return limitConf
}

// GetConfigs returns the rate limit configuration.
func (a *LimitConfigs) GetConfigs() *rls_config.RateLimitConfig {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.rateLimitConfig.Descriptors = nil
	a.rateLimitConfig.Descriptors = make([]*rls_config.RateLimitDescriptor, 0)
	for _, d := range a.descriptorMaps {
		a.rateLimitConfig.Descriptors = append(a.rateLimitConfig.Descriptors, d)
	}
	if len(a.rateLimitConfig.Descriptors) == 0 {
		return nil
	}
	return &a.rateLimitConfig
}

// AddToConfig adds rate limit policies to the configuration.
func (a *LimitConfigs) AddToConfig(policy HTTPProxyGlobalRateLimitPolicy) error {
	descriptors := convertToRateLimitDescriptors(policy.RateLimitsDescriptors)
	if descriptors != nil {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		for _, d := range descriptors {
			a.descriptorMaps[d.Key] = d
		}
		return nil
	}
	return fmt.Errorf("descriptors is empty")
}

// Delete removes rate limit policies associated with a namespace and name.
func (a *LimitConfigs) Delete(namespace string, name string) (deleted bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, d := range a.descriptorMaps {
		if isGenericKeyContainNameNamespace(d.Key, name, namespace) {
			delete(a.descriptorMaps, d.Key)
			deleted = true
		}
	}
	return deleted
}

// setDomain sets the domain and name of the rate limit configuration.
func (a *LimitConfigs) setDomain(name string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.rateLimitConfig.Domain = name
	a.rateLimitConfig.Name = name
}
