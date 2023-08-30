package data

// Cadence requires a mapping of string->string, which can be handled through json tags when marshalling.
// It also does not allow for null values, so we will be omitting them if empty
// Reference: https://docs.google.com/spreadsheets/d/1muUZowii0pqoyi6OPK1VPNJi7keSc8_5zI9vu_QvfOY/edit#gid=375111836
type PlayMetadata struct {
	FullName              string `json:",omitempty"`        // Not all plays have player information. Can be blank
	FirstName             string `json:",omitempty"`        // Not all plays have player information. Can be blank
	LastName              string `json:",omitempty"`        // Not all plays have player information. Can be blank
	Birthdate             string `json:",omitempty"`        // Not all plays have player information. Can be blank
	Birthplace            string `json:",omitempty"`        // Not all plays have player information. Can be blank
	JerseyNumber          string `json:",omitempty"`        // Not all plays have player information. Can be blank
	DraftTeam             string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	DraftYear             *int32 `json:",omitempty,string"` // Not all plays have draft information. Can be blank
	DraftSelection        string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	DraftRound            string `json:",omitempty"`        // Not all plays have draft information. Can be blank
	TeamAtMomentNBAID     string
	TeamAtMoment          string
	PrimaryPosition       string `json:",omitempty"`        // Not all plays have player information. Can be blank
	PlayerPosition        string `json:",omitempty"`        // Not all plays have player information. Can be blank
	Height                *int32 `json:",omitempty,string"` // Not all plays have player information. Can be blank
	Weight                *int32 `json:",omitempty,string"` // Not all plays have player information. Can be blank
	TotalYearsExperience  string `json:",omitempty"`        // Not all plays have player information. Can be blank
	NbaSeason             string `json:",omitempty"`        // Not all plays have player information. Can be blank
	DateOfMoment          string
	PlayCategory          string
	PlayType              string
	HomeTeamName          string
	AwayTeamName          string
	HomeTeamScore         *int32 `json:",string"`
	AwayTeamScore         *int32 `json:",string"`
	PlayerAutographType   string `json:",omitempty"`
	PlayerAutographDate   string `json:",omitempty"`
	PlayerAutographSigner string `json:",omitempty"`
	OverrideHeadline      string `json:",omitempty"`
	Tagline               string
}

// GenerateEmptyPlay generates a play with all its fields
// empty except for FullName for testing
func GenerateEmptyPlay(fullName string) PlayMetadata {
	num := int32(10)
	return PlayMetadata{FullName: fullName,
		FirstName:             "",
		LastName:              "",
		Birthdate:             "",
		Birthplace:            "",
		JerseyNumber:          "",
		TeamAtMomentNBAID:     "",
		TeamAtMoment:          "",
		PrimaryPosition:       "",
		PlayerPosition:        "",
		Height:                &num,
		Weight:                &num,
		TotalYearsExperience:  "",
		NbaSeason:             "",
		DateOfMoment:          "",
		PlayCategory:          "",
		PlayType:              "",
		HomeTeamName:          "",
		AwayTeamName:          "",
		HomeTeamScore:         &num,
		AwayTeamScore:         &num,
		PlayerAutographType:   "",
		PlayerAutographDate:   "",
		PlayerAutographSigner: "",
		OverrideHeadline:      "",
		Tagline:               "",
	}
}
