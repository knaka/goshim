package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// If env var is set, wait for a debugger to attach to this process
func waitForDebugger() {
	if os.Getenv("WAIT_FOR_DEBUGGER") == "" {
		return
	}
	pid := os.Getpid()
	log.Printf("Process %d is waiting\n", pid)
	for func() bool {
		time.Sleep(1 * time.Second)
		cmd := exec.Command("ps", "w")
		stdout, err := cmd.StdoutPipe()
		panicOn(err)
		defer func() { _ = stdout.Close() }()
		reader := bufio.NewReader(stdout)
		err = cmd.Start()
		panicOn(err)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF && len(line) == 0 {
				break
			}
			if strings.Contains(line, "dlv") &&
				strings.Contains(line, fmt.Sprintf("attach %d", pid)) {
				log.Println("Debugger connected")
				return false
			}
		}
		return true
	}() {
	}
}
