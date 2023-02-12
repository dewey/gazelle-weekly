package email

import (
	"context"
	"errors"
	"net/http"

	"github.com/mattevans/postmark-go"
)

type Client struct {
	client             *postmark.Client
	token              string
	fromEmail, toEmail string
}

func NewClient(token string, toEmail, fromEmail string) (*Client, error) {
	if token == "" || toEmail == "" || fromEmail == "" {
		return nil, errors.New("token, toEmail or fromEmail can't be empty on client initialization")
	}
	auth := &http.Client{
		Transport: &postmark.AuthTransport{Token: token},
	}
	client := postmark.NewClient(auth)
	return &Client{
		client:    client,
		token:     token,
		fromEmail: fromEmail,
		toEmail:   toEmail,
	}, nil
}

func (c *Client) Send(ctx context.Context, subject string, title string, variables []TemplateVariables) error {
	emailReq := &postmark.Email{
		From:       c.fromEmail,
		To:         c.toEmail,
		TemplateID: 30755957,
		TemplateModel: map[string]interface{}{
			"albums":  variables,
			"title":   title,
			"subject": subject,
		},
		TrackOpens: true,
	}

	_, _, err := c.client.Email.Send(emailReq)
	if err != nil {
		return err
	}
	return nil
}

// TemplateVariables contains all the variables we can set in the email template
type TemplateVariables struct {
	Artist      string   `json:"artist"`
	Album       string   `json:"album"`
	ReleaseYear int      `json:"release_year"`
	Tags        []string `json:"tags"`
	ArtworkURL  string   `json:"artwork_url"`
	// URL is the source URL
	URL string `json:"url"`
	// URLs contains the streaming service URLs
	URLs []StreamingURL `json:"urls"`
}

type StreamingURL struct {
	URL      string `json:"url"`
	LinkType string `json:"link_type"`
}
