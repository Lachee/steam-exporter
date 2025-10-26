package steam

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetOwnedGamesOptions contains the optional parameters for the GetOwnedGames API call
type GetOwnedGamesOptions struct {
	// IncludeAppInfo includes game name and logo information in the output.
	// The default is to return appids only.
	IncludeAppInfo bool

	// IncludePlayedFreeGames includes free games like Team Fortress 2 that the player has played.
	// By default, free games are excluded (as technically everyone owns them).
	IncludePlayedFreeGames bool

	// AppIDsFilter optionally filters the list to a set of appids.
	// Note: This requires using JSON format for the API call.
	AppIDsFilter []uint32
}

// GetOwnedGames fetches the owned games for a Steam user using the Steam Web API
func GetOwnedGames(apiKey, steamID string, options *GetOwnedGamesOptions) (*GetOwnedGamesResponse, error) {
	baseURL := "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/"

	// Build query parameters
	params := url.Values{}
	params.Set("key", apiKey)
	params.Set("steamid", steamID)
	params.Set("format", "json")

	// Add options
	if options != nil {
		if options.IncludeAppInfo {
			params.Set("include_appinfo", "true")
		}
		if options.IncludePlayedFreeGames {
			params.Set("include_played_free_games", "true")
		}

		// Convert appids to comma-separated string
		if len(options.AppIDsFilter) > 0 {
			appIDs := make([]string, len(options.AppIDsFilter))
			for i, appID := range options.AppIDsFilter {
				appIDs[i] = fmt.Sprintf("%d", appID)
			}
			params.Set("appids_filter", strings.Join(appIDs, ","))
		}
	}

	// Make the HTTP request
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	return ParseGetOwnedGamesResponse(body)
}
