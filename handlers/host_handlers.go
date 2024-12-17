package handlers

import (
	"idv_host/host"
	"log"

	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// make struct for a Virtual Machine specification in very detail
// the VM specification is used to create a new VM instance
// the VM specification contains the name, the CPU, the memory, the disk size, the network interface, and the OS

type VM struct {
	Name           string `json:"name"`
	CPU            int    `json:"cpu"`
	Memory         int    `json:"memory"`
	DiskSize       int    `json:"diskSize"`
	NetworkAdapter string `json:"networkAdapter"`
	OS             string `json:"os"`
}

// CreateVM creates a new VM, show the VM creation form
func CreateVM(c *gin.Context) {
	c.HTML(http.StatusOK, "create_vm.html", gin.H{})
}

// Createa a new VM instance
func CreateVMInstance(c *gin.Context) {

	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created successfully"})
}

// RestartHost restarts the host machine
func RestartHost(c *gin.Context) {
	log.Println("Restarting host")
	cmd := exec.Command("sudo", "reboot")
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restart host"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Host is restarting..."})
}

// ResetHost resets the host machine
func ResetHost(c *gin.Context) {
	log.Println("Resetting host")
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

//
