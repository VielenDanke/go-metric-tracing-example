package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/vielendanke/opentracing-example/internal/pkg/common"
	"github.com/vielendanke/opentracing-example/internal/pkg/sql"
	"github.com/vielendanke/opentracing-example/internal/pkg/trace"
	"github.com/vielendanke/opentracing-example/internal/users/handler"
	"github.com/vielendanke/opentracing-example/internal/users/service"
	"github.com/vielendanke/opentracing-example/internal/users/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := trace.ProviderConfig{
		JaegerEndpoint: "",
		ServiceName:    "service",
		ServiceVersion: "1.0.0",
		Environment:    "env",
		Disabled:       false,
	}
	exporter, exporterErr := trace.NewStdOutExporter(os.Stdout)
	if exporterErr != nil {
		log.Fatalln(exporterErr)
	}
	resource := trace.NewResource(cfg)
	pvd, providerErr := trace.NewProvider(cfg, exporter, resource)
	if providerErr != nil {
		log.Fatalln(providerErr)
	}
	defer pvd.Close(context.Background())

	conn, connErr := sql.OpenSQLConnection(
		"host=localhost port=5432 user=user password=password sslmode=disable dbname=users")

	if connErr != nil {
		log.Fatalln(connErr)
	}
	userRepo := storage.NewUserRepository(conn)

	svc := service.NewUserService(userRepo)

	mux := http.DefaultServeMux

	h := handler.NewUserHandler(svc)

	mux.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(h.FindAll, "controller_user_find_all"))
	mux.HandleFunc("/live", common.LiveReadyProbe)
	mux.HandleFunc("/ready", common.LiveReadyProbe)

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
