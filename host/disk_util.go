package host

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// GetDisks retrieves disk size and usage information use ubuntu command line.
func GetDisks() ([]DiskUsage, error) {
	var disks []DiskUsage

	// Get a list of all block devices
	devices, err := getBlockDevices()
	if err != nil {
		return nil, err
	}

	// Get the size and usage of each disk
	for _, device := range devices {
		size, err := getDiskSize(device)
		if err != nil {
			return nil, err
		}

		// Get the disk usage
		used, free, err := getDiskUsed(device)
		if err != nil {
			return nil, err
		}

		disks = append(disks, DiskUsage{
			Disk:  device,
			Total: size,
			Used:  used,
			Free:  free,
		})
	}
	log.Println("Disks", disks)
	return disks, nil
}

// getDiskUsed retrieves the used and free space of a disk in bytes.
func getDiskUsed(disk string) (uint64, uint64, error) {
	var stat unix.Statfs_t
	if err := unix.Statfs(disk, &stat); err != nil {
		return 0, 0, err
	}

	used := (stat.Blocks - stat.Bfree) * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)

	return used, free, nil
}

// getBlockDevices retrieves a list of block device names use command line.
func getBlockDevices() ([]string, error) {
	var devices []string

	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		// Skip the header line
		if fields[0] == "major" {
			continue
		}

		devices = append(devices, fmt.Sprintf("/dev/%s", fields[3]))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// getDiskSize retrieves the size of a disk in bytes.
func getDiskSize(disk string) (uint64, error) {
	fd, err := syscall.Open(disk, syscall.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(fd)

	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		return 0, err
	}

	return uint64(stat.Size), nil
}
