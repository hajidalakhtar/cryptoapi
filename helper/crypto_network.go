package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type Endpoint struct {
	URL    string `json:"url"`
	Failed bool   `json:"failed"`
}

var chainsFile = "chains.json"
var chains []Endpoint

// init loads the chains data from the JSON file when the package is initialized.
func init() {
	loadChainsFromFile()
}

// GetNetwork returns a random BSC testnet endpoint.
func GetNetwork() *ethclient.Client {

	testnet := getAvailableTestnetEndpoints()

	if len(testnet) == 0 {
		log.Warn().Msg("No available BSC testnet endpoints")
		return &ethclient.Client{}
	}

	for _, chain := range testnet {
		client, err := ethclient.Dial(chain.URL)
		chainId, err := client.ChainID(context.Background())
		if err != nil || chainId == nil {
			log.Warn().Err(err).Msg("Trying to connect to other endpoint")
			MarkEndpointAsFailed(chain.URL)

			continue
		}
		return client

	}

	log.Warn().Msg("No available BSC testnet endpoints")
	return &ethclient.Client{}

}

// MarkEndpointAsFailed adds the given endpoint URL to the list of failed endpoints and saves the list to the file.
func MarkEndpointAsFailed(endpointURL string) {
	for i, chain := range chains {
		if chain.URL == endpointURL {
			chains[i].Failed = true
			break
		}
	}
	saveChainsToFile()
}

// getAvailableTestnetEndpoints returns the list of available BSC testnet endpoints.
func getAvailableTestnetEndpoints() []Endpoint {
	var availableEndpoints []Endpoint
	for _, chain := range chains {
		if !chain.Failed {
			availableEndpoints = append(availableEndpoints, chain)
		}
	}
	return availableEndpoints
}

// saveChainsToFile saves the chains data to the JSON file.
func saveChainsToFile() {
	fmt.Println("Saving chains to file")
	dir, _ := os.Getwd()

	data, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		// Handle error (e.g., log or return an error)
		return
	}

	err = ioutil.WriteFile(dir+"/"+chainsFile, data, 0644)
	if err != nil {
		// Handle error (e.g., log or return an error)
	}
}

// loadChainsFromFile loads the chains data from the JSON file.
func loadChainsFromFile() {
	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(dir + "/" + chainsFile)
	if err != nil {
		// Handle error (e.g., log or return an error)
		return
	}

	err = json.Unmarshal(data, &chains)
	if err != nil {
		// Handle error (e.g., log or return an error)
		return
	}
}
