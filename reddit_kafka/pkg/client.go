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
	RateLimit sync.WaitGroup
}



func Init(config AuthConfig) (*RedditClient, error){
	client, err := Authenticate(config)
	if err != nil{
		return nil, err
	}
	client.ExpirationTicker = time.NewTicker(time.Second*time.Duration(client.Duration*0.75))
	client.Stream = Streaming{
		CommentListInterval: 1,
		PostListInterval:    1,
		PostListSlice:       1,
	}

	go func(config AuthConfig) {
		time.Sleep(time.Second*20)
		select{
		case t := <- client.ExpirationTicker.C:
			log.Printf("Client refreshed authentication at %s", t)
			temp, err := Authenticate(config)
			if err != nil{
				log.Fatal(err)
			}
			client.Token = temp.Token
		}
	}(client.Config)
	return client, err
}

