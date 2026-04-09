package models

import (
	"strings"
	"time"
)

const (
	CatMovie         = "Movie"
	CatSeries        = "Series"
	CatDocMovie      = "DocMovie"
	CatDocSeries     = "DocSeries"
	CatCartoonMovie  = "CartoonMovie"
	CatCartoonSeries = "CartoonSeries"
	CatTVShow        = "TVShow"
	CatAnime         = "Anime"
)

type TorrentDetails struct {
	Title        string
	Name         string
	Names        []string
	Categories   string
	Size         string
	CreateDate   time.Time
	Tracker      string
	Link         string
	Year         int
	Peer         int
	Seed         int
	Magnet       string
	Hash         string
	IMDBID       string
	VideoQuality int
	AudioQuality int
}

type TorrentFile struct {
	Name string
	Size int64
}

func (d TorrentDetails) GetNames() string {
	return strings.Join(d.Names, " ")
}
