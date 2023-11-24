package claim

type Claim struct {
	Identifier string `json:"identifier" bson:"_id"`
	UserID     string `json:"userId" bson:"userId"`
	Nonce      int64  `json:"nonce" bson:"nonce"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
	UpdateAt   int64  `json:"updatedAt" bson:"updatedAt"`
}
