package vm

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// VM represents a virtual machine with its name and status
type VM struct {
	ID     string
	Name   string
	Status string
}

// createVM function creates a new virtual machine with the given name and XML definition.
func createVM(vmName, xml string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	if xml == "" {
		return errors.New("xml cannot be empty")
	}

	// print log
	fmt.Printf("Creating VM with name: %s\n", vmName)

	cmd := exec.Command("virsh", "create", xml)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Created VM with name: %s\n", vmName)
	return nil
}

// getVMXMLPath retrieves the path to the XML configuration file of a VM using virsh dominfo.
func getVMXMLPath(vmName string) (string, error) {
	// Execute the virsh dominfo command.
	cmd := exec.Command("virsh", "dominfo", vmName)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute virsh dominfo: %w", err)
	}

	// Find the line containing "Config file:"
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Config file:") {
			// Extract the path using a regular expression.
			re := regexp.MustCompile(`Config file:\s*(.*)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 2 {
				return strings.TrimSpace(matches[1]), nil
			}
		}
	}

	return "", fmt.Errorf("config file path not found for VM %s", vmName)
}

// ListVMs returns a list of all available VMs with their status using the virsh command
func ListVMs() ([]VM, error) {
	cmd := exec.Command("virsh", "list", "--all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var vms []VM

	for _, line := range lines[2:] { // Skip the first two lines of the output
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		vms = append(vms, VM{
			ID:     fields[0],
			Name:   fields[1],
			Status: fields[2],
		})
	}

	return vms, nil
}

// StartVM starts a virtual machine given its name
func StartVM(vmName string) error {

	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	cmd := exec.Command("virsh", "start", vmName)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Started VM with name: %s\n", vmName)
	return nil
}

// RebootVM reboots a virtual machine given its name
func RebootVM(vmName string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	cmd := exec.Command("virsh", "reboot", vmName)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Rebooted VM with name: %s\n", vmName)
	return nil
}

// ResetVM resets a virtual machine given its name
func ResetVM(vmName string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	cmd := exec.Command("virsh", "reset", vmName)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Reset VM with name: %s\n", vmName)
	return nil
}

// ShutdownVM shuts down a virtual machine given its name
func ShutdownVM(vmName string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	cmd := exec.Command("virsh", "shutdown", vmName)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Shutdown VM with name: %s\n", vmName)
	return nil
}

// DestroyVM destroys a virtual machine given its name
func DestroyVM(vmName string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}
	// print log
	fmt.Printf("Destroying VM with name: %s\n", vmName)

	cmd := exec.Command("virsh", "destroy", vmName)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Destroyed VM with name: %s\n", vmName)
	return nil
}

// CreateVM creates a virtual machine given its name and XML definition
func CreateVM(vmName, xml string) error {
	if vmName == "" {
		return errors.New("vmName cannot be empty")
	}

	if xml == "" {
		return errors.New("xml cannot be empty")
	}

	// print log
	fmt.Printf("Creating VM with name: %s\n", vmName)

	cmd := exec.Command("virsh", "create", xml)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Created VM with name: %s\n", vmName)
	return nil
}
