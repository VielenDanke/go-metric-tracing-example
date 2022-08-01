package main

import (
	"context"
	"github.com/vielendanke/opentracing-example/internal/pkg/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/vielendanke/opentracing-example/internal/pkg/common"
	"github.com/vielendanke/opentracing-example/internal/pkg/trace"
	"github.com/vielendanke/opentracing-example/internal/users/handler"
	"github.com/vielendanke/opentracing-example/internal/users/service"
	"github.com/vielendanke/opentracing-example/internal/users/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	time.Sleep(5 * time.Second)

	cfg := trace.ProviderConfig{
		JaegerEndpoint: os.Getenv("JAEGER_ENDPOINT"),
		ServiceName:    "users",
		ServiceVersion: "1.0.0",
		Environment:    "integ",
		Disabled:       false,
	}
	exporter, exporterErr := trace.NewJaegerExporter(os.Getenv("JAEGER_ENDPOINT"))
	if exporterErr != nil {
		log.Println(exporterErr)
	}
	resource := trace.NewResource(cfg)
	pvd, providerErr := trace.NewProvider(cfg, exporter, resource)
	if providerErr != nil {
		log.Println(providerErr)
	}
	defer pvd.Close(context.Background())

	conn, connErr := sql.OpenSQLConnection(os.Getenv("DATABASE_URL"))

	if connErr != nil {
		log.Println(connErr)
	}
	migrErr := sql.SetupData(conn)

	if migrErr != nil {
		log.Println(migrErr)
	}
	userRepo := storage.NewUserRepository(conn)

	svc := service.NewUserService(userRepo)

	mux := http.DefaultServeMux

	h := handler.NewUserHandler(svc)

	mux.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(h.FindAll, "controller_user_find_all"))
	mux.HandleFunc("/live", common.LiveReadyProbe)
	mux.HandleFunc("/ready", common.LiveReadyProbe)

	log.Fatalln(http.ListenAndServe(":9090", mux))
}
