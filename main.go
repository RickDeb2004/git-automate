package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func updateCommitPush() {
	if err := executeCommand("git", "add", "-A"); err != nil {
		log.Fatalf("Failed to execute 'git add' command: %v", err)
	}

	name := nameGenerator()
	if err := executeCommand("git", "commit", "-m", name); err != nil {
		log.Fatalf("Failed to execute 'git commit' command: %v", err)
	}

	if err := executeCommand("git", "push", "origin", "master"); err != nil {
		log.Fatalf("Failed to execute 'git push' command: %v", err)
	}

	fmt.Println("Successfully added, committed, and pushed all the changes")
}

func nameGenerator() string {

	return "initialCommit"
}

func main() {
	updateCommitPush()
}
