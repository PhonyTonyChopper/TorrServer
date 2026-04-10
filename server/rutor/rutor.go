package rutor

import (
	"server/rutor/models"
)

// Search returns an empty result set. RuTor search is not available in this build.
func Search(_ string) []*models.TorrentDetails {
	return nil
}

// Start is a no-op stub.
func Start() {}

// Stop is a no-op stub.
func Stop() {}
