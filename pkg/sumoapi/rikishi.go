package sumoapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RikishiService service

type Rikishi struct {
	ID          int       `json:"id,omitempty"`
	SumoDBID    int       `json:"sumodbId,omitempty"`
	NSKID       int       `json:"nskId,omitempty"`
	Name        string    `json:"shikonaEn,omitempty"`
	CurrentRank string    `json:"currentRank,omitempty"`
	Heya        string    `json:"heya,omitempty"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	Shusshin    string    `json:"shusshin,omitempty"`
	Height      float64   `json:"height,omitempty"`
	Weight      float64   `json:"weight,omitempty"`
	Debut       string    `json:"debut,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	// Shikona     []byte    `json:"shikonaJp,omitempty"`
}

type RikishiGlobalStats struct {
	AbsenceByDivision map[string]int `json:"absenceByDivision,omitempty"`
	Basho             int            `json:"basho,omitempty"`
	BashoByDivision   map[string]int `json:"bashoByDivision,omitempty"`
	LossByDivision    map[string]int `json:"lossByDivision,omitempty"`
	Sansho            map[string]int `json:"sansho,omitempty"`
	TotalAbsences     int            `json:"totalAbsences,omitempy"`
	TotalByDivision   map[string]int `json:"totalByDivision,omitempty"`
	TotalLosses       int            `json:"totalLosses,omitempy"`
	TotalMatches      int            `json:"totalMatches,omitempy"`
	TotalWins         int            `json:"totalWins,omitempy"`
	WinsByDivision    map[string]int `json:"winsByDivision,omitempty"`
	Yusho             int            `json:"yusho,omitempy"`
	YushoByDivision   map[string]int `json:"yushoByDivision,omitempty"`
}

type Match struct {
	BashoID     string `json:"bashoId,omitempty"`
	Division    string `json:"division,omitempty"`
	Day         int    `json:"day,omitempty"`
	MatchNumber int    `json:"matchNo,omitempty"`
	EastID      int    `json:"eastId,omitempty"`
	EastName    string `json:"eastShikona,omitempty"`
	EastRank    string `json:"eastRank,omitempty"`
	WestID      int    `json:"westId,omitempty"`
	WestName    string `json:"westShikona,omitempty"`
	WestRank    string `json:"westRank,omitempty"`
	Kimarite    string `json:"kimarite,omitempty"`
	WinnerID    int    `json:"winnerId,omitempty"`
	WinnerName  string `json:"winnerEn,omitempty"`
	//WinnerShikona string `json:"winnerJp,omitempty"`
}

type MatchupStatistics struct {
	RikishiID      int
	KimariteLosses map[string]int `json:"kimariteLosses,omitempty"`
	KimariteWins   map[string]int `json:"kimariteWins,omitempty"`
	Matches        []Match        `json:"matches,omitempty"`
	OpponentWins   int            `json:"opponentWins,omitempty"`
	Wins           int            `json:"rikishiWins,omitempty"`
	Total          int            `json:"total,omitempty"`
}

func (r *RikishiService) List(ctx context.Context) ([]*Rikishi, *http.Response, error) {
	u := "api/rikishis"

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResponse apiBulkResponse

	resp, err := r.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, nil, err
	}

	var rikishi []*Rikishi

	err = json.Unmarshal(apiResponse.Records, &rikishi)
	if err != nil {
		return nil, nil, err
	}

	return rikishi, resp, nil
}

func (r *RikishiService) Get(ctx context.Context, id int) (*Rikishi, *http.Response, error) {
	u := fmt.Sprintf("api/rikishi/%d", id)

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rikishi Rikishi

	resp, err := r.client.Do(ctx, req, &rikishi)
	if err != nil {
		return nil, nil, err
	}

	return &rikishi, resp, nil
}

func (r *RikishiService) Stats(ctx context.Context, id int) (*RikishiGlobalStats, *http.Response, error) {
	u := fmt.Sprintf("api/rikishi/%d/stats", id)

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var gs RikishiGlobalStats

	resp, err := r.client.Do(ctx, req, &gs)
	if err != nil {
		return nil, nil, err
	}

	return &gs, resp, nil
}

func (r *RikishiService) Matches(ctx context.Context, id int) ([]*Match, *http.Response, error) {
	u := fmt.Sprintf("api/rikishi/%d/matches", id)

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResponse apiBulkResponse

	resp, err := r.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, nil, err
	}

	var matches []*Match

	err = json.Unmarshal(apiResponse.Records, &matches)
	if err != nil {
		return nil, nil, err
	}

	return matches, resp, nil
}

func (r *RikishiService) Matchup(ctx context.Context, subjectID int, opponentID int) (*MatchupStatistics, *http.Response, error) {
	u := fmt.Sprintf("api/rikishi/%d/matches/%d", subjectID, opponentID)

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var matchup MatchupStatistics

	resp, err := r.client.Do(ctx, req, &matchup)
	if err != nil {
		return nil, nil, err
	}

	matchup.RikishiID = subjectID

	return &matchup, resp, nil
}
