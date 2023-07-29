package helper

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Response struct {
	Name            string          `json:"name"`
	Symbol          string          `json:"symbol"`
	DetailPlatforms DetailPlatforms `json:"detail_platforms"`
	Image           Image           `json:"image"`
}

type DetailPlatforms struct {
	BinanceSmartChain BinanceSmartChain `json:"binance-smart-chain"`
}

type BinanceSmartChain struct {
	DecimalPlace    int    `json:"decimal_place"`
	ContractAddress string `json:"contract_address"`
}

type Image struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
	Large string `json:"large"`
}

func RandomBSCScanApiKey() string {
	apiKey := []string{"6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY", "6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY"}

	rand.Seed(time.Now().UnixNano())
	return apiKey[rand.Intn(len(apiKey))]
}

func GetCoingeckoTokenApiFromSmartContract(smartContact string) (Response, error) {
	apiURL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/bsc/contract/%s", smartContact)
	response, err := http.Get(apiURL)
	if err != nil {
		log.Warn().Err(err).Msg("Error making the HTTP request")
		return Response{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Error reading the response body")
		return Response{}, err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Warn().Err(err).Msg("Error unmarshalling the response body")
		return Response{}, err
	}

	return data, nil

}
