package main

import (
	"fmt"

	"github.com/urfave/cli"

	"strings"

	"github.com/sganon/ci-bot/gitlab"
	"github.com/sganon/ci-bot/slack"
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
		pj.Tag.Name = c.Args().Get(1)
		err = pj.FetchTagPipelines(glAPI)
		if err != nil {
			return fmt.Errorf("error getting pipelines: %v", err)
		}
		fmt.Println(pj.Tag.Pipelines)
		err = pj.FetchTag(glAPI)
		if err != nil {
			return fmt.Errorf("error getting tag: %v", err)
		}
		fmt.Println(pj.Tag.Release)

		text := strings.ReplaceAll(pj.Tag.Release.Description, "*", "•")
		attch := slack.Attachment{
			Fallback: "release text",
			Color:    "#008bd2",
			Pretext:  fmt.Sprintf("New release of project %s: *%s*", pj.NameWithNamespace, pj.Tag.Name),
			Title:    "Changelog",
			Text:     text,
		}
		attch.Send(cp.String("releasehook"))

		return err
	},
}
