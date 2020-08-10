# NBA Top Shot Go Packages

This directory conains packages for interacting with the NBA Top Shot
smart contracts from a Go programming environment.

# Package Guides

- `contracts`: Contains functions to generate the text of the contract code
for the contracts in the `/nba-smart-contracts/contracts` directory.
To generate the contracts:
1. Fetch the `contracts` package: `go get github.com/dapperlabs/nba-smart-contracts/contracts@v0.1.9`
2. Import the package at the top of your Go File: `import "github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"`
3. Call the `GenerateTopShotContract` and others to generate the full text of the contracts.
- `events`: Contains go definitions for the events that are emitted by
the Top Shot contracts so that these events can be monitored by applications.
- `templates`: Contains functions to return transaction templates
for common transactions and scripts for interacting with the Top Shot
smart contracts.
If you want to import the transactions in your Go programs
so you can submit them to interact with the NBA Top Shot smart contracts, 
you can do so with the `templates` package:
1. Fetch the `templates` package: `go get github.com/dapperlabs/nba-smart-contracts/templates@v0.1.10`
2. Import the package at the top of your Go File: `import "github.com/dapperlabs/nba-smart-contracts/lib/go/templates"`
3. Call the various functions in the `templates` package like `templates.GenerateTransferMomentScript()` and others to generate the full text of the templates that you can fill in with your arguments.
- `templates/data`: Contains go constructs for representing play metadata
for Top Shot plays on chain.
- `test`: Contains automated go tests for testing the functionality
of the Top Shot smart contracts.