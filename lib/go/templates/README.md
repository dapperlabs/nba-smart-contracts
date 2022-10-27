# NBA Top Shot Transaction Templates

This module contains transaction and script templates for the Top Shot contracts.

## Generated manifest files

The `manifest.mainnet.json` and `testnet.mainnet.json` files declare all transaction templates
in a portable format for mainnet and testnet respectively.

To update the manifest files:

- Add your desired templates to [cmd/manifest/manifest.go](./cmd/manifest/manifest.go).
- Run `make generate` in this directory.
