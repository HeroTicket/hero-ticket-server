package ticket

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
)

type TicketCollection struct {
	ID          string `json:"id" bson:"_id"`
	CreatorID   string `json:"creatorId" bson:"creatorId"`
	Address     string `json:"address" bson:"address"`
	Name        string `json:"name" bson:"name"`
	Symbol      string `json:"symbol" bson:"symbol"`
	Description string `json:"description" bson:"description"`
	Organizer   string `json:"organizer" bson:"organizer"`
	Location    string `json:"location" bson:"location"`
	Date        string `json:"date" bson:"date"`
	BannerImage string `json:"bannerImage" bson:"bannerImage"`
	TicketImage string `json:"ticketImage" bson:"ticketImage"`
	Price       int64  `json:"price" bson:"price"`
	TotalSupply int64  `json:"totalSupply" bson:"totalSupply"`
	Remaining   int64  `json:"remaining" bson:"remaining"`
	CreatedAt   int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt" bson:"updatedAt"`
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

type SaveTicketParams struct {
	Address      string `json:"address" bson:"address"`
	OwnerAddress string `json:"ownerAddress" bson:"ownerAddress"`
	TokenID      uint64 `json:"tokenId" bson:"tokenId"`
	BlockNumber  uint64 `json:"blockNumber" bson:"blockNumber"`
	PurchasedAt  int64  `json:"purchasedAt" bson:"purchasedAt"`
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
