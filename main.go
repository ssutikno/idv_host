package main

import (
	"encoding/json"
	"fmt"
	"idv_host/kvm"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", menuHandler)
	http.HandleFunc("/vms", listVMsHandler)
	http.HandleFunc("/vms/reboot/", rebootVMHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// make menu for the user to choose the action
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Choose an action:")
	fmt.Fprintln(w, "<ul>")
	fmt.Fprintln(w, "<li><a href='/vms'>List VMs</a></li>")
	fmt.Fprintln(w, "<li><a href='/vms/reboot/'>Reboot VM</a></li>")
	fmt.Fprintln(w, "</ul>")
	fmt.Fprintln(w, "Enter the name of the VM to reboot:")
	fmt.Fprintln(w, "<form method='post' action='/vms/reboot/'>")
	fmt.Fprintln(w, "<input type='text' name='vmName'>")
	fmt.Fprintln(w, "<input type='submit' value='Reboot'>")
	fmt.Fprintln(w, "</form>")

	if r.Method == http.MethodPost {

		vmName := r.FormValue("vmName")
		if vmName == "" {
			http.Error(w, "VM name is required", http.StatusBadRequest)
			return
		}

		err := kvm.RebootVM(vmName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "VM %q rebooted successfully", vmName)
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
