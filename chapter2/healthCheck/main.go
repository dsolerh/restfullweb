package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// HealthCheck API returns date time to client
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	// io.WriteString(w, currentTime.String())
	fmt.Fprintf(w, "%s", currentTime)
}

func main() {
	http.HandleFunc("/health", HealthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
