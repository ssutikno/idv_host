package handlers

import (
	"idv_host/host"
	"log"

	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// CreateVM creates a new VM, show the VM creation form
func CreateVM(c *gin.Context) {
	c.HTML(http.StatusOK, "create_vm.html", gin.H{})
}

// Createa a new VM instance
func CreateVMInstance(c *gin.Context) {

	// get the VM name and the VM image from the form
	vmName := c.PostForm("vm_name")
	vmImage := c.PostForm("vm_image")

	// create a new VM instance
	err := host.CreateVM(vmName, vmImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create VM"})
		return
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
