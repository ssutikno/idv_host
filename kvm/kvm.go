package kvm

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// ListVMs lists all virtual machines managed by libvirt.
func ListVMs() ([]string, error) {
	cmd := exec.Command("virsh", "list", "--all", "--json") // Add --json flag
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list VMs: %w", err)
	}

	var result struct {
		VMs []struct {
			Name string `json:"name"`
		} `json:"vms"`
	}
	err = json.Unmarshal(out, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	var vms []string
	for _, vm := range result.VMs {
		vms = append(vms, vm.Name)
	}

	return vms, nil
}

// RebootVM reboots the specified virtual machine.
func RebootVM(name string) error {
	cmd := exec.Command("virsh", "reboot", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to reboot VM %q: %w\nOutput: %s", name, err, out)
	}

	return nil
}
