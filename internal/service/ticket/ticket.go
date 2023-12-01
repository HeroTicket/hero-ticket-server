package ticket

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrTicketNotFound           = errors.New("ticket not found")
	ErrTicketCollectionNotFound = errors.New("ticket collection not found")
)

type TicketCollection struct {
	ID              string `json:"id" bson:"_id"`
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`
	IssuerAddress   string `json:"issuerAddress" bson:"issuerAddress"`
	Name            string `json:"name" bson:"name"`
	Symbol          string `json:"symbol" bson:"symbol"`
	Description     string `json:"description" bson:"description"`
	Organizer       string `json:"organizer" bson:"organizer"`
	Location        string `json:"location" bson:"location"`
	Date            string `json:"date" bson:"date"`
	BannerUrl       string `json:"bannerUrl" bson:"bannerUrl"`
	TicketUrl       string `json:"ticketUrl" bson:"ticketUrl"`
	EthPrice        string `json:"ethPrice" bson:"ethPrice"`
	TokenPrice      string `json:"tokenPrice" bson:"tokenPrice"`
	TotalSupply     string `json:"totalSupply" bson:"totalSupply"`
	Remaining       string `json:"remaining" bson:"remaining"`
	SaleStartAt     int64  `json:"saleStartAt" bson:"saleStartAt"`
	SaleEndAt       int64  `json:"saleEndAt" bson:"saleEndAt"`
	CreatedAt       int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt" bson:"updatedAt"`
}

type TicketCollectionDetail struct {
	TicketCollection
	UserHasTicket bool `json:"userHasTicket"`
}

type Ticket struct {
	ID           string `json:"id" bson:"_id"`
	Address      string `json:"address" bson:"address"`
	OwnerAddress string `json:"ownerAddress" bson:"ownerAddress"`
	TokenID      uint64 `json:"tokenId" bson:"tokenId"`
	Name         string `json:"name" bson:"name"`
	Symbol       string `json:"symbol" bson:"symbol"`
	Image        string `json:"image" bson:"image,omitempty"`
	PurchasedAt  int64  `json:"purchasedAt" bson:"purchasedAt"`
}

type OnchainTicketInfo struct {
	ContractAddress common.Address
	Issuer          common.Address
	Remaining       *big.Int
	EthPrice        *big.Int
	TokenPrice      *big.Int
	SaleStartAt     *big.Int
	SaleEndAt       *big.Int
}

type CreateTicketCollectionParams struct {
	ContractAddress string
	IssuerAddress   string
	Name            string
	Symbol          string
	Description     string
	Organizer       string
	Location        string
	Date            string
	BannerUrl       string
	TicketUrl       string
	EthPrice        string
	TokenPrice      string
	TotalSupply     string
	Remaining       string
	SaleStartAt     int64
	SaleEndAt       int64
}

type SaveTicketParams struct {
	Address      string
	OwnerAddress string
	TokenID      uint64
	BlockNumber  uint64
	PurchasedAt  int64
}

type IssueTicketParams struct {
	TicketName       string
	TicketSymbol     string
	TicketUri        string
	Issuer           common.Address
	TicketAmount     *big.Int
	TicketEthPrice   *big.Int
	TicketTokenPrice *big.Int
	SaleDuration     *big.Int
}

type TicketCollectionFilter struct {
	// TODO: add filter options
	Issuer string
}
