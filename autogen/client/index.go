// Code generated by goagen v1.4.3, DO NOT EDIT.
//
// API "feedpushr": index Resource Client
//
// Command:
// $ goagen
// --design=github.com/ncarlier/feedpushr/v3/design
// --out=/home/nicolas/workspace/fe/feedpushr/autogen
// --version=v1.4.3

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// GetIndexPath computes a request path to the get action of index.
func GetIndexPath() string {

	return fmt.Sprintf("/v2/")
}

// Get basic API information.
func (c *Client) GetIndex(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetIndexRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetIndexRequest create the request corresponding to the get action endpoint of the index resource.
func (c *Client) NewGetIndexRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
