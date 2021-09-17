package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/psiemens/sconfig"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
)

type Config struct {
	Network string `default:"mainnet" flag:"network" info:"Flow network to generate for"`
}

const envPrefix = "FLOW"

const (
	testnet = "testnet"
	mainnet = "mainnet"
)

const (
	testnetNonFungibleTokenAddress = "631e88ae7f1d7c20"
	testnetTopShotAddress     = "877931736ee77cff"
)

const (
	mainnetNonFungibleTokenAddress = "1d7e57aa55817448"
	mainnetTopShotAddress     = "0b2a3299cc857e29"
)

var conf Config

var cmd = &cobra.Command{
	Use:   "manifest <outfile>",
	Short: "Generate a JSON manifest of all core transaction templates",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, err := getEnv(conf)
		if err != nil {
			exit(err)
		}

		manifest := generateManifest(env)

		b, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			exit(err)
		}

		outfile := args[0]

		err = ioutil.WriteFile(outfile, b, 0777)
		if err != nil {
			exit(err)
		}
	},
}

func getEnv(conf Config) (templates.Environment, error) {

	if conf.Network == testnet {
		return templates.Environment{
			Network:              testnet,
			NFTAddress: testnetNonFungibleTokenAddress,
			TopShotAddress:     testnetTopShotAddress,
		}, nil
	}

	if conf.Network == mainnet {
		return templates.Environment{
			Network:              mainnet,
			NFTAddress: mainnetNonFungibleTokenAddress,
			TopShotAddress:     mainnetTopShotAddress,
		}, nil
	}

	return templates.Environment{}, fmt.Errorf("invalid network %s", conf.Network)
}

func init() {
	initConfig()
}

func initConfig() {
	err := sconfig.New(&conf).
		FromEnvironment(envPrefix).
		BindFlags(cmd.PersistentFlags()).
		Parse()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
