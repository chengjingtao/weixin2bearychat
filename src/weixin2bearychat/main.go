package main

import (
	"os"
	"weixin2bearychat/action"

	"weixin2bearychat/logger"

	"fmt"

	"github.com/codegangsta/cli"
)

var (
	version   = "unknown"
	buildDate = ""
)

var log = logger.New()

func main() {
	app := cli.NewApp()
	app.Name = "weixin2bearychat"
	app.Author = "chengjt"
	app.Version = fmt.Sprintf("%s (%s)", version, buildDate)
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "server",
			ShortName: "server",
			Usage:     "start weixin to bearychat server",
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
					Value: 80,
				},
				cli.StringFlag{
					Name:  "tmplpath",
					Usage: "msg template path",
					Value: "/etc/weixin2bearychat/tmpl/",
				},
				cli.StringFlag{
					Name:  "target,t",
					Usage: "which url that post msg to",
				},
			},
		},
	}
	fmt.Println("Version is " + version + ", Build Date is " + buildDate)
	err := app.Run(os.Args)
	if err != nil {
		log.Error("app run eror ", err.Error())
	}
}
