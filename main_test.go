package main

import "testing"
import (
	"github.com/stretchr/testify/assert"
	"go/build"
	"path/filepath"
)

func TestParseRawURL(t *testing.T) {
	// SSH remote style
	dir, remote := parseRawURL("git@github.com:brentnd/git-get.git")
	assert.Equal(t, filepath.FromSlash(build.Default.GOPATH+"/src/github.com/brentnd/git-get"), dir)
	assert.Equal(t, "git@github.com:brentnd/git-get.git", remote)

	// HTTPS remote style
	dir, remote = parseRawURL("https://github.com/brentnd/git-get.git")
	assert.Equal(t, filepath.FromSlash(build.Default.GOPATH+"/src/github.com/brentnd/git-get"), dir)
	assert.Equal(t, "https://github.com/brentnd/git-get.git", remote)

	// Golang package style
	dir, remote = parseRawURL("github.com/brentnd/git-get")
	assert.Equal(t, filepath.FromSlash(build.Default.GOPATH+"/src/github.com/brentnd/git-get"), dir)
	assert.Equal(t, "https://github.com/brentnd/git-get.git", remote)
}
