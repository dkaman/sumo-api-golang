package sumo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	BashoIDValidationErr = errors.New("basho id validation error")
)

type BashoService service

type Basho struct {
	Date          string       `json:"date,omitempty"`
	Location      string       `json:"location,omitempty"`
	StartDate     string       `json:"startDate"`
	EndDate       string       `json:"endDate"`
	Yusho         []prizeEntry `json:"yusho,omitempty"`
	SpecialPrizes []prizeEntry `json:"specialPrizes,omitempty"`
}

type prizeEntry struct {
	Type      string `json:"type,omitempty"`
	RikishiID int    `json:"rikishiId,omitempty"`
	Name      string `json:"shikonaEn,omitempty"`
	//Shikona string `json:"shikonaJp,omitempty"`
}

type Banzuke struct {
	BashoID  string         `json:"bashoId"`
	Division string         `json:"division"`
	East     []rikishiEntry `json:"east"`
	West     []rikishiEntry `json:"west"`
}

type rikishiEntry struct {
	Side      string `json:"side"`
	RikishiID int    `json:"rikishiID"`
	Name      string `json:"shikonaEn"`
	RankValue int    `json:"rankValue"`
	Rank      string `json:"rank"`
	Record    []Bout `json:"record"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
	Absences  int    `json:"absences"`
}

type Bout struct {
	Result            string `json:"result"`
	OpponentName string `json:"opponentShikonaEn"`
	OpponentID        int    `json:"opponentID"`
	Kimarite          string `json:"kimarite"`
	//OpponentShikona string `json:"opponentShikonaJp"`
}

type Torikumi struct {
	Basho

	Torikumi []Match `json:"torikumi,omitempty"`
}

type Division string
type Result string

const (
	Makuuchi  Division = "Makuuchi"
	Juryo     Division = "Juryo"
	Makushita Division = "Makushita"
	Sandanme  Division = "Sandanme"
	Jonidan   Division = "Jonidan"
	Jonokuchi Division = "Jonokuchi"

	Win    Result = "win"
	Loss   Result = "loss"
	Absent Result = "absent"
)

func validateID(id string) error {
	if len(id) != 6 {
		return fmt.Errorf("%w: must be 6 digits", BashoIDValidationErr)
	}

	yearPart := id[:4]
	monthPart := id[4:]

	year, err := strconv.Atoi(yearPart)
	if err != nil {
		return fmt.Errorf("%w: year is not a number", BashoIDValidationErr)
	}

	month, err := strconv.Atoi(monthPart)
	if err != nil {
		return fmt.Errorf("%w: month is not a number", BashoIDValidationErr)
	}

	if year < 1 || year > time.Now().Year() {
		return fmt.Errorf("%w: year '%d' is invalid", BashoIDValidationErr, year)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("%w: month '%d' is invalid", BashoIDValidationErr, month)
	}

	return nil
}

func (s *BashoService) Get(ctx context.Context, id string) (*Basho, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("api/basho/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var b Basho

	resp, err := s.client.Do(ctx, req, &b)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-200 http status code from server: %d", resp.StatusCode)
	}

	return &b, nil
}

func (s *BashoService) Banzuke(ctx context.Context, id string, division Division) (*Banzuke, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("api/basho/%s/banzuke/%s", id, division)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var b Banzuke

	resp, err := s.client.Do(ctx, req, &b)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-200 http status code from server: %d", resp.StatusCode)
	}

	return &b, nil
}

func (s *BashoService) Torikumi(ctx context.Context, id string, division Division, day int) (*Torikumi, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("api/basho/%s/torikumi/%s/%d", id, division, day)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var t Torikumi

	resp, err := s.client.Do(ctx, req, &t)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-200 http status code from server: %d", resp.StatusCode)
	}

	return &t, nil
}
