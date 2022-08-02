package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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
		ServiceName:    os.Getenv("APPLICATION_NAME"),
		ServiceVersion: os.Getenv("APPLICATION_VERSION"),
		Environment:    os.Getenv("ENVIRONMENT"),
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

	router := mux.NewRouter()

	h := handler.NewUserHandler(svc)

	router.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(h.FindAll, "controller_user_find_all")).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(h.Save, "controller_user_save")).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/{id}", trace.HTTPHandlerFunc(h.FindByID, "controller_user_find_by_id")).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(h.Update, "controller_user_update")).Methods(http.MethodPut)
	router.HandleFunc("/live", common.LiveReadyProbe)
	router.HandleFunc("/ready", common.LiveReadyProbe)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("CONTAINER_PORT")), router))
}
