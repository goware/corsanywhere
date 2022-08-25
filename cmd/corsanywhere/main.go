package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/goware/corsanywhere"
)

var (
	flags = flag.NewFlagSet("corsanywhere", flag.ExitOnError)
	fPort = flags.String("port", "8080", "Local port to listen for this corsanywhere service")
)

func main() {
	flags.Parse(os.Args[1:])

	port := *fPort

	fmt.Printf("CORS Anywhere started at http://localhost:%s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsanywhere.CORSAnywhereHandler())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
