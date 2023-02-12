package streamingservice

import (
	"context"
	"time"
)

type Repository interface {
	SearchAlbum(ctx context.Context, searchTerm string, releaseYear int) ([]Album, error)
}

// Album is our representation of an album
type Album struct {
	Artist           string         `json:"artist"`
	ArtistExternalID string         `json:"artist_external_id"`
	Album            string         `json:"album"`
	AlbumExternalID  string         `json:"album_external_id"`
	ArtworkURL       string         `json:"artwork_url"`
	ReleaseDate      time.Time      `json:"release_date"`
	URLs             []StreamingURL `json:"urls"`
}

// StreamingURL is a URL to an entity on a streaming service
type StreamingURL struct {
	URL      string `json:"url"`
	LinkType string `json:"link_type"`
}
