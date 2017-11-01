package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/containerum/solutions"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func gitClone(repoUrl, branch, destDir string) error {
	_, err := git.PlainClone(destDir, false, &git.CloneOptions{
		URL:           repoUrl,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Depth:         1,
	})
	return err
}

func githubDownload(user, repo, branch, destDir string, files []string) error {
	for _, file := range files {
		resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent/%s/%s/%s/%s", user, repo, branch, file))
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(destDir, file), content, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func fetchFiles(name, branch, destDir string, files []string) error {
	nameItems := strings.Split(name, "/")
	switch len(nameItems) {
	case 1: //containerum solution
		return githubDownload("containerum", name, branch, destDir, files)
	case 2: //3rd party solution on github
		return githubDownload(nameItems[0], nameItems[1], branch, destDir, files)
	default: //3rd party solution on any git hosting
		return gitClone(name, branch, destDir)
	}
	return nil
}

func DownloadSolution(name, solutionPath, branch string) error {
	if err := os.MkdirAll(solutionPath, os.ModePerm); err != nil {
		return err
	}

	// download config
	if err := fetchFiles(name, branch, solutionPath, []string{solutions.SolutionConfigFile}); err != nil {
		return err
	}

	// parse and download template files
	cfgFile, err := ioutil.ReadFile(path.Join(solutionPath, solutions.SolutionConfigFile))
	if err != nil {
		return err
	}
	var cfgObj solutions.SolutionConfig
	if err := json.Unmarshal(cfgFile, &cfgObj); err != nil {
		return err
	}

	var files []string
	for _, v := range cfgObj.Run {
		files = append(files, v.ConfigFile)
	}

	return fetchFiles(name, branch, solutionPath, files)
}
