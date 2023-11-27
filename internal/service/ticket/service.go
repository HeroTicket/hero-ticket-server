package ticket

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/heroticket/pkg/contracts/heroticket"
)

type Service interface {
	UpdateWhitelist(ctx context.Context, contractAddress, to string) error
	// CreateTBA
	// CreateTicketCollection
}

type TicketService struct {
	client *ethclient.Client
	hero   *heroticket.Heroticket
	pvk    *ecdsa.PrivateKey
	repo   Repository
}

func New(client *ethclient.Client, hero *heroticket.Heroticket, pvk *ecdsa.PrivateKey, repo Repository) *TicketService {
	return &TicketService{
		client: client,
		hero:   hero,
		pvk:    pvk,
		repo:   repo,
	}
}

func (s *TicketService) txOpts(ctx context.Context) (*bind.TransactOpts, error) {
	address := crypto.PubkeyToAddress(s.pvk.PublicKey)

	nonce, err := s.client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := s.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.pvk, chainID)
	if err != nil {
		return nil, err
	}

	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 3000000

	return auth, nil
}
