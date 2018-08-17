package main

import (
	"fmt"
	"os"

	"go/build"
	"gopkg.in/src-d/go-git.v4"
	"path/filepath"
)

func exitWithError(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

// Basic example of how to clone a repository using clone options.
func main() {
	if len(os.Args) < 2 {
		exitWithError("Usage: git get <repo>")
	}
	repo := os.Args[1]

	// TODO: strip possible scheme from url and .git?

	list := filepath.SplitList(build.Default.GOPATH)
	if len(list) == 0 {
		exitWithError("cannot download, $GOPATH not set. For more details see: 'go help gopath'")
	}
	srcRoot := filepath.Join(list[0], "src")
	root := filepath.Join(srcRoot, filepath.FromSlash(repo)) // TODO

	// Check that this is an appropriate place for the repo to be checked out.
	// The target directory must either not exist or have a repo checked out already.
	meta := filepath.Join(root, ".git")
	st, err := os.Stat(meta)
	if err == nil && !st.IsDir() {
		exitWithError(meta + " exists but is not a directory")
	}
	if err != nil {
		// .git directory does not exist. Prepare to checkout new copy.
		// Require the target directory not to exist to avoid stepping on existing work.
		if _, err := os.Stat(root); err == nil {
			exitWithError(fmt.Sprintf("%s exists but %s does not - stale checkout?", root, meta))
		}

		_, err := os.Stat(build.Default.GOPATH)
		gopathExisted := err == nil

		// Require the parent of the target to exist.
		parent, _ := filepath.Split(root)
		if err = os.MkdirAll(parent, 0777); err != nil {
			exitWithError("couldn't create directory " + parent)
		}
		if !gopathExisted {
			fmt.Fprintf(os.Stderr, "created GOPATH=%s; see 'go help gopath'\n", build.Default.GOPATH)
		}

		// Should this always be https://?
		url := fmt.Sprintf("https://%s.git", repo)
		_, err = git.PlainClone(root, false, &git.CloneOptions{URL: url})
		if err != nil {
			exitWithError(fmt.Sprintf("error: %s", err))
		}
	} else {
		// Metadata directory does exist; download incremental updates.
		// TODO: git pull?
	}
}
