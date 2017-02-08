package main

import (
	"os"
	"weixinGate/action"

	"weixinGate/logger"

	"github.com/codegangsta/cli"
)

var log = logger.New()

func main() {
	app := cli.NewApp()
	app.Name = "weixin gate"
	app.Author = "chengjt"
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "server",
			ShortName: "server",
			Usage:     "start weixin gate server",
			Action:    action.Serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host,H",
					Usage: "server bind host to",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "port,p",
					Usage: "bind port",
					Value: 8080,
				},
				cli.StringFlag{
					Name:  "target,t",
					Usage: "which url that post msg to",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Error("app run eror ", err.Error())
	}
}
