package events

import "github.com/onflow/cadence"

func NewCadenceString(str string) cadence.String {
	res, _ := cadence.NewString(str)
	return res
}
