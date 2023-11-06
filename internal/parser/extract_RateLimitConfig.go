package parser

import (
	"sync"

	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
)

var (
	ContourLimitConfigs = NewLimitConfigs("contour")
)

type LimitConfigs struct {
	rateLimitConfig rls_config.RateLimitConfig
	descriptorMaps  map[string]*rls_config.RateLimitDescriptor
	mutex           sync.Mutex
}

func NewLimitConfigs(domain string) *LimitConfigs {

	limitConf := &LimitConfigs{}
	limitConf.setDomain(domain)
	limitConf.descriptorMaps = make(map[string]*rls_config.RateLimitDescriptor)
	return limitConf

}

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
func (a *LimitConfigs) AddToConfig(policy HTTPProxyGlobalRateLimitPolicy) error {
	descriptors := convertToRateLimitDescriptors(policy.RateLimitsDescriptors)
	if descriptors != nil {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		for _, d := range descriptors {
			a.descriptorMaps[d.Key] = d
		}

	}
	//yml, _ := yaml.Marshal(&a.descriptorMaps)
	//fmt.Println(string(yml))
	return nil
}
func (a *LimitConfigs) Delete(namespace string, name string) (deleted bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, d := range a.descriptorMaps {
		if isGenericKeyContaineName_Namespace(d.Key, name, namespace) {
			delete(a.descriptorMaps, d.Key)
			deleted = true

		}
	}
	return deleted
}
func (a *LimitConfigs) setDomain(name string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.rateLimitConfig.Domain = name
	a.rateLimitConfig.Name = name

}
