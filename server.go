package main

import (
	"embed"
	"fmt"
	"net/http"

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
