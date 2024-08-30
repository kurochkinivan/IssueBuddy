package issueBuddy

import "time"

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type CreateUpdateIssue struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	// Open/Close
	State       string `json:"state,omitempty"`
	StateReason string `json:"state_reasons,omitempty"`
	// off-topic/too heated/resolved/spam
	LockReason string `json:"lock_reason,omitempty"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func NewIssue(title, body string) *Issue {
	return &Issue{
		Title: title,
		Body:  body,
	}
}
