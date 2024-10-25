package main

import (
	"idv_host/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/api/vms/start", handlers.StartVM)
	http.HandleFunc("/api/vms", handlers.ListVMs)
	http.HandleFunc("/api/vms/reboot", handlers.RebootVM)
	http.HandleFunc("/api/vms/reset", handlers.ResetVM)
	http.HandleFunc("/api/vms/shutdown", handlers.ShutdownVM)
	http.HandleFunc("/api/vms/poweroff", handlers.PowerOffVM)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
