// Copyright 2026 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// IssueAssigneesOption options for adding/removing issue assignees
type IssueAssigneesOption struct {
	Assignees []string `json:"assignees"`
}

// AddIssueAssignees add assignees to an issue
func (c *IssuesService) AddIssueAssignees(ctx context.Context, owner, repo string, index int64, opt IssueAssigneesOption) (*Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	issue := new(Issue)
	resp, err := c.getParsedResponse(ctx, "POST", fmt.Sprintf("/repos/%s/%s/issues/%d/assignees", owner, repo, index),
		jsonHeader, bytes.NewReader(body), issue)
	return issue, resp, err
}

// RemoveIssueAssignees remove assignees from an issue
func (c *IssuesService) DeleteIssueAssignees(ctx context.Context, owner, repo string, index int64, opt IssueAssigneesOption) (*Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	issue := new(Issue)
	resp, err := c.getParsedResponse(ctx, "DELETE", fmt.Sprintf("/repos/%s/%s/issues/%d/assignees", owner, repo, index),
		jsonHeader, bytes.NewReader(body), issue)
	return issue, resp, err
}

// CheckIssueAssignee check if a user can be assigned to an issue
func (c *IssuesService) CheckIssueAssignee(ctx context.Context, owner, repo, assignee string, index int64) (bool, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return false, nil, err
	}
	data, resp, err := c.getResponse(ctx, "GET", fmt.Sprintf("/repos/%s/%s/issues/%d/assignees/%s", owner, repo, index, assignee), nil, nil)

	// If the server returns 404, user is NOT assignable
	if resp != nil && resp.StatusCode == 404 {
		return false, resp, nil
	}

	// If any other error occurred, return it
	if err != nil {
		return false, resp, err
	}

	// 204 No Content → user is assignable
	if resp.StatusCode == 204 {
		return true, resp, nil
	}

	// Unexpected status code
	return false, resp, fmt.Errorf("unexpected status code %d, body: %s", resp.StatusCode, string(data))
}
