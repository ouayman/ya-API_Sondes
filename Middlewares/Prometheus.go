package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
)

const (
	reqsName    = "http_requests_total"
	latencyName = "http_request_duration_milliseconds"
)

type Prometheus struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewPrometheusMiddleware returns a new prometheus Middleware handler.
func NewPrometheus(name string, buckets ...float64) *Prometheus {
	var m Prometheus
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)

	return &m
}

func (m *Prometheus) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		nrw := newResponseWriter(w)
		next.ServeHTTP(nrw, r)

		m.reqs.WithLabelValues(strconv.Itoa(nrw.StatusCode), r.Method, r.URL.Path).Inc()
		m.latency.WithLabelValues(strconv.Itoa(nrw.StatusCode), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	})
}

/* Negroni middleware...
func (m *PrometheusMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)

	res := negroni.NewResponseWriter(rw)
	m.reqs.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Inc()
	m.latency.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
*/
