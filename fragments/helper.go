package fragments

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	promLoadFragmentsTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "micropuzzle_duration_load_milliseconds",
		Help:    "micropuzzle loading nanoseconds for microfrontends",
		Buckets: []float64{1, 5, 10, 30, 50, 80, 100, 1000},
	}, []string{"fragment", "frontend", "afterTimeout", "cached"})
)

func init() {
	prometheus.MustRegister(promLoadFragmentsTime)
}

func (sh *fragmentHandler) writePromMessage(options loadAsyncOptions, fromCache, insideTimeout bool, start time.Time) {
	promLoadFragmentsTime.WithLabelValues(options.FragmentName, options.Frontend, strconv.FormatBool(insideTimeout), strconv.FormatBool(fromCache)).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
func getLoggerWithUserInfo(logs *zap.SugaredLogger, user WebSocketUser) *zap.SugaredLogger {
	return logs.With("streamid", user.Id, "address", user.RemoteAddr)
}
