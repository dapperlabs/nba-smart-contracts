module github.com/dapperlabs/nba-smart-contracts/test

go 1.13

require (
	github.com/dapperlabs/flow-emulator v0.4.0
	github.com/dapperlabs/nba-smart-contracts/contracts v0.1.7
	github.com/dapperlabs/nba-smart-contracts/templates v0.1.6
	github.com/onflow/cadence v0.4.0
	github.com/onflow/flow-ft v0.1.0 // indirect
	github.com/onflow/flow-ft/contracts v0.1.2
	github.com/onflow/flow-ft/test v0.0.0-20200605203250-755c0ddcc598
	github.com/onflow/flow-go-sdk v0.4.0
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.28.0
)

replace github.com/dapperlabs/nba-smart-contracts/templates => ../templates

replace github.com/dapperlabs/nba-smart-contracts/contracts => ../contracts
