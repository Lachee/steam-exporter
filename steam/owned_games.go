package steam

import (
	"encoding/json"
	"fmt"
)

// GetOwnedGamesResponse represents the top-level response from the Steam Web API GetOwnedGames endpoint
type GetOwnedGamesResponse struct {
	Response OwnedGamesData `json:"response"`
}

// OwnedGamesData contains the actual owned games data within the response
type OwnedGamesData struct {
	GameCount int    `json:"game_count"`
	Games     []Game `json:"games"`
}

// Game represents a single game in the owned games list
type Game struct {
	// AppID is the unique identifier for the game
	AppID uint32 `json:"appid"`

	// Name is the name of the game (only present if include_appinfo was set to true)
	Name string `json:"name,omitempty"`

	// PlaytimeForever is the total number of minutes played "on record", since Steam began tracking
	// total playtime in early 2009
	PlaytimeForever        int `json:"playtime_forever"`
	PlaytimeLinuxForever   int `json:"playtime_linux_forever"`
	PlaytimeMacForever     int `json:"playtime_mac_forever"`
	PlaytimeDeckForever    int `json:"playtime_deck_forever"`
	PlaytimeWindowsForever int `json:"playtime_windows_forever"`

	LastPlayed int64 `json:"rtime_last_played"`

	// ImgIconURL is the filename of the icon image for the game
	// To construct the full URL: http://media.steampowered.com/steamcommunity/public/images/apps/{appid}/{hash}.jpg
	ImgIconURL string `json:"img_icon_url,omitempty"`

	// ImgLogoURL is the filename of the logo image for the game
	// To construct the full URL: http://media.steampowered.com/steamcommunity/public/images/apps/{appid}/{hash}.jpg
	ImgLogoURL string `json:"img_logo_url,omitempty"`

	// HasCommunityVisibleStats indicates there is a stats page with achievements or other game stats available
	// for this game. The uniform URL for accessing this data is:
	// http://steamcommunity.com/profiles/{steamid}/stats/{appid}
	HasCommunityVisibleStats bool `json:"has_community_visible_stats,omitempty"`
}

// GetImageURL constructs the full URL for a game's image given the image hash
func (g *Game) GetImageURL(imageHash string) string {
	if imageHash == "" {
		return ""
	}
	return fmt.Sprintf("http://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", g.AppID, imageHash)
}

// GetIconURL returns the full URL for the game's icon image
func (g *Game) GetIconURL() string {
	return g.GetImageURL(g.ImgIconURL)
}

// GetLogoURL returns the full URL for the game's logo image
func (g *Game) GetLogoURL() string {
	return g.GetImageURL(g.ImgLogoURL)
}

// GetStatsURL constructs the URL for the game's stats page for a given Steam ID
func (g *Game) GetStatsURL(steamID string) string {
	if !g.HasCommunityVisibleStats {
		return ""
	}
	return fmt.Sprintf("http://steamcommunity.com/profiles/%s/stats/%d", steamID, g.AppID)
}

// ParseGetOwnedGamesResponse parses a JSON response from the GetOwnedGames API
func ParseGetOwnedGamesResponse(data []byte) (*GetOwnedGamesResponse, error) {
	var response GetOwnedGamesResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
