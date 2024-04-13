package rlsparser

import (
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
)

// ParseGlobalRateLimit is a wrapper around the internal parser package.
func ParseGlobalRateLimit(httpProxy *contourv1.HTTPProxy) (hasRateLimitConfig bool, policies parser.HTTPProxyGlobalRateLimitPolicy, err error) {
	return parser.ExtractDescriptorsFromHTTPProxy(httpProxy)
}
