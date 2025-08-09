package sumo-api-golang

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
	BirthDate   time.Time `json:"birthDate"`
	Shusshin    string    `json:"shusshin,omitempty"`
	Height      float64   `json:"height,omitempty"`
	Weight      float64   `json:"weight,omitempty"`
	Debut       string    `json:"debut,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

func (s *RikishiService) List(ctx context.Context) ([]*Rikishi, error) {
	p := s.ListPager(ctx, 100)

	var out []*Rikishi

	for {
		page, err := p.NextPage()
		if errors.Is(err, PagerDone) {
			break
		}

		if err != nil {
			return nil, err
		}

		out = append(out, page...)
	}

	return out, nil
}

func (s *RikishiService) ListPager(ctx context.Context, pageSize int) *Pager[*Rikishi] {
	return newPager[*Rikishi](ctx, s.client, "GET", "api/rikishis", url.Values{}, pageSize)
}

func (s *RikishiService) Get(ctx context.Context, id int) (*Rikishi, error) {
	u := fmt.Sprintf("api/rikishi/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var rikishi Rikishi

	resp, err := s.client.Do(ctx, req, &rikishi)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300  {
		return nil, fmt.Errorf("received non-200 status code from server: %d", resp.StatusCode)
	}

	return &rikishi, nil
}

func (s *RikishiService) Stats(ctx context.Context, id int) (*RikishiGlobalStats, error) {
	u := fmt.Sprintf("api/rikishi/%d/stats", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var gs RikishiGlobalStats

	resp, err := s.client.Do(ctx, req, &gs)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300  {
		return nil, fmt.Errorf("received non-200 status code from server: %d", resp.StatusCode)
	}

	return &gs, nil
}

func (s *RikishiService) Matches(ctx context.Context, id int) ([]*Match, error) {
	p := s.MatchesPager(ctx, id, 100)

	var out []*Match

	for {
		page, err := p.NextPage()
		if errors.Is(err, PagerDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, page...)
	}

	return out, nil
}

func (s *RikishiService) MatchesPager(ctx context.Context, id int, pageSize int) *Pager[*Match] {
	path := fmt.Sprintf("api/rikishi/%d/matches", id)
	return newPager[*Match](ctx, s.client, "GET", path, url.Values{}, pageSize)
}

func (s *RikishiService) Matchup(ctx context.Context, subjectID int, opponentID int) (*MatchupStatistics, error) {
	u := fmt.Sprintf("api/rikishi/%d/matches/%d", subjectID, opponentID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var matchup MatchupStatistics

	resp, err := s.client.Do(ctx, req, &matchup)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300  {
		return nil, fmt.Errorf("received non-200 status code from server: %d", resp.StatusCode)
	}

	matchup.RikishiID = subjectID

	return &matchup, nil
}
