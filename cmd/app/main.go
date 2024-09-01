package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/kurochkinivan/IssueBuddy/internal/editors"
	issueBuddy "github.com/kurochkinivan/IssueBuddy/internal/issue_client"
	"github.com/kurochkinivan/IssueBuddy/internal/parsers"
	"github.com/kurochkinivan/IssueBuddy/internal/readers"
)

var (
	token string
	owner string
	repo  string

	state       = "close"
	stateReason = "i don't care"
	lockReason  = "off-topic"

	filename = "input.txt"
)

// TODO: think about log.Fatal

func main() {
	clearAll()

	mainMenu()
}

func mainMenu() {
	clearAll()
	fmt.Printf("Your data:\n[1] Your Token*: %s;\n[2] Repository owner*: %s;\n[3] Repository name*: %s;\n", token, owner, repo)
	fmt.Println("Actions:\n[4] Get all issues;\n[5] Get an issue;\n[6] Create new issue;\n[7] Edit an issue (includes closing);")
	fmt.Println("[8] Lock issue;\n[9] Unlock issue;\n[0] Refresh the menu;")

	for {
		var desicion string
		fmt.Scan(&desicion)

		switch desicion {
		case "1":
			fmt.Println("Enter your github token:")
			fmt.Scan(&token)
			mainMenu()
		case "2":
			fmt.Println("Enter the repo owner:")
			fmt.Scan(&owner)
			mainMenu()
		case "3":
			fmt.Println("Enter the repo name:")
			fmt.Scan(&repo)
			mainMenu()
		case "4":
			isEmpty()
			issues, err := issueBuddy.GetIssues(token, owner, repo)
			if err != nil {
				log.Fatalf("failed to get all issues, err: %v\n", err)
			}

			for _, issue := range issues {
				fmt.Printf("№%3d | Title: %40.40s... | state: %5q | user: %20s | user_url: %39s |\n",
					issue.Number, issue.Title, issue.State, issue.User.Login, issue.User.HTMLURL)
			}
		case "5":
			if empty := isEmpty(); empty {
				fmt.Println("cannot make a request: fields token, repo owner and repo name should not be empty;")
				refreshLoop()
			}

			var issueNumber int
			fmt.Println("Enter issue number:")
			fmt.Scan(&issueNumber)

			issue, err := issueBuddy.GetIssue(token, owner, repo, issueNumber)
			if err != nil {
				log.Fatalf("failed to get the issue by number, err: %v\n", err)
			}
			fmt.Printf("%v\n", issue)
		case "6":
			if empty := isEmpty(); empty {
				fmt.Println("cannot make a request: fields token, repo owner and repo name should not be empty;")
				refreshLoop()
			}

			fmt.Println("Enter the title (one line):")
			title := readers.ReadOneLine()

			fmt.Println("Enter the body (you can use texteditor (type ed) or continue using command line (type cl)'):")

			var body string
			for {
				var choice string
				fmt.Scan(&choice)
				if choice == "ed" {

					fmt.Println("choose editor: (Press enter if you want default)")
					var editor string
					fmt.Scanln(&editor)

					if editor == "" {
						editor = editors.DefaultEditor()
					}

					cmd := editors.OpenEditor(editor, filename)
					err := cmd.Start()
					if err != nil {
						fmt.Printf("failed to open editor, err: %v", err)
						continue
					}

					err = cmd.Wait()
					if err != nil {
						fmt.Printf("failed to wait for editor, err: %v", err)
						continue
					}

					body, err = readers.ReadFile(filename)
					if err != nil {
						fmt.Printf("failed to read file, err: %v", err)
						continue
					}

					break
				} else if choice == "cl" {
					body = readers.ReadOneLine()
					break
				} else {
					fmt.Println("invalid input, type 'cl' or 'ed'")
				}
			}

			fmt.Println("Enter the milestone (Press enter if you don't need it):")
			var milestone int
			fmt.Scanln(&milestone)

			fmt.Println("Enter the labels separated labels by comma ',' (Press enter if you don't need it):")
			labels := parsers.ParseLine(readers.ReadOneLine())

			fmt.Println("Enter the assignees separated by comma ',' (Press enter if you don't need it):")
			assignees := parsers.ParseLine(readers.ReadOneLine())

			createIssue := issueBuddy.CreateUpdateIssue{
				Title:     title,
				Body:      body,
				Milestone: milestone,
				Labels:    labels,
				Assignees: assignees,
			}

			fmt.Printf("Title: %s\nMilestone: %d\nLabels: %v\nAssignees: %v\nBody: %s\n",
				createIssue.Title, createIssue.Milestone, createIssue.Labels, createIssue.Assignees, createIssue.Body)
			fmt.Println("Do you confirm creating issue?(y/n)")

			if confirmed := confirmDialog(); confirmed {
				_, err := issueBuddy.CreateIssue(token, owner, repo, createIssue)
				if err != nil {
					log.Fatalf("failed to create the issue, err: %v\n", err)
				}
				break
			} else {
				mainMenu()
			}

			fmt.Println("New Issue was successfully created! (press 0 to refresh menu)")
		case "0":
			mainMenu()
		}
	}

}

func confirmDialog() bool {
	fmt.Println("Do you confirm creating issue?(y/n)")
	for {
		var yesNo string
		fmt.Scan(&yesNo)

		if strings.ToLower(yesNo) == "y" || strings.ToLower(yesNo) == "yes" {
			return true
		} else if strings.ToLower(yesNo) == "n" || strings.ToLower(yesNo) == "no" {
			return false
		} else {
			fmt.Println("invalid input, write y(yes) or n(no)")
		}
	}
}

func isEmpty() bool {
	if token == "" || repo == "" || owner == "" {
		return true
	}
	return false
}

func refreshLoop() {
	for {
		fmt.Println("In order to refresh terminal, enter '0'")
		var refresh string
		fmt.Scan(&refresh)
		if refresh == "0" {
			mainMenu()
		}
	}
}

func clearAll() {
	fmt.Print("\033[H\033[2J")
}

func test() {
	// // Get many
	// issues, err := issueBuddy.GetIssues(token, owner, repo)
	// if err != nil {
	// 	log.Fatalf("failed to get all issues, err: %v\n", err)
	// }
	// fmt.Printf("%v\n", issues)

	// // Get one
	// issue, err := issueBuddy.GetIssue(token, owner, repo, issueNumber)
	// if err != nil {
	// 	log.Fatalf("failed to get the issue by number, err: %v\n", err)
	// }
	// fmt.Printf("%v\n", issue)

	// // Create
	// createIssue := issueBuddy.CreateUpdateIssue{
	// 	Title:     title,
	// 	Body:      body,
	// 	Milestone: milestone,
	// 	Labels:    labels,
	// 	Assignees: assignees,
	// }
	// newIssue, err := issueBuddy.CreateIssue(token, owner, repo, createIssue)
	// if err != nil {
	// 	log.Fatalf("failed to create the issue, err: %v\n", err)
	// }
	// fmt.Printf("%v\n", newIssue)

	// // Update
	// updateIssue := issueBuddy.CreateUpdateIssue{
	// 	Title:       title,
	// 	Body:        body,
	// 	Milestone:   milestone,
	// 	Labels:      labels,
	// 	Assignees:   assignees,
	// 	State:       state,
	// 	StateReason: stateReason,
	// }
	// newIssue, err = issueBuddy.UpdateIssue(token, owner, repo, issueNumber, updateIssue)
	// if err != nil {
	// 	log.Fatalf("failed to update the issue, err: %v\n", err)
	// }
	// fmt.Printf("%v\n", newIssue)

	// // Lock
	// lockIssue := issueBuddy.CreateUpdateIssue{
	// 	LockReason: lockReason,
	// }
	// err = issueBuddy.LockIssue(token, owner, repo, issueNumber, lockIssue)
	// if err != nil {
	// 	log.Fatalf("failed to lock the issue, err: %v\n", err)
	// }

	// // Unlock
	// err = issueBuddy.UnlockIssue(token, owner, repo, issueNumber)
	// if err != nil {
	// 	log.Fatalf("failed to unlock the issue, err: %v\n", err)
	// }
	fmt.Println("")
}
