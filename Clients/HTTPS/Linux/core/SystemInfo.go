package core

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func MachineID() string {
	if runtime.GOOS == "linux" {
		// Standard paths
		id, err := os.ReadFile("/var/lib/dbus/machine-id")
		if err != nil {
			id, err = os.ReadFile("/etc/machine-id")
		}
		// Fallback for BusyBox/Minimal systems without machine-id
		if err != nil {
			// Try to get first MAC address as ID
			ifas, _ := net.Interfaces()
			for _, ifa := range ifas {
				if ifa.HardwareAddr.String() != "" {
					return MD5Hash(ifa.HardwareAddr.String())
				}
			}
			// Last resort: Hostname
			host, _ := os.Hostname()
			return MD5Hash(host)
		}
		return strings.TrimSpace(strings.Trim(string(id), "\n"))
	} else if runtime.GOOS == "darwin" {
		buf := &bytes.Buffer{}
		err := run(buf, os.Stderr, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		if err != nil {
			return "ERROR"
		}
		id := extractID(buf.String())
		if err != nil {
			return "ERROR"
		}
		return strings.TrimSpace(strings.Trim(id, "\n"))
	} else if runtime.GOOS == "bsd" {
		buf, err := ioutil.ReadFile("/etc/hostid")
		if err != nil {
			buf := &bytes.Buffer{}
			err := run(buf, os.Stderr, "kenv", "-q", "smbios.system.uuid")
			if err != nil {
				return ""
			}
			return strings.TrimSpace(strings.Trim(buf.String(), "\n"))
		}
		return strings.TrimSpace(strings.Trim(string(buf), "\n"))
	}
	return "ERROR"
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func extractID(lines string) string {
	for _, line := range strings.Split(lines, "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.SplitAfter(line, `" = "`)
			if len(parts) == 2 {
				return strings.TrimRight(parts[1], `"`)
			}
		}
	}
	return ""
}

func GetRAM(mb int) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	rmb := uint64(mb)
	ram := m.TotalAlloc / 1024 / 1024

	return ram < rmb
}

func GetCPUCores(cores int) bool {
	x := false
	num_procs := runtime.NumCPU()
	if !(num_procs >= cores) {
		x = true
	}
	return x
}

func GetClientIP() string {
	rsp, _ := http.Get("https://checkip.amazonaws.com/")
	if rsp.StatusCode == 200 {
		defer rsp.Body.Close()
		buf, _ := ioutil.ReadAll(rsp.Body)
		return string(bytes.TrimSpace(buf))
	}
	return "127.0.0.1"
}

func GetClientPath() string {
	return os.Args[0]
}
