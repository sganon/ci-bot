package main

import (
	"github.com/urfave/cli"

	"github.com/sganon/code-bot/api"
	"github.com/sganon/code-bot/gitlab"
)

var apiCmd = cli.Command{
	Name:  "api",
	Usage: "launch api handling slash command and other request from slack",
	Flags: []cli.Flag{
		apiHostFlag, apiPortFlag,
		signinSecretFlag,
	},
	Action: func(c *cli.Context) error {
		cp := c.Parent()
		glAPI := gitlab.NewAPI(cp.String(gitlabAddrFlag.Name), cp.String(gitlabTokenFlag.Name))
		api := api.New(
			c.String(apiHostFlag.Name), c.String(apiPortFlag.Name),
			c.String(signinSecretFlag.Name), glAPI)
		api.Serve()
		return nil
	},
}
