package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// OrphanInfo describes a detected orphaned picoclaw-agents process.
type OrphanInfo struct {
	PID  int
	PPID int
	Cmd  string
}

// FindOrphans scans for picoclaw-agents processes whose parent PID is 1.
func FindOrphans() ([]OrphanInfo, error) {
	if runtime.GOOS == "windows" {
		return nil, nil
	}

	cmd := exec.Command("ps", "aux")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ps failed: %w", err)
	}

	var orphans []OrphanInfo
	myPID := os.Getpid()

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		pid, _ := strconv.Atoi(fields[1])
		ppid, _ := strconv.Atoi(fields[2])

		if pid == myPID {
			continue
		}

		cmdName := strings.Join(fields[10:], " ")
		if strings.Contains(cmdName, "picoclaw-agents") && ppid == 1 {
			orphans = append(orphans, OrphanInfo{PID: pid, PPID: ppid, Cmd: cmdName})
		}
	}

	return orphans, nil
}

// KillOrphans sends SIGTERM to each orphan process.
func KillOrphans() ([]OrphanInfo, error) {
	orphans, err := FindOrphans()
	if err != nil {
		return nil, err
	}

	killed := make([]OrphanInfo, 0, len(orphans))
	for _, o := range orphans {
		proc, err := os.FindProcess(o.PID)
		if err != nil {
			continue
		}
		if err := proc.Signal(os.Interrupt); err == nil {
			killed = append(killed, o)
		}
	}
	return killed, nil
}
