package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	// Command to list VMs using virsh
	cmd := exec.Command("virsh", "list", "--all")

	// Execute the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	// print the output
	fmt.Println(string(out))

	// Split the output into lines
	lines := strings.Split(string(out), "\n")

	// skip below if there are no VMs
	if len(lines) < 3 {
		fmt.Println("No VMs found")
		return
	}
	// Skip the first two lines (header) and the last line (empty)
	for _, line := range lines[2 : len(lines)-1] {
		// Split each line into fields
		fields := strings.Fields(line)

		// Extract VM ID, Name, and State
		vmID := fields[0]
		vmName := fields[1]
		vmState := fields[2]

		// Print the VM information
		fmt.Printf("ID: %s, Name: %s, State: %s\n", vmID, vmName, vmState)
	}
}
