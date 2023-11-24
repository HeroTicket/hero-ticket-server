package did

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrRequestNotFound = errors.New("request not found")
)

var (
	DefaultClient      = http.DefaultClient
	DefaultCacheExpiry = 10 * time.Minute
)

type Membership struct {
	Identifier string    `json:"identifier" bson:"_id"`
	ClaimID    string    `json:"claimId,omitempty" bson:"claimId,omitempty"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TicketOwnership struct {
	ID         string    `json:"id" bson:"_id"`
	Identifier string    `json:"identifier" bson:"identifier"`
	ClaimID    string    `json:"claimId" bson:"claimId"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

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