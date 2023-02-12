package main

import (
	"context"
	"flag"
	"fmt"
	"gazelle-weekly/repositories/email"
	"gazelle-weekly/repositories/gazelle"
	"gazelle-weekly/repositories/streamingservice"
	"gazelle-weekly/repositories/streamingservice/applemusic"
	"gazelle-weekly/services/sending"
	"os"

	"github.com/peterbourgon/ff/v3"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	ctx := context.Background()
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	fs := flag.NewFlagSet("gazelle-weekly", flag.ContinueOnError)
	var (
		postmarkAPIToken = fs.String("postmark-api-token", "", "The Postmark.com API token")
		fromEmail        = fs.String("from-email", "", "The email address the weekly update should be sent from. This has to match the Postmark verified address.")
		gazelleAPIToken  = fs.String("gazelle-api-token", "", "The Gazelle API token")
		gazelleBaseURL   = fs.String("gazelle-base-url", "", "The Gazelle base url")
		toEmail          = fs.String("to-email", "", "The email address the weekly update should be sent to.")
	)
	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		level.Error(logger).Log("msg", "error parsing flags", "err", err)
		return
	}

	// Initialize clients and services
	gazelleClient := gazelle.NewClient(nil, *gazelleAPIToken, *gazelleBaseURL)
	amClient := applemusic.NewClient(ctx, logger, nil)
	emailerClient := email.NewClient(*postmarkAPIToken, *toEmail, *fromEmail)
	sendingService := sending.NewSendingService(logger, emailerClient, gazelleClient)

	result, err := gazelleClient.GetUniqueTop10(ctx, "week")
	if err != nil {
		level.Error(logger).Log("msg", "error getting top 10 from gazelle", "err", err)
		return
	}
	var albums []sending.DecoratedResult
	for _, searchResult := range result {
		dr := sending.DecoratedResult{
			Result: searchResult,
		}
		var results []streamingservice.Album
		res, err := amClient.SearchAlbum(ctx, "DE", fmt.Sprintf("%s %s", searchResult.Artist, searchResult.GroupName), searchResult.Year, true)
		if err != nil {
			level.Error(logger).Log("msg", "error getting album search results from apple music", "err", err)
			continue
		}
		results = res
		// If we don't get a hit, we just search for the artist and their latest album so we at least have a link to the artist page
		if len(results) < 1 {
			res, err := amClient.LatestAlbum(ctx, "DE", searchResult.Artist, searchResult.Year)
			if err != nil {
				level.Error(logger).Log("msg", "error getting artist search results from apple music", "err", err)
				continue
			}
			results = res
		}
		for _, response := range results {
			dr.Albums = append(dr.Albums, response)
		}
		albums = append(albums, dr)
		level.Info(logger).Log("msg", "album prepared for sending", "artist", dr.Result.Artist, "album", dr.Result.GroupName, "streaming_urls_count", len(dr.Albums))
	}

	if err := sendingService.Send(ctx, albums); err != nil {
		level.Error(logger).Log("msg", "error sending notification", "err", err)
		return
	}
}
