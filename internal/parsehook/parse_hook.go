//Package parsehook parses a github-webhook. Push-Events and Pull-Request-Events are handled.
package parsehook

import (
	"context"
	"fmt"
	"main/internal/authentication"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	githubWebhook "gopkg.in/go-playground/webhooks.v5/github"
)

type PushResult struct {
	UserName      string
	OwnerName     string
	RepoName      string
	AddedFiles    []string
	ModifiedFiles []string
	Sha           string
	Branch        string
}

type PullResult struct {
	UserName  string
	RepoName  string
	OwnerName string
	Branch    string
	Sha       string
	Number    int64
}

type ParseResult struct {
	Event string
	Push  *PushResult
	Pull  *PullResult
}

//parseHookPullRequest gets a githubWebhook.PullRequestPayload and checks for .yml and .yaml-files
func parseHookPullRequest(payload githubWebhook.PullRequestPayload, client *authentication.Client) (*PullResult, error) { //TODO sollte Pointer zurückgeben, damit nil returnt werden kann
	var result PullResult
	if payload.Action == "opened" ||
		payload.Action == "edited" ||
		payload.Action == "reopened" ||
		payload.Action == "synchronized" {

		result.UserName = payload.PullRequest.User.Login
		result.OwnerName = payload.PullRequest.Head.Repo.Owner.Login
		result.RepoName = payload.PullRequest.Head.Repo.Name
		result.Branch = payload.PullRequest.Head.Ref
		result.Sha = payload.PullRequest.Head.Sha
		result.Number = payload.Number

		//client wird von oben durchgereicht
		//var ghClient := authentication.CreateClient()
		//githubClient := authentication.GetGithubClient()

		var options = github.ListOptions{}

		files, response, err := client.GithubClient.PullRequests.ListFiles(context.Background(), payload.PullRequest.Head.Repo.Owner.Login, payload.Repository.Name, int(payload.Number), &options)
		if err != nil {
			fmt.Println("Error while getting filenames:\n", err, "\n", response)
			return nil, err
		}

		for _, file := range files {
			if strings.Contains(*file.Filename, ".yml") || strings.Contains(*file.Filename, "yaml") {
				fmt.Println("Found modified or added yamls in PullRequest.")
				return &result, nil
			}
		}
	}
	return nil, nil
}

//parseHookPush gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
func parseHookPush(payload githubWebhook.PushPayload, client *authentication.Client) (*PushResult, error) {
	var result = PushResult{}
	modifiedFilenames := lookForYamlInArray(payload.HeadCommit.Modified)
	addedFilenames := lookForYamlInArray(payload.HeadCommit.Added)
	if len(modifiedFilenames) == 0 && len(addedFilenames) == 0 {
		return &result, nil
	} else {
		commitSha := payload.HeadCommit.ID
		branchRef := *&payload.Ref

		result.AddedFiles = addedFilenames
		result.ModifiedFiles = modifiedFilenames
		result.RepoName = payload.Repository.Name
		result.OwnerName = payload.Repository.Owner.Login
		result.UserName = payload.Pusher.Name
		result.Sha = commitSha
		result.Branch = branchRef

		return &result, nil
	}
}

//lookForYamlInArray looks for .yaml or .yml-files, adds them to a string-array and returns it.
func lookForYamlInArray(filesInCommit []string) []string {
	var yamlFilenames []string
	for i := 0; i < len(filesInCommit); i++ {
		if strings.Contains(filesInCommit[i], ".yaml") ||
			strings.Contains(filesInCommit[i], ".yml") {
			yamlFilenames = append(yamlFilenames, filesInCommit[i])
		}
	}
	return yamlFilenames
}

//ParseHook checks the hook for githubWebhook.PushPayload or githubWebhook.PullRequestPayload
//and passes the payloads to the appropriate methods. It ultimately returns
//a list of modified files, a list of added files, and the commit-SHA.
func ParseHook(r *http.Request, secret string, client *authentication.Client) (*ParseResult, error) { //([]string, []string, string, string, PrSourceBranchInformation) {
	hook, err := githubWebhook.New(githubWebhook.Options.Secret(secret))
	if err != nil {
		return nil, err
	}

	payload, err := hook.Parse(r, githubWebhook.PushEvent, githubWebhook.PullRequestEvent)
	if err != nil {
		if err == githubWebhook.ErrEventNotFound {
			//This happens if the webhook sends an event that is not push or pull-request.
			fmt.Println("This event is neither a Push nor a Pull-request.\n", err)
		}
		return nil, err
	}

	var pushRes *PushResult
	var pullRes *PullResult

	var result ParseResult
	result.Event = "none"

	switch payload.(type) {

	case githubWebhook.PushPayload:
		fmt.Println("Receiving Push-Payload:")
		commit := payload.(githubWebhook.PushPayload)
		pushRes, err = parseHookPush(commit, client)
		if err != nil {
			return nil, err
		}
		result.Push = pushRes
		result.Event = "push"

	case githubWebhook.PullRequestPayload:
		fmt.Println("Receiving Pull-Request-Payload:")
		pullRequest := payload.(githubWebhook.PullRequestPayload)
		pullRes, err = parseHookPullRequest(pullRequest, client)
		if err != nil {
			return nil, err
		}
		result.Pull = pullRes
		result.Event = "pull"
		if result.Pull.Sha == "" {
			result.Event = "none"
		}
	}

	result.Push = pushRes
	result.Pull = pullRes
	fmt.Println("ParseResult:", result)
	return &result, nil
}
