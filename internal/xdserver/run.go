package xdserver

import (
	"context"
	"flag"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

var (
	logger Logger
	port   uint
	nodeID string
)

func init() {
	logger = Logger{}

	flag.BoolVar(&logger.Debug, "debug", false, "Enable xDS server debug logging")
	flag.UintVar(&port, "port", 18000, "xDS management server port")
	flag.StringVar(&nodeID, "nodeID", "test-node-id", "Node ID")

}

var snapshotCache = cache.NewSnapshotCache(true, cache.IDHash{}, logger)

func CreateNewSnapshot(version string) {
	flag.Parse()

	// Create a snapshotCache

	snapshot := GenerateSnapshot(version)
	if err := snapshot.Consistent(); err != nil {
		logger.Errorf("Snapshot is inconsistent: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	logger.Debugf("Will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := snapshotCache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		logger.Errorf("Snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

}
