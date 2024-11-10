// handlers/host.go
package handlers

import (
	"idv_host/host"

	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// RestartHost restarts the host machine
func RestartHost(c *gin.Context) {
	cmd := exec.Command("sudo", "reboot")
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restart host"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Host is restarting..."})
}

// ResetHost resets the host machine
func ResetHost(c *gin.Context) {
	cmd := exec.Command("sudo", "shutdown", "-r", "now")
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset host"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Host is resetting..."})
}

func GetNetworkData(c *gin.Context) {
	network, err := host.GetNetworkInterfaces() // call GetNetworkInterfaces
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network data"})
		return
	}

	c.JSON(http.StatusOK, network)
}
