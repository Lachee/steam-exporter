package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lachee/steam-exporter/steam"
	"github.com/lpernett/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

var (
	steamGamesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "steam_games_total",
			Help: "Total number of games owned by the Steam user",
		},
		[]string{"steam_id"},
	)
	steamGamePlaytime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "steam_games_playtime_total_minutes",
			Help: "Total playtime of games.",
		},
		[]string{"steam_id", "app_id", "name", "platform"},
	)
	steamGameLastPlayed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "steam_game_last_played_timestamp",
			Help: "Unix timestamp of when the game was last played",
		},
		[]string{"steam_id", "app_id", "name"},
	)
)

func init() {
	prometheus.MustRegister(steamGamesCount, steamGamePlaytime, steamGameLastPlayed)

	defaultGatherer := prometheus.DefaultGatherer
	gatherers := prometheus.Gatherers{
		prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) {
			collectSteamMetrics()
			return defaultGatherer.Gather()
		}),
	}

	prometheus.DefaultGatherer = prometheus.Gatherers(gatherers)
}

func collectSteamMetrics() {
	apiKey := os.Getenv("STEAM_WEB_API_KEY")
	steamID := os.Getenv("STEAM_ID64")

	if apiKey == "" || steamID == "" {
		log.Printf("Steam API key or Steam ID not configured")
		return
	}

	// Fetch the owned games stats
	response, err := steam.GetOwnedGames(apiKey, steamID, &steam.GetOwnedGamesOptions{
		IncludeAppInfo:         true,
		IncludePlayedFreeGames: true,
	})
	if err != nil {
		log.Printf("Error fetching owned games: %v", err)
		return
	}
	log.Printf("Updated steam_games_total metric: %d games for user %s", response.Response.GameCount, steamID)
	steamGamesCount.WithLabelValues(steamID).Set(float64(response.Response.GameCount))

	// Update the metrics
	for _, game := range response.Response.Games {
		app_id := fmt.Sprintf("%d", game.AppID)
		name := game.Name
		steamGamePlaytime.WithLabelValues(steamID, app_id, name, "any").Set(float64(game.PlaytimeForever))
		steamGamePlaytime.WithLabelValues(steamID, app_id, name, "windows").Set(float64(game.PlaytimeWindowsForever))
		steamGamePlaytime.WithLabelValues(steamID, app_id, name, "linux").Set(float64(game.PlaytimeLinuxForever))
		steamGamePlaytime.WithLabelValues(steamID, app_id, name, "mac").Set(float64(game.PlaytimeMacForever))
		steamGamePlaytime.WithLabelValues(steamID, app_id, name, "deck").Set(float64(game.PlaytimeDeckForever))

		steamGameLastPlayed.WithLabelValues(steamID, app_id, name).Set(float64(game.LastPlayed))
	}

}
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
	}

	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Steam exporter on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
