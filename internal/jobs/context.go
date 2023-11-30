package jobs

import (
	"os/exec"
	"strings"
)

type Context struct {
	RepositoryURL string
	CommitHash    string
}

func NewContext(repositoryURL string, commitHash string) Context {
	return Context{repositoryURL, commitHash}
}

func NewCurrentContext() (Context, error) {
	gitRepo, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return Context{}, err
	}

	commitHash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return Context{}, err
	}
	repo, _ := strings.CutSuffix(string(gitRepo), "\n")
	hash, _ := strings.CutSuffix(string(commitHash), "\n")

	return NewContext(string(repo), string(hash)), nil
}
