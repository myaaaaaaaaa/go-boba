package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var r = rand.Intn(1000000)

func main() {
	// Generate a random socket path
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("mpv-socket-%d", r))

	// Prepare arguments for mpv
	mpvArgs := []string{
		"--input-ipc-server=" + socketPath,
	}
	mpvArgs = append(mpvArgs, os.Args[1:]...)

	cmd := exec.Command("mpv", mpvArgs...)
	fmt.Println(mpvArgs)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start mpv: %v\n", err)
		os.Exit(1)
	}

	// Ensure cleanup of the process and socket
	defer func() {
		// We ignore the error here because the process might have already exited
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		os.Remove(socketPath)
	}()

	// Connect to the IPC socket
	var conn net.Conn

	for {
		var err error
		conn, err = net.Dial("unix", socketPath)
		if err == nil {
			defer conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Fprintln(conn, `{"command":["observe_property",1,"time-pos"]}`)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		parseTimePosEvent(scanner.Text(), func(timePos float64) {
			fmt.Printf("%.3f\n", timePos)
		})
	}

	// Wait for mpv to finish
	cmd.Wait()
}

func parseTimePosEvent(line string, timePosEvent func(timePos float64)) {
	var event struct {
		Event string
		Name  string
		Data  *float64
	}

	if err := json.Unmarshal([]byte(line), &event); err != nil {
		fmt.Println("invalid message:", line)
	}

	switch {
	case event.Event != "property-change":
	case event.Name != "time-pos":
	case event.Data == nil:
	default:
		timePosEvent(*event.Data)
	}
}
