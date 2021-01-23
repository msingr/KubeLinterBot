//Package getcommit is used to download all folders with .yaml and .yml-files.
package getcommit

import (
	"context"
	"fmt"
	"io"
	"main/internal/authentication"
	"net/http"
	"os"
	"regexp"

	"github.com/google/go-github/github"
)

// var personalAccessToken string

const mainDir = "./downloadedYaml/"

// type TokenSource struct {
// 	AccessToken string
// }

// //Token creates the oauth2.Token for oauth.
// func (t *TokenSource) Token() (*oauth2.Token, error) {
// 	token := &oauth2.Token{
// 		AccessToken: t.AccessToken,
// 	}
// 	return token, nil
// }

//DownloadCommit authenticates with oauth and downloads all folders with .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(ownername string, reponame string, commitSha string, branchRef string, addedFiles []string, modifiedFiles []string) bool {
	var downloadStatus = false

	githubClient := authentication.GetGithubClient()
	//oAuthClient := authentication.GetOAuthClient()

	_, err := downloadFolder(ownername, reponame, "", commitSha, branchRef, githubClient) //TODO path not hardcoded
	if err != nil {
		fmt.Println("Error while creating folder DownloadCommit", err)
	} // } else {
	// 	fmt.Println(folder)
	// }

	downloadStatus = true
	return downloadStatus
}

//downloadFolder downloads all files in a folder, creating subfolders as necessary.
// func downloadFolder(path string, username string, reponame string, commitSha string, client *github.Client) ([]*github.RepositoryContent, error) {
// 	//var options = github.RepositoryContentGetOptions{}
// 	// _, folder, _, err := client.Repositories.GetContents(oauth2.NoContext,
// 	// 	username,
// 	// 	reponame,
// 	// 	path, //TODO path not hardcoded
// 	// 	&options)
// 	commit, r, err := client.Repositories.GetCommit(oauth2.NoContext,
// 		username,
// 		reponame,
// 		commitSha)
// 	if err != nil {
// 		return _, err
// 	} else {
// 		for _, file := range commit.Files {
// 			if string(file.GetType()) == "dir" {
// 				err := os.MkdirAll(string(mainDir+commitSha+"/"+file.GetPath()), 0755)
// 				if err != nil {
// 					fmt.Println("Error while creating folder.", err)
// 					//return downloadStatus
// 				} else {
// 					fmt.Println("Folder created:", file.GetPath())
// 				}
// 				downloadFolder(file.GetPath(), username, reponame, commitSha, client)
// 			} else if string(file.GetType()) == "file" {
// 				downloadFile(file.GetDownloadURL(), commitSha+"/"+file.GetPath())
// 				fmt.Println(mainDir)
// 				fmt.Println(mainDir + commitSha)
// 				fmt.Println(mainDir + commitSha + "/" + file.GetPath())
// 			}
// 		}
// 		return folder, err
// 	}
// }

//downloadFolder downloads all files in a folder, creating subfolders as necessary.
func downloadFolder(ownername string, reponame string, subpath string, commitSha string, branchRef string, client *github.Client) ([]*github.RepositoryContent, error) {
	fmt.Println("subpath", subpath, "branchref:", branchRef)
	var commitDir = mainDir + commitSha + "/"

	branch, _, err := client.Repositories.GetBranch(context.Background(),
		ownername,
		reponame,
		branchRef)

	var options = github.RepositoryContentGetOptions{Ref: branch.GetName()}
	_, folder, _, err := client.Repositories.GetContents(context.Background(),
		ownername,
		reponame,
		subpath, //TODO path not hardcoded
		&options)

	errDir := os.MkdirAll(string(commitDir), 0777)
	fmt.Println("commitdir:", commitDir)
	if errDir != nil {
		fmt.Println("Error while creating folder.", err)
	} else {
		fmt.Println("Folder created:", commitDir)
	}

	if err != nil {
		return folder, err
	} else {
		for _, file := range folder {
			if string(file.GetType()) == "dir" {
				err := os.MkdirAll(string(commitDir+file.GetPath()), 0777)
				fmt.Println("commitdir+getpath", commitDir+file.GetPath())
				if err != nil {
					fmt.Println("Error while creating folder.", err)
					//return downloadStatus
				} else {
					fmt.Println("Folder created:", file.GetPath())
				}
				fmt.Println("commitdir+subpath+getpath", commitDir+file.GetPath()+"/")
				downloadFolder(ownername, reponame, file.GetPath(), commitSha, branchRef, client)
			} else if string(file.GetType()) == "file" {
				fmt.Println("filedownload: commitdir+filegetpath", commitDir+file.GetPath())
				downloadFile(file.GetDownloadURL(), commitDir+file.GetPath())
			}
		}
		return folder, err
	}
}

//downloadFile downloads a single file.
func downloadFile(url string, filepath string) error {
	fmt.Println("Downloading file: " + url + "\n")
	fmt.Println("filepath:", filepath)
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println("create", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("get", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("copy", err)
		return err
	}
	return nil
}

// //downloadFolder downloads all files in a folder, creating subfolders as necessary.
// func downloadFolderNeu(ownername string, reponame string, commitSha string, client *github.Client) (*github.RepositoryCommit, error) {
// 	fmt.Println("Ownername:", ownername)
// 	fmt.Println("reponame:", reponame)
// 	fmt.Println("commitSha:", commitSha)
// 	commit, _, err := client.Repositories.GetCommit(context.Background(),
// 		ownername,
// 		reponame,
// 		commitSha)
// 	if err != nil {
// 		fmt.Println("commit:", commit)
// 		return nil, err
// 	} else {
// 		fmt.Println("commit:", commit)
// 		for _, file := range commit.Files {
// 			fmt.Println("Datei-----:", file)
// 			path := file.GetFilename()
// 			directory := truncateFilepath(path)
// 			//mainDirectory := truncateMainpath(path)
// 			// err := os.MkdirAll(string(mainDir+commitSha+"/"+directory), 0700)
// 			// if err != nil {
// 			// 	fmt.Println("Error while creating folder.", err)
// 			// 	//return downloadStatus
// 			// } else {
// 			// 	fmt.Println("Folder created:", directory)
// 			// }
// 			var options = github.RepositoryContentGetOptions{}
// 			_, dir, _, err := client.Repositories.GetContents(context.Background(),
// 				ownername,
// 				reponame,
// 				directory, //mainDirectory, //TODO path not hardcoded
// 				&options)
// 			fmt.Println(dir)
// 			for _, f := range dir {
// 				fmt.Println(f.GetDownloadURL())
// 				if f.GetType() == "dir" {
// 					err := os.MkdirAll(string(mainDir+commitSha+"/"+f.GetPath()), 0700)
// 					if err != nil {
// 						fmt.Println("Error while creating folder.", err)
// 						//return downloadStatus
// 					}
// 				} else if f.GetType() == "file" {
// 					err3 := downloadFile(f.GetDownloadURL(), commitSha, file.GetFilename())
// 					if err3 != nil {
// 						fmt.Println("downloadFile:", err, err3)
// 					}
// 				}
// 			}

// 			// err2 := downloadFile(file.GetContentsURL(), commitSha, file.GetFilename())
// 			// if err2 != nil {
// 			// 	fmt.Println("Error while creating file.", err2)
// 			// }
// 			// if string(file.) == "dir" {
// 			// 	err := os.MkdirAll(string(mainDir+commitSha+"/"+file.GetPath()), 0755)
// 			// 	if err != nil {
// 			// 		fmt.Println("Error while creating folder.", err)
// 			// 		//return downloadStatus
// 			// 	} else {
// 			// 		fmt.Println("Folder created:", file.GetPath())
// 			// 	}
// 			// 	downloadFolder(file.GetPath(), ownername, reponame, commitSha, client)
// 			// } else if string(file.GetType()) == "file" {
// 			// 	downloadFile(file.GetDownloadURL(), commitSha+"/"+file.GetPath())
// 			// 	fmt.Println(mainDir)
// 			// 	fmt.Println(mainDir + commitSha)
// 			// 	fmt.Println(mainDir + commitSha + "/" + file.GetPath())
// 			// }
// 		}
// 		return commit, err
// 	}
// }

func truncateFilepath(path string) string {
	var regex = `^(.*[\\\/])` //`/^(.*[\\\/])/` //last slash in string
	reg := regexp.MustCompile(regex)
	index := reg.FindStringIndex(path)
	fmt.Println("Index", index)
	if index == nil {
		return "/"
	}
	var truncated string = path[index[0]:index[1]]
	fmt.Println("Filepath was:", path, "\ntruncated:", truncated)
	return truncated
}

//downloadFile downloads a single file.
func downloadFileNeu(url string, commitSha string, filename string) error {
	fmt.Println("Downloading file: " + url + "\n")
	out, err := os.Create(mainDir + commitSha + "/" + filename)
	// fmt.Println("mainDir:", mainDir)
	// fmt.Println("filename:", filename)
	if err != nil {
		// fmt.Println("create", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println("get", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// fmt.Println("copy", err)
		return err
	}
	return nil
}
