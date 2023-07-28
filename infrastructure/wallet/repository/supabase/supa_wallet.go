package supabase

import (
	"context"
	"cryptoapi/domain"
	"encoding/json"
	"fmt"
	"github.com/nedpals/supabase-go"
	"io/ioutil"
	"net/http"
)

type SupabaseWalletRepository struct {
	client *supabase.Client
}

type Response struct {
	ID int `json:"id"`
	// Add other fields if needed, depending on the API response.
}

func NewSupabaseWalletRepository(supaClient *supabase.Client) domain.SupabaseWalletRepository {
	return &SupabaseWalletRepository{
		client: supaClient,
	}
}

func (s SupabaseWalletRepository) GetToken(ctx context.Context, tokenAddr []string) ([]domain.Token, error) {
	var token []domain.Token
	err := s.client.DB.From("token").Select("*").In("smart_contract_address", tokenAddr).Execute(&token)
	if err != nil {
		panic(err)
	}
	return token, nil
}

func (s SupabaseWalletRepository) AddToken(ctx context.Context, strings []string) ([]domain.Token, error) {

	// URL of the API endpoint you want to query.
	apiURL := "https://api.coingecko.com/api/v3/coins/bsc/contract/0x6b175474e89094c44da98b954eedeac495271d0f"

	// Send the HTTP GET request to the API.
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making the HTTP request:", err)
	}
	defer response.Body.Close()

	// Read the response body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
	}

	// Parse the JSON response into our custom struct.
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
	}

	// Access the ID from the struct and print it.
	fmt.Println("ID:", data.ID)

	return nil, nil

}
