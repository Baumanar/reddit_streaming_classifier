package main

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_api_streaming/api_models"
	"github.com/Baumanar/reddit_proj/reddit_storage/data"
	"github.com/gocql/gocql"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"github.com/Baumanar/reddit_proj/cli"
	"time"
)

func main() {
	//Connect to cassandra cluster
	cluster := gocql.NewCluster("cassandra:9042")
	//cluster.Authenticator = gocql.PasswordAuthenticator{
	//	Username: "cassandra",
	//	Password: "cassandra",
	//}
	cluster.Keyspace = "reddit_storage"
	cluster.Consistency = gocql.Quorum

	Session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("ERR",err)
	}
	defer Session.Close()


	broker := cli.SetBroker(os.Args)
	fmt.Printf("Using Broker: %v\n--------------------------\n\n", broker)

	// Create consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker.String(),
		"group.id":          "storageGroup"})

	if err != nil {
		log.Fatal("err", err)
	}
	err = consumer.SubscribeTopics([]string{"reddit_stream_comments", "reddit_stream_submissions", "reddit_classification"}, nil)
	if err != nil{
		log.Fatal("☠️ Uh oh, there was an error subscribing to the topic :\n\t%v", err)
	}
	//tm.Clear()

	//Get database stats every 10 seconds
	go func() {
		start := time.Now()
		for {
			//tm.MoveCursor(1, 1)
			log.Println("Uptime: ", time.Now().Sub(start).String())
			var commentCount string
			iter := Session.Query(`SELECT count(*) FROM comments`).Iter()
			for iter.Scan(&commentCount) {
				log.Printf("Number of comments in database: %s\n", commentCount)
			}

			var submissionCount string
			iter = Session.Query(`SELECT count(*) FROM submissions`).Iter()
			for iter.Scan(&submissionCount) {
				log.Printf("Number of submissions in database: %s\n", submissionCount)
			}

			var classificationCount string
			iter = Session.Query(`SELECT count(*) FROM classifications`).Iter()
			_, err = iter.RowData()
			if err != nil{
				log.Fatal(err)
			}
			for iter.Scan(&classificationCount) {
				log.Printf("Number of classifications in database: %s\n", classificationCount)
			}
			//tm.Flush()
			time.Sleep(10 * time.Second)

		}
	}()

	// Consume messages coming from kafka
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {

			if *msg.TopicPartition.Topic == "reddit_stream_comments" {

				sub := api_models.Comment{}
				err := json.Unmarshal(msg.Value, &sub)
				if err != nil {
					log.Fatal(err)
				}
				//log.Println("got com")
				err = data.CreateComment(&sub, Session)
				if err != nil {
					log.Fatal(err)
				}
			} else if *msg.TopicPartition.Topic == "reddit_stream_submissions" {
				sub := api_models.Submission{}
				err := json.Unmarshal(msg.Value, &sub)
				if err != nil {
					log.Fatal(err)
				}
				//log.Println("got sub")

				err = data.CreateSubmission(&sub, Session)
				if err != nil {
					log.Fatal(err)
				}
			} else if *msg.TopicPartition.Topic == "reddit_classification" {
				classif := data.Classification{}
				err := json.Unmarshal(msg.Value, &classif)
				if err != nil {
					log.Fatal(err)
				}
					err = data.CreateClassification(&classif, Session)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("KK error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()

}
