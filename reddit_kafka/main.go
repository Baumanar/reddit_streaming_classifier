package main

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_streaming_classifier/cli"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/pkg"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
)

func main() {

	// Authenticate with the reddit client
	authConfig := pkg.GetConfigByFile("auth.conf")
	client := &pkg.RedditClient{}
	client, err := pkg.Init(authConfig)
	if err != nil {
		panic(err)
	}
	// Stream comments and submissions from some subreddits
	c1Sub, _ := client.StreamSubredditSubmissions("politics", "new", 30)
	c2Sub, _ := client.StreamSubredditSubmissions("memes", "new", 10)
	c3Sub, _ := client.StreamSubredditSubmissions("funny", "new", 20)
	c4Sub, _ := client.StreamSubredditSubmissions("gaming", "new", 20)
	c5Sub, _ := client.StreamSubredditSubmissions("movies", "new", 20)

	c1Com, _ := client.StreamSubredditComments("memes", 10)
	c2Com, _ := client.StreamSubredditComments("politics", 10)
	c3Com, _ := client.StreamSubredditComments("funny", 10)
	c4Com, _ := client.StreamSubredditComments("gaming", 10)
	c5Com, _ := client.StreamSubredditComments("movies", 10)

	// Merges all messages in two channels
	submissionChan := pkg.MergeSubmissionChannels(c1Sub, c2Sub, c3Sub, c4Sub, c5Sub)
	commentChan := pkg.MergeCommentChannels(c1Com, c2Com, c3Com, c4Com, c5Com)

	broker := cli.SetBroker(os.Args)
	fmt.Printf("Using Broker: %v\n--------------------------\n\n", broker)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker.String()})
	if err != nil {
		log.Fatal(err)
	}

	defer p.Close()

	submissionTopic := "reddit_stream_submissions"
	commentTopic := "reddit_stream_comments"

	// Send new submissions and comments to kafka broker
	for {
		for {
			select {
			case submission := <-submissionChan:
				payload, err := json.Marshal(submission)
				if err != nil {
					log.Fatal(err)
				}
				err = p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &submissionTopic, Partition: kafka.PartitionAny},
					Value:          payload,
				}, nil)
				//log.Printf("Sent a submission: %s\n", submission.Subreddit)
				if err != nil {
					log.Fatal(err)
				}

			case comment := <-commentChan:
				payload, err := json.Marshal(comment)
				if err != nil {
					log.Fatal(err)
				}

				err = p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &commentTopic, Partition: kafka.PartitionAny},
					Value:          payload,
				}, nil)
				//log.Printf("Sent a comment: %s\n", comment.Subreddit)

				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}

}
