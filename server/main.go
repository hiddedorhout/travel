package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hiddedorhout/travel/server/common"
)

var (
	port string
	db   string
)

func init() {
	flag.StringVar(&port, "port", "8000", "the local port")
	flag.StringVar(&db, "db", "serviceDB", "The db name")
	flag.Parse()
}

func main() {

	service, err := common.New(fmt.Sprintf("%s.db", db))
	if err != nil {
		log.Fatal(err)
	}

	service.SetupRoutes()

	log.Println(fmt.Sprintf("Service running on port: %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
