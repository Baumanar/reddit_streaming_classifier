package pkg

import "github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/models"

// MergeSubmissionChannels merges multiple input submissions channels and returns a single submission channel
func MergeSubmissionChannels(channels ...<-chan models.Submission) chan models.Submission {
	outChan := make(chan models.Submission)
	for _, ch := range channels {
		go func(ch <-chan models.Submission) {
			for v := range ch {
				outChan <- v
			}
		}(ch)
	}
	return outChan
}

// MergeCommentChannels merges multiple input comments channels and returns a single comment channel
func MergeCommentChannels(channels ...<-chan models.Comment) chan models.Comment {
	outChan := make(chan models.Comment)
	for _, ch := range channels {
		go func(ch <-chan models.Comment) {
			for v := range ch {
				outChan <- v
			}
		}(ch)
	}
	return outChan
}
