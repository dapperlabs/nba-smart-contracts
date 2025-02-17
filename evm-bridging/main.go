package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/bjartek/overflow/v2"
)

/**
 * This script is used to deploy TopShot contracts on EVM networks
 */

// Overflow prefixes signer names with the current network - e.g. "emulator-topshot-signer"
// Ensure accounts in flow.json are named accordingly
var networks = []string{"emulator", "testnet", "mainnet"}

// Addresses by network
type addressesByNetwork struct {
	topShotFlow                  string
	topshotCoa                   string
	bridgedDeployedTopshotERC721 string
	flowEvmBridgeCoa             string
	transferValidator            string
	royaltyRecipient             string
	proxyContract                string
}

const (
	placeholderEvmAddress = "0x1234567890abcdef1234567890abcdef12345678"
)

// Addresses by network
var addresses = map[string]addressesByNetwork{
	"emulator": {
		topShotFlow:                  "abcdef1234567890",
		flowEvmBridgeCoa:             placeholderEvmAddress,
		bridgedDeployedTopshotERC721: placeholderEvmAddress,
		transferValidator:            "0xA000027A9B2802E1ddf7000061001e5c005A0000",
		royaltyRecipient:             placeholderEvmAddress,
	},
	"testnet": {
		topShotFlow:                  "877931736ee77cff",
		flowEvmBridgeCoa:             "0x0000000000000000000000023f946ffbc8829bfd",
		bridgedDeployedTopshotERC721: "",
		transferValidator:            "0xA000027A9B2802E1ddf7000061001e5c005A0000", // StrictAuthorizedTransferSecurityRegistry
		royaltyRecipient:             placeholderEvmAddress,
	},
	"mainnet": {
		topShotFlow:                  "0b2a3299cc857e29",
		flowEvmBridgeCoa:             "0x00000000000000000000000249250a5c27ecab3b",
		bridgedDeployedTopshotERC721: "",
		// TODO: confirm StrictAuthorizedTransferSecurityRegistry or CreatorTokenTransferValidator (0x721C0078c2328597Ca70F5451ffF5A7B38D4E947)
		transferValidator: "0xA000027A9B2802E1ddf7000061001e5c005A0000", // StrictAuthorizedTransferSecurityRegistry
		// TODO: get royalty recipient
		royaltyRecipient: placeholderEvmAddress,
	},
}

// Provider struct
type provider struct {
	Dir                string
	TopshotAccountName string
	Network            string
	Addresses          addressesByNetwork
	GasLimit           int
	*OverflowState
}

// To run, execute 'go run main.go $NETWORK'
func main() {
	// Check prerequisites
	for _, prerequisite := range []string{"forge", "cast", "flow"} {
		if _, err := exec.LookPath(prerequisite); err != nil {
			log.Fatalf("Please install %s", prerequisite)
		}
	}

	// Get current directory
	dir, err := os.Getwd()
	checkNoErr(err)

	// Get network from command line argument
	network := getSpecifiedNetwork()

	// Initialize provider
	p := provider{
		Dir:                dir,
		TopshotAccountName: "topshot-signer",
		Network:            network,
		Addresses:          addresses[network],
		GasLimit:           15000000,
		OverflowState: Overflow(
			WithNetwork(network),
			WithFlowConfig("cadence/transactions/admin/deploy/flow.json"),
			WithTransactionFolderName("cadence/transactions"),
			WithScriptFolderName("cadence/scripts"),
			WithGlobalPrintOptions(WithTransactionUrl()),
			WithLogNone(),
		),
	}
	log.Printf("Provider initialized%s", separatorString())

	// Deploy implementation and proxy contracts
	p.retrieveOrCreateCOA()
	recompileContracts()
	p.deployContracts()

	// Verify contracts
	if p.Network == "testnet" || p.Network == "mainnet" {
		// TODO
	}

	// Set up royalty management
	p.setRoyaltyManagement()

	log.Printf("\n\nSETUP COMPLETE!")
}

// Retrieve or create a COA
func (p *provider) retrieveOrCreateCOA() {
	log.Printf("\t...retrieving COA")
	topshotCOAHex, err := p.OverflowState.Script("get_evm_address_string",
		WithArg("flowAddress", p.OverflowState.Address(p.TopshotAccountName)),
	).GetAsInterface()
	checkNoErr(err)
	if topshotCOAHex != nil {
		log.Printf("Using existing COA with EVM address: %s", topshotCOAHex)
	} else {
		log.Printf("\t...creating new COA")
		createCoaResult := p.OverflowState.Tx("admin/deploy/create_coa",
			WithSigner(p.TopshotAccountName),
			WithArg("amount", 1.0),
		)
		checkNoErr(createCoaResult.Err)
		topshotCOAHex, err = p.OverflowState.Script("get_evm_address_string",
			WithArg("flowAddress", p.OverflowState.Address(p.TopshotAccountName)),
		).GetAsInterface()
		checkNoErr(err)
		log.Printf("Created new COA with EVM address: %s%s", topshotCOAHex, separatorString())
	}
	p.Addresses.topshotCoa = fmt.Sprintf("0x%s", topshotCOAHex.(string))
}

// Deploy implementation and proxy contracts
func (p *provider) deployContracts() {
	// Deploy implementation contract
	log.Printf("\t...deploying implementation contract")
	deployImplementationResult := p.OverflowState.Tx("admin/deploy/deploy_contract",
		WithSigner(p.TopshotAccountName),
		WithArg("bytecode", p.getContractBytecodeFromABI("BridgedTopShotMoments")),
		WithArg("gasLimit", p.GasLimit),
	)
	checkNoErr(deployImplementationResult.Err)
	implementationAddr := getContractAddressFromEVMEvent(deployImplementationResult)
	log.Printf("Implementation contract deployed to address: %s%s", implementationAddr, separatorString())

	// Generate initialize data
	initializeData := generateAbiEncodedInitializerData(
		p.Addresses.topshotCoa,
		p.Addresses.bridgedDeployedTopshotERC721,
		p.Addresses.flowEvmBridgeCoa,
		`"NBA Top Shot"`,
		"TOPSHOT",
		// TODO: replace with actual baseTokenURI
		"https://api.cryptokitties.co/tokenuri/",
		p.Addresses.topShotFlow,
		fmt.Sprintf("A.%s.TopShot.NFT", p.Addresses.topShotFlow),
		// TODO: replace with actual contract metadata
		`data:application/json;utf8,{\"name\": \"Name of NFT\",\"description\":\"Description of NFT\"}`,
	)

	// Deploy proxy contract
	log.Printf("\t...deploying proxy contract")
	deployProxyResult := p.OverflowState.Tx("admin/deploy/deploy_contract",
		WithSigner(p.TopshotAccountName),
		WithArg("bytecode", fmt.Sprintf("%s%s",
			p.getContractBytecodeFromABI("ERC1967Proxy"),
			initializeData,
		)),
		WithArg("gasLimit", p.GasLimit),
		WithArg("implementation", implementationAddr),
	)
	checkNoErr(deployProxyResult.Err)
	proxyAddr := getContractAddressFromEVMEvent(deployProxyResult)
	p.Addresses.proxyContract = fmt.Sprintf("0x%s", proxyAddr)
	log.Printf("Proxy contract deployed to address: %s%s", proxyAddr, separatorString())
}

// Set up royalty management
func (p *provider) setRoyaltyManagement() {
	if p.Addresses.proxyContract == "" {
		log.Fatal("Proxy contract not deployed")
	}

	log.Printf("\t...setting up royalty management")
	setRoyaltyManagementResult := p.OverflowState.Tx("admin/set_up_royalty_management",
		WithSigner(p.TopshotAccountName),
		WithArg("erc721C", p.Addresses.proxyContract),
		WithArg("validator", p.Addresses.transferValidator),
		WithArg("royaltyRecipient", p.Addresses.royaltyRecipient),
		WithArg("royaltyBasisPoints", 500),
	)
	checkNoErr(setRoyaltyManagementResult.Err)
	log.Printf("Royalty management set up%s", separatorString())
}

// Parse network argument from command line
func getSpecifiedNetwork() string {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a network as an argument: ", networks)
	}
	network := os.Args[1]

	if !slices.Contains(networks, network) {
		log.Fatal("Please provide a valid network as an argument: ", networks)
	}
	return network
}

// Extract deployed contract address from TransactionExecuted event
func getContractAddressFromEVMEvent(res *OverflowResult) string {
	evts := res.GetEventsWithName("TransactionExecuted")
	contractAddr := evts[0].Fields["contractAddress"]
	if contractAddr == nil {
		log.Fatal("Contract address not found in event")
	}
	return strings.ToLower(strings.Split(contractAddr.(string), "x")[1])
}

// Read contract ABI file and return the bytecode object
func (p *provider) getContractBytecodeFromABI(contractName string) string {
	// Read ABI
	abiFile, err := os.ReadFile(filepath.Join(p.Dir, fmt.Sprintf("out/%s.sol/%s.json", contractName, contractName)))
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Parse and return bytecode object without 0x prefix
	var abi struct {
		Bytecode struct {
			Object string `json:"object"`
		} `json:"bytecode"`
	}
	if err := json.Unmarshal(abiFile, &abi); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return abi.Bytecode.Object[2:]
}

// Compile contracts by running 'forge clean' and 'forge build'
func recompileContracts() {
	// Run 'forge clean'
	cleanCmd := exec.Command("forge", "clean")
	log.Println("Executing command:", cleanCmd.String())
	if err := cleanCmd.Run(); err != nil {
		log.Fatalf("Failed to run 'forge clean': %v", err)
	}

	// Run 'forge build'
	buildCmd := exec.Command("forge", "build")
	log.Println("Executing command:", buildCmd.String())
	buildCmdOutput, err := buildCmd.Output()
	if err != nil {
		log.Fatalf("Failed to run 'forge build': %v", err)
	}
	log.Println("Output:\n", string(buildCmdOutput))
}

// Generate ABI encoded initializer data for proxy contract
func generateAbiEncodedInitializerData(
	owner,
	underlyingNftContractAddress,
	vmBridgeAddress,
	name,
	symbol,
	baseTokenURI,
	cadenceNFTAddress,
	cadenceNFTIdentifier,
	contractMetadata string,
) string {
	//Run 'cast abi-encode'
	abiEncodeCmd := exec.Command(
		"cast",
		"abi-encode",
		`initialize(address,address,address,string,string,string,string,string,string)`,
		owner,
		underlyingNftContractAddress,
		vmBridgeAddress,
		name,
		symbol,
		baseTokenURI,
		cadenceNFTAddress,
		cadenceNFTIdentifier,
		contractMetadata,
	)
	// Print command for logging
	log.Println("Executing command:", abiEncodeCmd.String())

	// Run and return output without 0x prefix
	output, err := abiEncodeCmd.Output()
	if err != nil {
		log.Fatalf("Failed to run 'cast abi-encode' and get output from command: %v", err)
	}
	return string(output)[2:]
}

func checkNoErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func separatorString() string {
	return "\n--------------------------------\n"
}
