package pkg

import "github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"

func MergeSubmissionChannels(channels ...<-chan api_models.Submission) chan api_models.Submission {
	outChan := make(chan api_models.Submission)
	for _, ch := range channels {
		go func(ch <-chan api_models.Submission) {
			for v := range ch {
				outChan <- v
			}
		}(ch)
	}
	return outChan
}

func MergeCommentChannels(channels ...<-chan api_models.Comment) chan api_models.Comment {
	outChan := make(chan api_models.Comment)
	for _, ch := range channels {
		go func(ch <-chan api_models.Comment) {
			for v := range ch {
				outChan <- v
			}
		}(ch)
	}
	return outChan
}
