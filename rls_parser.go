package main

import (
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
)

// RLSParser is used to extract rate limit configs from httpproxy.
type RlsParser struct{}

// ExtractRlsDescriptors is a wrapper around the internal parser package.
func (r *RlsParser) ExtractRlsDescriptors(httpProxy *contourv1.HTTPProxy) (hasRateLimitConfig bool, policies parser.HTTPProxyGlobalRateLimitPolicy, err error) {
	return parser.ExtractDescriptorsFromHTTPProxy(httpProxy)
}
