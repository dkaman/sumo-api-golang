package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sumo "github.com/dkaman/sumo-api-golang"
)

func main() {
	ctx := context.Background()
	client := sumo.NewClient(nil)

	// ---- Rikishi ----
	fmt.Println("=== Rikishi List (first page) ===")
	rikishiPage, err := client.Rikishi.ListPager(ctx, 10).NextPage()
	if err != nil {
		log.Fatalf("error getting rikishi list: %v", err)
	}
	for _, r := range rikishiPage {
		fmt.Printf("ID: %d | Name: %s | Rank: %s | Heya: %s\n", r.ID, r.Name, r.CurrentRank, r.Heya)
	}

	// Pick a sample Rikishi ID
	var rikishiID int
	if len(rikishiPage) > 0 {
		rikishiID = rikishiPage[0].ID
	} else {
		fmt.Println("No rikishi found to test other calls")
		os.Exit(0)
	}

	fmt.Printf("\n=== Rikishi Get (%d) ===\n", rikishiID)
	rikishi, err := client.Rikishi.Get(ctx, rikishiID)
	if err != nil {
		log.Fatalf("error getting rikishi: %v", err)
	}
	fmt.Printf("%+v\n", rikishi)

	fmt.Printf("\n=== Rikishi Stats (%d) ===\n", rikishiID)
	stats, err := client.Rikishi.Stats(ctx, rikishiID)
	if err != nil {
		log.Fatalf("error getting rikishi stats: %v", err)
	}
	fmt.Printf("%+v\n", stats)

	fmt.Printf("\n=== Rikishi Matches (%d, first page) ===\n", rikishiID)
	matches, err := client.Rikishi.MatchesPager(ctx, rikishiID, 5).NextPage()
	if err != nil {
		log.Fatalf("error getting rikishi matches: %v", err)
	}
	for _, m := range matches {
		fmt.Printf("Basho: %s | Day %d | %s vs %s | Winner: %s\n",
			m.BashoID, m.Day, m.EastName, m.WestName, m.WinnerName)
	}

	if len(matches) > 0 {
		opponentID := matches[0].WestID
		if opponentID == rikishiID {
			opponentID = matches[0].EastID
		}
		fmt.Printf("\n=== Rikishi Matchup (%d vs %d) ===\n", rikishiID, opponentID)
		matchup, err := client.Rikishi.Matchup(ctx, rikishiID, opponentID)
		if err != nil {
			log.Fatalf("error getting matchup: %v", err)
		}
		fmt.Printf("%+v\n", matchup)
	}

	// ---- Basho ----
	// Use a recent basho ID (YYYYMM format)
	bashoID := "202501"
	fmt.Printf("\n=== Basho Get (%s) ===\n", bashoID)
	basho, err := client.Basho.Get(ctx, bashoID)
	if err != nil {
		log.Printf("error getting basho: %v", err)
	} else {
		fmt.Printf("%+v\n", basho)
	}

	fmt.Printf("\n=== Basho Banzuke (%s, Makuuchi) ===\n", bashoID)
	banzuke, err := client.Basho.Banzuke(ctx, bashoID, sumo.Makuuchi)
	if err != nil {
		log.Printf("error getting banzuke: %v", err)
	} else {
		fmt.Printf("%+v\n", banzuke)
	}

	fmt.Printf("\n=== Basho Torikumi (%s, Makuuchi, Day 1) ===\n", bashoID)
	torikumi, err := client.Basho.Torikumi(ctx, bashoID, sumo.Makuuchi, 1)
	if err != nil {
		log.Printf("error getting torikumi: %v", err)
	} else {
		fmt.Printf("%+v\n", torikumi)
	}
}
