package sumoapi

import "time"

type BashoService service

type prizeEntry struct {
	Type      string `json:"type,omitempty"`
	RikishiID int    `json:"rikishiId,omitempty"`
	Name      string `json:"shikonaEn,omitempty"`
	//Shikona string `json:"shikonaJp,omitempty"`
}

type Basho struct {
	Date          string       `json:"date,omitempty"`
	Location      string       `json:"location,omitempty"`
	StartDate     time.Time    `json:"startDate,omitempty"`
	EndDate       time.Time    `json:"endDate,omitempty"`
	Yusho         []prizeEntry `json:"yusho,omitempty"`
	SpecialPrizes []prizeEntry `json:"specialPrizes,omitempty"`
}

type Torikumi struct {
	Basho

	Torikumi []Match `json:"torikumi,omitempty"`
}
