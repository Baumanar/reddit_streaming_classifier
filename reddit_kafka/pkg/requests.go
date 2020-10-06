package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Request holds info for a request to the reddit api
type Request struct {
	SubReddit string            // subreddit name
	Method    string            // http method
	Path      string            // url path
	Payload   map[string]string // request parameters
}

// Request makes a request to the reddit api
// returns the body of the response
func (c RedditClient) Request(request Request) ([]byte, error) {
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

	c.checkRateLimit(resp)

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err := findRedditError(data); err != nil {
		return nil, err
	}
	return data, nil
}

// checkRateLimit checks in the response headers how many remaining requests we can do in this time interval
func (c *RedditClient) checkRateLimit(response *http.Response) {
	temp, err := strconv.ParseFloat(response.Header.Get("X-Ratelimit-Remaining"), 64)
	if err != nil {
		log.Println("checkRateLimit failed while parsing X-Ratelimit-Remaining, pass...")
		return
	}
	rateLimitRemaining := int(temp)
	rateLimitReset, _ := strconv.Atoi(response.Header.Get("X-Ratelimit-Reset"))
	// If the rateLimit is dangerously low, block all streams by using the waitGroup
	if rateLimitRemaining <= 10 {
		c.Stream.RateLimit.Add(1)
		fmt.Printf("Too many request, waiting %ds\n", rateLimitReset)
		// Wait for the interval to reset
		time.Sleep(time.Duration(rateLimitReset) * time.Second)
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

// GetSubredditSubmissions sends a request for submissions in specified subreddit, sort and limit (number of items)
func (c *RedditClient) GetSubredditSubmissions(subreddit string, sort string, tdur string, limit int) ([]models.Submission, error) {
	target := redditOauth + "/r/" + subreddit + "/" + sort + ".json"
	// Create associated request struct
	req := Request{
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
	var ret models.SubmissionListing
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

// GetSubredditSubmissionsAfter sends a request for submissions after last submission (full name) in specified subreddit,
// sort and limit (number of items)
func (c *RedditClient) GetSubredditSubmissionsAfter(subreddit string, last string, limit int) ([]models.Submission, error) {
	target := redditOauth + "/r/" + subreddit + "/new.json"
	// Create associated request struct
	req := Request{
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit":  strconv.Itoa(limit),
			"before": last,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret models.SubmissionListing
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

// GetSubredditComments sends a request for comments in specified subreddit, sort and limit
func (c *RedditClient) GetSubredditComments(subreddit string, sort string, tdur string, limit int) ([]models.Comment, error) {
	target := redditOauth + "/r/" + subreddit + "/comments.json"
	// Create associated request struct
	req := Request{
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit": strconv.Itoa(limit),
			"sort":  sort,
			"t":     tdur,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret models.CommentListing
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

// GetSubredditCommentsAfter sends a request for comments after last commment in specified subreddit, sort and limit
func (c *RedditClient) GetSubredditCommentsAfter(subreddit string, sort string, last string, limit int) ([]models.Comment, error) {
	target := redditOauth + "/r/" + subreddit + "/comments.json"
	req := Request{
		SubReddit: subreddit,
		Method:    "GET",
		Path:      target,
		Payload: map[string]string{
			"limit":  strconv.Itoa(limit),
			"sort":   sort,
			"before": last,
		},
	}
	ans, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	var ret models.CommentListing
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

// GetSubmission sends a request to get a specific submission
func (c *RedditClient) GetSubmission(subreddit string, id string) (*models.Submission, error) {
	got, err := c.Request(Request{
		SubReddit: subreddit,
		Method:    "GET",
		Path:      redditOauth + "/r/" + subreddit + "/comments/" + id + ".json",
		Payload:   nil,
	})
	if err != nil {
		log.Fatal(err)
	}
	var ret []models.SubmissionListing
	err = json.Unmarshal(got, &ret)
	if err != nil {
		return nil, err
	}
	sumbmissionArr, err := ret[0].UnwrapData()
	return &sumbmissionArr[0], nil
}

// GetComment sends a request to get a specific comment
func (c *RedditClient) GetComment(id string) (*models.Comment, error) {
	got, err := c.Request(Request{
		SubReddit: "",
		Method:    "GET",
		Path:      redditOauth + "/api/info",
		Payload: map[string]string{
			"id": id,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	var ret models.CommentListing
	err = json.Unmarshal(got, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Data.Children[0].Data, nil
}

// IsDeletedComment checks if comment is deleted
func (c *RedditClient) IsDeletedComment(id string) (bool, error) {
	comment, err := c.GetComment(id)
	if err != nil {
		return false, err
	}
	matched, err := regexp.MatchString(`\[effacé\]|\[deleted\]`, comment.Body)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// IsDeletedSubmission checks if submission is deleted
func (c *RedditClient) IsDeletedSubmission(subreddit string, id string) (bool, error) {
	submission, err := c.GetSubmission(subreddit, id)
	if err != nil {
		return false, err
	}
	return submission.RemovedByCategory != "", nil
}

// GetSubmissionAnchor gets latest submission in specified subreddit, returns its full Name (lastID),
// its id, and an error
func (c *RedditClient) GetSubmissionAnchor(subreddit string, sort string) (last *string, lastID *string, err error) {
	anchor, err := c.GetSubredditSubmissions(subreddit, sort, "hour", 1)
	if err != nil {
		return nil, nil, err
	}
	if len(anchor) > 0 {
		last = &anchor[0].Name
		lastID = &anchor[0].ID
	}
	return last, lastID, err
}

// GetCommentAnchor gets latest comment in specified subreddit, returns its full Name (lastID),
// its id, and an error
func (c *RedditClient) GetCommentAnchor(subreddit string) (last *string, lastID *string, err error) {
	anchor, err := c.GetSubredditComments(subreddit, "new", "hour", 1)
	if err != nil {
		return nil, nil, err
	}
	if len(anchor) > 0 {
		last = &anchor[0].Name
		lastID = &anchor[0].ID
	}
	return last, lastID, err
}
