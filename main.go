package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	CliFallbackLoader        = "fallbackloader"
	CliMicrofrontends        = "microfrontends"
	CliPublicFolder          = "publicfoder"
	CliTimeout               = "timeoutms"
	CliLogLevel              = "logLevel"
	CliPort                  = "port"
	CliManagementPort        = "managementport"
	CliRedisAddress          = "redisaddr"
	CliRedisUser             = "redisuser"
	CliRedisPassword         = "redispassword"
	CliRedisDb               = "redisdb"
	EnvPrefix         string = "MICROPUZZLE_"
)

const (
	SOCKET_PATH      = "socket"
	SOCKET_ENDPOINT  = "/micro-puzzle"
	LIB_ENDPOINT     = "/micro-lib/*"
	METRICS_ENDPOINT = "/metrics"
)

func getFlagEnvByFlagName(flagName string) string {
	return EnvPrefix + strings.ToUpper(flagName)
}

func main() {
	app := cli.NewApp()
	app.Name = "Micropuzzle"

	app.Description = "Application to combine Server Side Include and Afterloading"
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
			Value:   "3300",
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
			Value:   "./config/frontends.yaml",
			Usage:   "A yaml file to describe available Frontends",
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
		&cli.Int64Flag{
			Name:    CliManagementPort,
			EnvVars: []string{getFlagEnvByFlagName(CliManagementPort)},
			Usage:   "Port to get data not needed from client",
			Value:   3301,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}

type Frontends map[string]map[string]Frontend

type Frontend struct {
	Url            string `yaml:"url"`
	StripUrlPrefix bool   `yaml:"stripUrlPrefix"`
}

//go:embed micro-lib/*.js
var embeddedLib embed.FS

type Runner struct{}

func (ru *Runner) Run(c *cli.Context) error {
	logs, err := logger.Initialize(c.String(CliLogLevel))
	if err != nil {
		return err
	}
	frontends, err := loadFrontends(c.String(CliMicrofrontends))
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	r.Use(chiprometheus.NewMiddleware("micropuzzle"))
	r.Use(middleware.Compress(5))
	registerReverseProxy(frontends, r)

	cache, err := NewRedisHandler(&redis.Options{
		Addr:     c.String(CliRedisAddress),
		DB:       c.Int(CliRedisDb),
		Username: c.String(CliRedisUser),
		Password: c.String(CliRedisPassword),
	})
	if err != nil {
		return err
	}
	websocketHandler := NewWebSocketHandler(cache, c.Duration(CliTimeout), frontends, c.String(CliFallbackLoader))

	r.HandleFunc("/"+SOCKET_PATH, websocketHandler.Handle)
	f := FileHandler{
		server:    &websocketHandler,
		socketUrl: SOCKET_PATH,
	}
	r.Get(SOCKET_ENDPOINT, websocketHandler.LoadFragmentHandler)
	r.Handle(LIB_ENDPOINT, http.FileServer(http.FS(embeddedLib)))
	f.ChiFileServer(r, "/", http.Dir(c.String(CliPublicFolder)))

	logs.Infof("Start Server on Port :%s and Management on port %s", c.String(CliPort), c.String(CliManagementPort))
	managementR := chi.NewRouter()
	managementR.Handle(METRICS_ENDPOINT, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%s", c.String(CliManagementPort)), managementR)
	return http.ListenAndServe(fmt.Sprintf(":%s", c.String(CliPort)), r)

}

func loadFrontends(frontendsPath string) (Frontends, error) {
	frontendsBody, err := ioutil.ReadFile(frontendsPath)
	if err != nil {
		return nil, err
	}
	var frontends Frontends
	err = yaml.Unmarshal(frontendsBody, &frontends)
	if err != nil {
		return nil, err
	}
	return frontends, nil
}

func registerReverseProxy(frontends Frontends, r chi.Router) {
	for key, one := range frontends {
		for oneK, oneV := range one {
			prefix := ""
			if key != "global" {
				prefix = key + "."
			}
			err := registerMicrofrontendProxy(r, prefix+oneK, oneV)
			if err != nil {
				logger.Get().Warnw(fmt.Sprintf("Error by setting Reverseproxy for destination %s", prefix+oneK), "error", err)
			}
		}
	}
}

func registerMicrofrontendProxy(r chi.Router, name string, frontend Frontend) error {
	url, err := url.Parse(frontend.Url)
	if err != nil {
		return err
	}
	logger.Get().Infof("Register endpoint /%s/* for frontend %s", name, name)
	r.HandleFunc(fmt.Sprintf("/%s/*", name), func(w http.ResponseWriter, r *http.Request) {
		if frontend.StripUrlPrefix {
			path := strings.Replace(r.URL.String(), "/"+name, "", 1)
			r.URL, err = url.Parse(path)
		}
		p := httputil.NewSingleHostReverseProxy(url)
		p.ModifyResponse = rewriteBodyHandler("/" + name)
		p.ServeHTTP(w, r)
	})
	return nil
}
func rewriteBodyHandler(prefix string) func(*http.Response) error {
	return func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body) //Read html
		if err != nil {
			return err
		}
		log.Println("rewriteBodyHandler", resp.Request.URL.String(), resp.StatusCode)
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		var res string
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			res, err = ChangePathOfRessources(string(b), prefix)
			if err != nil {
				return err
			}
		} else if strings.Contains(resp.Header.Get("Content-Type"), "text/css") {
			res = ChangePathOfRessourcesCss(string(b), prefix)
		} else {
			res = string(b)
		}
		body := ioutil.NopCloser(bytes.NewReader([]byte(res)))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		return nil
	}
}

func mimeTypeForFile(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"

	default:
		return mime.TypeByExtension(ext)
	}
}

type FileHandler struct {
	server    *WebSocketHandler
	socketUrl string
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

	handler, err := NewTemplateHandler(r, filehandler.socketUrl, id, filehandler.server)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
