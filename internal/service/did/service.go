package did

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/heroticket/internal/cache"
)

type Service interface {
	CreateIdentity(ctx context.Context, identity CreateIdentityRequest) (*CreateIdentityResponse, error)
	CreateClaim(ctx context.Context, identifier string, claim CreateClaimRequest) (*CreateClaimResponse, error)
	GetClaimQrCode(ctx context.Context, identifier string, claimId string) (*GetClaimQrCodeResponse, error)
}

type DidServiceConfig struct {
	RPCUrl    string
	IssuerUrl string
	Username  string
	Password  string
	QrCache   cache.Cache
	Client    *http.Client
}

type DidService struct {
	rpcUrl    string
	issuerUrl string
	username  string
	password  string

	qrCache cache.Cache
	client  *http.Client
}

func New(cfg DidServiceConfig) Service {
	svc := &DidService{
		rpcUrl:    cfg.RPCUrl,
		issuerUrl: cfg.IssuerUrl,
		username:  cfg.Username,
		password:  cfg.Password,
		qrCache:   cfg.QrCache,
		client:    http.DefaultClient,
	}

	if cfg.Client != nil {
		svc.client = cfg.Client
	}

	return svc
}

func (s *DidService) CreateIdentity(ctx context.Context, identity CreateIdentityRequest) (*CreateIdentityResponse, error) {
	url := fmt.Sprintf("%s/v1/identities", s.issuerUrl)

	body, err := json.Marshal(identity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	s.setAuthorizationHeader(req)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, errorFromResponse(res)
	}

	var createIdentityResponse CreateIdentityResponse

	err = json.NewDecoder(res.Body).Decode(&createIdentityResponse)
	if err != nil {
		return nil, err
	}

	if createIdentityResponse.Identifier == "" {
		return nil, fmt.Errorf("identity id is empty")
	}

	return &createIdentityResponse, nil
}

func (s *DidService) CreateClaim(ctx context.Context, identifier string, claim CreateClaimRequest) (*CreateClaimResponse, error) {
	url := fmt.Sprintf("%s/v1/%s/claims", identifier, s.issuerUrl)

	body, err := json.Marshal(claim)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	s.setAuthorizationHeader(req)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, errorFromResponse(res)
	}

	var createClaimResponse CreateClaimResponse

	err = json.NewDecoder(res.Body).Decode(&createClaimResponse)
	if err != nil {
		return nil, err
	}

	return &createClaimResponse, nil
}

func (s *DidService) GetClaimQrCode(ctx context.Context, identifier string, claimId string) (*GetClaimQrCodeResponse, error) {
	// check if qrcode exists in cache
	var qrcode GetClaimQrCodeResponse

	err := s.qrCache.Get(ctx, claimId, &qrcode)
	if err == nil {
		return &qrcode, nil
	}

	url := fmt.Sprintf("%s/v1/%s/claims/%s/qrcode", s.issuerUrl, identifier, claimId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	s.setAuthorizationHeader(req)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res)
	}

	var getClaimQrCodeResponse GetClaimQrCodeResponse

	err = json.NewDecoder(res.Body).Decode(&getClaimQrCodeResponse)
	if err != nil {
		return nil, err
	}

	// save qrcode to cache
	err = s.qrCache.Set(ctx, claimId, getClaimQrCodeResponse, time.Hour)
	if err != nil {
		return nil, err
	}

	return &getClaimQrCodeResponse, nil
}

func (s *DidService) setAuthorizationHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.username, s.password)))))
}

func errorFromResponse(res *http.Response) error {
	var data map[string]interface{}

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	if msg, ok := data["message"].(string); ok {
		return fmt.Errorf(msg)
	}

	return fmt.Errorf("unknown error")
}
