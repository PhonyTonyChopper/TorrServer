package torrServer

import (
	server "server"
	"server/settings"
)

func StartTorrentServer(pathdb string) {
	settings.Path = pathdb
	settings.Args = &settings.ExecArgs{}
	server.Start()
}

func WaitTorrentServer() {
	server.WaitServer()
}

func StopTorrentServer() {
	server.Stop()
}

func AddTrackers(trackers string) {
	server.AddTrackers(trackers)
}
