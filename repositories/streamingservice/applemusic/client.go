package applemusic

import (
	"context"
	"encoding/json"
	"fmt"
	"gazelle-weekly/repositories/streamingservice"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Repository struct {
	client *http.Client
	l      log.Logger
}

func NewClient(ctx context.Context, l log.Logger, client *http.Client) *Repository {
	if client == nil {
		client = http.DefaultClient
	}
	return &Repository{
		client: client,
		l:      l,
	}
}

// LatestAlbum returns the latest album in a list of albums.
func (s *Repository) LatestAlbum(ctx context.Context, countryISO string, searchTerm string, releaseYear int) ([]streamingservice.Album, error) {
	var latest streamingservice.Album
	res, err := s.SearchAlbum(ctx, countryISO, searchTerm, releaseYear, false)
	if err != nil {
		return nil, err
	}
	for _, re := range res {
		if re.ReleaseDate.After(latest.ReleaseDate) {
			latest = re
		}
	}
	return []streamingservice.Album{latest}, nil
}

func (s *Repository) SearchAlbum(ctx context.Context, countryISO string, searchTerm string, releaseYear int, releaseYearFilter bool) ([]streamingservice.Album, error) {
	u, err := url.Parse("https://itunes.apple.com/search")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("term", searchTerm)
	q.Add("country", countryISO)
	q.Add("entity", "album")
	q.Add("media", "apple_music")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	level.Debug(s.l).Log("msg", "searching streaming service", "term", searchTerm, "country", countryISO)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code. got %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}
	var albums []streamingservice.Album
	for _, i := range apiResponse.Results {
		if releaseYearFilter && !releaseYearWithinRange(i.ReleaseDate.Year(), releaseYear, 3) {
			level.Info(s.l).Log("msg", "release year not within filter range", "release_year", i.ReleaseDate.Year(), "input_release_year", releaseYear)
			continue
		}
		a := streamingservice.Album{
			Artist:           i.ArtistName,
			ArtistExternalID: strconv.Itoa(i.ArtistId),
			Album:            i.CollectionName,
			AlbumExternalID:  strconv.Itoa(i.CollectionId),
			ArtworkURL:       i.ArtworkUrl100,
			ReleaseDate:      i.ReleaseDate,
		}
		var urls []streamingservice.StreamingURL
		urls = append(urls, streamingservice.StreamingURL{
			URL:      fmt.Sprintf("https://music.apple.com/%s/album/%d?l=en", strings.ToLower(countryISO), i.CollectionId),
			LinkType: "web_streaming",
		})
		a.URLs = urls
		albums = append(albums, a)
	}
	return albums, nil
}

// Response is a generic top level response from the Apple iTunes API
type Response struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

// Result contains an album result from the Apple iTunes API
type Result struct {
	WrapperType            string    `json:"wrapperType"`
	CollectionType         string    `json:"collectionType"`
	ArtistId               int       `json:"artistId"`
	CollectionId           int       `json:"collectionId"`
	AmgArtistId            int       `json:"amgArtistId"`
	ArtistName             string    `json:"artistName"`
	CollectionName         string    `json:"collectionName"`
	CollectionCensoredName string    `json:"collectionCensoredName"`
	ArtistViewUrl          string    `json:"artistViewUrl"`
	CollectionViewUrl      string    `json:"collectionViewUrl"`
	ArtworkUrl60           string    `json:"artworkUrl60"`
	ArtworkUrl100          string    `json:"artworkUrl100"`
	CollectionPrice        float64   `json:"collectionPrice"`
	CollectionExplicitness string    `json:"collectionExplicitness"`
	ContentAdvisoryRating  string    `json:"contentAdvisoryRating"`
	TrackCount             int       `json:"trackCount"`
	Copyright              string    `json:"copyright"`
	Country                string    `json:"country"`
	Currency               string    `json:"currency"`
	ReleaseDate            time.Time `json:"releaseDate"`
	PrimaryGenreName       string    `json:"primaryGenreName"`
}
