package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/MessageDream/webIM/modules/setting"
	"github.com/MessageDream/webIM/modules/socket"
	soc "github.com/MessageDream/webIM/routers/chat/socket"
)

var CmdSocketServer = cli.Command{
	Name:  "socketServer",
	Usage: "Start IM socketServer",
	Description: `IM socketServer is the only thing you need to run, 
and it takes care of all the other things for you`,
	Action: runSocket,
	Flags:  []cli.Flag{},
}

func runSocket(*cli.Context) {
	server := socket.NewServer()

	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.TcpPort)
	fmt.Println(listenAddr)
	server.ListenTCP(listenAddr)
	server.OnConnected = soc.OnConnected
	server.OnMessage = soc.OnMessage
	server.OnDisconnected = soc.OnDisconnected
	server.Boot()
	server.Wait()
}
