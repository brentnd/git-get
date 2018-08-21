package main

import (
	"fmt"
	"os"

	"go/build"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"
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

func parseRawURL(rawurl string) (directory, remote string) {
	if !strings.Contains(rawurl, "//") {
		// "//" needed to provide protocol to make git@ seen as user
		rawurl = "//" + rawurl
	}
	repoURL, err := url.Parse(rawurl)
	if err != nil {
		exitWithError("syntax error in repo: %s", err)
	}
	repoURL.Path = strings.Replace(repoURL.Path, ".git", "", 1)
	directory = filepath.Join(gopath(), "src", repoURL.Hostname(), repoURL.Port(), filepath.FromSlash(repoURL.Path))

	if repoURL.Scheme == "//" {
		repoURL.Scheme = ""
	}
	if !repoURL.IsAbs() && repoURL.User == nil {
		repoURL.Scheme = "https"
	}
	repoURL.Path += ".git"
	remote = strings.TrimPrefix(repoURL.String(), "//")
	return
}

func remoteRepoExists(remote string) error {
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return err
	}
	rem, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})
	if err != nil {
		return err
	}
	_, err = rem.List(&git.ListOptions{})
	return err
}

func main() {
	if len(os.Args) < 2 {
		exitWithError("fatal: You must specify a repository to clone.\n\nusage: git get <repo>")
	}

	directory, remote := parseRawURL(os.Args[1])

	if err := remoteRepoExists(remote); err != nil {
		exitWithError("fatal: repository '%s' does not exist", os.Args[1])
	}

	// Check that this is an appropriate place for the repo to be checked out.
	// The target directory must either not exist or have a repo checked out already.
	meta := filepath.Join(directory, ".git")
	st, err := os.Stat(meta)
	if err == nil && !st.IsDir() {
		exitWithError("%s exists but is not a directory", meta)
	}
	if err != nil {
		// .git directory does not exist. Prepare to checkout new copy.
		// Require the target directory not to exist to avoid stepping on existing work.
		if _, err := os.Stat(directory); err == nil {
			exitWithError("%s exists but %s does not - stale checkout?", directory, meta)
		}
		// Require the parent of the target to exist.
		parent, _ := filepath.Split(directory)
		if err = os.MkdirAll(parent, 0777); err != nil {
			exitWithError("couldn't create directory %s", parent)
		}
		_, err = git.PlainClone(directory, false, &git.CloneOptions{URL: remote})
		if err != nil {
			exitWithError("fatal: %s", err)
		}
	} else {
		fmt.Printf("repo already exists at %s\n", directory)
	}
}
