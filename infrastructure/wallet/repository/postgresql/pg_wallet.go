package postgresql

import (
	"context"
	"cryptoapi/domain"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

type PostgresqlWalletRepository struct {
	conn *pgx.Conn
}

func NewPgWalletRepository(conn *pgx.Conn) domain.PostgresqlWalletRepository {
	return &PostgresqlWalletRepository{
		conn: conn,
	}
}

func (s PostgresqlWalletRepository) GetTokens(ctx context.Context, tokenAddr []string) ([]domain.Token, error) {
	tokenAddressesStr := "'" + strings.Join(tokenAddr, "', '") + "'"
	query := fmt.Sprintf(`SELECT id, name, symbol, smart_contract_address, decimal, logo_url, token_abi FROM crypto_tokens WHERE smart_contract_address IN (%s)`, tokenAddressesStr)
	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		log.Warn().Err(err).Msg("Error executing the query")
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.Token

	for rows.Next() {
		var token domain.Token
		var id uuid.UUID

		err := rows.Scan(
			&id,
			&token.Name,
			&token.Symbol,
			&token.SmartContractAddress,
			&token.Decimal,
			&token.LogoURL,
			&token.TokenABI,
		)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot scan a token row")
			return nil, err
		}

		token.Id = id
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Warn().Err(err).Msg("Error iterating over rows")
		return nil, err
	}
	return tokens, nil
}

func (s PostgresqlWalletRepository) AddToken(ctx context.Context, token domain.Token) (domain.Token, error) {

	query := `
		INSERT INTO crypto_tokens (id, name, symbol, decimal, smart_contract_address, logo_url, token_abi)
		VALUES ($1, $2, $3, $4, $5, $6,$7);
	`

	_, err := s.conn.Exec(
		ctx,
		query,
		uuid.New(),
		token.Name,
		token.Symbol,
		token.Decimal,
		token.SmartContractAddress,
		token.LogoURL,
		token.TokenABI,
	)

	if err != nil {
		log.Warn().Err(err).Msg("Error executing the query")
		return domain.Token{}, err
	}

	return domain.Token{}, nil

}
