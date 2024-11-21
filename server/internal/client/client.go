package client

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stepan-pokladov/csw/server/internal/handlers"
	mw "github.com/stepan-pokladov/csw/server/internal/handlers/middleware"
	"github.com/stepan-pokladov/csw/server/internal/queue"
	"github.com/stepan-pokladov/csw/server/internal/report_processor/av_processor"
)

type Service struct {
	q queue.QueueProducer
}

func NewService(q queue.QueueProducer) *Service {
	return &Service{
		q: q,
	}
}

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors (non-2xx responses)",
		},
		[]string{"method", "status"},
	)
)

// init registers the prometheus metrics
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(errorCount)
}

// recordMetrics records metrics for each request
func recordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.Method, fmt.Sprint(rr.statusCode)).Inc()
		requestDuration.WithLabelValues(r.Method, fmt.Sprint(rr.statusCode)).Observe(duration)

		if rr.statusCode >= 400 {
			errorCount.WithLabelValues(r.Method, fmt.Sprint(rr.statusCode)).Inc()
		}
	})
}

// responseRecorder is a wrapper around http.ResponseWriter that records the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader records the status code
func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (s *Service) Run(port string) {
	rp := av_processor.NewAVProcessor(s.q)
	r := chi.NewRouter()
	r.Use(recordMetrics)
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Handle("/metrics", promhttp.Handler())

	r.With(mw.GzipDecompressorMiddleware).Post("/api/visit/v1", handlers.HandlerForRoute("visit", rp))
	r.With(mw.GzipDecompressorMiddleware).Post("/api/activity/v1", handlers.HandlerForRoute("activity", rp))
	fmt.Println("Server started at " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
