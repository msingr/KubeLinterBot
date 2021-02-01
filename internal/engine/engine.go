package engine

import (
	"fmt"
	"main/internal/authentication"
	"main/internal/callkubelinter"
	"main/internal/config"
	"main/internal/getcommit"
	"main/internal/handleresult"
	"main/internal/parsehook"
	"net/http"
)

type analysisEngine struct{}

func GetEngine() *analysisEngine {
	var ae analysisEngine
	return &ae
}

func (ae *analysisEngine) Analyse(r *http.Request, cfg config.Config) {
	fmt.Println("test")
	// var added []string
	// var modified []string
	var commitSha string
	var token string = cfg.User.AccessToken
	client := authentication.CreateClient(token) //TODO

	// var ownerName = cfg.Repositories[0].Owner
	// var repoName = cfg.Repositories[0].Name

	result, err := parsehook.ParseHook(r, cfg.Repositories[0].Webhook.Secret, client)
	if err != nil {
		fmt.Println("Error while parsing hook:\n", err)
	}

	//make prettier. commitSha should be named dl-directory or something
	if result.Event == "push" {
		commitSha = result.Push.Sha
	} else if result.Event == "pull" {
		commitSha = result.Pull.Sha
	}
	fmt.Println("ParseResult:", result)
	if result.Event != "none" {
		getcommit.GetCommit(result, *client)

		var lintResult, exitCode = callkubelinter.CallKubelinter()
		handleresult.Handle(result, lintResult, exitCode, commitSha, client)
	} else {
		fmt.Println("No need to lint, as no .yml or .yaml were changed.\nKubeLinterBot is listening for Webhooks...")
	}
}