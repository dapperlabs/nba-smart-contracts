package main

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
)

type manifest struct {
	Network   string     `json:"network"`
	Templates []template `json:"templates"`
}

func (m *manifest) addTemplate(t template) {
	m.Templates = append(m.Templates, t)
}

type template struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Source    string     `json:"source"`
	Arguments []argument `json:"arguments"`
	Network   string     `json:"network"`
	Hash      string     `json:"hash"`
}

type argument struct {
	Type         string         `json:"type"`
	Name         string         `json:"name"`
	Label        string         `json:"label"`
	SampleValues []cadenceValue `json:"sampleValues"`
}

type cadenceValue struct {
	cadence.Value
}

func (v cadenceValue) MarshalJSON() ([]byte, error) {
	return jsoncdc.Encode(v.Value)
}

func (v cadenceValue) UnmarshalJSON(bytes []byte) (err error) {
	v.Value, err = jsoncdc.Decode(bytes)
	if err != nil {
		return err
	}

	return nil
}

type templateGenerator func(env templates.Environment) []byte

func generateTemplate(
	id, name string,
	env templates.Environment,
	generator templateGenerator,
	arguments []argument,
) template {
	source := generator(env)

	h := sha256.New()
	h.Write(source)
	hash := h.Sum(nil)

	return template{
		ID:        id,
		Name:      name,
		Source:    string(source),
		Arguments: arguments,
		Network:   env.Network,
		Hash:      hex.EncodeToString(hash),
	}
}

func generateManifest(env templates.Environment) *manifest {
	m := &manifest{
		Network: env.Network,
	}

	sampleMomentID := cadenceValue{
		cadence.NewUInt64(42),
	}

	m.addTemplate(generateTemplate(
		"TS.01", "Set up Top Shot Collection",
		env,
		templates.GenerateSetupAccountScript,
		[]argument{},
	))

	m.addTemplate(generateTemplate(
		"TS.02", "Transfer Top Shot Moment",
		env,
		templates.GenerateTransferMomentScript,
		[]argument{
			{
				Type:         "Address",
				Name:         "recipient",
				Label:        "Recipient",
				SampleValues: []cadenceValue{sampleAddress(env.Network)},
			},
			{
				Type:         "UInt64",
				Name:         "withdrawID",
				Label:        "Moment ID",
				SampleValues: []cadenceValue{sampleMomentID},
			},
		},
	))

	return m
}

func sampleAddress(network string) cadenceValue {
	var address flow.Address

	switch network {
	case testnet:
		address = flow.NewAddressGenerator(flow.Testnet).NextAddress()
	case mainnet:
		address = flow.NewAddressGenerator(flow.Mainnet).NextAddress()
	}

	return cadenceValue{cadence.Address(address)}
}
