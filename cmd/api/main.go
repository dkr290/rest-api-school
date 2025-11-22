package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/cmd/router"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/middleware"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/repository/sqlconnect"
)

var port = ":8082"

func main() {
	conf := config.LoadConfig()

	db, err := sqlconnect.ConnectDB(conf)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	llogger := logging.Init(conf.Debug)

	router := router.Router(db, *conf)

	rl := middleware.NewRateLimit(200, time.Minute)
	server := &http.Server{
		Addr: port,
		Handler: rl.Middleware(middleware.ResponseTimeMiddleware(
			middleware.SecurityHeaders(
				middleware.Cors(middleware.JWTMiddleware(router, *conf, *llogger)),
			),
		),
		),
	}

	fmt.Println("The server is starting on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Error Starting the server", err)
	}
}
