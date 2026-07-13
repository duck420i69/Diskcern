package ui

import (
	"bufio"
	"fmt"
	"os/exec"
)

func ExecuteCommand(command string, outputCh chan<- string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputCh <- scanner.Text()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputCh <- "[ERROR] " + scanner.Text()
		}
	}()

	go func() {
		err := cmd.Wait()
		if err != nil {
			outputCh <- fmt.Sprintf("[EXIT] Process finished with error: %v", err)
		} else {
			outputCh <- "[EXIT] Process finished successfully."
		}
		close(outputCh)
	}()

	return nil
}
