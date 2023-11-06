package xdserver

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
)

// MakeRlsConfig creates rate limit config resources.
func MakeRlsConfig() []types.Resource {
	return []types.Resource{
		parser.ContourLimitConfigs.GetConfigs(),
	}
}

// GenerateSnapshot generates an xDS snapshot for the given version.
func GenerateSnapshot(version string) *cache.Snapshot {
	snap, _ := cache.NewSnapshot(version,
		map[resource.Type][]types.Resource{
			resource.RateLimitConfigType: MakeRlsConfig(),
		},
	)
	return snap
}

var snapshotCache = cache.NewSnapshotCache(true, cache.IDHash{}, logger)
var snapShotVersion int
var snapShotMutex sync.Mutex

// CreateNewSnapshot creates and serves a new snapshot.
func CreateNewSnapshot() {
	version := func() (version string) {
		snapShotMutex.Lock()
		snapShotVersion++
		version = fmt.Sprint(snapShotVersion)
		snapShotMutex.Unlock()
		return

	}()

	// Generate a snapshot
	snapshot := GenerateSnapshot(version)

	// Check if the snapshot is empty for "contour"
	if m := snapshot.GetResources(resource.RateLimitConfigType); m["contour"] == nil {
		logger.Errorf("Snapshot is empty for : %+v\n%+v", NodeID)
		return
	}

	// Check if the snapshot is consistent
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
	logger.Infof("New version of Snapshot Created: ", version)
}
