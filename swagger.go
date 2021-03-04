package main

import (
	"flag"
	"fmt"
	"github.com/smoug25/go-swagger-ui/server"
	"log"
	"net/http"
)

var (
	// Build of git, got by LDFLAGS on build
	Build = "-unknown-"
	// Version of git, got by LDFLAGS on build
	Version = "-unknown-"
)

var (
	serverAddr  = flag.String("l", ":8080", "server's listening Address")
	swaggerPath  = flag.String("p", "/", "swagger url path")
	swaggerFile = flag.String("f",
		"http://petstore.swagger.io/v2/swagger.json",
		"swagger url or local file path")
	enableTopbar    = flag.Bool("b", false, "enable the topbar")
)

func main() {
	flag.Parse()
	fmt.Printf("Server listening on %s and path %s\n", *serverAddr, *swaggerPath)
	server.SwaggerPath = *swaggerPath

	// test if swagger file is a local one
	server.SetSwaggerFile(*swaggerFile)
	if server.IsNativeSwaggerFile {
		fmt.Printf("Using default local swagger file %s\n", *swaggerFile)
	} else {
		fmt.Printf("Using default online swagger file %s\n", *swaggerFile)
	}
	if *enableTopbar {
		server.EnableTopbar = true
		fmt.Println("Topbar enabled")
	} else {
		fmt.Println("Topbar disabled")
	}
	fmt.Println("Swagger UI version", Version, ", build", Build)
	http.HandleFunc(*swaggerPath, server.Serv)
	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}
