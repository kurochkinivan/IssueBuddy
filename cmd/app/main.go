package main

import (
	"flag"
	"fmt"
	"strings"

	cuserr "github.com/kurochkinivan/IssueBuddy/internal/custome_erros"
	"github.com/kurochkinivan/IssueBuddy/internal/editors"
	issueBuddy "github.com/kurochkinivan/IssueBuddy/internal/issue_client"
	"github.com/kurochkinivan/IssueBuddy/internal/parsers"
	"github.com/kurochkinivan/IssueBuddy/internal/readers"
)

var (
	token string
	owner string
	repo  string

	filename = "input.txt"
)

func main() {
	flag.StringVar(&token, "token", "", "Github access token")
	flag.StringVar(&owner, "owner", "", "Repository owner")
	flag.StringVar(&repo, "repo", "", "Repository name")
	flag.Parse()

	if token == "" {
		fmt.Println("Enter your github access token:")
		fmt.Scan(&token)
	}
	if owner == "" {
		fmt.Println("Enter the repo owner:")
		fmt.Scan(&owner)
	}
	if repo == "" {
		fmt.Println("Enter the repo name:")
		fmt.Scan(&repo)
	}

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
				fmt.Printf("failed to get all issues, err: %v\n", err)
				refreshLoop()
			}

			for _, issue := range issues {
				fmt.Printf("№%3d | Title: %40.40s... | state: %5q | user: %20s | user_url: %39s |\n",
					issue.Number, issue.Title, issue.State, issue.User.Login, issue.User.HTMLURL)
			}
		case "5":
			if empty := isEmpty(); empty {
				fmt.Println(cuserr.CantMakeRequest + cuserr.EmptyError)
				refreshLoop()
			}

			var issueNumber int
			fmt.Println("Enter issue number:")
			fmt.Scan(&issueNumber)

			issue, err := issueBuddy.GetIssue(token, owner, repo, issueNumber)
			if err != nil {
				fmt.Printf("failed to get the issue by number, err: %v\n", err)
				refreshLoop()
			}

			printIssue(issue)
		case "6":
			if empty := isEmpty(); empty {
				fmt.Println(cuserr.CantMakeRequest + cuserr.EmptyError)
				refreshLoop()
			}

			title, body, milestone, labels, assignees := createUpdateDialog()

			createIssue := issueBuddy.CreateUpdateIssue{
				Title:     title,
				Body:      body,
				Milestone: milestone,
				Labels:    labels,
				Assignees: assignees,
			}

			fmt.Println("----------------------------------------------------")
			fmt.Printf("\nTitle: %s\nBody: %s\nMilestone: %d\nLabels: %v\nAssignees: %v\n",
				createIssue.Title, createIssue.Body, createIssue.Milestone, createIssue.Labels, createIssue.Assignees)
			fmt.Println("----------------------------------------------------")

			fmt.Println("Do you confirm creating issue?(y/n)")

			if confirmed := confirmDialog(); confirmed {
				_, err := issueBuddy.CreateIssue(token, owner, repo, createIssue)
				if err != nil {
					fmt.Printf("failed to create the issue, err: %v\n", err)
					refreshLoop()
				}
			} else {
				mainMenu()
			}

			fmt.Println("New Issue was successfully created!")
		case "7":
			if empty := isEmpty(); empty {
				fmt.Println(cuserr.CantMakeRequest + cuserr.EmptyError)
				refreshLoop()
			}

			var issueNumber int
			fmt.Println("Enter issue number:")
			fmt.Scan(&issueNumber)

			oldIssue, err := issueBuddy.GetIssue(token, owner, repo, issueNumber)
			if err != nil {
				fmt.Printf("failed to get the issue by number, err: %v\n", err)
				refreshLoop()
			}
			printIssue(oldIssue)

			title, body, milestone, labels, assignees := createUpdateDialog()

			fmt.Println("Enter the state(close/open)", cuserr.SkipFieldInstruction)
			var state string
			for {
				fmt.Scanln(&state)
				if state == "close" || state == "open" || state == "" {
					break
				} else {
					fmt.Println(cuserr.InvalidInput)
				}
			}

			fmt.Println("Enter the state reason", cuserr.SkipFieldInstruction)
			stateReason, err := readers.ReadOneLine()
			if err != nil {
				fmt.Println(err)
			}

			updateIssue := issueBuddy.CreateUpdateIssue{
				Title:       title,
				Body:        body,
				Milestone:   milestone,
				Labels:      labels,
				Assignees:   assignees,
				State:       state,
				StateReason: stateReason,
			}

			fmt.Println("----------------------------------------------------")
			fmt.Println("Old Issue:")
			fmt.Printf("\nTitle: %s\nBody: %s\nMilestone: %d\nLabels: %v\nAssignees: %v\nState: %s\n",
				oldIssue.Title, oldIssue.Body, oldIssue.Milestone, oldIssue.Labels, oldIssue.Assignees, oldIssue.State)
			fmt.Println("----------------------------------------------------")
			fmt.Println("----------------------------------------------------")
			fmt.Printf("\nTitle: %s\nBody: %s\nMilestone: %d\nLabels: %v\nAssignees: %v\nState: %s\nState reason: %s\n",
				updateIssue.Title, updateIssue.Body, updateIssue.Milestone, updateIssue.Labels, updateIssue.Assignees, updateIssue.State, updateIssue.StateReason)
			fmt.Println("----------------------------------------------------")

			fmt.Println("Do you confirm updating issue?(y/n)")

			if confirmed := confirmDialog(); confirmed {
				_, err := issueBuddy.UpdateIssue(token, owner, repo, issueNumber, updateIssue)
				if err != nil {
					fmt.Printf("failed to update the issue, err: %v\n", err)
					refreshLoop()
				}
			} else {
				mainMenu()
			}

			fmt.Println("Issue was successfully updated!")
		case "8":
			if empty := isEmpty(); empty {
				fmt.Println(cuserr.CantMakeRequest + cuserr.EmptyError)
				refreshLoop()
			}

			fmt.Println("Enter issue number:")
			var issueNumber int
			fmt.Scan(&issueNumber)

			fmt.Println("Enter lock reason(off-topic/too heated/resolved/spam):")
			var lockReason string
			for {
				lockReason, _ = readers.ReadOneLine()
				if lockReason == "off-topic" || lockReason == "too heated" || lockReason == "resolved" || lockReason == "spam" {
					break
				} else {
					fmt.Println(cuserr.InvalidInput)
				}
			}

			lockIssue := issueBuddy.CreateUpdateIssue{
				LockReason: lockReason,
			}
			err := issueBuddy.LockIssue(token, owner, repo, issueNumber, lockIssue)
			if err != nil {
				fmt.Printf("failed to lock the issue, err: %v\n", err)
				refreshLoop()
			}

			fmt.Println("Issue was successfully locked!")
		case "9":
			if empty := isEmpty(); empty {
				fmt.Println(cuserr.CantMakeRequest + cuserr.EmptyError)
				refreshLoop()
			}

			fmt.Println("Enter issue number:")
			var issueNumber int
			fmt.Scan(&issueNumber)

			err := issueBuddy.UnlockIssue(token, owner, repo, issueNumber)
			if err != nil {
				fmt.Printf("failed to unlock the issue, err: %v\n", err)
			}

			fmt.Println("Issue was successfully unlocked!")
		case "0":
			mainMenu()
		default:
			fmt.Println(cuserr.InvalidInput)
			continue
		}
		refreshLoop()
	}
}

func printIssue(issue issueBuddy.Issue) {
	fmt.Println("----------------------------------------------------")
	fmt.Printf("№%d\nCreated at: %v\nState: %q\nUser: %s\nUser_url: %s\nTitle: %s\nBody: %s\n",
		issue.Number, issue.CreatedAt.Format("2006-01-02 15:04:05"), issue.State, issue.User.Login, issue.User.HTMLURL, issue.Title, issue.Body)
	fmt.Println("----------------------------------------------------")
}

func createUpdateDialog() (title string, body string, milestone int, labels []string, assignees []string) {
	fmt.Println("Enter the title (one line)", cuserr.SkipFieldInstruction)
	title, err := readers.ReadOneLine()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Enter the body (using texteditor(type ed)/using command line(type cl))", cuserr.SkipFieldInstruction)

	for {
		var choice string
		fmt.Scanln(&choice)
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
				fmt.Println(err)
				continue
			}

			break
		} else if choice == "cl" {
			line, err := readers.ReadOneLine()
			if err != nil {
				fmt.Println(err)
				fmt.Println("Enter the body (using texteditor(type ed)/using command line(type cl))")
				continue
			}
			body = strings.TrimSpace(line)
			break
		} else if choice == "" {
			body = ""
			continue
		} else {
			fmt.Println("invalid input, type 'cl' or 'ed'")
		}
	}

	fmt.Println("Enter the milestone", cuserr.SkipFieldInstruction)
	fmt.Scanln(&milestone)

	fmt.Println("Enter the labels separated by comma ','", cuserr.SkipFieldInstruction)
	line, err := readers.ReadOneLine()
	if err != nil {
		fmt.Println(err)
	}
	labels = parsers.ParseLine(line)

	fmt.Println("Enter the assignees separated by comma ','", cuserr.SkipFieldInstruction)
	line, err = readers.ReadOneLine()
	if err != nil {
		fmt.Println(err)
	}
	assignees = parsers.ParseLine(line)

	return title, body, milestone, labels, assignees
}

func confirmDialog() bool {
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
