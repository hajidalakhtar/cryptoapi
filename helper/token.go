package helper

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ABIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

//func getTokenABI(tokenAddr string) string {
//	apiKey := RandomBSCScanApiKey()
//	resp, err := http.Get("https://api.bscscan.com/api?module=contract&action=getabi&address=" + tokenAddr + "&apikey=" + apiKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//	// Read the response body
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Failed to read response body:", err)
//	}
//
//	// Parse the JSON response
//	var abiResponse ABIResponse
//	err = json.Unmarshal(body, &abiResponse)
//	if err != nil {
//		fmt.Println("Failed to parse JSON response:", err)
//	}
//
//	// Check if the request was successful
//	if abiResponse.Status != "1" {
//		fmt.Println("Failed to fetch ABI:", abiResponse.Message)
//	}
//
//	return abiResponse.Result
//}

// Token is the Go binding of the solidity interface for your token
type Token struct {
	contract *abi.ABI
	address  common.Address
	callOpts *bind.CallOpts
	client   *ethclient.Client
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, tokenABI string, client *ethclient.Client) (*Token, error) {
	TokenABI := tokenABI
	parsedABI, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}

	return &Token{
		contract: &parsedABI,
		address:  address,
		callOpts: &bind.CallOpts{},
		client:   client,
	}, nil
}

func (t *Token) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	if !common.IsHexAddress(who.Hex()) {
		return nil, errors.New("invalid Ethereum address")
	}

	packedData, err := t.contract.Pack("balanceOf", who)
	if err != nil {
		return nil, err
	}

	// Create the Ethereum call messag
	callMsg := ethereum.CallMsg{
		To:   &t.address,
		Data: packedData,
	}

	output, err := t.client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	results := new(big.Int)
	if err := t.contract.UnpackIntoInterface(&results, "balanceOf", output); err != nil {
		return nil, err
	}

	return results, nil
}

func GetTokenBalance(client *ethclient.Client, addr string, tokenAddr string, tokenABI string) (*big.Float, error) {
	tokenAddress := common.HexToAddress(tokenAddr)
	address := common.HexToAddress(addr)
	tokenContract, err := NewToken(tokenAddress, tokenABI, client)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}

	balance, err := tokenContract.BalanceOf(nil, address)

	if err != nil {
		log.Fatalf("Failed to retrieve token balance: %v", err)
	}

	fbalance := new(big.Float).SetInt(balance)
	value := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return value, nil
}
