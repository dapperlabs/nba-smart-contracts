module github.com/dapperlabs/nba-smart-contracts/lib/go/test

go 1.13

require (
	github.com/dapperlabs/flow-emulator v0.4.0
	github.com/dapperlabs/nba-smart-contracts/lib/go/contracts v0.0.0-00010101000000-000000000000
	github.com/dapperlabs/nba-smart-contracts/lib/go/templates v0.0.0-00010101000000-000000000000
	github.com/onflow/cadence v0.4.0
	github.com/onflow/flow-ft/contracts v0.1.3
	github.com/onflow/flow-ft/test v0.0.0-20200619173914-64c953134397
	github.com/onflow/flow-go-sdk v0.4.0
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.28.0
)

replace github.com/dapperlabs/nba-smart-contracts/lib/go/templates => ../templates

replace github.com/dapperlabs/nba-smart-contracts/lib/go/contracts => ../contracts
