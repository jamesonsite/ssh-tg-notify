package follow

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os/exec"
)

// Journal streams journalctl stdout lines until ctx is cancelled.
func Journal(ctx context.Context, binary string, args []string, lines chan<- string) error {
	if binary == "" {
		binary = "journalctl"
	}
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Stderr = io.Discard
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		_ = cmd.Process.Kill()
	}()
	sc := bufio.NewScanner(stdout)
	// Some syslog lines can be long.
	const max = 256 * 1024
	sc.Buffer(make([]byte, 0, 64*1024), max)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			_ = cmd.Wait()
			return ctx.Err()
		case lines <- sc.Text():
		}
	}
	if err := sc.Err(); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("journalctl read: %v", err)
	}
	waitErr := cmd.Wait()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if waitErr != nil && !errors.Is(waitErr, context.Canceled) {
		return waitErr
	}
	return nil
}
