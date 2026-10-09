package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FYI https://api-public-docs.cs-prod.leetify.com/#/player/get_v3_profile

// type LeetifyProfile struct {
// 	PrivacyMode    string      `json:"privacy_mode"`
// 	WinRate        json.Number `json:"winrate"`
// 	TotalMatches   json.Number `json:"total_matches"`
// 	FirstMatchDate string      `json:"first_match_date"`
// 	UserName       string      `json:"name"`
// 	Bans           []struct {
// 		Platform         string `json:"platform"`
// 		PlatformNickname string `json:"platform_nickname"`
// 		BannedSince      string `json:"banned_since"`
// 	} `json:"bans"`
// 	SteamID string `json:"steam64_id"`
// 	ID      string `json:"id"`
// 	Ranks   struct {
// 		Leetify     json.Number `json:"leetify"`
// 		Premier     json.Number `json:"premier"`
// 		FaceIT      json.Number `json:"faceit"`
// 		FaceITELO   json.Number `json:"faceit_elo"`
// 		Wingman     json.Number `json:"wingman"`
// 		Renown      json.Number `json:"renown"`
// 		Competitive []struct {
// 			MapName string      `json:"map_name"`
// 			Rank    json.Number `json:"rank"`
// 		} `json:"competitive"`
// 	} `json:"ranks"`
// 	Rating struct {
// 		Aim         json.Number `json:"aim"`
// 		Positioning json.Number `json:"positioning"`
// 		Utility     json.Number `json:"utility"`
// 		Clutch      json.Number `json:"clutch"`
// 		Opening     json.Number `json:"opening"`
// 		CtLeetify   json.Number `json:"ct_leetify"`
// 		Tleetify    json.Number `json:"t_leetify"`
// 	} `json:"rating"`
// 	Stats struct {
// 		AccuracyEnemySpotted           json.Number `json:"accuracy_enemy_spotted"`
// 		AccuracyHead                   json.Number `json:"accuracy_head"`
// 		CounterStrafingGoodShotsRatio  json.Number `json:"counter_strafing_good_shots_ratio"`
// 		CTOpeningAggressionSuccessRate json.Number `json:"ct_opening_aggression_success_rate"`
// 		CTOpeningDuelSuccessPercentage json.Number `json:"ct_opening_duel_success_percentage"`
// 		FlashbangHitFoeAvgDuration     json.Number `json:"flashbang_hit_foe_avg_duration"`
// 		FlashbangHitFoePerFlashbang    json.Number `json:"flashbang_hit_foe_per_flashbang"`
// 		FlashbangHitFriendPerFlashbang json.Number `json:"flashbang_hit_friend_per_flashbang"`
// 		FlashbangLeadingToKill         json.Number `json:"flashbang_leading_to_kill"`
// 		FlashbangThrown                json.Number `json:"flashbang_thrown"`
// 		HEFoesDamageAvg                json.Number `json:"he_foes_damage_avg"`
// 		HEFriendsDamageAvg             json.Number `json:"he_friends_damage_avg"`
// 		Preaim                         json.Number `json:"preaim"`
// 		ReactionTimeMS                 json.Number `json:"reaction_time_ms"`
// 		SprayAccuracy                  json.Number `json:"spray_accuracy"`
// 		TOpeningAggressionSuccessRate  json.Number `json:"t_opening_aggression_success_rate"`
// 		TOpeningDuelSuccessPercentage  json.Number `json:"t_opening_duel_success_percentage"`
// 		TradedDeathsSuccessPercentage  json.Number `json:"traded_deaths_success_percentage"`
// 		TradeKillOpportunitiesPerRound json.Number `json:"trade_kill_opportunities_per_round"`
// 		TradeKillsSuccessPercentage    json.Number `json:"trade_kills_success_percentage"`
// 		UtilityOnDeathAvg              json.Number `json:"utility_on_death_avg"`
// 	} `json:"stats"`
// 	RecentMatches []struct {
// 		ID                   string        `json:"id"`
// 		FinishedAt           string        `json:"finished_at"`
// 		DataSource           string        `json:"data_source"`
// 		Outcome              string        `json:"outcome"`
// 		Rank                 json.Number   `json:"rank"`
// 		RankType             json.Number   `json:"rank_type"`
// 		MapName              string        `json:"map_name"`
// 		LeetifyRating        json.Number   `json:"leetify_rating"`
// 		Score                []json.Number `json:"score"`
// 		Preaim               json.Number   `json:"preaim"`
// 		ReactionTimeMS       json.Number   `json:"reaction_time_ms"`
// 		AccuracyEnemySpotted json.Number   `json:"accuracy_enemy_spotted"`
// 		AccuracyHead         json.Number   `json:"accuracy_head"`
// 		SprayAccuracy        json.Number   `json:"spray_accuracy"`
// 	} `json:"recent_matches"`
// 	RecentTeammates []struct {
// 		SteamID            string      `json:"steam64_id"`
// 		RecentMatchesCount json.Number `json:"recent_matches_count"`
// 	} `json:"recent_teammates"`
// }

type CSMetricsProfile struct {
	SteamID       string      `json:"steamID"`
	UserName      string      `json:"name"`
	Kills         json.Number `json:"kills"`
	Deaths        json.Number `json:"deaths"`
	Assists       json.Number `json:"assists"`
	HeadshotKills json.Number `json:"headshotKills"`
	Damage        json.Number `json:"damage"`
	RoundsPlayed  json.Number `json:"roundsPlayed"`
	RoundsWon     json.Number `json:"roundsWon"`
	Mvps          json.Number `json:"mvps"`
	KastRounds    json.Number `json:"kastRounds"`
	MatchesPlayed json.Number `json:"matchesPlayed"`
	MatchesWon    json.Number `json:"matchesWon"`
	FirstKills    json.Number `json:"firstKills"`
	FirstDeaths   json.Number `json:"firstDeaths"`
	TradeKills    json.Number `json:"tradeKills"`
	TradeDeaths   json.Number `json:"tradeDeaths"`
	Kdr           json.Number `json:"kdr"`
	Adr           json.Number `json:"adr"`
	HsPercent     json.Number `json:"hsPercent"`
	KastPercent   json.Number `json:"kastPercent"`
	WinRate       json.Number `json:"winRate"`

	Ranks struct {
		Premier struct {
			Current json.Number `json:"current"`
			Peak    json.Number `json:"peak"`
			Wins    json.Number `json:"wins"`
		} `json:"premier"`

		Competitive struct {
			Current json.Number `json:"current"`
			Peak    json.Number `json:"peak"`
			Wins    json.Number `json:"wins"`
			Maps    []struct {
				Map     string      `json:"map"`
				Current json.Number `json:"current"`
				Peak    json.Number `json:"peak"`
				Wins    json.Number `json:"wins"`
			} `json:"maps"`
		} `json:"competitive"`

		Wingman struct {
			Current json.Number `json:"current"`
			Peak    json.Number `json:"peak"`
			Wins    json.Number `json:"wins"`
		} `json:"wingman"`

		FaceIT struct {
			Current json.Number `json:"current"`
			Peak    json.Number `json:"peak"`
			Wins    json.Number `json:"wins"`
			Linked  bool        `json:"linked"`
		} `json:"faceit"`
	} `json:"ranks"`

	Matches []struct {
		MatchUUID           string      `json:"matchUUID"`
		MatchID             string      `json:"matchID"`
		Map                 string      `json:"map"`
		Mode                string      `json:"mode"`
		MatchDate           string      `json:"matchDate"`
		Duration            json.Number `json:"duration"`
		CTScore             json.Number `json:"ctScore"`
		TScore              json.Number `json:"tScore"`
		Team                json.Number `json:"team"`
		Kills               json.Number `json:"kills"`
		Deaths              json.Number `json:"deaths"`
		Assists             json.Number `json:"assists"`
		HeadshotKills       json.Number `json:"headshotKills"`
		Damage              json.Number `json:"damage"`
		RoundsPlayed        json.Number `json:"roundsPlayed"`
		RoundsWon           json.Number `json:"roundsWon"`
		KASTRounds          json.Number `json:"kastRounds"`
		MVPs                json.Number `json:"mvps"`
		FirstKills          json.Number `json:"firstKills"`
		FirstDeaths         json.Number `json:"firstDeaths"`
		TradeKills          json.Number `json:"tradeKills"`
		TradeDeaths         json.Number `json:"tradeDeaths"`
		AverageTimeToShot   json.Number `json:"averageTimeToShot"`
		AverageTimeToDamage json.Number `json:"averageTimeToDamage"`
		AverageTimeToKill   json.Number `json:"averageTimeToKill"`
		MedianCrosshair     json.Number `json:"medianCrosshair"`
		Rank                json.Number `json:"rank"`
		UserParticipated    bool        `json:"userParticipated"`
	} `json:"matches"`

	Tracked bool `json:"tracked"`
}

// gets all csmetrics stats
func getCSMStats(steamID string) (CSMetricsProfile, error) {

	time.Sleep(1 * time.Second)

	if steamID == "" {
		return CSMetricsProfile{}, fmt.Errorf("empty steamID")
	}

	// url := fmt.Sprintf(
	// 	"https://api-public.cs-prod.leetify.com/v3/profile?steam64_id=%s",
	// 	steamID,
	// )

	url := fmt.Sprintf(
		"https://csmetrics.app/api/v1/player/%s",
		steamID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error while getting response from csmetrics: %s\n", err)
		return CSMetricsProfile{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error while getting response from csmetrics: %s\n", err)
		return CSMetricsProfile{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error while reading from csmetrics: %s\n", err)
		return CSMetricsProfile{}, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return CSMetricsProfile{}, fmt.Errorf("csmetrics profile not found")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("csmetrics returned status %d\n", resp.StatusCode)
		fmt.Println(string(body))
		return CSMetricsProfile{}, fmt.Errorf("non-200 response")
	}

	var profile CSMetricsProfile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		fmt.Printf("error while unmarshalling json from csmetrics: %s\n", err)
		return CSMetricsProfile{}, err
	}
	return profile, nil
}

func loopFunc() {
	fSteamIDList := linkGetAllSteamIds()
	dSteamIDList := dLinkGetAllSteamIds()

	for _, steamID := range fSteamIDList {
		playerstats, _ := getCSMStats(steamID)
		// looped stat functions go here

		csMatchListener(playerstats)
		updatePremierRatingRole(playerstats)
	}

	for _, steamID := range dSteamIDList {
		playerstats, _ := getCSMStats(steamID)
		// looped stat functions go here

		csMatchListener(playerstats)
		dUpdatePremierRatingRole(playerstats)
	}

	// checks all steamIDs, even if there are duplicates between the two lists
}

// starts with bot and runs stats functions every 5 mins
func initCSMStatsLoop() {
	loopFunc()
	for range time.Tick(time.Minute * 5) {
		loopFunc()
	}
}
