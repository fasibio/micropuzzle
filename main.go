package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		next.ServeHTTP(w, r)
	})
}

const (
	CliFallbackLoader        = "fallbackloader"
	CliMicrofrontends        = "microfrontends"
	CliPublicFolder          = "publicfoder"
	CliTimeout               = "timeoutms"
	CliLogLevel              = "logLevel"
	CliPort                  = "port"
	CliRedisAddress          = "redisaddr"
	CliRedisUser             = "redisuser"
	CliRedisPassword         = "redispassword"
	CliRedisDb               = "redisdb"
	EnvPrefix         string = "MICROPUZZLE_"
)

func getFlagEnvByFlagName(flagName string) string {
	return EnvPrefix + strings.ToUpper(flagName)
}

func main() {
	app := cli.NewApp()
	app.Name = "BoulderDB API"
	app.Description = "Service for low level Data-Communication"
	runner := Runner{}
	app.Action = runner.Run
	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:    CliTimeout,
			EnvVars: []string{getFlagEnvByFlagName(CliTimeout)},
			Usage:   "Timeout for loading Microfrontends (for all slower, it will be use streaming to bring it to the client)",
			Value:   45 * time.Millisecond,
		},
		&cli.StringFlag{
			Name:    CliLogLevel,
			EnvVars: []string{getFlagEnvByFlagName(CliLogLevel)},
			Usage:   "Loglevel debug, info, warn, error",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    CliPort,
			EnvVars: []string{getFlagEnvByFlagName(CliPort)},
			Usage:   "port where server will be started",
			Value:   "3000",
		},
		&cli.StringFlag{
			Name:    CliPublicFolder,
			EnvVars: []string{getFlagEnvByFlagName(CliPublicFolder)},
			Value:   "./public",
			Usage:   "Folder where all html js css from server directly will be foundable (Public folder for the web)",
		},
		&cli.StringFlag{
			Name:    CliMicrofrontends,
			EnvVars: []string{getFlagEnvByFlagName(CliMicrofrontends)},
			Value:   "./config/frontends.ini",
			Usage:   "A ini file (key=value) key is for logic name of microfrontend. value is the url where to fetch the content (groups are . seperated by using)",
		},
		&cli.StringFlag{
			Name:    CliFallbackLoader,
			EnvVars: []string{getFlagEnvByFlagName(CliFallbackLoader)},
			Usage:   "key of inifile where to find fallbackhtml which will shown if microfrontend is lower than timeout",
			Value:   "fallback",
		},
		&cli.StringFlag{
			Name:    CliRedisAddress,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisAddress)},
			Usage:   "The domian/ip:port of redis",
			Value:   "localhost:6379",
		},
		&cli.StringFlag{
			Name:    CliRedisUser,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisUser)},
			Usage:   "Username to connect to redis",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    CliRedisPassword,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisPassword)},
			Usage:   "Password to connect to redis",
			Value:   "",
		},
		&cli.Int64Flag{
			Name:    CliRedisDb,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisDb)},
			Usage:   "Db to use by redis",
			Value:   0,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}

type Runner struct{}

func (ru *Runner) Run(c *cli.Context) error {
	logger.Initialize(c.String(CliLogLevel))
	iniF, err := ini.Load(c.String(CliMicrofrontends))
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	cache, err := NewRedisHandler(&redis.Options{
		Addr:     c.String(CliRedisAddress),
		DB:       c.Int(CliRedisDb),
		Username: c.String(CliRedisUser),
		Password: c.String(CliRedisPassword),
	})
	if err != nil {
		return err
	}

	// sockerHandler := NewSocketHandler(cache, c.Duration(CliTimeout), iniF, c.String(CliFallbackLoader), &socketio.RedisAdapterOptions{
	// 	Addr:    c.String(CliRedisAddress),
	// 	Prefix:  "MICROPUZZLE_SOCKET.IO",
	// 	Network: "tcp",
	// })

	websocketHandler := NewWebSocketHandler(cache, c.Duration(CliTimeout), iniF, c.String(CliFallbackLoader))
	// defer sockerHandler.Server.Close()
	// r.Handle("/socket.io/", sockerHandler.Server)
	r.HandleFunc("/socket", websocketHandler.Handle)
	f := FileHandler{
		server: &websocketHandler,
	}
	f.ChiFileServer(r, "/", http.Dir(c.String(CliPublicFolder)))

	logger.Get().Infof("Start Server on Port :%s", c.String(CliPort))
	return http.ListenAndServe(fmt.Sprintf(":%s", c.String(CliPort)), r)

}

func mimeTypeForFile(file string) string {
	// We use a built in table of the common types since the system
	// TypeByExtension might be unreliable. But if we don't know, we delegate
	// to the system.
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"

		// ...

	default:
		return mime.TypeByExtension(ext)
	}
}

type FileHandler struct {
	server *WebSocketHandler
}

func (filehandler *FileHandler) ChiFileServer(r chi.Router, path string, root http.FileSystem) {

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := root.Open(path)
		mimetype := mimeTypeForFile(path)
		w.Header().Set("Content-Type", mimetype)
		if err == nil {
			if mimetype == "application/javascript" {
				io.Copy(w, f)
				return
			}
			err := filehandler.handleTemplate(f, w, r)
			if err != nil {
				logger.Get().Warnw("Error handle template", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (filehandler *FileHandler) handleTemplate(f http.File, dst io.Writer, r *http.Request) error {

	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := NewTemplateHandler(r, "ws://localhost:3000/socket", id, filehandler.server)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
