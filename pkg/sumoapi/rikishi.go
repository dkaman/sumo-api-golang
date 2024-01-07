package sumoapi

import (
	"context"
	"net/http"
	"time"
)

type RikishiService service

type Rikishi struct {
	ID          int       `json:"id,omitempty"`
	NSKID       int       `json:"nskId,omitempty"`
	Shikona     string    `json:"shikonaEn,omitempty"`
	CurrentRank string    `json:"currentRank,omitempty"`
	Heya        string    `json:"heya,omitempty"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	Shusshin    string    `json:"shusshin,omitempty"`
	Height      int       `json:"height,omitempty"`
	Weight      int       `json:"weight,omitempty"`
	Debut       time.Time `json:"debut,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type RikishiResponse struct {
	Limit   int        `json:"limit,omitempty"`
	Skip    int        `json:"skip,omitempty"`
	Total   int        `json:"total,omitempty"`
	Records []*Rikishi `json:"records,omitempty"`
}

func (r *RikishiService) List(ctx context.Context) ([]*Rikishi, *http.Response, error) {
	u := "api/rikishis"

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResponse RikishiResponse
	resp, err := r.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, nil, err
	}

	return apiResponse.Records, resp, nil
}
