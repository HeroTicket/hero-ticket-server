package did

import (
	"errors"
	"time"
)

var (
	ErrRequestNotFound = errors.New("request not found")
	ErrClaimNotFound   = errors.New("claim not found")
)

var DefaultCacheExpiry = 1 * time.Hour

type CreateIdentityRequest struct {
	DidMetadata DidMetadata `json:"didMetadata"`
}

type DidMetadata struct {
	Blockchain string                               `json:"blockchain"`
	Method     string                               `json:"method"`
	Network    string                               `json:"network"`
	Type       CreateIdentityRequestDidMetadataType `json:"type"`
}

type CreateIdentityRequestDidMetadataType string

const (
	BJJ CreateIdentityRequestDidMetadataType = "BJJ"
	ETH CreateIdentityRequestDidMetadataType = "ETH"
)

type CreateIdentityResponse struct {
	Address    string `json:"address"`
	Identifier string `json:"identifier"`
}

type CreateClaimRequest struct {
	CredentialSchema      string                 `json:"credentialSchema"`
	CredentialSubject     map[string]interface{} `json:"credentialSubject"`
	Expiration            *int64                 `json:"expiration,omitempty"`
	MerklizedRootPosition *string                `json:"merklizedRootPosition,omitempty"`
	RevNonce              *uint64                `json:"revNonce,omitempty"`
	SubjectPosition       *string                `json:"subjectPosition,omitempty"`
	Type                  string                 `json:"type"`
	Version               *uint32                `json:"version,omitempty"`
}

type CreateClaimResponse struct {
	ID string `json:"id"`
}

type GetClaimQrCodeResponse struct {
	Body struct {
		Credentials []struct {
			Description string `json:"description"`
			Id          string `json:"id"`
		} `json:"credentials"`
		Url string `json:"url"`
	} `json:"body"`
	From string `json:"from"`
	Id   string `json:"id"`
	Thid string `json:"thid"`
	To   string `json:"to"`
	Typ  string `json:"typ"`
	Type string `json:"type"`
}

type SaveClaimParams struct {
	ID              string
	UserID          string
	ContractAddress string
}

type Claim struct {
	ID              string `json:"id" bson:"_id"`
	UserID          string `json:"userId" bson:"userId"`
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`
	CreatedAt       int64  `json:"createdAt" bson:"createdAt"`
	UpdateAt        int64  `json:"updatedAt" bson:"updatedAt"`
}
