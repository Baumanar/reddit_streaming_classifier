package pkg

import (
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"
	"log"
	"time"
)

func (r *RedditClient) StreamSubredditComments(subreddit string, refresh int) (<-chan api_models.Comment, error) {
	c := make(chan api_models.Comment, 100)
	last, _, err := r.GetCommentAnchor(subreddit)
	if err != nil {
		log.Fatal("error at GetSubmissionAnchor ", err)
	}

	go func() {
		for {
			new, err := r.GetSubredditCommentsAfter(subreddit, "new", *last, 100)
			if err != nil {
				log.Fatal("error at GetSubredditCommentsAfter ", err)
			}
			if len(new) < 1 {
				log.Printf("No new comment found in: %s %s, sleeping for %ds\n", subreddit, *last, refresh)
				isDeletedComment, err := r.IsDeletedComment(*last)
				if isDeletedComment{
					log.Printf("last comment got deleted, updating anchor: %s %s", subreddit, *last)
					last, _, err = r.GetCommentAnchor(subreddit)
					if err != nil {
						log.Fatal("error at GetSubmissionAnchor in nested loop ", err)
					}
				}
				time.Sleep(time.Duration(refresh) * time.Second)
				continue
			} else {
				log.Printf("Found %d new comments in: %s %s\n", len(new), subreddit, *last)
			}
			last = &new[0].Name
			for _, v := range new {
				c <- v
			}
			time.Sleep(time.Duration(refresh) * time.Second)
			// Wait if the request rate is above the limit rate
			r.Stream.RateLimit.Wait()
		}
	}()
	return c, nil
}

func (r *RedditClient) StreamSubredditSubmissions(subreddit string, sort string, refresh int) (<-chan api_models.Submission, error) {
	c := make(chan api_models.Submission, 100)
	last, lastId, err := r.GetSubmissionAnchor(subreddit, sort)
	if err != nil {
		log.Fatal("error at GetSubmissionAnchor ", err)
	}
	go func() {
		for {
			new, err := r.GetSubredditSubmissionsAfter(subreddit, *last, 100)
			if err != nil {
				log.Fatal("error at GetSubredditSubmissionsAfter ", err)
			}
			if len(new) < 1 {
				log.Printf("No new submission found in: %s %s, sleeping for %ds\n", subreddit, *last, refresh)
				isDeletedSub, err := r.IsDeletedSubmission(subreddit, *lastId)
				if isDeletedSub{
					log.Printf("last submission got deleted, updating anchor: %s %s", subreddit, *lastId)
					last, lastId, err = r.GetSubmissionAnchor(subreddit, sort)
					if err != nil {
						log.Fatal("error at GetSubmissionAnchor in nested loop ", err)
					}
				}
				time.Sleep(time.Duration(refresh) * time.Second)
				continue
			} else {
				log.Printf("Found %d new submissions in: %s %s\n", len(new), subreddit, *last)
			}
			last = &new[0].Name
			lastId = &new[0].ID

			for i := range new {
				c <- new[len(new)-i-1]
			}
			time.Sleep(time.Duration(refresh) * time.Second)
			// Wait if the request rate is above the limit rate
			r.Stream.RateLimit.Wait()
		}
	}()
	return c, nil
}
