package main

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_api_streaming/pkg"
	"github.com/Baumanar/reddit_classif_stream/cli"
	channels "github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/pkg"
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
	c1_sub, err := client.StreamSubredditSubmissions("politics", "new", 10)
	c2_sub, err := client.StreamSubredditSubmissions("memes", "new", 10)
	c1_com, err := client.StreamSubredditComments("memes", 10)
	c2_com, err := client.StreamSubredditComments("politics", 10)

	// Merges all messages in two channels
	submissionChan := channels.MergeSubmissionChannels(c1_sub, c2_sub)
	commentChan := channels.MergeCommentChannels(c1_com, c2_com)

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
