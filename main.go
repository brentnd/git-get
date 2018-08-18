package main

import (
	"fmt"
	"os"

	"go/build"
	"gopkg.in/src-d/go-git.v4"
	"net/url"
	"path/filepath"
	"strings"
)

func exitWithError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func gopath() string {
	list := filepath.SplitList(build.Default.GOPATH)
	if len(list) == 0 {
		exitWithError("cannot download, $GOPATH not set. For more details see: 'go help gopath'")
	}
	_, err := os.Stat(build.Default.GOPATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "created GOPATH=%s; see 'go help gopath'\n", build.Default.GOPATH)
	}
	return list[0]
}

// Basic example of how to clone a repository using clone options.
func main() {
	if len(os.Args) < 2 {
		exitWithError("usage: git get <repo>")
	}
	rawurl := os.Args[1]
	if !strings.Contains(rawurl, "//") {
		// "//" needed to provide protocol to make git@ seen as user
		rawurl = "//" + rawurl
	}
	repoUrl, err := url.Parse(rawurl)
	if err != nil {
		exitWithError("syntax error in repo: %s", err)
	}
	repoUrl.Path = strings.Replace(repoUrl.Path, ".git", "", 1)
	destinationDir := filepath.Join(gopath(), "src", repoUrl.Hostname(), repoUrl.Port(), filepath.FromSlash(repoUrl.Path))

	if repoUrl.Scheme == "//" {
		repoUrl.Scheme = ""
	}
	if !repoUrl.IsAbs() && repoUrl.User == nil {
		repoUrl.Scheme = "https"
	}
	repoUrl.Path += ".git"

	// Check that this is an appropriate place for the repo to be checked out.
	// The target directory must either not exist or have a repo checked out already.
	meta := filepath.Join(destinationDir, ".git")
	st, err := os.Stat(meta)
	if err == nil && !st.IsDir() {
		exitWithError("%s exists but is not a directory", meta)
	}
	if err != nil {
		// .git directory does not exist. Prepare to checkout new copy.
		// Require the target directory not to exist to avoid stepping on existing work.
		if _, err := os.Stat(destinationDir); err == nil {
			exitWithError("%s exists but %s does not - stale checkout?", destinationDir, meta)
		}
		// Require the parent of the target to exist.
		parent, _ := filepath.Split(destinationDir)
		if err = os.MkdirAll(parent, 0777); err != nil {
			exitWithError("couldn't create directory %s", parent)
		}
		remoteUrl := strings.TrimPrefix(repoUrl.String(), "//")
		_, err = git.PlainClone(destinationDir, false, &git.CloneOptions{URL: remoteUrl})
		if err != nil {
			exitWithError("clone failed: %s", err)
		}
	} else {
		fmt.Printf("repo already exists at %s\n", destinationDir)
	}
}
