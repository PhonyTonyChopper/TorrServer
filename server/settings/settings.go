package settings

import (
	"os"
	"path/filepath"

	"server/log"
)

var (
	tdb      *TDB
	Path     string
	IP       string
	Port     string
	Ssl      bool
	SslPort  string
	ReadOnly bool
	HttpAuth bool
	SearchWA bool
	PubIPv4  string
	PubIPv6  string
	TorAddr  string
	MaxSize  int64
)

func InitSets(readOnly, searchWA bool) {
	ReadOnly = readOnly
	SearchWA = searchWA
	tdb = NewTDB()
	if tdb == nil {
		log.TLogln("Error open db:", filepath.Join(Path, "config.db"))
		os.Exit(1)
	}
	loadBTSets()
	Migrate()
}

func CloseDB() {
	tdb.CloseDB()
}

// GetStoragePreferences returns the current storage configuration.
// In this build, all data is stored in bbolt.
func GetStoragePreferences() map[string]interface{} {
	return map[string]interface{}{
		"settings": "bbolt",
		"viewed":   "bbolt",
	}
}

// SetStoragePreferences is a no-op in this build.
// Storage backend switching is not supported.
func SetStoragePreferences(_ map[string]interface{}) error {
	return nil
}
