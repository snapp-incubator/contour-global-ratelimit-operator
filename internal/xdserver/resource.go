package xdserver

import (
	"context"
	"os"

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

var snapshotCache = cache.NewSnapshotCache(true, cache.IDHash{}, logger)

func CreateNewSnapshot(version string) {

	// Create a snapshotCache

	snapshot := GenerateSnapshot(version)
	if err := snapshot.Consistent(); err != nil {
		logger.Errorf("Snapshot is inconsistent: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	logger.Debugf("Will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := snapshotCache.SetSnapshot(context.Background(), NodeID, snapshot); err != nil {
		logger.Errorf("Snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

}
