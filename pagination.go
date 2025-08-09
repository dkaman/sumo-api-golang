package sumo-api-golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

var PagerDone = errors.New("no more pages to iterate")

type PageInfo struct {
	Limit   int
	Skip    int
	HasNext bool
	HasPrev bool
}

type pageIterator[T any] interface {
	NextPage() ([]T, error)
}

// Pager is a generic, offset/limit-based pager over SUMO endpoints that return apiBulkResponse.
type Pager[T any] struct {
	ctx    context.Context
	client *Client
	method string
	path   string
	query  url.Values

	limit     int
	skip      int
	inited    bool
	exhausted bool
	lastLen   int
	lastErr   error
}

func newPager[T any](ctx context.Context, client *Client, method, path string, q url.Values, size int) *Pager[T] {
	if size <= 0 {
		size = 50
	}

	// Make a safe copy of the query and set limit/skip
	q2 := url.Values{}
	for k, v := range q {
		q2[k] = append([]string(nil), v...)
	}

	q2.Set("limit", fmt.Sprint(size))

	return &Pager[T]{
		ctx:    ctx,
		client: client,
		method: method,
		path:   path,
		query:  q2,
		limit:  size,
		skip:   0,
	}
}

func (p *Pager[T]) PageInfo() PageInfo {
	return PageInfo{
		Limit:   p.limit,
		Skip:    p.skip,
		HasNext: !p.exhausted || !p.inited,
		HasPrev: p.skip > 0,
	}
}

// NextPage fetches the next page (or the first, if not started).
func (p *Pager[T]) NextPage() ([]T, error) {
	// If we've already failed, keep returning the same error (typical iterator behavior)
	if p.lastErr != nil && !errors.Is(p.lastErr, PagerDone) {
		return nil, p.lastErr
	}

	// If initialized and we know we've reached/exceeded the end, signal done
	if p.inited && p.exhausted {
		p.lastErr = PagerDone
		return nil, p.lastErr
	}

	items, err := p.fetchCurrent()
	if err != nil {
		p.lastErr = err
		return nil, err
	}


	p.lastLen = len(items)
	p.inited = true
	p.exhausted = p.lastLen < p.limit

	// advance cursor by the configured window (common offset/limit pattern)
	p.skip += p.limit
	p.query.Set("skip", fmt.Sprint(p.skip))

	// If the API ever returns an empty page, treat as exhausted immediately.
	if p.lastLen == 0 {
		p.lastErr = PagerDone
		return nil, p.lastErr
	}

	return items, nil
}

// PrevPage fetches the previous page.
func (p *Pager[T]) PrevPage() ([]T, error) {
	if p.lastErr != nil && !errors.Is(p.lastErr, PagerDone) {
		return nil, p.lastErr
	}
	// Move the window back one page; if we're already at start, done.
	if !p.inited || p.skip-p.limit < 0 {
		return nil, PagerDone
	}

	p.skip -= p.limit
	if p.skip < 0 {
		p.skip = 0
	}
	p.query.Set("skip", fmt.Sprint(p.skip))

	items, err := p.fetchCurrent()
	if err != nil {
		p.lastErr = err
		return nil, err
	}

	p.lastLen = len(items)
	p.inited = true
	// Going backward means there can be a next page again.
	p.exhausted = false

	return items, nil
}

func (p *Pager[T]) fetchCurrent() ([]T, error) {
	qs := ""
	if len(p.query) > 0 {
		qs = "?" + p.query.Encode()
	}
	req, err := p.client.NewRequest(p.method, p.path+qs, nil)
	if err != nil {
		return nil, err
	}

	var bulk apiBulkResponse
	_, err = p.client.Do(p.ctx, req, &bulk)
	if err != nil {
		return nil, err
	}

	// honor server-provided limit if present; keep our cursor window consistent
	if bulk.Limit > 0 && bulk.Limit != p.limit {
		p.limit = bulk.Limit
		p.query.Set("limit", fmt.Sprint(p.limit))
	}

	var out []T
	if len(bulk.Records) > 0 {
		if err := json.Unmarshal(bulk.Records, &out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

type PageScanner[T any] struct {
	Pager pageIterator[T]
	list  []T
}

func (s *PageScanner[T]) Next() (item T, err error) {
	for len(s.list) == 0 {
		s.list, err = s.Pager.NextPage()
		if err != nil {
			var zero T
			return zero, err
		}
	}
	item = s.list[0]
	s.list = s.list[1:]
	return item, nil
}
