// make host information data available to the template.
package host

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Host represents the host information
type Host struct {
	Hostname string
	IP       string
	CPU      string
	Memory   Memory
	Disks    []DiskUsage
	UpTime   string
}

// DiskUsage contains information about disk space
type DiskUsage struct {
	Path        string  `json:"path"`
	Disk        string  `json:"disks"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// make network interfaces available to the template.
// the content of the network interfaces is displayed in the web page.
// values of network interface are including the name of the interface and the IPs Information
// NetWorkInterface contains information about network interfaces including the name, ips, mac address and status

type NetWorkInterface struct {
	Name   string   `json:"name"`
	IPs    []string `json:"ips"`
	MAC    string   `json:"mac"`
	Status string   `json:"status"`
}

// Process represents the process information
type Process struct {
	PID       int    `json:"pid"`
	Name      string `json:"name"`
	Cmd       string `json:"cmd"`
	StartTime string `json:"startTime"`
}

// Memory represents the memory information
type Memory struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Available   uint64  `json:"available"`
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

	memory, err := GetMemory()
	if err != nil {
		return Host{}, err
	}

	disks, err := GetDiskUsage()
	if err != nil {
		return Host{}, err
	}
	log.Printf("Disks: %v", disks)
	cmd = exec.Command("uptime", "-p")
	uptime, err := cmd.Output()
	if err != nil {
		return Host{}, err
	}

	host := Host{
		Hostname: strings.TrimSpace(string(hostname)),
		IP:       strings.TrimSpace(string(ip)),
		CPU:      cpu,
		Memory:   memory,
		Disks:    disks,
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

func GetMemory() (Memory, error) {
	var memory Memory
	cmd := exec.Command("free", "-b")
	output, err := cmd.Output()
	if err != nil {
		return memory, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Skip the header line
	scanner.Scan()
	if scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 7 {
			return memory, fmt.Errorf("invalid free output")
		}

		memory.Total = parseSize(parts[1])
		memory.Used = parseSize(parts[2])
		memory.Free = parseSize(parts[3])
		memory.Available = parseSize(parts[6])
		memory.UsedPercent = float64(memory.Used) / float64(memory.Total) * 100
	}
	return memory, nil
}

// function GetDiskUsage in array of DiskUsage struct
func GetDiskUsage() ([]DiskUsage, error) {
	var disk []DiskUsage
	var partition DiskUsage

	cmd := exec.Command("df")
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

		// Parse the path information into DiskUsage struct which contains the path is started with '/dev/' and the total, used, free and used percent information
		if strings.HasPrefix(parts[0], "/dev/") {
			partition.Path = parts[0]
			partition.Total = parseSize(parts[1])
			partition.Used = parseSize(parts[2])
			partition.Free = parseSize(parts[3])
			partition.UsedPercent = parsePercent(parts[4])

			disk = append(disk, partition)
		}
	}
	return disk, nil
}

// function parseSize
func parseSize(size string) uint64 {
	var value uint64
	fmt.Sscanf(size, "%d", &value)
	return value
}

// function parsePID
func parsePID(pidStr string) (int, error) {
	var pid int
	_, err := fmt.Sscanf(pidStr, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// function parsePercent
func parsePercent(percent string) float64 {
	var value float64
	fmt.Sscanf(percent, "%f%%", &value)
	return value
}

func GetNetworkInterfaces() ([]NetWorkInterface, error) {
	var interfaces []NetWorkInterface
	cmd := exec.Command("ip", "a")
	output, err := cmd.Output()
	if err != nil {
		return interfaces, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var iface NetWorkInterface
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mtu") {
			if iface.Name != "" {
				interfaces = append(interfaces, iface)
			}
			iface = NetWorkInterface{}
			iface.Name = strings.Split(line, ":")[1]
		} else if strings.Contains(line, "inet") {
			parts := strings.Fields(line)
			iface.IPs = append(iface.IPs, parts[1])
		}
		// the mac address is in the same line as the name of the interface
		if strings.Contains(line, "link/ether") {
			parts := strings.Fields(line)
			iface.MAC = parts[1]
		}
		// the status of the interface is in the same line as the name of the interface
		if strings.Contains(line, "state") {
			parts := strings.Fields(line)
			// if "state" is in the 7th position, the status is in the 8th position. else, the status is in the 10th position
			if parts[7] == "state" {
				iface.Status = parts[8]
			} else {
				iface.Status = parts[10]
			}
			iface.Status = parts[8]
		}
	}

	if err := scanner.Err(); err != nil {
		return interfaces, err
	}

	return interfaces, nil
}

func GetProcesses() ([]Process, error) {
	var processes []Process
	cmd := exec.Command("ps", "-e", "-o", "pid,comm,lstart")
	output, err := cmd.Output()
	if err != nil {
		return processes, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Skip the header line
	scanner.Scan()
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 3 {
			continue
		}

		pid, err := parsePID(parts[0])
		if err != nil {
			continue
		}

		process := Process{
			PID:       pid,
			Name:      parts[1],
			StartTime: strings.Join(parts[2:], " "),
		}
		processes = append(processes, process)
	}
	return processes, nil
}

func KillProcess(pid string) error {
	cmd := exec.Command("kill", pid)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
