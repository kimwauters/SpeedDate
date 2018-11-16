package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/proepkes/speeddate/usersvc"
	repositorysvr "github.com/proepkes/speeddate/usersvc/gen/http/repository/server"
	swaggersvr "github.com/proepkes/speeddate/usersvc/gen/http/swagger/server"
	repository "github.com/proepkes/speeddate/usersvc/gen/repository"
	goahttp "goa.design/goa/http"
	"goa.design/goa/http/middleware"
)

type config struct {
	DbIP       string `env:"DB_IP" envDefault:"127.0.0.1"`
	DbPort     string `env:"DB_PORT" envDefault:"8888"`
	DbName     string `env:"DB_NAME" envDefault:"speeddate"`
	DbUser     string `env:"DB_USER" envDefault:"speeddateuser"`
	DbPassword string `env:"DB_PASSWORD" envDefault:""`
}

func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		//TODO: configurable listen address
		addr = flag.String("listen", ":8001", "HTTP listen `address`")
		dbg  = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	cfg := config{}
	env.Parse(&cfg)

	// Setup logger and goa log adapter. Replace logger with your own using
	// your log package of choice. The goa.design/middleware/logging/...
	// packages define log adapters for common log packages.
	var (
		adapter middleware.Logger
		logger  *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[usersvc] ", log.Ltime)
		adapter = middleware.NewLogger(logger)
	}

	// Initialize service dependencies such as databases.
	var (
		db *gorm.DB
	)
	{
		var err error
		//TODO: secure dbpass + configurable sslmode
		address := "postgresql://" + cfg.DbUser + "@" + cfg.DbIP + ":" + cfg.DbPort + "/" + cfg.DbName + "?sslmode=disable"
		log.Println("Connecting to " + address)

		db, err = gorm.Open("postgres", address)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db.LogMode(true)

		// Automatically create the "stored_users" table based on the StoredUser model.
		db.AutoMigrate(&repository.StoredUser{})
	}

	// Create the structs that implement the services.
	var (
		repositorySvc repository.Service
	)
	{
		var err error
		repositorySvc, err = usersvc.NewRepository(db, logger)
		if err != nil {
			logger.Fatalf("error creating database: %s", err)
		}
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
		repositoryEndpoints *repository.Endpoints
	)
	{
		repositoryEndpoints = repository.NewEndpoints(repositorySvc)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		repositoryServer *repositorysvr.Server
		swaggerServer    *swaggersvr.Server
	)
	{
		eh := ErrorHandler(logger)
		repositoryServer = repositorysvr.New(repositoryEndpoints, mux, dec, enc, eh)
		swaggerServer = swaggersvr.New(nil, mux, dec, enc, eh)
	}

	// Configure the mux.
	repositorysvr.Mount(mux, repositoryServer)
	swaggersvr.Mount(mux)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		if *dbg {
			handler = middleware.Debug(mux, os.Stdout)(handler)
		}
		handler = middleware.Log(adapter)(handler)
		handler = middleware.RequestID()(handler)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the service to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	srvHealth := http.NewServeMux()
	srvHealth.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	go func() {
		http.ListenAndServe(":8002", srvHealth)
	}()
	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: *addr, Handler: handler}
	go func() {
		for _, m := range repositoryServer.Mounts {
			logger.Printf("method %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
		}
		for _, m := range swaggerServer.Mounts {
			logger.Printf("file %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
		}
		logger.Printf("listening on %s", *addr)
		errc <- srv.ListenAndServe()
	}()

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	logger.Println("exited")
}

// ErrorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func ErrorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}
