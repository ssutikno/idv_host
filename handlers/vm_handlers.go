package handlers

import (
	"encoding/json"
	"html/template"
	"idv_host/host"
	"idv_host/vm"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/home.html"))
var index = template.Must(template.ParseFiles("templates/index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	vms, err := vm.ListVMs()
	if err != nil {
		http.Error(w, "Failed to list VMs", http.StatusInternalServerError)
		return
	}

	// log.Println("Listing VMs", vms)
	hostdata, err := host.GetHostData()
	if err != nil {
		http.Error(w, "Failed to get host data", http.StatusInternalServerError)
		return
	}

	data := struct {
		VMs  []vm.VM
		Host host.Host
	}{
		VMs:  vms,
		Host: hostdata,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func ListVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := vm.ListVMs()
	if err != nil {
		http.Error(w, "Failed to list VMs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vms)
}

func StartVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}

	err := vm.StartVM(vmName)
	if err != nil {
		http.Error(w, "Failed to start VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM started successfully"))
}

func RebootVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}

	err := vm.RebootVM(vmName)
	if err != nil {
		http.Error(w, "Failed to reboot VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM rebooted successfully"))
}

func ResetVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}

	err := vm.ResetVM(vmName)
	if err != nil {
		http.Error(w, "Failed to reset VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM reset successfully"))
}

func ShutdownVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}

	err := vm.ShutdownVM(vmName)
	if err != nil {
		http.Error(w, "Failed to shutdown VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM shutdown successfully"))
}

func PowerOffVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}

	err := vm.DestroyVM(vmName)
	if err != nil {
		http.Error(w, "Failed to power off VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM powered off successfully"))
}

func CreateVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing VM name", http.StatusBadRequest)
		return
	}
	vmXML := r.URL.Query().Get("xml")
	if vmXML == "" {
		http.Error(w, "Missing VM XML", http.StatusBadRequest)
		return
	}

	err := vm.CreateVM(vmName, vmXML)
	if err != nil {
		http.Error(w, "Failed to create VM", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("VM created successfully"))
}
