package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/fasibio/micropuzzle/cache"
	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/filehandler"
	"github.com/fasibio/micropuzzle/fragments"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

/* Customer endpoints */
const (
	SOCKET_PATH     = "socket"
	SOCKET_ENDPOINT = "/micro-puzzle"
	LIB_ENDPOINT    = "/micro-lib/*"
)

/* Management endpoints*/
const (
	METRICS_ENDPOINT   = "/metrics"
	HEALTH_ENDPOINT    = "/health"
	FRONTENDS_ENDPOINT = "/frontends"
)

//go:embed micro-lib/*.js
var embeddedLib embed.FS

type runner struct{}

func NewRunner() *runner {
	return &runner{}
}

func (ru *runner) Run(c *cli.Context) error {
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
	fragmentHandler := fragments.NewFragmentHandler(cache, cache, c.Duration(CliTimeout), frontends, c.String(CliFallbackLoader))
	fragmentHandler.RegisterHandler(r, SOCKET_PATH, SOCKET_ENDPOINT)
	f := filehandler.NewFileHandler(&fragmentHandler, SOCKET_PATH)

	r.Handle(LIB_ENDPOINT, http.FileServer(http.FS(embeddedLib)))
	f.RegisterFileHandler(r, "/", http.Dir(c.String(CliPublicFolder)))

	logs.Infof("Start Server on Port :%s and Management on port %s", c.String(CliPort), c.String(CliManagementPort))
	go ru.StartManagementEndpoint(c.String(CliManagementPort), frontends)
	return http.ListenAndServe(fmt.Sprintf(":%s", c.String(CliPort)), r)

}

type FrontedsManagementResult struct {
	Frontends []string `json:"frontends"`
}

type HealthCheckObj struct {
	GitHash     string `json:"git_hash"`
	Version     string `json:"version"`
	ServiceName string `json:"service_name"`
}

func (r *runner) StartManagementEndpoint(port string, frontends configloader.Frontends) {
	managementR := chi.NewRouter()
	managementR.Handle(METRICS_ENDPOINT, promhttp.Handler())
	managementR.Get(HEALTH_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthCheckObj{
			GitHash:     os.Getenv("COMMIT_SHA"),
			Version:     os.Getenv("APPLICATION_BUILD_ID"),
			ServiceName: "mircopuzzle",
		})
	})
	managementR.Get(FRONTENDS_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FrontedsManagementResult{
			Frontends: frontends.GetKeyList(),
		})
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), managementR)
}
