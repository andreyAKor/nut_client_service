package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/nut_client_service/internal/http/clients/nut"
	handlerCommand "github.com/andreyAKor/nut_client_service/internal/http/server/handlers/command"
	handlerGet "github.com/andreyAKor/nut_client_service/internal/http/server/handlers/get"
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "nut_client_service",
	Name:      "http_response_time_seconds",
	Help:      "Duration of HTTP requests.",
	Buckets:   []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
}, []string{"method", "path"})

var (
	ErrServerNotInit  = errors.New("server not init")
	ErrInvalidRequest = errors.New("the request body can’t be parsed as valid data")

	_ io.Closer = (*Server)(nil)
)

type Server struct {
	host      string
	port      int
	bodyLimit int

	nutClient *nut.Client

	server *http.Server
	ctx    context.Context
}

func New(host string, port int, bodyLimit int, nutClient *nut.Client) (*Server, error) {
	return &Server{
		host:      host,
		port:      port,
		bodyLimit: bodyLimit,
		nutClient: nutClient,
	}, nil
}

// Run Running http-server.
func (s *Server) Run(ctx context.Context) error {
	s.ctx = ctx

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/get", s.method(s.toJSON(handlerGet.New(s.nutClient).Handle()), "GET"))
	mux.HandleFunc("/command", s.method(s.toJSON(handlerCommand.New(s.nutClient).Handle()), "POST"))

	// middlewares
	handler := s.metrics(mux)
	handler = s.headers(handler)
	handler = s.body(handler)
	handler = s.logger(handler)

	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.host, strconv.Itoa(s.port)),
		Handler: handler,
	}
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http-server listen fail")
	}

	return nil
}

func (s *Server) Close() error {
	if s.server == nil {
		//nolint:wrapcheck
		return ErrServerNotInit
	}

	return s.server.Shutdown(s.ctx)
}

// metrics Middleware sets metrics to prometheus.
func (s Server) metrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		handler.ServeHTTP(w, r)
	})
}

// headers Middleware sets http-headers for response.
func (s Server) headers(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST")
		}

		// For OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		// JSON header
		w.Header().Set("Content-Type", "application/json")

		handler.ServeHTTP(w, r)
	})
}

// logger Middleware logger output log info of request, e.g.: r.Method, r.URL etc.
func (s Server) logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newAppResponseWriter(w)

		start := time.Now()
		defer func() {
			i := log.Info()

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				i.Err(err)
			}

			i.Str("ip", host).
				Str("startAt", start.String()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("proto", r.Proto).
				Int("status", rw.statusCode).
				TimeDiff("latency", time.Now(), start)

			if len(r.UserAgent()) > 0 {
				i.Str("userAgent", r.UserAgent())
			}

			i.Msg("http-request")
		}()

		handler.ServeHTTP(rw, r)
	})
}

// body Middleware preparing body request.
func (s Server) body(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, int64(s.bodyLimit)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("body read fail")

			if err := s.writeJSON(Response{Error: ErrInvalidRequest.Error()}, w); err != nil {
				log.Error().Err(err).Msg("writeJSON fail")
			}

			return
		}

		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("body close fail")

			if err := s.writeJSON(Response{Error: ErrInvalidRequest.Error()}, w); err != nil {
				log.Error().Err(err).Msg("writeJSON fail")
			}

			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		handler.ServeHTTP(w, r)
	})
}

// method Checking allowed method for endpoint.
func (s Server) method(handler http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		handler(w, r)
	}
}

// toJSON Converting Response from endpoint to json-response.
func (s Server) toJSON(h func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rs Response

		data, err := h(w, r)
		if err != nil {
			rs.Error = err.Error()
		} else {
			rs.Data = data
		}

		if err := s.writeJSON(rs, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("writeJSON fail")

			return
		}
	}
}

// writeJSON Writing Response structure to json http-response.
func (s Server) writeJSON(rs Response, w io.Writer) error {
	res, err := json.Marshal(&rs)
	if err != nil {
		return errors.Wrap(err, "JSON-marshal fail")
	}

	if _, err := w.Write(res); err != nil {
		return errors.Wrap(err, "write fail")
	}

	return nil
}

var _ http.ResponseWriter = (*appResponseWriter)(nil)

// appResponseWriter App wrapper over http.ResponseWriter.
type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newAppResponseWriter(w http.ResponseWriter) *appResponseWriter {
	return &appResponseWriter{w, http.StatusOK}
}

func (a *appResponseWriter) WriteHeader(code int) {
	a.statusCode = code
	a.ResponseWriter.WriteHeader(code)
}
