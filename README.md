# git get

A simple plugin for git which allows cloning repos from any directory into their directory based on how
gopath works.

## Getting Started

All source cloned with `git get` will be installed in `$GOPATH/src/...`.


### Prerequisites

What things you need to install the software and how to install them:
* `git`
* `go`
* `dep`

You will need to have go installed and setup with `$GOPATH`. Your
`$GOPATH/bin` should be in your path.

### Installing

Follow these steps to get the `git get` command working.

Install the git-get command into `$GOPATH/bin`.

```bash
go get github.com/brentnd/git-get
```

### Usage

As a git plugin, the git-get command is run with:
```bash
git get git@github.com:brentnd/git-get.git       # ssh
git get https://github.com/brentnd/git-get.git   # https
git get github.com/brentnd/git-get               # go package style
```

To clone a private repo, be sure your SSH key is added to the authentication agent.
In most cases done with:
```bash
ssh-add [key-file]
```

## Built With

* [go-git](https://github.com/src-d/go-git) - go-git library

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

