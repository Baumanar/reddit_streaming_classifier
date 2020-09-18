package pkg

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type RedditClient struct {
	Token            string  `json:"access_token"`
	Duration         float64 `json:"expires_in"`
	ExpirationTicker *time.Ticker
	Config           AuthConfig
	UserAgent        string
	Stream           Streaming
	Client           *http.Client
}

type Streaming struct {
	CommentListInterval int
	PostListInterval    int
	PostListSlice       int
	RateLimit           sync.WaitGroup
}

func Init(config AuthConfig) (*RedditClient, error) {
	client, err := Authenticate(&config)
	if err != nil {
		return nil, err
	}
	client.ExpirationTicker = time.NewTicker(45 * time.Minute)
	client.Stream = Streaming{
		CommentListInterval: 1,
		PostListInterval:    1,
		PostListSlice:       1,
	}
	
	go client.auto_refresh()
	return client, err
}

func (c *RedditClient) update_creds(){
	temp, _ := Authenticate(&c.Config)
	c.Token = temp.Token
}


func (c *RedditClient) auto_refresh() {
	for {
		select {
		case <- c.ExpirationTicker.C:
			log.Println("refresh authentication")
			c.update_creds()
		}
	}
}