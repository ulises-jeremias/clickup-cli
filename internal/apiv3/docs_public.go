package apiv3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/triptechtravel/clickup-cli/internal/api"
)

// DocCore holds fields common to list and detail responses for workspace docs.
type DocCore struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Deleted    bool   `json:"deleted"`
	Archived   bool   `json:"archived"`
	Visibility string `json:"visibility"`
	Creator    struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"creator"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	Parent      struct {
		ID   string `json:"id"`
		Type int    `json:"type"`
	} `json:"parent"`
	Workspace struct {
		ID string `json:"id"`
	} `json:"workspace"`
}

// DocsListResponse is the JSON body for listing docs in a workspace.
type DocsListResponse struct {
	Docs       []DocCore `json:"docs"`
	NextCursor string    `json:"next_cursor"`
}

// PageRef is a summary entry in a page listing response.
type PageRef struct {
	ID         string    `json:"id"`
	DocID      string    `json:"doc_id"`
	Name       string    `json:"name"`
	SubTitle   string    `json:"sub_title"`
	OrderIndex int       `json:"order_index"`
	Pages      []PageRef `json:"pages"`
}

// PagesListResponse is the JSON body for listing pages under a doc.
type PagesListResponse struct {
	Pages []PageRef `json:"pages"`
}

// PageDetail is the full page returned by the get-page endpoint.
type PageDetail struct {
	ID            string `json:"id"`
	DocID         string `json:"doc_id"`
	Name          string `json:"name"`
	SubTitle      string `json:"sub_title"`
	Content       string `json:"content"`
	ContentFormat string `json:"content_format"`
	OrderIndex    int    `json:"order_index"`
	DateCreated   string `json:"date_created"`
	DateUpdated   string `json:"date_updated"`
	Creator       struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	} `json:"creator"`
	Pages []PageRef `json:"pages"`
}

// SearchDocsPublicParams holds optional query parameters for listing docs.
type SearchDocsPublicParams struct {
	Deleted    bool
	Archived   bool
	Creator    int
	ParentID   string
	ParentType int
	Limit      int
	Cursor     string
}

// CreateDocPublicRequest is the JSON body for creating a doc.
type CreateDocPublicRequest struct {
	Name       string     `json:"name"`
	CreatePage bool       `json:"create_page"`
	Parent     *DocParent `json:"parent,omitempty"`
	Visibility string     `json:"visibility,omitempty"`
}

// DocParent identifies where to create a doc in the hierarchy.
type DocParent struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
}

// CreatePagePublicRequest is the JSON body for creating a page.
type CreatePagePublicRequest struct {
	Name          string `json:"name"`
	ParentPageID  string `json:"parent_page_id,omitempty"`
	SubTitle      string `json:"sub_title,omitempty"`
	Content       string `json:"content,omitempty"`
	ContentFormat string `json:"content_format,omitempty"`
}

// EditPagePublicRequest is the JSON body for editing a page (partial update).
type EditPagePublicRequest struct {
	Name            string `json:"name,omitempty"`
	SubTitle        string `json:"sub_title,omitempty"`
	Content         string `json:"content,omitempty"`
	ContentFormat   string `json:"content_format,omitempty"`
	ContentEditMode string `json:"content_edit_mode,omitempty"`
}

func doJSON(ctx context.Context, client *api.Client, method, rawURL string, body any, result any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return api.APIErrorFromResponseBody(resp.StatusCode, respBody)
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func v3Root(client *api.Client) string {
	return strings.TrimSuffix(client.APIV3BaseURL(), "/")
}

func docsBaseURL(client *api.Client, workspaceID string) string {
	return v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs"
}

// SearchDocsPublic lists docs in a workspace (GET .../workspaces/{id}/docs).
func SearchDocsPublic(ctx context.Context, client *api.Client, workspaceID string, params SearchDocsPublicParams) (*DocsListResponse, error) {
	base, err := url.Parse(docsBaseURL(client, workspaceID))
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params.Deleted {
		q.Set("deleted", "true")
	}
	if params.Archived {
		q.Set("archived", "true")
	}
	if params.Creator != 0 {
		q.Set("creator", strconv.Itoa(params.Creator))
	}
	if params.ParentID != "" {
		q.Set("parent_id", params.ParentID)
		if params.ParentType != 0 {
			q.Set("parent_type", strconv.Itoa(params.ParentType))
		}
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	base.RawQuery = q.Encode()

	var out DocsListResponse
	if err := doJSON(ctx, client, http.MethodGet, base.String(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDocPublic fetches a single doc.
func GetDocPublic(ctx context.Context, client *api.Client, workspaceID, docID string) (*DocCore, error) {
	u := v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs/" + url.PathEscape(docID)
	var out DocCore
	if err := doJSON(ctx, client, http.MethodGet, u, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDocPublic creates a doc in the workspace.
func CreateDocPublic(ctx context.Context, client *api.Client, workspaceID string, req *CreateDocPublicRequest) (*DocCore, error) {
	u := docsBaseURL(client, workspaceID)
	var out DocCore
	if err := doJSON(ctx, client, http.MethodPost, u, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDocPagesPublic lists pages for a doc.
func GetDocPagesPublic(ctx context.Context, client *api.Client, workspaceID, docID string, maxPageDepth int) (*PagesListResponse, error) {
	u, err := url.Parse(v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs/" + url.PathEscape(docID) + "/pages")
	if err != nil {
		return nil, err
	}
	if maxPageDepth >= 0 {
		q := url.Values{}
		q.Set("max_page_depth", strconv.Itoa(maxPageDepth))
		u.RawQuery = q.Encode()
	}
	var out PagesListResponse
	if err := doJSON(ctx, client, http.MethodGet, u.String(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreatePagePublic creates a page under a doc.
func CreatePagePublic(ctx context.Context, client *api.Client, workspaceID, docID string, req *CreatePagePublicRequest) (*PageDetail, error) {
	u := v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs/" + url.PathEscape(docID) + "/pages"
	var out PageDetail
	if err := doJSON(ctx, client, http.MethodPost, u, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditPagePublic updates a page.
func EditPagePublic(ctx context.Context, client *api.Client, workspaceID, docID, pageID string, req *EditPagePublicRequest) (*PageDetail, error) {
	u := v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs/" + url.PathEscape(docID) + "/pages/" + url.PathEscape(pageID)
	var out PageDetail
	if err := doJSON(ctx, client, http.MethodPut, u, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPagePublic fetches a single page; contentFormat is optional (e.g. text/md).
func GetPagePublic(ctx context.Context, client *api.Client, workspaceID, docID, pageID, contentFormat string) (*PageDetail, error) {
	u, err := url.Parse(v3Root(client) + "/workspaces/" + url.PathEscape(workspaceID) + "/docs/" + url.PathEscape(docID) + "/pages/" + url.PathEscape(pageID))
	if err != nil {
		return nil, err
	}
	if contentFormat != "" {
		q := url.Values{}
		q.Set("content_format", contentFormat)
		u.RawQuery = q.Encode()
	}
	var out PageDetail
	if err := doJSON(ctx, client, http.MethodGet, u.String(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
