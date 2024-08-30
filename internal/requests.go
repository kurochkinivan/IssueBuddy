package issueBuddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var client = &http.Client{}

func checkResponseStatus(resp *http.Response, validStatusCode int) error {
	if resp.StatusCode != validStatusCode {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read data from response body, err: %v", err)
		}

		var buf bytes.Buffer
		json.Indent(&buf, data, "", " ")
		return fmt.Errorf("something went wrong, json response: %v", buf.String())
	}

	return nil
}

func GetIssues(token, owner, repo string) ([]Issue, error) {
	req := &http.Request{
		Method: http.MethodGet,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues", owner, repo),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Issue{}, fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusOK); err != nil {
		return []Issue{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Issue{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var issues []Issue
	json.Unmarshal(body, &issues)

	return issues, nil
}

func GetIssue(token, owner, repo string, issueNumber int) (Issue, error) {
	req := &http.Request{
		Method: http.MethodGet,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, issueNumber),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusOK); err != nil {
		return Issue{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var issue Issue
	json.Unmarshal(body, &issue)

	return issue, nil
}

func CreateIssue(token, owner, repo string, issue CreateUpdateIssue) (CreateUpdateIssue, error) {
	data, err := json.Marshal(issue)
	if err != nil {
		return CreateUpdateIssue{}, fmt.Errorf("failed to marshal data, err: %v", err)
	}

	req := &http.Request{
		Method: http.MethodPost,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues", owner, repo),
		},
		Body: io.NopCloser(bytes.NewBuffer(data)),
	}

	resp, err := client.Do(req)
	if err != nil {
		return CreateUpdateIssue{}, fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusCreated); err != nil {
		return CreateUpdateIssue{}, err
	}

	return issue, nil
}

func UpdateIssue(token, owner, repo string, issueNumber int, issue CreateUpdateIssue) (CreateUpdateIssue, error) {
	data, err := json.Marshal(issue)
	if err != nil {
		return CreateUpdateIssue{}, fmt.Errorf("failed to marshal data, err: %v", err)
	}

	req := &http.Request{
		Method: http.MethodPatch,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, issueNumber),
		},
		Body: io.NopCloser(bytes.NewBuffer(data)),
	}

	resp, err := client.Do(req)
	if err != nil {
		return CreateUpdateIssue{}, fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusOK); err != nil {
		return CreateUpdateIssue{}, err
	}

	return issue, nil
}

func LockIssue(token, owner, repo string, issueNumber int, issue CreateUpdateIssue) error {
	data, err := json.Marshal(issue)
	if err != nil {
		return fmt.Errorf("failed to marshal data, err: %v", err)
	}

	req := &http.Request{
		Method: http.MethodPut,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues/%d/lock", owner, repo, issueNumber),
		},
		Body: io.NopCloser(bytes.NewBuffer(data)),
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

func UnlockIssue(token, owner, repo string, issueNumber int) error {
	req := &http.Request{
		Method: http.MethodDelete,
		Header: map[string][]string{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {fmt.Sprintf("Bearer %s", token)},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
			Path:   fmt.Sprintf("/repos/%s/%s/issues/%d/lock", owner, repo, issueNumber),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a request, err: %v", err)
	}
	defer resp.Body.Close()

	if err = checkResponseStatus(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}
