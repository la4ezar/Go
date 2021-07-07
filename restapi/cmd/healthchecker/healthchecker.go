// Checks if the server is running
package main // import "github.com/la4ezar/restapi"

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "8080", "port on localhost to check")
	flag.Parse()

	resp, err := http.Get("http://localhost:" + *port + "/api/health")

	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}

	os.Exit(0)
}
