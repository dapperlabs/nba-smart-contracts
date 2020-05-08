package data

// Cadence requires a mapping of string->string, which can be handled through json tags when marshalling.
// It also does not allow for null values, so we will be omitting them if empty
type PlayMetadata struct {
	// VideoURLs  VideoUrls `sql:"video_urls"`
	// Images     Images    `sql:"images"`

	// Sport Radar Stats
	PlayerID                  string `json:",omitempty"`
	FullName                  string `json:",omitempty"`
	JerseyNumber              string `json:",omitempty"`
	TeamAtMoment              string `json:",omitempty"`
	PrimaryPosition           string `json:",omitempty"`
	TotalYearsExperience      string `json:",omitempty"`
	AwayTeamScore             *int32 `json:",omitempty,string"`
	AwayTeamName              string `json:",omitempty"`
	HomeTeamName              string `json:",omitempty"`
	HomeTeamScore             *int32 `json:",omitempty,string"`
	DateOfMoment              string `json:",omitempty"`
	TeamAtMomentNBAID         string `json:",omitempty"`
	CurrentTeam               string `json:",omitempty"`
	CurrentTeamID             string `json:",omitempty"`
	Height                    *int32 `json:",omitempty,string"`
	Weight                    *int32 `json:",omitempty,string"`
	PlayerGameScores          string `json:",omitempty"`
	PlayerSeasonAverageScores string `json:",omitempty"`
	HomeTeamNbaID             string `json:",omitempty"`
	AwayTeamNbaID             string `json:",omitempty"`
	NbaSeason                 string `json:",omitempty"`
	DraftYear                 *int32 `json:",omitempty,string"`
	DraftSelection            string `json:",omitempty"`
	DraftRound                string `json:",omitempty"`
	Birthplace                string `json:",omitempty"`
	Birthdate                 string `json:",omitempty"`
	DraftTeam                 string `json:",omitempty"`
	DraftTeamNbaID            string `json:",omitempty"`
	PlayType                  string `json:",omitempty"`
	PlayCategory              string `json:",omitempty"`
	Quarter                   string `json:",omitempty"`
	HomeTeamScoresByQuarter   string `json:",omitempty"`
	AwayTeamScoresByQuarter   string `json:",omitempty"`
}
