package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/fasibio/micropuzzle/cache"
	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/fragments"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
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

//go:embed micro-lib/*.js
var embeddedLib embed.FS

type Runner struct{}

func (ru *Runner) Run(c *cli.Context) error {
	logs, err := logger.Initialize(c.String(CliLogLevel))
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	frontends, err := configloader.LoadFrontends(c.String(CliMicrofrontends))
	if err != nil {
		return err
	}
	r.Use(chiprometheus.NewMiddleware("micropuzzle"))
	r.Use(middleware.Compress(5))
	proxy.RegisterReverseProxy(r, frontends)

	cache, err := cache.NewRedisHandler(&redis.Options{
		Addr:     c.String(CliRedisAddress),
		DB:       c.Int(CliRedisDb),
		Username: c.String(CliRedisUser),
		Password: c.String(CliRedisPassword),
	})
	if err != nil {
		return err
	}
	fragmentHandler := fragments.NewFragmentHandler(cache, c.Duration(CliTimeout), frontends, c.String(CliFallbackLoader))
	fragmentHandler.RegisterHandler(r, SOCKET_PATH, SOCKET_ENDPOINT)
	f := FileHandler{
		server:    &fragmentHandler,
		socketUrl: SOCKET_PATH,
	}
	r.Handle(LIB_ENDPOINT, http.FileServer(http.FS(embeddedLib)))
	f.ChiFileServer(r, "/", http.Dir(c.String(CliPublicFolder)))

	logs.Infof("Start Server on Port :%s and Management on port %s", c.String(CliPort), c.String(CliManagementPort))
	go ru.StartManagementEndpoint(c.String(CliManagementPort))
	return http.ListenAndServe(fmt.Sprintf(":%s", c.String(CliPort)), r)

}

func (r *Runner) StartManagementEndpoint(port string) {
	managementR := chi.NewRouter()
	managementR.Handle(METRICS_ENDPOINT, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", port), managementR)
}
