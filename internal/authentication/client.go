// Package authentication is responsible for registering KubeLinterBot to a github-Repository.
// It also handles functions related to the oauth-token like serializing it or reading it again.
package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client lets KLB login to github.
type Client struct {
	jsonToken           string
	personalAccessToken string
	oauthClient         *http.Client
	GithubClient        *github.Client
}

// GetToken returns the access-token as a string, without bearer and expiry.
func (ao *Client) GetToken() string {
	return ao.personalAccessToken
}

// GetJSONToken returns the access-token as a JSON-string, with bearer and expiry.
func (ao *Client) GetJSONToken() string {
	return ao.jsonToken
}

// SetJSONToken sets the access-token as a JSON-string, with bearer and expiry.
func (ao *Client) SetJSONToken(token string) {
	ao.jsonToken = token
}

func (ao *Client) getGithubClient() *github.Client {
	return ao.GithubClient
}

func (ao *Client) getOAuthClient() *http.Client {
	return ao.oauthClient
}

// TokenSource must be implemented for oauth.
type TokenSource struct {
	AccessToken string
}

// Token must be implemented for oauth.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// CreateClient creates and returns the Client-Object needed to login to github.
func CreateClient(token string) *Client {
	var c Client
	c.personalAccessToken = extractTokenStringFromJSONToken(token)

	tokenSource := &TokenSource{
		AccessToken: c.personalAccessToken,
	}

	c.oauthClient = oauth2.NewClient(context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{}), tokenSource)
	c.GithubClient = github.NewClient(c.oauthClient)
	//fmt.Println("Client:", c)
	return &c
}

type jsonTokenStruct struct {
	Access_token string
	Token_type   string
	Expiry       string
}

func extractTokenStringFromJSONToken(completeToken string) string {
	var tokenStruct jsonTokenStruct
	json.Unmarshal([]byte(completeToken), &tokenStruct)
	var tokenString string = tokenStruct.Access_token
	if tokenString == "" {
		return completeToken
	} else {
		return tokenString
	}
}
