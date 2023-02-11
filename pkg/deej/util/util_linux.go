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

	// Use xprop to get currently active window's PID and use that to output just the process name with ps.
	// It most likely won't work with sandboxed applications
	cmd, err := exec.Command("bash", "-c", "ps -q $(xprop -id $(xprop -root _NET_ACTIVE_WINDOW | cut -d ' ' -f 5) _NET_WM_PID | awk '{print $NF}') -o comm=").Output()

	if err != nil {
		return nil, fmt.Errorf("error getting current process name %w", err)
	}

	result := []string{strings.TrimSpace(string(cmd))}

	lastGetCurrentWindowResult = result
	return result, nil
}
