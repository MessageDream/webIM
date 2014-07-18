package app

import (
	"github.com/MessageDream/webIM/modules/base"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/middleware"
)

const (
	WELCOME base.TplName = "app/welcome"
)

func Welcome(ctx *middleware.Context) {
	ctx.HTML(200, WELCOME)
}

// Join method handles POST requests
func Join(ctx *middleware.Context) {
	// Get form value.
	uname := ctx.Query("uname")
	roomid := ctx.Query("roomid")
	tech := ctx.Query("tech")

	log.Warn("%s", tech)
	// Check valid.
	if len(uname) == 0 || len(roomid) == 0 || len(tech) == 0 {
		ctx.Redirect("/", 302)
		return
	}

	switch tech {
	case "Long Polling", "长轮询":
		ctx.Redirect("/lp?uname="+uname+"&roomid="+roomid, 302)
	case "WebSocket":
		ctx.Redirect("/ws?uname="+uname+"&roomid="+roomid, 302)
	default:
		ctx.Redirect("/", 302)
	}

	// Usually put return after redirect.
	return
}
