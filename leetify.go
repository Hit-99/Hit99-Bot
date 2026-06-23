package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// FYI https://api-public-docs.cs-prod.leetify.com/#/player/get_v3_profile

type LeetifyProfile struct {
	PrivacyMode    string      `json:"privacy_mode"`
	WinRate        json.Number `json:"winrate"`
	TotalMatches   json.Number `json:"total_matches"`
	FirstMatchDate string      `json:"first_match_date"`
	UserName       string      `json:"name"`
	Bans           []struct {
		Platform         string `json:"platform"`
		PlatformNickname string `json:"platform_nickname"`
		BannedSince      string `json:"banned_since"`
	} `json:"bans"`
	SteamID string `json:"steam64_id"`
	ID      string `json:"id"`
	Ranks   struct {
		Leetify     json.Number `json:"leetify"`
		Premier     json.Number `json:"premier"`
		FaceIT      json.Number `json:"faceit"`
		FaceITELO   json.Number `json:"faceit_elo"`
		Wingman     json.Number `json:"wingman"`
		Renown      json.Number `json:"renown"`
		Competitive []struct {
			MapName string      `json:"map_name"`
			Rank    json.Number `json:"rank"`
		} `json:"competitive"`
	} `json:"ranks"`
	Rating struct {
		Aim         json.Number `json:"aim"`
		Positioning json.Number `json:"positioning"`
		Utility     json.Number `json:"utility"`
		Clutch      json.Number `json:"clutch"`
		Opening     json.Number `json:"opening"`
		CtLeetify   json.Number `json:"ct_leetify"`
		Tleetify    json.Number `json:"t_leetify"`
	} `json:"rating"`
	Stats struct {
		AccuracyEnemySpotted           json.Number `json:"accuracy_enemy_spotted"`
		AccuracyHead                   json.Number `json:"accuracy_head"`
		CounterStrafingGoodShotsRatio  json.Number `json:"counter_strafing_good_shots_ratio"`
		CTOpeningAggressionSuccessRate json.Number `json:"ct_opening_aggression_success_rate"`
		CTOpeningDuelSuccessPercentage json.Number `json:"ct_opening_duel_success_percentage"`
		FlashbangHitFoeAvgDuration     json.Number `json:"flashbang_hit_foe_avg_duration"`
		FlashbangHitFoePerFlashbang    json.Number `json:"flashbang_hit_foe_per_flashbang"`
		FlashbangHitFriendPerFlashbang json.Number `json:"flashbang_hit_friend_per_flashbang"`
		FlashbangLeadingToKill         json.Number `json:"flashbang_leading_to_kill"`
		FlashbangThrown                json.Number `json:"flashbang_thrown"`
		HEFoesDamageAvg                json.Number `json:"he_foes_damage_avg"`
		HEFriendsDamageAvg             json.Number `json:"he_friends_damage_avg"`
		Preaim                         json.Number `json:"preaim"`
		ReactionTimeMS                 json.Number `json:"reaction_time_ms"`
		SprayAccuracy                  json.Number `json:"spray_accuracy"`
		TOpeningAggressionSuccessRate  json.Number `json:"t_opening_aggression_success_rate"`
		TOpeningDuelSuccessPercentage  json.Number `json:"t_opening_duel_success_percentage"`
		TradedDeathsSuccessPercentage  json.Number `json:"traded_deaths_success_percentage"`
		TradeKillOpportunitiesPerRound json.Number `json:"trade_kill_opportunities_per_round"`
		TradeKillsSuccessPercentage    json.Number `json:"trade_kills_success_percentage"`
		UtilityOnDeathAvg              json.Number `json:"utility_on_death_avg"`
	} `json:"stats"`
	RecentMatches []struct {
		ID                   string        `json:"id"`
		FinshedAt            string        `json:"finished_at"`
		DataSource           string        `json:"data_source"`
		Outcome              string        `json:"outcome"`
		Rank                 json.Number   `json:"rank"`
		RankType             json.Number   `json:"rank_type"`
		MapName              string        `json:"map_name"`
		LeetifyRating        json.Number   `json:"leetify_rating"`
		Score                []json.Number `json:"score"`
		Preaim               json.Number   `json:"preaim"`
		ReactionTimeMS       json.Number   `json:"reaction_time_ms"`
		AccuracyEnemySpotted json.Number   `json:"accuracy_enemy_spotted"`
		AccuracyHead         json.Number   `json:"accuracy_head"`
		SprayAccuracy        json.Number   `json:"spray_accuracy"`
	} `json:"recent_matches"`
	RecentTeammates []struct {
		SteamID            string      `json:"steam64_id"`
		RecentMatchesCount json.Number `json:"recent_matches_count"`
	} `json:"recent_teammates"`
}

// gets all leetify stats
func getLeetifyStats(steamID string) (LeetifyProfile, error) {

	time.Sleep(1 * time.Second)

	if steamID == "" {
		return LeetifyProfile{}, fmt.Errorf("empty steamID")
	}

	url := fmt.Sprintf(
		"https://api-public.cs-prod.leetify.com/v3/profile?steam64_id=%s",
		steamID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error while getting response from leetify: %s\n", err)
		return LeetifyProfile{}, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("LEETIFY_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error while getting response from leetify: %s\n", err)
		return LeetifyProfile{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error while reading from leetify: %s\n", err)
		return LeetifyProfile{}, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return LeetifyProfile{}, fmt.Errorf("leetify profile not found")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("leetify returned status %d\n", resp.StatusCode)
		fmt.Println(string(body))
		return LeetifyProfile{}, fmt.Errorf("non-200 response")
	}

	var profile LeetifyProfile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		fmt.Printf("error while unmarshalling json from leetify: %s\n", err)
		return LeetifyProfile{}, err
	}
	return profile, nil
}

func loopFunc() {
	fSteamIDList := linkGetAllSteamIds()
	dSteamIDList := dLinkGetAllSteamIds()

	for _, steamID := range fSteamIDList {
		playerstats, _ := getLeetifyStats(steamID)
		// looped stat functions go here

		csMatchListener(playerstats)
		updatePremierRatingRole(playerstats)
	}

	for _, steamID := range dSteamIDList {
		playerstats, _ := getLeetifyStats(steamID)
		// looped stat functions go here

		csMatchListener(playerstats)
		dUpdatePremierRatingRole(playerstats)
	}

	// checks all steamIDs, even if there are duplicates between the two lists
}

// starts with bot and runs stats functions every 5 mins
func initLeetifyStatsLoop() {
	loopFunc()
	for range time.Tick(time.Minute * 5) {
		loopFunc()
	}
}
