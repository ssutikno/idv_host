package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ssutikno/idv_host/modules/kvm"
)

func main() {
	http.HandleFunc("/", menuHandler)
	http.HandleFunc("/vms", listVMsHandler)
	http.HandleFunc("/vms/reboot/", rebootVMHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "KVM Management API")
	fmt.Fprintln(w, "-----------------")
	fmt.Fprintln(w, "Available actions:")
	fmt.Fprintln(w, "1. List all VMs")
	fmt.Fprintln(w, "2. Reboot a VM")
	fmt.Fprintln(w, "Enter your choice: ")

	reader := bufio.NewReader(r.Body)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	switch choice {
	case "1":
		listVMsHandler(w, r)
	case "2":
		fmt.Fprintln(w, "Enter the name of the VM to reboot: ")
		vmName, _ := reader.ReadString('\n')
		vmName = strings.TrimSpace(vmName)
		r.URL.Path = "/vms/reboot/" + vmName // Simulate URL path for reboot handler
		rebootVMHandler(w, r)
	default:
		fmt.Fprintln(w, "Invalid choice.")
	}
}

func listVMsHandler(w http.ResponseWriter, r *http.Request) {
	vms, err := kvm.ListVMs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vms)
}

func rebootVMHandler(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Path[len("/vms/reboot/"):]
	if vmName == "" {
		http.Error(w, "VM name is required", http.StatusBadRequest)
		return
	}

	err := kvm.RebootVM(vmName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "VM rebooted successfully")
}
