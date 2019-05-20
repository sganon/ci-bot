package main

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/sganon/ci-bot/gitlab"
)

var releaseCmd = cli.Command{
	Name:      "release",
	Usage:     "Send release message to channel for a given project and tag",
	UsageText: "release [GROUP]/PROJECT TAG",
	Action: func(c *cli.Context) error {
		cp := c.Parent()
		glAPI := gitlab.NewAPI(cp.String(gitlabAddrFlag.Name), cp.String(gitlabTokenFlag.Name))
		fmt.Println(c.Args())
		if c.NArg() != 2 {
			return fmt.Errorf("usage error: you need to provide [GROUP]/PROJECT TAG")
		}
		pj, err := gitlab.GetProjectByName(glAPI, c.Args().Get(0))
		if err != nil {
			return fmt.Errorf("error getting project: %v", err)
		}
		fmt.Println(pj)
		pj.Ref.Value = c.Args().Get(1)
		err = pj.FetchRefPipelines(glAPI)
		if err != nil {
			return fmt.Errorf("error getting pipelines: %v", err)
		}
		fmt.Println(pj.Ref.Pipelines)
		return err
	},
}
