package routers

import (
	//"path"

	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/middleware"
	"github.com/MessageDream/webIM/modules/setting"
)

func NewServices() {
	setting.NewServices()
}
func GlobalInit() {
	setting.NewConfigContext()
	log.Trace("Custom path: %s", setting.CustomPath)
	log.Trace("Log path: %s", setting.LogRootPath)
	NewServices()
	//if setting.InstallLock {
	//	if err := models.NewEngine(); err != nil {
	//		log.Fatal("Fail to initialize ORM engine: %v", err)
	//	}

	//	log.NewGitLogger(path.Join(setting.LogRootPath, "http.log"))
	//}
	//if models.EnableSQLite3 {
	//	log.Info("SQLite3 Enabled")
	//}
	//checkRunMode()
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Data["PageIsNotFound"] = true
	ctx.Handle(404, "home.NotFound", nil)
}
