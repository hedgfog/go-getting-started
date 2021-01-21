package gitsource

type UserInfo struct {
	ID        string
	LoginName string
	Email     string
}

type RepoInfo struct {
	ID           string
	Path         string
	HTMLURL      string
	SSHCloneURL  string
	HTTPCloneURL string
}
