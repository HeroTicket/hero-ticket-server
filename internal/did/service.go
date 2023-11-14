package did

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	"github.com/iden3/iden3comm/v2/protocol"
)

type Service interface {
	LoginRequest(id, audience, callbackUrl string) (protocol.AuthorizationRequestMessage, error)
	LoginCallback(ctx context.Context, id, token string) (*protocol.AuthorizationResponseMessage, error)
	// TODO: add more methods
	// CreateCredential()
	// VerifyCredential()
	// VerifyCredentialCallback()
	// RevokeCredential()
}

type didService struct {
	rpcUrl string

	loginRequestMap map[string]interface{}

	mu sync.RWMutex
}

func New(rpcUrl string) Service {
	svc := &didService{
		rpcUrl:          rpcUrl,
		loginRequestMap: make(map[string]interface{}),
		mu:              sync.RWMutex{},
	}

	return svc
}

func (s *didService) LoginRequest(id string, audience string, callbackUrl string) (protocol.AuthorizationRequestMessage, error) {
	var request protocol.AuthorizationRequestMessage = auth.CreateAuthorizationRequestWithMessage(
		"Login to Hero Ticket",
		"Scan the QR code to login to Hero Ticket",
		audience,
		callbackUrl,
	)

	request.ID = id
	request.ThreadID = id

	s.mu.Lock()
	s.loginRequestMap[id] = request
	s.mu.Unlock()

	return request, nil
}

func (s *didService) LoginCallback(ctx context.Context, id string, token string) (*protocol.AuthorizationResponseMessage, error) {
	s.mu.RLock()
	request, ok := s.loginRequestMap[id]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrRequestNotFound
	}

	ipfsUrl := "https://ipfs.io"
	contractAddress := "134B1BE34911E39A8397ec6289782989729807a4"
	resolverPrefix := "polygon:mumbai"
	ketDir := "./keys"

	var verificationKeyLoader = &loaders.FSKeyLoader{
		Dir: ketDir,
	}

	resolver := state.ETHResolver{
		RPCUrl:          s.rpcUrl,
		ContractAddress: common.HexToAddress(contractAddress),
	}

	resolvers := map[string]pubsignals.StateResolver{
		resolverPrefix: &resolver,
	}

	verifier, err := auth.NewVerifier(
		verificationKeyLoader,
		resolvers,
		auth.WithIPFSGateway(ipfsUrl),
	)
	if err != nil {
		return nil, err
	}

	response, err := verifier.FullVerify(
		ctx,
		token,
		request.(protocol.AuthorizationRequestMessage),
		pubsignals.WithAcceptedProofGenerationDelay(time.Minute*5),
	)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	delete(s.loginRequestMap, id)
	s.mu.Unlock()

	return response, nil
}
