package model

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"KeyboardLayout", cfg.KeyboardLayout, "us"},
		{"Hostname", cfg.Hostname, "arch-server"},
		{"NetworkDHCP", cfg.NetworkDHCP, true},
		{"FilesystemType", cfg.FilesystemType, "ext4"},
		{"BootloaderType", cfg.BootloaderType, "grub"},
		{"UEFIMode", cfg.UEFIMode, true},
		{"Locale", cfg.Locale, "en_US.UTF-8"},
		{"KernelType", cfg.KernelType, "linux"},
		{"SSHPort", cfg.SSHPort, 22},
		{"EnableSSH", cfg.EnableSSH, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		{"valid simple", "server", true},
		{"valid with hyphen", "my-server", true},
		{"valid with number", "server01", true},
		{"valid with domain", "server.example.com", true},
		{"empty", "", false},
		{"too long", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", false},
		{"with space", "my server", false},
		{"with special char", "server$", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateHostname(tt.hostname)
			if got != tt.want {
				t.Errorf("ValidateHostname(%q) = %v, want %v", tt.hostname, got, tt.want)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"valid IPv4", "192.168.1.1", true},
		{"valid IPv4", "10.0.0.1", true},
		{"valid IPv4", "172.16.0.1", true},
		{"valid IPv4", "8.8.8.8", true},
		{"empty", "", false},
		{"too many octets", "1.2.3.4.5", false},
		{"invalid octet", "256.1.2.3", false},
		{"with letters", "abc.def.ghi.jkl", false},
		{"missing octet", "1.2.3", false},
		{"leading zeros allowed", "192.168.001.001", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateIP(tt.ip)
			if got != tt.want {
				t.Errorf("ValidateIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestDefaultConfigPointer(t *testing.T) {
	cfg1 := DefaultConfig()
	cfg2 := DefaultConfig()

	// Configs should be independent (different pointers)
	if cfg1 == cfg2 {
		t.Error("DefaultConfig() returned same pointer twice")
	}

	// Modifying one should not affect the other
	cfg1.Hostname = "server1"
	if cfg2.Hostname == "server1" {
		t.Error("Configs are not independent")
	}
}