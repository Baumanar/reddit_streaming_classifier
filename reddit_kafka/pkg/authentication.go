package pkg

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	tokenURL    = "https://www.reddit.com/api/v1/access_token"
	redditOauth = "https://oauth.reddit.com"
)

// AuthConfig is the struct storing the reddit credentials
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	UserAgent    string
	Username     string
	Password     string
}

// GetConfigByEnv gets the reddit credentials from the env vars
func GetConfigByEnv() AuthConfig {
	return AuthConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		UserAgent:    os.Getenv("USER_AGENT"),
		Username:     os.Getenv("USERNAME"),
		Password:     os.Getenv("PASSWORD"),
	}
}

// GetConfigByFile gets the reddit credentials from a config file
func GetConfigByFile(filePath string) AuthConfig {
	ClientID, _ := regexp.Compile(`CLIENT_ID\s*=\s*(.+)`)
	ClientSecret, _ := regexp.Compile(`CLIENT_SECRET\s*=\s*(.+)`)
	Username, _ := regexp.Compile(`USERNAME\s*=\s*(.+)`)
	Password, _ := regexp.Compile(`PASSWORD\s*=\s*(.+)`)
	UserAgent, _ := regexp.Compile(`USER_AGENT\s*=\s*(.+)`)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return AuthConfig{}
	}
	s := string(data)
	creds := AuthConfig{
		ClientID.FindStringSubmatch(s)[1],
		ClientSecret.FindStringSubmatch(s)[1],
		UserAgent.FindStringSubmatch(s)[1],
		Username.FindStringSubmatch(s)[1],
		Password.FindStringSubmatch(s)[1],
	}
	return creds
}

// Authenticate creates a new RedditClient from the provided config credentials
func Authenticate(config *AuthConfig) (*RedditClient, error) {

	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", config.Username)
	form.Add("password", config.Password)

	raw := config.ClientID + ":" + config.ClientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Authorization", "Basic "+encoded)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	auth := RedditClient{}
	// Unmarshall the body response in the client struct
	json.Unmarshal(body, &auth)
	auth.Client = client
	auth.Config = *config
	auth.UserAgent = config.UserAgent
	return &auth, err
}
