package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateDynamicCommitMessage() string {

	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute 'git diff --cached': %v", err)
	}

	lines := strings.Split(string(output), "\n")
	commitMessage := "chore: Update\n\n"

	for _, line := range lines {
		if strings.HasPrefix(line, "+") {

		} else if strings.HasPrefix(line, "-") {

			commitMessage += "fix: " + strings.TrimSpace(line[1:]) + "\n"
		}
	}

	if commitMessage != "chore: Update\n\n" {
		return commitMessage
	}

	return "chore: No significant changes\n\n"
}

func updateCommitPush(branchName string) {

	commitMessage := generateDynamicCommitMessage()

	if err := executeCommand("git", "checkout", "-b", branchName); err != nil {
		log.Fatalf("Failed to execute 'git checkout -b' command '%s': %v", branchName, err)
	}
	if err := executeCommand("git", "add", "-A"); err != nil {
		log.Fatalf("Failed to execute 'git add' command: %v", err)
	}

	if err := executeCommand("git", "commit", "-m", commitMessage); err != nil {
		log.Fatalf("Failed to execute 'git commit' command: %v", err)
	}

	if err := executeCommand("git", "push", "origin", branchName); err != nil {
		log.Fatalf("Failed to execute 'git push' command: %v", err)
	}

	fmt.Println("Successfully added, committed, and pushed all the changes with dynamic commit message")
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the branch name: ")
	branchName, _ := reader.ReadString('\n')
	branchName = strings.TrimSpace(branchName)

	if branchName == "" {
		log.Fatal("branch name cannot be empty")
	}
	updateCommitPush(branchName)
}
