package githubsource

import (
	"cilintservice/servcies/gitsource"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type Client struct {
	client           *github.Client
	oauth2HTTPClient *http.Client
	APIURL           string
	WebURL           string
	oauth2ClientID   string
	oauth2Secret     string
}

func New(oAuthClient *http.Client) *Client {
	client := github.NewClient(oAuthClient)

	return &Client{
		client:           client,
		oauth2HTTPClient: oAuthClient,
	}
}

func (c *Client) CreateRepoWebhook(repopath, url, secret string) error {
	owner, reponame, err := parseRepoPath(repopath)
	if err != nil {
		return err
	}

	hook := &github.Hook{
		Events: []string{"push", "pull_request"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          url,
			"content_type": "json",
			"secret":       secret,
		},
	}

	if _, _, err := c.client.Repositories.CreateHook(context.Background(), owner, reponame, hook); err != nil {
		return fmt.Errorf("error creating repository webhook: %w", err)
	}

	return nil
}

func (c *Client) ListUserRepos() ([]*gitsource.RepoInfo, error) {
	remoteRepos := []*github.Repository{}

	opt := &github.RepositoryListOptions{}
	for {
		pRemoteRepos, resp, err := c.client.Repositories.List(context.Background(), "", opt)
		if err != nil {
			return nil, err
		}
		remoteRepos = append(remoteRepos, pRemoteRepos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	repos := []*gitsource.RepoInfo{}

	for _, rr := range remoteRepos {
		// keep only repos with admin permissions
		if rr.Permissions != nil {
			if !(*rr.Permissions)["admin"] {
				continue
			}
			repos = append(repos, fromGithubRepo(rr))
		}
	}

	return repos, nil
}

func (c *Client) GetUserInfo() (*gitsource.UserInfo, error) {
	user, _, err := c.client.Users.Get(context.Background(), "")
	if err != nil {
		return nil, err
	}

	userInfo := &gitsource.UserInfo{
		ID:        strconv.FormatInt(*user.ID, 10),
		LoginName: *user.Login,
	}
	if user.Email != nil {
		userInfo.Email = *user.Email
	}

	return userInfo, nil
}

func fromGithubRepo(rr *github.Repository) *gitsource.RepoInfo {
	return &gitsource.RepoInfo{
		ID:           strconv.FormatInt(*rr.ID, 10),
		Path:         path.Join(*rr.Owner.Login, *rr.Name),
		HTMLURL:      *rr.HTMLURL,
		SSHCloneURL:  *rr.SSHURL,
		HTTPCloneURL: *rr.CloneURL,
	}
}

func parseRepoPath(repopath string) (string, string, error) {
	parts := strings.Split(repopath, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("wrong github repo path: %q", repopath)
	}
	return parts[0], parts[1], nil
}
