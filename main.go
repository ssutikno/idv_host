package main

import (
	"idv_host/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/index", handlers.IndexHandler)
	http.HandleFunc("/api/vms/start", handlers.StartVM)
	http.HandleFunc("/api/vms", handlers.ListVMs)
	http.HandleFunc("/api/vms/reboot", handlers.RebootVM)
	http.HandleFunc("/api/vms/reset", handlers.ResetVM)
	http.HandleFunc("/api/vms/shutdown", handlers.ShutdownVM)
	http.HandleFunc("/api/vms/poweroff", handlers.PowerOffVM)
	http.HandleFunc("/api/vms/create", handlers.CreateVM)

	http.HandleFunc("/api/host/network", handlers.GetNetworkData)

	http.HandleFunc("/api/host/restart", handlers.RestartHost)
	http.HandleFunc("/api/host/reset", handlers.ResetHost)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
