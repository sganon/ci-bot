package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	log.SetLevel(log.DebugLevel)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		gitlabAddrFlag, gitlabTokenFlag,
		releaseSlackHookFlag,
	}
	app.Commands = []cli.Command{
		releaseCmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
