package utils

import (
	"fmt"
	"os"
	"runtime"
)

// ArchEnvironment contains detected information about the current system.
type ArchEnvironment struct {
	IsArchISO       bool
	IsRoot          bool
	Architecture    string
	HasInternet     bool
	VirtualMemoryMB int64
	DiskInfo        []DiskInfo
}

// DiskInfo contains basic information about a block device.
type DiskInfo struct {
	Device string
	Size   string
	Model  string
}

// DetectEnvironment probes the current system and returns environment info.
func DetectEnvironment() *ArchEnvironment {
	env := &ArchEnvironment{
		Architecture:    runtime.GOARCH,
		IsRoot:          os.Geteuid() == 0,
		IsArchISO:       checkArchISO(),
		HasInternet:     checkInternet(),
		VirtualMemoryMB: getVirtualMemory(),
		DiskInfo:        getDiskInfo(),
	}
	return env
}

// Summary returns a human-readable summary of the environment.
func (e *ArchEnvironment) Summary() string {
	s := fmt.Sprintf("Architecture: %s\n", e.Architecture)
	s += fmt.Sprintf("Root: %v\n", e.IsRoot)
	s += fmt.Sprintf("Arch ISO: %v\n", e.IsArchISO)
	s += fmt.Sprintf("Internet: %v\n", e.HasInternet)
	s += fmt.Sprintf("Memory: %d MB\n", e.VirtualMemoryMB)
	s += "Disks:\n"
	for _, d := range e.DiskInfo {
		s += fmt.Sprintf("  %s (%s)\n", d.Device, d.Size)
	}
	return s
}

// checkArchISO returns true if running on an Arch Linux live ISO.
func checkArchISO() bool {
	// Check for common indicators of Arch ISO
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return true
	}
	if _, err := os.Stat("/run/archiso"); err == nil {
		return true
	}
	return false
}

// checkInternet pings a known host to check connectivity.
func checkInternet() bool {
	// Try common DNS servers
	hosts := []string{"1.1.1.1", "8.8.8.8", "223.5.5.5"}
	for _, host := range hosts {
		if ping(host) {
			return true
		}
	}
	return false
}

// ping checks if a host is reachable (best-effort).
func ping(host string) bool {
	_, err := runCommand("ping", "-c", "1", "-W", "2", host)
	return err == nil
}

// getVirtualMemory returns total virtual memory in MB.
func getVirtualMemory() int64 {
	out, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	var totalKB int64
	_, err = fmt.Sscanf(string(out), "MemTotal: %d kB", &totalKB)
	if err != nil {
		return 0
	}
	return totalKB / 1024
}

// getDiskInfo returns information about available block devices.
func getDiskInfo() []DiskInfo {
	var disks []DiskInfo
	out, err := runCommand("lsblk", "-d", "-o", "NAME,SIZE,MODEL", "-n")
	if err != nil {
		return disks
	}

	lines := splitLines(out)
	for _, line := range lines {
		fields := split(line)
		if len(fields) >= 2 {
			d := DiskInfo{
				Device: "/dev/" + fields[0],
				Size:   fields[1],
			}
			if len(fields) >= 3 {
				d.Model = fields[2]
			}
			disks = append(disks, d)
		}
	}
	return disks
}