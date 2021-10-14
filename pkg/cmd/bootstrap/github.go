package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/scrape"
	"github.com/google/go-github/v39/github"
)

// generateManifest generate manifest from the given options
func generateManifest(opts *bootstrapOpts) ([]byte, error) {
	sc := scrape.AppManifest{
		Name:           github.String(opts.GithubApplicationName),
		URL:            github.String(opts.GithubApplicationURL),
		HookAttributes: map[string]string{"url": opts.RouteName},
		RedirectURL:    github.String(fmt.Sprintf("http://localhost:%d", opts.webserverPort)),
		Description:    github.String("Pipilines as Code Application"),
		Public:         github.Bool(true),
		DefaultEvents: []string{
			"commit_comment",
			"issue_comment",
			"pull_request",
			"pull_request_review",
			"pull_request_review_comment",
			"push",
		},
		DefaultPermissions: &github.InstallationPermissions{
			Checks:           github.String("write"),
			Contents:         github.String("write"),
			Issues:           github.String("write"),
			Members:          github.String("read"),
			Metadata:         github.String("read"),
			OrganizationPlan: github.String("read"),
			PullRequests:     github.String("write"),
		},
	}
	return json.Marshal(sc)
}

// getGHClient get github client
func getGHClient(opts *bootstrapOpts) (*github.Client, error) {
	if opts.GithubAPIURL == defaultPublicGithub {
		return github.NewClient(nil), nil
	}

	gvcs, err := github.NewEnterpriseClient(opts.GithubAPIURL, "", nil)
	if err != nil {
		return nil, err
	}
	return gvcs, nil
}
