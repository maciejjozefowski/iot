package main

import (
	"IoT_device_discovery/internal/devices"
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

func main() {

	var port int
	var username string = "postgres"
	var password string = "1234"
	var dbname string = "devices"
	var sslMode string = "disable"
	var dbPort int = 5432
	var dbHost string = "localhost"

	ctx := context.Background()
	connPool, err := pgxpool.NewWithConfig(ctx, PgxPoolConfig(username, password, dbname, sslMode, dbPort, dbHost))
	if err != nil {
		panic(err)
	}
	connection, err := connPool.Acquire(ctx)
	if err != nil {
		log.Fatal("Failed to acquire a connection pool to the database!!", err)
	}

	defer connection.Release()
	err = connection.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping the database!!")
	}
	fmt.Println("Connected to the database!!")

	flag.IntVar(&port, "port", 3000, "port to run the server on")
	fmt.Println("Server running on port: ", port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/devices", registerDevicesEndpoints(connPool))

	http.ListenAndServe(":3000", r)
}

func registerDevicesEndpoints(conn *pgxpool.Pool) chi.Router {
	router := chi.NewRouter()
	deviceQuerries := devices.New(conn)
	handler := devices.NewDeviceHandler(deviceQuerries)

	router.Post("/", handler.CreteDevice)
	router.Get("/", handler.GetDevicePage)
	router.Delete("/{id}", handler.DeleteDevice)
	router.Put("/{id}", handler.UpdateDevice)
	router.Get("/{id}", handler.GetDeviceByID)
	//router.Options("/", handler.Options)
	return router
}

func PgxPoolConfig(username string, password string, dbname string, sslMode string, dbPort int, dbHost string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5
	dbConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, username, password, dbname, sslMode)
	dbConfig, err := pgxpool.ParseConfig(dbConnString)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}
