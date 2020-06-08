package data

// Cadence requires a mapping of string->string, which can be handled through json tags when marshalling.
// It also does not allow for null values, so we will be omitting them if empty
// Reference: https://docs.google.com/spreadsheets/d/1muUZowii0pqoyi6OPK1VPNJi7keSc8_5zI9vu_QvfOY/edit#gid=375111836
type PlayMetadata struct {
	FullName             string
	FirstName            string
	LastName             string
	Birthdate            string
	Birthplace           string
	JerseyNumber         string
	DraftTeam            string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	DraftYear            *int32 `json:",omitempty,string"` // Not all plays have draft information. Can be blank
	DraftSelection       string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	DraftRound           string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	TeamAtMomentNBAID    string
	CurrentTeamID        string
	TeamAtMoment         string
	CurrentTeam          string
	PrimaryPosition      string
	PlayerPosition       string
	Height               *int32 `json:",string"`
	Weight               *int32 `json:",string"`
	TotalYearsExperience string
	NbaSeason            string
	DateOfMoment         string
	PlayCategory         string
	PlayType             string
	HomeTeamName         string
	AwayTeamName         string
	HomeTeamScore        *int32 `json:",string"`
	AwayTeamScore        *int32 `json:",string"`
}
