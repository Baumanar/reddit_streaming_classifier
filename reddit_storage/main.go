package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_storage/data"
	"github.com/gocql/gocql"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"time"
)

func main() {
	//Connect to cassandra cluster

	var kafkaHost = flag.String("kafka", "my-release-kafka.default.svc.cluster.local:9092", "kafka host")
	var cassHost = flag.String("cass", "cassandra.default.svc.cluster.local:9042", "cassandra host")
	var cassUser = flag.String("u", "cassandra", "cassandra user")
	var cassPassword = flag.String("p", "cassandra", "cassandra password")

	flag.Parse()

	//broker := cli.SetBroker(os.Args)
	fmt.Printf("Using Broker: %v\n--------------------------\n\n", *kafkaHost)
	fmt.Printf("Using Cassandra: %v\n--------------------------\n\n", *cassHost)

	cluster := gocql.NewCluster(*cassHost)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: *cassUser,
		Password: *cassPassword,
	}
	fmt.Println("OKAY LETS GO")
	cluster.Keyspace = "reddit_storage"
	//cluster.Consistency = gocql.Quorum

	Session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("ERR ", err)
	}
	defer Session.Close()

	// Create consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": *kafkaHost,
		"group.id":          "storageGroup"})

	if err != nil {
		log.Fatal("err", err)
	}
	err = consumer.SubscribeTopics([]string{"reddit_stream_comments", "reddit_stream_submissions", "reddit_classification"}, nil)
	if err != nil {
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

			var submissionClassifCount string
			iter = Session.Query(`SELECT count(*) FROM submissions WHERE is_classified = true ALLOW FILTERING`).Iter()
			for iter.Scan(&submissionClassifCount) {
				log.Printf("Number of submissions classified in database: %s\n", submissionClassifCount)
			}

			var commentClassifCount string
			iter = Session.Query(`SELECT count(*) FROM comments WHERE is_classified = true ALLOW FILTERING`).Iter()
			for iter.Scan(&commentClassifCount) {
				log.Printf("Number of comments classified in database: %s\n", commentClassifCount)
			}
			//iter.Scan(&submissionClassifCount) {
			//	log.Printf("Number of submissions in database: %s\n", submissionCount)
			//}

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

				err = data.UpdateClassification(classif.Type, classif.Name, true, classif.Class==1, classif.ProbaHateful, classif.ProbaNotHateful, Session)
				//err = data.CreateClassification(&classif, Session)
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
