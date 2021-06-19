package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CreateTag creates an annotated tag
func CreateTag(tag, msg string) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
	fmt.Printf("\nCreating annotated tag %v with message '%v'...\n", tag, msg)

	_, err = exec.Command("git", "tag", "-a", tag, "-m", msg).CombinedOutput()
	return
}

// CommitFiles adds the given files to staging and commits them with the given message
func CommitFiles(sList []string, m string) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
	fmt.Printf("\nCommiting files...\n")

	fmt.Printf("\tStaging files...\n")
	for _, s := range sList {
		if _, err = exec.Command("git", "add", s).CombinedOutput(); err != nil {
			_ = RemoveFromStaging(sList, true)
			return
		}
	}
	fmt.Printf("\tCommiting with '%v' as message...\n", m)
	if _, err = exec.Command("git", "commit", "-m", m).CombinedOutput(); err != nil {
		_ = RemoveFromStaging(sList, true)
		return
	}

	return
}

// GetLatestTag Returns the latest tag for the git repo related to the working directory
func GetLatestTag(noFetch bool) (tag string, err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
	fmt.Printf("\nGetting latest git tag...\n")

	if !noFetch {
		fmt.Printf("\tFetching...\n")
		if _, err = exec.Command("git", "fetch", "--tags").CombinedOutput(); err != nil {
			fmt.Println("...fetching failed!")
		}
	}

	cmd1 := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	output1 := &bytes.Buffer{}
	cmd1.Stdout = output1
	if err = cmd1.Run(); err != nil {
		return
	}
	hash := string(output1.Bytes())
	hash = strings.TrimSuffix(hash, "\n")

	cmd2 := exec.Command("git", "describe", "--tags", hash)
	output2 := &bytes.Buffer{}
	cmd2.Stdout = output2
	if err = cmd2.Run(); err != nil {
		return
	}

	tag = string(output2.Bytes())
	tag = strings.TrimSuffix(tag, "\n")

	return
}

// HasGit checks if the git command exists in PATH
func HasGit() bool {
	s, err := exec.LookPath("git")
	return s != "" && err == nil
}

// RemoveFromStaging removes the given files from the stagin area
func RemoveFromStaging(sList []string, ignoreErrors bool) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}

	for _, s := range sList {
		if _, err = exec.Command("git", "reset", s).CombinedOutput(); err != nil {
			if !ignoreErrors {
				return
			}
		}
	}
	return
}
