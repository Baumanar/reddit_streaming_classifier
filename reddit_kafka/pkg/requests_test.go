package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/models"
	"log"
	"reflect"
	"testing"
)

func TestRedditClient_GetSubredditCommentsAfter(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		subreddit string
		sort      string
		last      string
		limit     int
	}
	tests := []struct {
		name            string
		args            args
		wantNumCOmments int
		wantErr         bool
	}{
		{name: "deletedComment1", args: args{
			subreddit: "memes",
			sort:      "new",
			last:      "t1_g5hfj9c",
			limit:     100,
		}, wantNumCOmments: 100, wantErr: false},
		{name: "deletedComment2", args: args{
			subreddit: "comments",
			sort:      "new",
			last:      "t1_g5ki4kk",
			limit:     100,
		}, wantNumCOmments: 100, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}

			got, err := c.GetSubredditCommentsAfter(tt.args.subreddit, tt.args.sort, tt.args.last, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubredditCommentsAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(tt.args.last, len(got))
		})
	}
}

func TestRedditClient_GetSubredditSubmissionsAfter(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		subreddit string
		sort      string
		last      string
		limit     int
	}
	tests := []struct {
		name       string
		args       args
		wantDelete bool
	}{
		{name: "deletedSub1", args: args{
			subreddit: "politics",
			last:      "t3_ityvmg",
			limit:     100,
		}, wantDelete: true},
		{name: "deletedSub2", args: args{
			subreddit: "politics",
			last:      "t3_itvxvd",
			limit:     100,
		}, wantDelete: false},
		{name: "deletedSub3", args: args{
			subreddit: "politics",
			last:      "t3_iu0k3b",
			limit:     100,
		}, wantDelete: true},
		{name: "deletedSub3", args: args{
			subreddit: "politics",
			last:      "t3_iufpej",
			limit:     100,
		}, wantDelete: true},
		{name: "okSub1", args: args{
			subreddit: "memes",
			last:      "t3_itx0r1",
			limit:     100,
		}, wantDelete: true},
		{name: "okSub2", args: args{
			subreddit: "memes",
			last:      "t3_itydkn",
			limit:     100,
		}, wantDelete: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}

			got, err := c.GetSubredditSubmissionsAfter(tt.args.subreddit, tt.args.last, tt.args.limit)
			if len(got) > 0 == tt.wantDelete {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantDelete)
			}
			fmt.Println(tt.args.last, len(got))
		})
	}
}

func TestRedditClient_Request(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		request Request
	}
	tests := []struct {
		name        string
		args        args
		wantRemoved string
	}{
		{name: "reqDeletedModSubmission1", args: args{request: Request{
			SubReddit: "politics",
			Method:    "GET",
			Path:      redditOauth + "/r/politics/comments/itx07n.json",
			Payload:   nil,
		}}, wantRemoved: "deleted"},

		{name: "reqDeletedModSubmission2", args: args{request: Request{
			SubReddit: "memes",
			Method:    "GET",
			Path:      redditOauth + "/r/memes/comments/itxkhf.json",
			Payload:   nil,
		}}, wantRemoved: "deleted"},

		{name: "reqDeletedModSubmission2", args: args{request: Request{
			SubReddit: "memes",
			Method:    "GET",
			Path:      redditOauth + "/r/politics/comments/iufpej.json",
			Payload:   nil,
		}}, wantRemoved: "deleted"},

		{name: "reqOkSubmission", args: args{request: Request{
			SubReddit: "memes",
			Method:    "GET",
			Path:      redditOauth + "/r/memes/comments/itxf7n.json",
			Payload:   nil,
		}}, wantRemoved: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}

			got, err := c.Request(tt.args.request)
			var ret []models.SubmissionListing
			err = json.Unmarshal(got, &ret)
			if err != nil {
				log.Fatal(err)
			}
			sumbmissionArr, err := ret[0].UnwrapData()
			removedByCategory := sumbmissionArr[0].RemovedByCategory
			fmt.Println(removedByCategory)
			if removedByCategory != tt.wantRemoved {
				t.Errorf("Request() got = %v, want %v", removedByCategory, tt.wantRemoved)
			}
		})
	}
}

func TestRedditClient_GetSubmission(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		subreddit string
		id        string
	}

	tests := []struct {
		name        string
		args        args
		wantRemoved string
	}{
		{name: "testDeletedSub", args: args{
			subreddit: "politics",
			id:        "iufpej",
		}, wantRemoved: "deleted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}

			got, err := c.GetSubmission(tt.args.subreddit, tt.args.id)
			if err != nil {
				t.Errorf("Request() error = %v", err)
			}
			if !reflect.DeepEqual(got.ID, tt.args.id) {
				t.Errorf("GetSubmission() = %v, want %v", got, tt.args.id)
			}
			if !reflect.DeepEqual(got.RemovedByCategory, tt.wantRemoved) {
				t.Errorf("GetSubmission() = %v, want %v", got, tt.wantRemoved)
			}
		})
	}
}

func TestRedditClient_GetComment(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		subreddit string
		id        string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "testDelComment1", args: args{
			id: "t1_g5hfz2w",
		}},
		{name: "testDelComment2", args: args{
			id: "t1_g5ki4kk",
		}},
		{name: "testDelComment2", args: args{
			id: "t1_g5kiz50",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}
			comment, err := c.GetComment(tt.args.id)
			//fmt.Println(comment.Subreddit)
			//fmt.Println(comment.Body)
			//fmt.Printf("%+v\n", comment)
			if !reflect.DeepEqual(comment.Name, tt.args.id) {
				t.Errorf("GetSubmission() = %v, want %v", comment.Name, tt.args.id)
			}
		})
	}
}

func TestRedditClient_IsDeletedComment(t *testing.T) {
	authConfig := GetConfigByFile("../auth.conf")

	type args struct {
		subreddit string
		id        string
	}
	tests := []struct {
		name      string
		args      args
		isDeleted bool
	}{
		{name: "testDelComment1", args: args{
			id: "t1_g5hfz2w",
		}, isDeleted: false},
		{name: "testDelComment2", args: args{
			id: "t1_g5ki4kk",
		}, isDeleted: true},
		{name: "testDelComment2", args: args{
			id: "t1_g5kiz50",
		}, isDeleted: false},
		{name: "testDelComment2", args: args{
			id: "t1_g5km26m",
		}, isDeleted: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RedditClient{}
			c, err := Init(authConfig)
			if err != nil {
				panic(err)
			}
			gotDeleted, err := c.IsDeletedComment(tt.args.id)
			if !reflect.DeepEqual(gotDeleted, tt.isDeleted) {
				t.Errorf("GetSubmission() = %v, want %v", gotDeleted, tt.isDeleted)
			}
		})
	}
}
