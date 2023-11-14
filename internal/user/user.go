package user

type User struct {
	DID           string `json:"did" bson:"_id"`
	WalletAddress string `json:"wallet_address" bson:"wallet_address"`
	Name          string `json:"name" bson:"name"`
	IsAdmin       bool   `json:"is_admin" bson:"is_admin"`
	CreatedAt     int64  `json:"created_at" bson:"created_at"`
	UpdatedAt     int64  `json:"updated_at" bson:"updated_at"`
}
