package main

import (
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
	// Use "git diff --cached" to get the staged changes
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute 'git diff --cached': %v", err)
	}

	// Process the output to create a commit message
	lines := strings.Split(string(output), "\n")
	commitMessage := "chore: Update\n\n" // Default message if no changes

	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			// This line represents an addition, so include it in the commit message
			commitMessage += "feat: " + strings.TrimSpace(line[1:]) + "\n"
		} else if strings.HasPrefix(line, "-") {
			// This line represents a deletion, so include it in the commit message
			commitMessage += "fix: " + strings.TrimSpace(line[1:]) + "\n"
		}
	}

	// Remove the default line if changes were detected
	if commitMessage != "chore: Update\n\n" {
		return commitMessage
	}

	// No changes detected
	return "chore: No significant changes\n\n"
}

// func updateCommitPush() {
// 	if err := executeCommand("git", "add", "-A"); err != nil {
// 		log.Fatalf("Failed to execute 'git add' command: %v", err)
// 	}

// 	name := nameGenerator()
// 	if err := executeCommand("git", "commit", "-m", name); err != nil {
// 		log.Fatalf("Failed to execute 'git commit' command: %v", err)
// 	}

// 	if err := executeCommand("git", "push", "origin", "master"); err != nil {
// 		log.Fatalf("Failed to execute 'git push' command: %v", err)
// 	}

// 	fmt.Println("Successfully added, committed, and pushed all the changes")
// }

// func nameGenerator() string {

// 	return "initialCommit"
// }

//	func main() {
//		updateCommitPush()
//	}
func updateCommitPush() {
	// Get the dynamic commit message
	commitMessage := generateDynamicCommitMessage()

	// Perform the Git add, commit, and push with the dynamic commit message
	if err := executeCommand("git", "add", "-A"); err != nil {
		log.Fatalf("Failed to execute 'git add' command: %v", err)
	}

	if err := executeCommand("git", "commit", "-m", commitMessage); err != nil {
		log.Fatalf("Failed to execute 'git commit' command: %v", err)
	}

	if err := executeCommand("git", "push", "origin", "master"); err != nil {
		log.Fatalf("Failed to execute 'git push' command: %v", err)
	}

	fmt.Println("Successfully added, committed, and pushed all the changes with dynamic commit message")
}
func main() {
	updateCommitPush()
}
