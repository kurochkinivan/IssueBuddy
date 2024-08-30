package main

import (
	"fmt"
	"log"

	issueBuddy "github.com/kurochkinivan/IssueBuddy/internal"
)

const (
	token       = ""
	owner       = "liminfinity"
	repo        = "FarmFood"
	issueNumber = 2

	title       = "Futher dev"
	body        = "How are you planning to develop this project?"
	milestone   = 1
	state       = "close"
	stateReason = "i don't care"
	lockReason = "off-topic"
)

var (
	labels    = []string{"bug", "question"}
	assignees = []string{"liminfinity"}
)

func main() {
	// Get many
	issues, err := issueBuddy.GetIssues(token, owner, repo)
	if err != nil {
		log.Fatalf("failed to get all issues, err: %v\n", err)
	}
	fmt.Printf("%v\n", issues)

	// Get one
	issue, err := issueBuddy.GetIssue(token, owner, repo, issueNumber)
	if err != nil {
		log.Fatalf("failed to get the issue by number, err: %v\n", err)
	}
	fmt.Printf("%v\n", issue)

	// Create
	createIssue := issueBuddy.CreateUpdateIssue{
		Title:     title,
		Body:      body,
		Milestone: milestone,
		Labels:    labels,
		Assignees: assignees,
	}
	newIssue, err := issueBuddy.CreateIssue(token, owner, repo, createIssue)
	if err != nil {
		log.Fatalf("failed to create the issue, err: %v\n", err)
	}
	fmt.Printf("%v\n", newIssue)

	// Update
	updateIssue := issueBuddy.CreateUpdateIssue{
		Title:       title,
		Body:        body,
		Milestone:   milestone,
		Labels:      labels,
		Assignees:   assignees,
		State:       state,
		StateReason: stateReason,
	}
	newIssue, err = issueBuddy.UpdateIssue(token, owner, repo, issueNumber, updateIssue)
	if err != nil {
		log.Fatalf("failed to update the issue, err: %v\n", err)
	}
	fmt.Printf("%v\n", newIssue)

	// Lock
	lockIssue := issueBuddy.CreateUpdateIssue{
		LockReason: lockReason,
	}
	err = issueBuddy.LockIssue(token, owner, repo, issueNumber, lockIssue)
	if err != nil {
		log.Fatalf("failed to lock the issue, err: %v\n", err)
	}

	// Unlock
	err = issueBuddy.UnlockIssue(token, owner, repo, issueNumber)
	if err != nil {
		log.Fatalf("failed to unlock the issue, err: %v\n", err)
	}

}
