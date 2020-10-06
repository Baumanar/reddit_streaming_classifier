package pkg

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// RedditClient is the client that makes requests to the api
type RedditClient struct {
	Token            string  `json:"access_token"`
	Duration         float64 `json:"expires_in"`
	ExpirationTicker *time.Ticker
	Config           AuthConfig
	UserAgent        string
	Stream           Streaming
	Client           *http.Client
}

// Streaming is for streaming parameters
type Streaming struct {
	RateLimit sync.WaitGroup
}

// Init initializes the client
func Init(config AuthConfig) (*RedditClient, error) {
	client, err := Authenticate(&config)
	if err != nil {
		return nil, err
	}
	client.ExpirationTicker = time.NewTicker(45 * time.Minute)
	client.Stream = Streaming{}

	go client.autoRefresh()
	return client, err
}

// updateCreds gets a new auth token
func (c *RedditClient) updateCreds() {
	temp, _ := Authenticate(&c.Config)
	c.Token = temp.Token
}

// autoRefresh the client auth
func (c *RedditClient) autoRefresh() {
	for {
		select {
		case <-c.ExpirationTicker.C:
			log.Println("refresh authentication")
			c.updateCreds()
		}
	}
}
