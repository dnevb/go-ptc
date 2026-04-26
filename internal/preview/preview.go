package preview

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Preview(dir string) error {
	if _, err := exec.LookPath("plymouthd"); err != nil {
		return fmt.Errorf("plymouthd not found: %w", err)
	}
	cmdVer := exec.Command("plymouthd", "--version")
	if out, err := cmdVer.CombinedOutput(); err != nil {
		return fmt.Errorf("plymouthd --version failed: %w: %s", err, string(out))
	}

	cmdDaemon := exec.Command("plymouthd", "--debug", "--daemonize")
	if out, err := cmdDaemon.CombinedOutput(); err != nil {
		return fmt.Errorf("plymouthd failed: %w: %s", err, string(out))
	}

	cmdShow := exec.Command("plymouth", "--show-splash")
	if out, err := cmdShow.CombinedOutput(); err != nil {
		_ = exec.Command("plymouth", "--quit").Run()
		return fmt.Errorf("plymouth --show-splash failed: %w: %s", err, string(out))
	}

	fmt.Fprintln(os.Stderr, "preview running for 5 seconds...")
	time.Sleep(5 * time.Second)

	cmdQuit := exec.Command("plymouth", "--quit")
	if out, err := cmdQuit.CombinedOutput(); err != nil {
		return fmt.Errorf("plymouth --quit failed: %w: %s", err, string(out))
	}
	return nil
}
