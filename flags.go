package main

import (
	"github.com/urfave/cli"
)

var gitlabAddrFlag = cli.StringFlag{
	EnvVar: "GITLAB_ADDR",
	Name:   "gladdr",
}

var gitlabTokenFlag = cli.StringFlag{
	EnvVar: "GITLAB_API_TOKEN",
	Name:   "gltoken",
}

var releaseSlackHookFlag = cli.StringFlag{
	EnvVar: "RELEASE_SLACK_HOOK",
	Name:   "releasehook",
}

var apiHostFlag = cli.StringFlag{
	Name:  "apihost",
	Value: "0.0.0.0",
}

var apiPortFlag = cli.StringFlag{
	Name:  "apiport",
	Value: "8080",
}

var signinSecretFlag = cli.StringFlag{
	Name:   "signinsecret",
	EnvVar: "SLACK_SIGNIN_SECRET",
}
