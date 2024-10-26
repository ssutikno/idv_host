// handlers/host.go
package handlers

import (
	"net/http"
	"os/exec"
)

// RestartHost restarts the host machine
func RestartHost(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sudo", "reboot")
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to restart host", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Host is restarting..."))
}

// ResetHost resets the host machine
func ResetHost(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sudo", "shutdown", "-r", "now")
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to reset host", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Host is resetting..."))
}
