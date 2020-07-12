module github.com/dapperlabs/nba-smart-contracts/lib/go/test

go 1.13

require (
	github.com/dapperlabs/flow-emulator v0.5.0
	github.com/dapperlabs/nba-smart-contracts/lib/go/contracts v0.0.0-00010101000000-000000000000
	github.com/dapperlabs/nba-smart-contracts/lib/go/templates v0.0.0-00010101000000-000000000000
	github.com/onflow/cadence v0.5.1
	github.com/onflow/flow-ft v0.1.3-0.20200710214526-873e655e7c37 // indirect
	github.com/onflow/flow-ft/lib/go/contracts v0.0.0-20200629211940-37a9fc480521
	github.com/onflow/flow-ft/lib/go/templates v0.0.0-20200629211940-37a9fc480521
	github.com/onflow/flow-go-sdk v0.7.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.28.0
)

replace github.com/dapperlabs/nba-smart-contracts/lib/go/templates => ../templates

replace github.com/dapperlabs/nba-smart-contracts/lib/go/contracts => ../contracts
