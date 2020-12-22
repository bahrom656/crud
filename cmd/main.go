package main

import (
	"context"
	"github.com/bahrom656/crud/cmd/app"
	"github.com/bahrom656/crud/pkg/customers"
	"github.com/bahrom656/crud/pkg/managers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	postgres := "postgres://app:pass@localhost:5432/db"
	err := execute(host, port, postgres)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
func execute(host, port, postgres string, ) (err error) {
	dependencies := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			conn, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(conn, postgres)
		},
		customers.NewService,
		managers.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    host + ":" + port,
				Handler: server,
			}
		},
	}

	container := dig.New()

	for _, v := range dependencies {
		err = container.Provide(v)
		if err != nil {
			return err
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})

	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})

}
