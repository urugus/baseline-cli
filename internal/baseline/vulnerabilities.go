package baseline

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type ListVulnerabilitiesOptions struct {
	Page    int
	PerPage int
	AssetID string
}

func (c *Client) ListVulnerabilities(ctx context.Context, opts ListVulnerabilitiesOptions) (PageResponse[Vulnerability], []byte, error) {
	query := url.Values{}
	if opts.Page > 0 {
		query.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		query.Set("perPage", strconv.Itoa(opts.PerPage))
	}
	if opts.AssetID != "" {
		query.Set("asset", opts.AssetID)
	}

	body, err := c.GetRaw(ctx, "/api/v1/vulnerabilities", query)
	if err != nil {
		return PageResponse[Vulnerability]{}, nil, err
	}

	var response PageResponse[Vulnerability]
	if err := json.Unmarshal(body, &response); err != nil {
		return PageResponse[Vulnerability]{}, nil, err
	}
	return response, body, nil
}

func (c *Client) GetVulnerability(ctx context.Context, id string) (EntityResponse[Vulnerability], []byte, error) {
	body, err := c.GetRaw(ctx, "/api/v1/vulnerabilities/"+url.PathEscape(id), nil)
	if err != nil {
		return EntityResponse[Vulnerability]{}, nil, err
	}

	var response EntityResponse[Vulnerability]
	if err := json.Unmarshal(body, &response); err != nil {
		return EntityResponse[Vulnerability]{}, nil, err
	}
	return response, body, nil
}
