package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var logger *log.Logger

func init() {
	logFile, err := os.OpenFile("git_workflow.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}
func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func logAction(action string) {
	logger.Printf("Action: %s\n", action)
}

func generateDynamicCommitMessage(issueNumber string) string {

	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute 'git diff --cached': %v", err)
	}

	lines := strings.Split(string(output), "\n")
	commitMessage := "chore: Update\n\n"

	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			fmt.Printf("Do you want to include this addition in the commit? (y/n) \n")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if strings.ToLower(answer) == "y" {
				commitMessage += "feat: " + strings.TrimSpace(line[1:]) + "\n"
			}
		} else if strings.HasPrefix(line, "-") {
			fmt.Printf("Do you want to include this deletion in the commit? (y/n) \n")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if strings.ToLower(answer) == "y" {
				commitMessage += "fix: " + strings.TrimSpace(line[1:]) + "\n"
			}
		}
	}

	if commitMessage != "chore: Update\n\n" {
		commitMessage += fmt.Sprintf("Issue: %s\n", issueNumber)
		return commitMessage
	}

	return "chore: No significant changes\n\n"
}
func promptForIssueNumber() string {
	fmt.Printf("Enter the issue number : ")
	reader := bufio.NewReader(os.Stdin)
	issueNumber, _ := reader.ReadString('\n')
	return strings.TrimSpace(issueNumber)
}
func updateCommitPush(branchName string, issueNumber string) {

	commitMessage := generateDynamicCommitMessage(issueNumber)
	logAction("Switching to branch: " + branchName)
	if err := executeCommand("git", "checkout", "-b", branchName); err != nil {
		log.Fatalf("Failed to execute 'git checkout -b' command '%s': %v", branchName, err)
	}
	logAction("Adding changes to the commit")
	if err := executeCommand("git", "add", "-A"); err != nil {
		log.Fatalf("Failed to execute 'git add' command: %v", err)
	}
	logAction("Committing changes with message: " + commitMessage)
	if err := executeCommand("git", "commit", "-m", commitMessage); err != nil {
		log.Fatalf("Failed to execute 'git commit' command: %v", err)
	}
	logAction("Pushing changes to branch: " + branchName)
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
	issueNumber := promptForIssueNumber()
	updateCommitPush(branchName, issueNumber)
}
