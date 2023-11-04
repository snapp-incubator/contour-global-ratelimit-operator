package xdserver

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
)

func makeRlsConfig() []types.Resource {
	return []types.Resource{
		parser.ContourLimitConfigs.GetConfigs(),
	}

}

func GenerateSnapshot(version string) *cache.Snapshot {
	snap, _ := cache.NewSnapshot(version,
		map[resource.Type][]types.Resource{
			resource.RateLimitConfigType: makeRlsConfig(),
		},
	)
	return snap
}
