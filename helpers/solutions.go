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
	"github.com/libgit2/git2go"
)

func gitFetch(remoteUrl, destDir string) error {
	repo, err := git.InitRepository(destDir, true)
	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Create("origin", remoteUrl)
	if err != nil {
		return err
	}

	return remote.Fetch(nil, nil, "")
}

func gitCheckout(checkoutTarget, destDir string, files []string) error {
	repo, err := git.OpenRepository(destDir)
	if err != nil {
		return err
	}

	ref, err := repo.References.Lookup(checkoutTarget)
	if err != nil {
		return err
	}

	tree, err := repo.LookupTree(ref.Target())
	if err != nil {
		return err
	}

	return repo.CheckoutTree(tree, &git.CheckoutOpts{Paths: files})
}

func isGitRepo(destDir string) bool {
	_, err := git.OpenRepository(destDir)
	return err != nil
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
		if !isGitRepo(destDir) {
			if err := gitFetch(name, destDir); err != nil {
				return err
			}
		}

		if err := gitCheckout(branch, destDir, files); err != nil {
			return err
		}
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
