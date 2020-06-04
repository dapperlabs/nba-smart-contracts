package data

// Cadence requires a mapping of string->string, which can be handled through json tags when marshalling.
// It also does not allow for null values, so we will be omitting them if empty
type SetMetadata struct {
	ID               string
	FlowId           *uint32
	FlowSeriesNumber *uint32
	FlowName         string
}
