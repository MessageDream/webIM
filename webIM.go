package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/MessageDream/webIM/cmd"
	"github.com/MessageDream/webIM/modules/setting"
)

const APP_VER = "0.0.1"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "WebIM"
	app.Usage = "WebIM Service"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdApp,
		//cmd.CmdSocketServer,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
