package cmd

import (
	"html/template"
	"io/ioutil"
	"net/http"
	//"os"
	"fmt"
	"path"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"

	"github.com/MessageDream/webIM/modules/base"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/middleware"
	//"github.com/MessageDream/webIM/modules/middleware/binding"
	"github.com/MessageDream/webIM/modules/setting"
	"github.com/MessageDream/webIM/routers"
	"github.com/MessageDream/webIM/routers/app"
	"github.com/MessageDream/webIM/routers/chat/longpolling"
	"github.com/MessageDream/webIM/routers/chat/websocket"
)

var CmdApp = cli.Command{
	Name:  "app",
	Usage: "Start IM server",
	Description: `IM server is the only thing you need to run, 
and it takes care of all the other things for you`,
	Action: runApp,
	Flags:  []cli.Flag{},
}

// checkVersion checks if binary matches the version of temolate files.
func checkVersion() {
	data, err := ioutil.ReadFile(path.Join(setting.StaticRootPath, "templates/VERSION"))
	if err != nil {
		log.Fatal("Fail to read 'templates/VERSION': %v", err)
	}
	if string(data) != setting.AppVer {
		log.Fatal("Binary and template file version does not match, did you forget to recompile?")
	}
}

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(middleware.Logger())
	m.Use(martini.Recovery())
	m.Use(middleware.Static("public",
		middleware.StaticOptions{SkipLogging: !setting.DisableRouterLog}))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func runApp(*cli.Context) {
	routers.GlobalInit()
	checkVersion()

	m := newMartini()

	// Middlewares.
	m.Use(middleware.Renderer(middleware.RenderOptions{
		Directory:  path.Join(setting.StaticRootPath, "templates"),
		Funcs:      []template.FuncMap{base.TemplateFuncs},
		IndentJSON: true,
	}))
	m.Use(middleware.InitContext())

	//reqSignIn := middleware.Toggle(&middleware.ToggleOptions{SignInRequire: true})
	//ignSignIn := middleware.Toggle(&middleware.ToggleOptions{SignInRequire: setting.Service.RequireSignInView})
	//ignSignInAndCsrf := middleware.Toggle(&middleware.ToggleOptions{DisableCsrf: true})

	//reqSignOut := middleware.Toggle(&middleware.ToggleOptions{SignOutRequire: true})

	//	bindIgnErr := binding.BindIgnErr

	// Routers.
	//m.Get("/", ignSignIn, routers.Home)
	//m.Get("/install", bindIgnErr(auth.InstallForm{}), routers.Install)
	//m.Post("/install", bindIgnErr(auth.InstallForm{}), routers.InstallPost)
	//m.Group("", func(r martini.Router) {
	//	r.Get("/issues", user.Issues)
	//	r.Get("/pulls", user.Pulls)
	//	r.Get("/stars", user.Stars)
	//}, reqSignIn)

	//m.Group("/api", func(_ martini.Router) {
	//	m.Group("/v1", func(r martini.Router) {
	//		// Miscellaneous.
	//		r.Post("/markdown", bindIgnErr(apiv1.MarkdownForm{}), v1.Markdown)
	//		r.Post("/markdown/raw", v1.MarkdownRaw)

	//		// Users.
	//		r.Get("/users/search", v1.SearchUser)

	//		r.Any("**", func(ctx *middleware.Context) {
	//			ctx.JSON(404, &base.ApiJsonErr{"Not Found", v1.DOC_URL})
	//		})
	//	})
	//})

	//m := martini.Classic()
	//	m.Get("/", &controllers.AppController{})
	//	// Indicate AppController.Join method to handle POST requests.
	//	m.Post("/join", &controllers.AppController{}, "post:Join")

	//	// Long polling.
	//	m.Get("/lp", &controllers.LongPollingController{}, "get:Join")
	//	m.Post("/lp/post", &controllers.LongPollingController{})
	//	m.Post("/lp/fetch", &controllers.LongPollingController{}, "get:Fetch")

	//	// WebSocket.
	//	m.Get("/ws", &controllers.WebSocketController{})
	//	m.Get("/ws/join", &controllers.WebSocketController{}, "get:Join")

	m.Get("/", app.Welcome)
	m.Post("/join", app.Join)

	m.Get("/lp", longpolling.Join)
	m.Group("/lp", func(r martini.Router) {
		m.Post("/post", longpolling.Post)
		m.Get("/fetch", longpolling.Fetch)
	})

	m.Get("/ws", websocket.Get)
	m.Get("/ws/join", websocket.Join)
	//Not found handler.
	m.NotFound(routers.NotFound)

	var err error
	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	log.Info("Listen: %v://%s", setting.Protocol, listenAddr)
	switch setting.Protocol {
	case setting.HTTP:
		err = http.ListenAndServe(listenAddr, m)
	case setting.HTTPS:
		err = http.ListenAndServeTLS(listenAddr, setting.CertFile, setting.KeyFile, m)
	default:
		log.Fatal("Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Fatal("Fail to start server: %v", err)
	}
}
