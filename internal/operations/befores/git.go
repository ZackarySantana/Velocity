package befores

import (
	"github.com/urfave/cli/v2"
)

type GitInfo struct {
	Owner string
	Repo  string
}

func Git(c *cli.Context) error {
	// TODO: Run command to get git info
	owner := "owner"
	repo := "repo"
	c.App.Metadata["git_owner"] = owner
	c.App.Metadata["git_repository"] = repo
	return nil
}

func GetGit(c *cli.Context) (*GitInfo, error) {
	owner, ownerOk := c.App.Metadata["git_owner"].(string)
	repo, repoOk := c.App.Metadata["git_repository"].(string)
	if !ownerOk || !repoOk {
		return nil, cli.Exit("git repository not found", 1)
	}
	return &GitInfo{
		Owner: owner,
		Repo:  repo,
	}, nil
}
