# NBA Top Shot Go Packages

This directory conains packages for interacting with the NBA Top Shot
smart contracts from a Go programming environment.

# Package Guides

- `contracts`: Contains functions to generate the text of the contract code
for the contracts in the `/nba-smart-contracts/contracts` directory
- `events`: Contains go definitions for the events that are emitted by
the Top Shot contracts so that these events can be monitored by applications.
- `templates`: Contains functions to return transaction templates
for common transactions and scripts for interacting with the Top Shot
smart contracts.
- `templates/data`: Contains go constructs for representing play metadata
for Top Shot plays on chain.
- `test`: Contains automated go tests for testing the functionality
of the Top Shot smart contracts.