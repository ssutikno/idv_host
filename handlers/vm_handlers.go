package handlers

import (
	"html/template"
	"idv_host/host"
	"idv_host/vm"
	"log"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
)

// var tmpl = template.Must(template.New("home").Funcs(template.FuncMap{"humanizeBytes": humanizeBytes}).ParseFiles("templates/home.html"))

func humanizeBytes(size uint64) string {
	log.Println("humanizeBytes")
	return humanize.Bytes(size)
}

var tmpl = template.Must(template.ParseFiles("templates/home.html"))

var index = template.Must(template.ParseFiles("templates/index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func HomeHandler(c *gin.Context) {
	vms, err := vm.ListVMs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list VMs"})
		return
	}

	hostdata, err := host.GetHostData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get host data"})
		return
	}

	data := struct {
		VMs  []vm.VM
		Host host.Host
	}{
		VMs:  vms,
		Host: hostdata,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
	}
}

func ListVMs(c *gin.Context) {
	vms, err := vm.ListVMs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list VMs"})
		return
	}

	c.JSON(http.StatusOK, vms)
}

func StartVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}

	err := vm.StartVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM started successfully"})
}

func RebootVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}

	err := vm.RebootVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reboot VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM rebooted successfully"})
}

func ResetVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}

	err := vm.ResetVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM reset successfully"})
}

func ShutdownVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}

	err := vm.ShutdownVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shutdown VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM shutdown successfully"})
}

func PowerOffVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}

	err := vm.DestroyVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to power off VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM powered off successfully"})
}

func CreateVM(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM name"})
		return
	}
	vmXML := c.Query("xml")
	if vmXML == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing VM XML"})
		return
	}

	err := vm.CreateVM(vmName, vmXML)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create VM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created successfully"})
}
