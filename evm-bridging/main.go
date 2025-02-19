package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onflow/go-ethereum/accounts/abi"

	. "github.com/bjartek/overflow/v2"
)

/**
 * This script is used to deploy TopShot contracts on EVM networks
 */

// Overflow prefixes signer names with the current network - e.g. "emulator-topshot-signer"
// Ensure accounts in flow.json are named accordingly
var networks = []string{"emulator", "testnet", "mainnet"}
var scriptTypes = []string{"setup", "tests"}

// Config by network
type config struct {
	topShotFlowAddr                 string
	topshotCoaAddr                  string
	bridgeDeployedTopshotERC721Addr string
	flowEvmBridgeCoaAddr            string
	transferValidatorAddr           string
	royaltyRecipientAddr            string
	proxyContractAddr               string
	rpcUrl                          string
	verifierUrl                     string
	verifierProvider                string
}

const (
	placeholderEvmAddress = "0x1234567890abcdef1234567890abcdef12345678"
)

// Addresses by network
var configByNetwork = map[string]config{
	"emulator": {
		topShotFlowAddr:                 "abcdef1234567890",
		flowEvmBridgeCoaAddr:            placeholderEvmAddress,
		bridgeDeployedTopshotERC721Addr: placeholderEvmAddress,
		transferValidatorAddr:           "0xA000027A9B2802E1ddf7000061001e5c005A0000",
		royaltyRecipientAddr:            placeholderEvmAddress,
	},
	"testnet": {
		topShotFlowAddr:                 "877931736ee77cff",
		flowEvmBridgeCoaAddr:            "0x0000000000000000000000023f946ffbc8829bfd",
		bridgeDeployedTopshotERC721Addr: "0xB3627E6f7F1cC981217f789D7737B1f3a93EC519",
		transferValidatorAddr:           "0x721C0078c2328597Ca70F5451ffF5A7B38D4E947", // CreatorTokenTransferValidator
		royaltyRecipientAddr:            placeholderEvmAddress,
		rpcUrl:                          "https://testnet.evm.nodes.onflow.org",
		verifierUrl:                     "https://evm-testnet.flowscan.io/api",
		verifierProvider:                "blockscout",
	},
	"mainnet": {
		topShotFlowAddr:                 "0b2a3299cc857e29",
		flowEvmBridgeCoaAddr:            "0x00000000000000000000000249250a5c27ecab3b",
		bridgeDeployedTopshotERC721Addr: "0x50AB3a827aD268e9D5A24D340108FAD5C25dAD5f",
		transferValidatorAddr:           "0x721C0078c2328597Ca70F5451ffF5A7B38D4E947", // CreatorTokenTransferValidator
		// TODO: get royalty recipient
		royaltyRecipientAddr: placeholderEvmAddress,
		rpcUrl:               "https://mainnet.evm.nodes.onflow.org",
		verifierUrl:          "https://evm.flowscan.io/api",
		verifierProvider:     "blockscout",
	},
}

// Provider struct
type provider struct {
	Dir                string
	TopshotAccountName string
	Network            string
	Config             config
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

	// Get network and script type from command line argument
	scriptType, network := getSpecifiedNetworkAndScriptType()

	// Initialize provider
	p := provider{
		Dir:                dir,
		TopshotAccountName: "topshot-signer",
		Network:            network,
		Config:             configByNetwork[network],
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

	// Retrieve or create COA
	p.retrieveOrCreateCOA()

	// Deploy test contract
	switch scriptType {
	case "setup":
		p.setupProject()
	case "tests":
		p.tests()
	}
}

func (p *provider) setupProject() {
	// Compile contracts
	recompileContracts()

	// Deploy implementation contract
	implementationAddr := p.deployContract("BridgedTopShotMoments", "")

	// Generate encoded initialize function call
	encodedInitializeFunctionCall := generateEncodedInitializeFunctionCall(
		p.Config.topshotCoaAddr,
		p.Config.bridgeDeployedTopshotERC721Addr,
		p.Config.flowEvmBridgeCoaAddr,
		"NBA Top Shot",
		"TOPSHOT",
		// TODO: replace with actual baseTokenURI
		"https://api.cryptokitties.co/tokenuri/",
		p.Config.topShotFlowAddr,
		fmt.Sprintf("A.%s.TopShot.NFT", p.Config.topShotFlowAddr),
		// TODO: replace with actual contract metadata
		`{"name": "Name of NFT","description":"Description of NFT"}`,
	)

	// Deploy proxy contract
	proxyAddr := p.deployContract("ERC1967Proxy",
		generateProxyEncodedConstructorData(
			fmt.Sprintf("0x%s", implementationAddr),
			fmt.Sprintf("0x%s", encodedInitializeFunctionCall),
		),
	)
	p.Config.proxyContractAddr = fmt.Sprintf("0x%s", proxyAddr)

	// Verify contracts
	if p.Network == "testnet" || p.Network == "mainnet" {
		p.verifyContract(implementationAddr)
		p.verifyContract(proxyAddr)
	}

	// Set up royalty management
	p.setRoyaltyManagement()

	log.Printf("\n\nSETUP COMPLETE!")
}

func (p *provider) tests() {
	if p.Network != "emulator" {
		log.Fatal("Test script can only be run on the emulator network")
	}
	testContractAddr := p.deployContract("TestContract", "")

	log.Printf("\t...running test uint array encoding tx")
	result := p.OverflowState.Tx("tests/test_uint_array_encoding",
		WithSigner(p.TopshotAccountName),
		WithArg("evmContractAddress", testContractAddr),
	)
	checkNoErr(result.Err)
	log.Printf("Tx executed%s", separatorString())
}

func (p *provider) deployContract(name, encodedConstructorData string) string {
	log.Printf("\t...deploying %s contract", name)

	// Get contract bytecode
	bytecode := p.getContractBytecodeFromABIFile(name)
	// Debug log the bytecode
	log.Printf("Original bytecode: %s", bytecode)

	// NEW: Properly handle the constructor data
	if encodedConstructorData != "" {
		// Remove 0x prefix if present
		encodedConstructorData = strings.TrimPrefix(encodedConstructorData, "0x")
		// Ensure both parts are valid hex strings
		if _, err := hex.DecodeString(bytecode); err != nil {
			log.Fatalf("Invalid bytecode hex: %v", err)
		}
		if _, err := hex.DecodeString(encodedConstructorData); err != nil {
			log.Fatalf("Invalid constructor data hex: %v", err)
		}
		bytecode = bytecode + encodedConstructorData
	}

	// Deploy contract and return address
	result := p.OverflowState.Tx("admin/deploy/deploy_contract",
		WithSigner(p.TopshotAccountName),
		WithArg("bytecode", bytecode),
		WithArg("gasLimit", p.GasLimit),
	)
	checkNoErr(result.Err)
	address := getContractAddressFromEVMEvent(result)
	log.Printf("%s contract deployed to address: %s%s", name, address, separatorString())
	return address
}

// Read contract ABI file and return the bytecode object
func (p *provider) getContractBytecodeFromABIFile(contractName string) string {
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

func generateProxyEncodedConstructorData(implementationAddr, abiEncodedInitializeFunctionCall string) string {
	implementationAddr = strings.TrimPrefix(implementationAddr, "0x")
	abiEncodedInitializeFunctionCall = strings.TrimPrefix(abiEncodedInitializeFunctionCall, "0x")
	fmt.Printf("abiEncodedInitializeFunctionCall is like this: %+v\n", abiEncodedInitializeFunctionCall)
	//Run 'cast abi-encode'
	initCallBytes, err := hex.DecodeString(abiEncodedInitializeFunctionCall)
	if err != nil {
		log.Fatalf("Failed to decode init call hex: %v", err)
	}
	const contractABI = `[{
		"type": "constructor",
		"inputs": [
			{"name": "implementation", "type": "address"},
			{"name": "initializeCall", "type": "bytes"}
		]
}]`

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	implementation := common.HexToAddress(implementationAddr)
	// Print command for logging
	//encodedBytes := common.FromHex(abiEncodedInitializeFunctionCall)
	data, err := parsedABI.Pack("", implementation, initCallBytes)
	if err != nil {
		log.Fatalf("Failed to pack proxy ABI: %v", err)
	}

	// Convert to hex string
	hexData := hex.EncodeToString(data)

	// Debug log
	log.Printf("Constructor data hex: %s", hexData)
	return hexData
}

func (p *provider) verifyContract(contractAddr string) string {
	//Run 'cast abi-encode'
	cmd := exec.Command(
		"forge",
		"verify-contract",
		"--rpc-url",
		p.Config.rpcUrl,
		"--verifier",
		p.Config.verifierProvider,
		"--verifier-url",
		p.Config.verifierUrl,
		contractAddr,
	)
	// Print command for logging
	log.Println("Executing command:", cmd.String())

	// Get error output
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("abiEncode: Failed to run 'cast abi-encode' and get output from command: %v; %s", err, stderr.String())
	}

	return out.String()[2:]
}

// Generate ABI encoded initializer data for proxy contract
func generateEncodedInitializeFunctionCall(
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
	const contractABI = `[{
	"name": "initialize",
	"type": "function",
	"inputs": [
		{"name": "owner", "type": "address"},
		{"name": "underlyingNftContractAddress", "type": "address"},
		{"name": "vmBridgeAddress", "type": "address"},
		{"name": "name", "type": "string"},
		{"name": "symbol", "type": "string"},
		{"name": "baseTokenURI", "type": "string"},
		{"name": "cadenceNFTAddress", "type": "string"},
		{"name": "cadenceNFTIdentifier", "type": "string"},
		{"name": "contractMetadata", "type": "string"}
	]
}]` // Print command for logging
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
	owner = ensureHexPrefix(owner)
	underlyingNftContractAddress = ensureHexPrefix(underlyingNftContractAddress)
	vmBridgeAddress = ensureHexPrefix(vmBridgeAddress)

	ownerAddr := common.HexToAddress(owner)
	nftContractAddr := common.HexToAddress(underlyingNftContractAddress)
	bridgeAddr := common.HexToAddress(vmBridgeAddress)

	// Run and return output without 0x prefix
	data, err := parsedABI.Pack(
		"initialize",
		ownerAddr,
		nftContractAddr,
		bridgeAddr,
		name,
		symbol,
		baseTokenURI,
		cadenceNFTAddress,
		cadenceNFTIdentifier,
		contractMetadata,
	)
	if err != nil {
		log.Fatalf("Failed to pack ABI: %v", err)
	}

	log.Printf("Encoded ABI: 0x%x", data) // Output the hex-encoded transaction data
	return hex.EncodeToString(data)
}

// Retrieve or create a COA
func (p *provider) retrieveOrCreateCOA() {
	log.Printf("\t...retrieving COA")
	log.Printf(p.TopshotAccountName)
	topshotCOAHex, err := p.OverflowState.Script("get_evm_address_string",
		WithArg("flowAddress", p.OverflowState.Address(p.TopshotAccountName)),
	).GetAsInterface()
	log.Printf(fmt.Sprintf("%v", topshotCOAHex), err)
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
	p.Config.topshotCoaAddr = fmt.Sprintf("0x%s", topshotCOAHex.(string))
}

// Set up royalty management
func (p *provider) setRoyaltyManagement() {
	if p.Config.proxyContractAddr == "" {
		log.Fatal("Proxy contract not deployed")
	}

	log.Printf("\t...setting up royalty management")
	setRoyaltyManagementResult := p.OverflowState.Tx("admin/set_up_royalty_management",
		WithSigner(p.TopshotAccountName),
		WithArg("erc721C", p.Config.proxyContractAddr),
		WithArg("validator", p.Config.transferValidatorAddr),
		WithArg("royaltyRecipient", p.Config.royaltyRecipientAddr),
		WithArg("royaltyBasisPoints", 500),
	)
	checkNoErr(setRoyaltyManagementResult.Err)
	log.Printf("Royalty management set up%s", separatorString())
}

// Parse script type and network argument from command line
func getSpecifiedNetworkAndScriptType() (string, string) {
	if len(os.Args) < 3 {
		log.Fatal("Please provide a network and script type as an argument: ", networks)
	}
	scriptType := os.Args[1]
	network := os.Args[2]

	if !slices.Contains(scriptTypes, scriptType) {
		log.Fatal("Please provide a valid script type as an argument: ", scriptTypes)
	}

	if !slices.Contains(networks, network) {
		log.Fatal("Please provide a valid network as an argument: ", networks)
	}

	return scriptType, network
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

func checkNoErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func separatorString() string {
	return "\n--------------------------------\n"
}

// NEW: Added helper function
func ensureHexPrefix(addr string) string {
	if !strings.HasPrefix(addr, "0x") {
		return "0x" + addr
	}
	return addr
}
