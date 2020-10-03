package pkg

import (
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"
	"log"
	"time"
)

// StreamSubredditComments streams comments from a specified subreddit
// Returns an output channel on which new comments are sent
func (r *RedditClient) StreamSubredditComments(subreddit string, refresh int) (<-chan api_models.Comment, error) {
	// Create an output channel
	c := make(chan api_models.Comment, 100)
	// Get the latest comment and start from this point
	last, _, err := r.GetCommentAnchor(subreddit)
	if err != nil {
		log.Fatal("error at GetSubmissionAnchor: ", err)
	}
	// Start streaming
	go func() {
		for {
			// Get comments that have been posted after the last stored comment
			new, err := r.GetSubredditCommentsAfter(subreddit, "new", *last, 100)
			if err != nil {
				log.Fatal("error at GetSubredditCommentsAfter: ", err)
			}
			if len(new) < 1 {
				// No new comments
				log.Printf("No new comment found in: %s %s, sleeping for %ds\n", subreddit, *last, refresh)
				// Check if the comment has been deleted
				// If the latest comment has been deleted, the api wont be able to find newer comments and will
				// indefinitely send an empty list of new comments
				isDeletedComment, err := r.IsDeletedComment(*last)
				if isDeletedComment{
					log.Printf("last comment got deleted, updating anchor: %s %s", subreddit, *last)
					// Get latest comment sent
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
			// send new comments
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

// StreamSubredditSubmissions streams comments from a specified subreddit
// Returns an output channel on which new submissions are sent
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
