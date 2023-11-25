package ticket

type TicketCollection struct {
	ID              string `json:"id" bson:"_id"`
	CreatorID       string `json:"creatorId" bson:"creatorId"`
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`
	Name            string `json:"name" bson:"name"`
	Symbol          string `json:"symbol" bson:"symbol"`
	Description     string `json:"description" bson:"description"`
	Organizer       string `json:"organizer" bson:"organizer"`
	Location        string `json:"location" bson:"location"`
	Date            string `json:"date" bson:"date"`
	BannerImage     string `json:"bannerImage" bson:"bannerImage"`
	TicketImage     string `json:"ticketImage" bson:"ticketImage"`
	Price           int64  `json:"price" bson:"price"`
	TotalSupply     int64  `json:"totalSupply" bson:"totalSupply"`
	Remaining       int64  `json:"remaining" bson:"remaining"`
	CreatedAt       int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt" bson:"updatedAt"`
}

type Ticket struct {
	ID                string `json:"id" bson:"_id"`
	Name              string `json:"name" bson:"name"`
	TokenID           string `json:"tokenId" bson:"tokenId"`
	CollectionAddress string `json:"collectionAddress" bson:"collectionAddress"`
	OwnerAddress      string `json:"ownerAddress" bson:"ownerAddress"`
	Image             string `json:"image" bson:"image,omitempty"`
	PurchasedAt       int64  `json:"purchasedAt" bson:"purchasedAt"`
}
