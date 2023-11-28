package jobs

type Context struct {
	RepositoryURL string
	CommitHash    string
}

func NewContext(repositoryURL string, commitHash string) Context {
	return Context{repositoryURL, commitHash}
}
