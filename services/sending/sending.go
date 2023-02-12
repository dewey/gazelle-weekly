package sending

import (
	"context"
	"fmt"
	"gazelle-weekly/repositories/email"
	"gazelle-weekly/repositories/gazelle"
	"gazelle-weekly/repositories/streamingservice"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Service struct {
	l       log.Logger
	emailer *email.Client
	gazelle *gazelle.Client
	token   string
}

func NewSendingService(l log.Logger, emailer *email.Client, gazelle *gazelle.Client) *Service {
	return &Service{
		l:       l,
		emailer: emailer,
		gazelle: gazelle,
	}
}

func (s *Service) Send(ctx context.Context, payload []DecoratedResult) error {
	var albums []email.TemplateVariables
	for _, pl := range payload {
		vars := email.TemplateVariables{
			Album:       pl.Result.GroupName,
			Artist:      pl.Result.Artist,
			ReleaseYear: pl.Result.GroupYear,
			Tags:        pl.Result.Tags,
			URL:         fmt.Sprintf("https://%s/torrents.php?id=%d&torrentid=%d", s.gazelle.BaseURL, pl.Result.GroupId, pl.Result.TorrentId),
			ArtworkURL:  pl.Result.WikiImage,
		}
		var urls []email.StreamingURL
		for _, album := range pl.Albums {
			for _, l := range album.URLs {
				urls = append(urls, email.StreamingURL{
					URL:      l.URL,
					LinkType: l.LinkType,
				})
			}
		}
		vars.URLs = urls
		albums = append(albums, vars)
	}

	level.Info(s.l).Log("msg", "sending new email", "item_count", len(albums))
	if err := s.emailer.Send(ctx,
		fmt.Sprintf("Redacted Weekly: Top 10 on %s", time.Now().Format("2 January 2006")),
		fmt.Sprintf("Redacted Weekly, %s", time.Now().Format("2 January 2006")),
		albums,
	); err != nil {
		return err
	}
	return nil
}

// DecoratedResult is an album and it's streaming urls
type DecoratedResult struct {
	Result gazelle.Result
	Albums []streamingservice.Album
}
