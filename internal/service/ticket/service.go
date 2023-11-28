package ticket

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/heroticket/pkg/contracts/heroticket"
)

type Service interface {
	UpdateWhitelist(ctx context.Context, contractAddress, to string) (*types.Transaction, error)
	CreateTBA(ctx context.Context, to common.Address, tokenURI string) (*types.Transaction, error)
	IssueTicket(ctx context.Context, tokenContractAddress common.Address, ticketName, ticketSymbol, ticketURI string, initialOwner common.Address, ticketAmount, ticketPrice int64) (*types.Transaction, error)
	BuyTicket(ctx context.Context, contractAddress common.Address) (*types.Transaction, error)
	BuyTicketByEther(ctx context.Context, TicketContractAddress, adminAddress common.Address, ticketPrice int64) (*types.Transaction, error)
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

func (s *TicketService) UpdateWhitelist(ctx context.Context, contractAddress, to common.Address, ticketPrice int64) (*types.Transaction, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	return s.hero.UpdateWhiteList(auth, contractAddress, to)
}

func (s *TicketService) CreateTBA(ctx context.Context, to common.Address, tokenURI string) (*types.Transaction, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	return s.hero.Mint(auth, to, tokenURI)
}

func (s *TicketService) IssueTicket(ctx context.Context, tokenContractAddress common.Address, ticketName, ticketSymbol, ticketURI string, initialOwner common.Address, ticketAmount, ticketPrice int64) (*types.Transaction, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	return s.hero.IssueTicket(auth, tokenContractAddress, ticketName, ticketSymbol, ticketURI, initialOwner, big.NewInt(ticketAmount), big.NewInt(ticketPrice))
}

func (s *TicketService) BuyTicket(ctx context.Context, contractAddress common.Address) (*types.Transaction, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	return s.hero.BuyTicket(auth, contractAddress)
}

func (s *TicketService) BuyTicketByEther(ctx context.Context, TicketContractAddress, adminAddress common.Address, ticketPrice int64) (*types.Transaction, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}
	return s.hero.BuyTicketByEther(auth, TicketContractAddress, adminAddress, big.NewInt(ticketPrice))
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
