package sumoapi

import (
	"context"
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

	var apiResponse RikishiResponse
	resp, err := r.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, nil, err
	}

	return apiResponse.Records, resp, nil
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
