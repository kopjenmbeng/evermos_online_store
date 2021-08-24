package api

import (
	"net/http"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/db_context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/cors"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware"
	newrelic "github.com/newrelic/go-agent"
)

type Server struct {
	address    string
	httpServer *http.Server

	log     *log.Logger
	router  *chi.Mux
	apm     newrelic.Application
	readDb  *sqlx.DB
	writeDb *sqlx.DB
	jwe     *jwe_auth.JWE
}

type Option func(*Server)

func NewServer(address string, log *log.Logger, apm newrelic.Application, readDb *sqlx.DB, writeDb *sqlx.DB ,options ...Option) *Server {
	s := &Server{address: address, log: log, apm: apm, readDb: readDb, writeDb: writeDb}
	for _, o := range options {
		o(s)
	}
	s.Setup()
	return s
}

func (srv *Server) Setup() {
	srv.router = srv.Route()
	srv.httpServer = &http.Server{Addr: srv.address, Handler: srv.router}
}

func (srv *Server) Route() *chi.Mux {
	r := chi.NewRouter()
	// compressor := chim.NewCompressor(flate.DefaultCompression)
	r.Use(
		chim.NoCache,
		chim.RedirectSlashes,
		chim.Heartbeat("/ping"),
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestLogger(srv.log),
		// middleware.Header(),
		// telemetry.Telemetry(srv.apm),
		db_context.DBR(srv.readDb),
		db_context.DBW(srv.writeDb),
		// db_context.Redis(srv.redisClient),
		jwe_auth.InitJWE(srv.jwe),
		// auth.InitJWERefresh(srv.jw_refresh),
		chim.RequestID,
		chim.RealIP,
		chim.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			// MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	routes(r)
	return r
}

func (server *Server) Serve() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Infof("%s %s", method, route)
		return nil
	}
	if err := chi.Walk(server.router, walkFunc); err != nil {
		log.Panicln(errors.Cause(err))
	}
	// server.log.Println(server.address)

	log.Fatal(http.ListenAndServe(server.address, server.router))
	// server.Route()
	// Router()
}
