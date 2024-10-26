// make host information data available to the template.
package host

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Host represents the host information
type Host struct {
	Hostname string
	IP       string
	CPU      string
	Memory   string
	Disk     []DiskUsage
	UpTime   string
}

// DiskUsage contains information about disk space
type DiskUsage struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// GetHostData returns the host information using linux commands
func GetHostData() (Host, error) {
	cmd := exec.Command("hostname")
	hostname, err := cmd.Output()
	if err != nil {
		return Host{}, err
	}

	cmd = exec.Command("hostname", "-I")
	ip, err := cmd.Output()
	if err != nil {
		return Host{}, err
	}

	cpu, err := GetCPUModel()
	if err != nil {
		return Host{}, err
	}

	cmd = exec.Command("free", "-h")
	memory, err := cmd.Output()
	if err != nil {
		return Host{}, err
	}

	disk, err := GetDiskUsage()
	if err != nil {
		return Host{}, err
	}

	cmd = exec.Command("uptime", "-p")
	uptime, err := cmd.Output()
	if err != nil {
		return Host{}, err
	}

	host := Host{
		Hostname: strings.TrimSpace(string(hostname)),
		IP:       strings.TrimSpace(string(ip)),
		CPU:      cpu,
		Memory:   string(memory),
		Disk:     disk,
		UpTime:   string(uptime),
	}
	return host, nil
}

// function getCPUModel
func GetCPUModel() (string, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			// Split the line and get the CPU model name
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("CPU model name not found")
}

// function GetDiskUsage in array of DiskUsage struct
func GetDiskUsage() ([]DiskUsage, error) {
	var disk []DiskUsage
	var partition DiskUsage

	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return disk, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Skip the header line
	scanner.Scan()
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 6 {
			continue
		}

		// Parse the path information into DiskUsage struct
		partition.Path = parts[0]
		partition.Total = parseSize(parts[1])
		partition.Used = parseSize(parts[2])
		partition.Free = parseSize(parts[3])
		partition.UsedPercent = parsePercent(parts[4])

		disk = append(disk, partition)
	}
	return disk, nil
}

// function parseSize
func parseSize(size string) uint64 {
	var value uint64
	fmt.Sscanf(size, "%d", &value)
	return value
}

// function parsePercent
func parsePercent(percent string) float64 {
	var value float64
	fmt.Sscanf(percent, "%f%%", &value)
	return value
}
