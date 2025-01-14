package provider

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	testRetestAllRegex    = regexp.MustCompile(`(?m)^(/retest|/test)\s*$`)
	testRetestSingleRegex = regexp.MustCompile(`(?m)^(/test|/retest)[ \t]+\S+`)
	oktotestRegex         = regexp.MustCompile(`(?m)^/ok-to-test\s*$`)
)

const (
	ProviderGitHubApp     = "GitHubApp"
	ProviderGitHubWebhook = "GitHubWebhook"
	ProviderGitLabWebhook = "GitLabWebhook"
)

func Valid(value string, validValues []string) bool {
	for _, v := range validValues {
		if v == value {
			return true
		}
	}
	return false
}

func IsTestRetestComment(comment string) bool {
	return testRetestSingleRegex.MatchString(comment) || testRetestAllRegex.MatchString(comment)
}

func IsOkToTestComment(comment string) bool {
	return oktotestRegex.MatchString(comment)
}

func GetPipelineRunFromComment(comment string) string {
	var splitTest []string
	if strings.Contains(comment, "/test") {
		splitTest = strings.Split(comment, "/test")
	} else {
		splitTest = strings.Split(comment, "/retest")
	}
	// now get the first line
	getFirstLine := strings.Split(splitTest[1], "\n")
	// trim spaces
	return strings.TrimSpace(getFirstLine[0])
}

// CompareHostOfURLS compares the host of two parsed URLs and returns true if
// they are
func CompareHostOfURLS(uri1, uri2 string) bool {
	u1, err := url.Parse(uri1)
	if err != nil || u1.Host == "" {
		return false
	}
	u2, err := url.Parse(uri2)
	if err != nil || u2.Host == "" {
		return false
	}
	return u1.Host == u2.Host
}
