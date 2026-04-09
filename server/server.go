package server

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"server/log"
	"server/settings"
	"server/torr/utils"
	"server/web"
)

func Start(pathdb, port string, roSets, searchWA bool) {
	settings.Path = pathdb
	settings.InitSets(roSets, searchWA)
	if roSets {
		log.TLogln("Enabled Read-only DB mode!")
	}

	if port == "" {
		port = "8090"
	}
	log.TLogln("Check web port", port)
	l, err := net.Listen("tcp", ":"+port)
	if l != nil {
		l.Close()
	}
	if err != nil {
		log.TLogln("Port", port, "already in use! Please set different port for HTTP. Abort")
		os.Exit(1)
	}
	// remove old disk caches
	go cleanCache()
	// set http port and start web server
	settings.Port = port
	web.Start()
}

func cleanCache() {
	if !settings.BTsets.UseDisk || settings.BTsets.TorrentsSavePath == "/" || settings.BTsets.TorrentsSavePath == "" {
		return
	}

	dirs, err := os.ReadDir(settings.BTsets.TorrentsSavePath)
	if err != nil {
		return
	}

	torrs := settings.ListTorrent()

	log.TLogln("Remove unused cache in dir:", settings.BTsets.TorrentsSavePath)
	keep := map[string]bool{}
	for _, d := range dirs {
		if len(d.Name()) != 40 {
			// Not a hash
			continue
		}

		if !settings.BTsets.RemoveCacheOnDrop {
			keep[d.Name()] = true
			for _, t := range torrs {
				if d.IsDir() && d.Name() == t.InfoHash.HexString() {
					keep[d.Name()] = false
					break
				}
			}
			for hash, del := range keep {
				if del && hash == d.Name() {
					log.TLogln("Remove unused cache:", d.Name())
					removeAllFiles(filepath.Join(settings.BTsets.TorrentsSavePath, d.Name()))
				}
			}
		} else {
			if d.IsDir() {
				log.TLogln("Remove unused cache:", d.Name())
				removeAllFiles(filepath.Join(settings.BTsets.TorrentsSavePath, d.Name()))
			}
		}
	}
}

func removeAllFiles(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, f := range files {
		name := filepath.Join(path, f.Name())
		os.Remove(name)
	}
	os.Remove(path)
}

func WaitServer() string {
	err := web.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func Stop() {
	web.Stop()
	settings.CloseDB()
}

func AddTrackers(trackers string) {
	lines := strings.Split(trackers, "\n")
	var tracks []string
	for _, l := range lines {
		l = strings.Trim(l, " ,\r")
		if l != "" {
			tracks = append(tracks, l)
		}
	}
	utils.SetDefTrackers(tracks)
}

