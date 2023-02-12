package gazelle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
)

type Client struct {
	client  *http.Client
	token   string
	BaseURL string
}

func NewClient(client *http.Client, token string, baseURL string) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		client:  client,
		token:   token,
		BaseURL: baseURL,
	}
}

// GetUniqueTop10 is deduplicating the Top 10, as different encodings are distinct entries in the list
func (c *Client) GetUniqueTop10(ctx context.Context, period string) ([]Result, error) {
	res, err := c.GetTop10(ctx, period)
	if err != nil {
		return nil, err
	}
	uniqueEntries := map[int]struct{}{}
	var uniqueResults []Result
	for _, searchResult := range res {
		// We don't want singles, remixes, bootlegs etc.
		switch searchResult.ReleaseType {
		case "1", "3", "5", "6", "11", "18", "19":
		default:
			continue
		}
		if _, ok := uniqueEntries[searchResult.GroupId]; !ok {
			uniqueResults = append(uniqueResults, searchResult)
			uniqueEntries[searchResult.GroupId] = struct{}{}
		}
		if len(uniqueResults) == 10 {
			return uniqueResults, nil
		}
	}
	return nil, errors.New("couldn't gather required number of unique top 10 entries")
}

func (c *Client) GetTop10(ctx context.Context, period string) ([]Result, error) {
	// Sanitizing available input values according to API
	var details string
	switch period {
	case "day", "week", "month", "year":
		details = period
	default:
		details = "day"
	}
	u, err := url.Parse(fmt.Sprintf("https://%s/ajax.php", c.BaseURL))
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("details", details)
	q.Add("action", "top10")
	q.Add("limit", "100")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Header.Add("Authorization", c.token)
	req.Header.Set("User-Agent", "https://github.com/dewey/gazelle-weekly")

	resp, err := c.client.Do(req)
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
	if apiResponse.Status != "success" {
		return nil, fmt.Errorf("unexpected status message. got %s, expected %s", apiResponse.Status, "success")
	}
	if len(apiResponse.Response) == 1 {
		return apiResponse.Response[0].Results, nil
	}
	return nil, fmt.Errorf("unexpected result count. got %d, expected %d", len(apiResponse.Response), 1)
}

// UnmarshalJSON is a custom Unmarshaler we use as Gazelle has encoding problems depending on which API route you use.
func (r *Result) UnmarshalJSON(bytes []byte) error {
	type CleanResult Result
	var result CleanResult

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*r = Result(result)
	r.GroupName = html.UnescapeString(result.GroupName)
	r.Artist = html.UnescapeString(result.Artist)
	return nil
}

// Response is a generic API response from Gazelle
type Response struct {
	Status   string `json:"status"`
	Response []struct {
		Caption string   `json:"caption"`
		Tag     string   `json:"tag"`
		Limit   int      `json:"limit"`
		Results []Result `json:"results"`
	} `json:"response"`
}

// Result is a single search result
type Result struct {
	TorrentId     int      `json:"torrentId"`
	GroupId       int      `json:"groupId"`
	Artist        string   `json:"artist"`
	GroupName     string   `json:"groupName"`
	GroupCategory int      `json:"groupCategory"`
	GroupYear     int      `json:"groupYear"`
	RemasterTitle string   `json:"remasterTitle"`
	Format        string   `json:"format"`
	Encoding      string   `json:"encoding"`
	HasLog        bool     `json:"hasLog"`
	HasCue        bool     `json:"hasCue"`
	Media         string   `json:"media"`
	Scene         bool     `json:"scene"`
	Year          int      `json:"year"`
	Tags          []string `json:"tags"`
	Snatched      int      `json:"snatched"`
	Seeders       int      `json:"seeders"`
	Leechers      int      `json:"leechers"`
	Data          int64    `json:"data"`
	Size          int      `json:"size"`
	WikiImage     string   `json:"wikiImage"`
	ReleaseType   string   `json:"releaseType"`
}
