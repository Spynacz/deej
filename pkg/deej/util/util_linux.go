package util

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	getCurrentWindowInternalCooldown = time.Millisecond * 350
)

var (
	lastGetCurrentWindowResult []string
	lastGetCurrentWindowCall   = time.Now()
)

func getCurrentWindowProcessNames() ([]string, error) {
	// Timeout just like the one in windows util. Might not be necessary, but shouldn't hurt either.
	now := time.Now()
	if lastGetCurrentWindowCall.Add(getCurrentWindowInternalCooldown).After(now) {
		return lastGetCurrentWindowResult, nil
	}

	lastGetCurrentWindowCall = now

	// Use xprop to get currently active window's PID and use that to get it's process's binary from /proc.
	// Messy and probably could be replaced with some proper go api but works for now.
	cmd, err := exec.Command("sh", "-c",
		"basename $(readlink /proc/$(xprop -id $(xprop -root _NET_ACTIVE_WINDOW | cut -d ' ' -f 5) "+
			"_NET_WM_PID | awk '{print $NF}')/exe)").Output()

	if err != nil {
		return nil, fmt.Errorf("Error getting current process name %w", err)
	}

	result := []string{strings.TrimSpace(string(cmd))}

	lastGetCurrentWindowResult = result
	return result, nil
}
