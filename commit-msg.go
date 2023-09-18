package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	// The first argument is the path to the temporary file that contains the commit message
	if len(os.Args) < 2 {
		log.Fatal("No commit message file provided")
	}
	commitMessageFile := os.Args[1]

	// Get the current branch name
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	branchName := strings.TrimSpace(out.String())

	// Check if the branch name matches the pattern
	pattern := `^(feature|bugfix)/([a-zA-Z]+-\d+)-.*$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(branchName)
	if matches == nil {
		log.Println("Branch name does not match the pattern")
		return
	}

	// Get the JIRA issue key
	jiraIssueKey := matches[2]

	// Read the commit message
	commitMessage, err := ioutil.ReadFile(commitMessageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Prepend the JIRA issue key to the commit message
	newCommitMessage := jiraIssueKey + ": " + string(commitMessage)

	// Write the new commit message back to the file
	err = ioutil.WriteFile(commitMessageFile, []byte(newCommitMessage), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
