package base

import (
	"html/template"
	"runtime"

	"github.com/MessageDream/webIM/modules/setting"
)

var TemplateFuncs template.FuncMap = map[string]interface{}{
	"GoVer": func() string {
		return runtime.Version()
	},
	"AppName": func() string {
		return setting.AppName
	},
	"AppVer": func() string {
		return setting.AppVer
	},
	"AppDomain": func() string {
		return setting.Domain
	},
}
