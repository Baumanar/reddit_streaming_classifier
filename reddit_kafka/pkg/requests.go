package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ReqType string

const(
	Submission       ReqType = "submission"
	SubredditComment ReqType = "subreddit_comments"
)

type Request struct {
	ReqType   ReqType
	SubReddit string
	Method    string
	Path      string
	Payload   map[string]string
}


func (c RedditClient) Request(request Request) ([]byte, error){
	values := "?"
	for i, v := range request.Payload {
		v = url.QueryEscape(v)
		values += fmt.Sprintf("%s=%s&", i, v)
	}
	values = values[:len(values)-1]
	req, err := http.NewRequest(request.Method, request.Path+values, nil)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "bearer "+c.Token)

	resp, err := c.Client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Request done for %s,    %s\n", request.ReqType, request.SubReddit)

	//c.checkRateLimit(resp)

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err := findRedditError(data); err != nil {
		return nil, err
	}
	return data, nil
}


func (c *RedditClient) checkRateLimit(response *http.Response){
	//rateLimitUsed, _ := strconv.Atoi(response.Header.Get("X-Ratelimit-Used"))
	temp, err := strconv.ParseFloat(response.Header.Get("X-Ratelimit-Remaining"), 64)
	if err != nil{
		log.Fatal(err)
	}
	rateLimitRemaining := int(temp)
	rateLimitReset, _ := strconv.Atoi(response.Header.Get("X-Ratelimit-Reset"))
	if rateLimitRemaining <= 10 {
		c.Stream.RateLimit.Add(1)
		fmt.Printf("Too many request, waiting %ds\n", rateLimitReset)
		time.Sleep(time.Duration(rateLimitReset)  * time.Second)
		c.Stream.RateLimit.Done()
	}

}


// RedditErr is a struct to store reddit error messages
type RedditErr struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func findRedditError(data []byte) error {
	object := &RedditErr{}
	json.Unmarshal(data, object)
	if object.Message != "" || object.Error != "" {
		return fmt.Errorf("%s | error code: %s", object.Message, object.Error)
	}
	return nil
}

func (c *RedditClient) GetSubredditSubmissions(subreddit string, sort string, tdur string, limit int)  ([]api_models.Submission, error){
	target := RedditOauth + "/r/" + subreddit + "/" + sort + ".json"
	req := Request{
		ReqType:   Submission,
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit": strconv.Itoa(limit),
			"t":     tdur,
			},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret api_models.SubmissionListing
	err = json.Unmarshal(ans, &ret)
	if err != nil {
		return nil, err
	}
	sumbmissions, err := ret.UnwrapData()
	if err != nil {
		return nil, err
	}
	return sumbmissions, nil
}

func (c *RedditClient) GetSubredditSubmissionsAfter(subreddit string, last string, limit int) ([]api_models.Submission, error){

	target := RedditOauth + "/r/" + subreddit + "/new.json"
	req := Request{
		ReqType:   Submission,
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit": strconv.Itoa(limit),
			"before":     last,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret api_models.SubmissionListing
	err = json.Unmarshal(ans, &ret)
	if err != nil {
		return nil, err
	}
	sumbmissions, err := ret.UnwrapData()
	if err != nil {
		return nil, err
	}
	return sumbmissions, nil
}

func (c *RedditClient) GetSubredditComments(subreddit string, sort string, tdur string, limit int) ([]api_models.Comment, error){
	target := RedditOauth + "/r/" + subreddit + "/comments.json"
	req := Request{
		ReqType:   SubredditComment,
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit": strconv.Itoa(limit),
			"sort":     sort,
			"t":     tdur,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret api_models.CommentListing
	err = json.Unmarshal(ans, &ret)
	if err != nil {
		return nil, err
	}
	sumbmissions, err := ret.UnwrapData()
	if err != nil {
		return nil, err
	}
	return sumbmissions, nil
}


func (c *RedditClient) GetSubredditCommentsAfter(subreddit string, sort string, last string, limit int) ([]api_models.Comment, error) {
	target := RedditOauth + "/r/" + subreddit + "/comments.json"
	req := Request{
		ReqType:   SubredditComment,
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit": strconv.Itoa(limit),
			"sort":  sort,
			"before":     last,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret api_models.CommentListing
	err = json.Unmarshal(ans, &ret)
	if err != nil {
		return nil, err
	}
	sumbmissions, err := ret.UnwrapData()
	if err != nil {
		return nil, err
	}
	return sumbmissions, nil
}

