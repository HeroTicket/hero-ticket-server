package auth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/heroticket/internal/cache"
	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	"github.com/iden3/iden3comm/v2/protocol"
)

type Service interface {
	AuthorizationRequest(ctx context.Context, params AuthorizationRequestParams) (protocol.AuthorizationRequestMessage, error)
	AuthorizationCallback(ctx context.Context, id, token string, deleteOnSuccess bool) (*protocol.AuthorizationResponseMessage, error)
}

type AuthServiceConfig struct {
	IPFSUrl         string
	RPCUrl          string
	ContractAddress string
	ResolverPrefix  string
	KeyDir          string
	ReqCache        cache.Cache
}

type AuthService struct {
	ipfsUrl         string
	rpcUrl          string
	contractAddress string
	resolverPrefix  string
	keyDir          string
	reqCache        cache.Cache
}

func New(config AuthServiceConfig) Service {
	return &AuthService{
		ipfsUrl:         config.IPFSUrl,
		rpcUrl:          config.RPCUrl,
		contractAddress: config.ContractAddress,
		resolverPrefix:  config.ResolverPrefix,
		keyDir:          config.KeyDir,
		reqCache:        config.ReqCache,
	}
}

func (s *AuthService) AuthorizationRequest(ctx context.Context, params AuthorizationRequestParams) (protocol.AuthorizationRequestMessage, error) {
	var req protocol.AuthorizationRequestMessage = auth.CreateAuthorizationRequestWithMessage(
		params.Reason,
		params.Message,
		params.Sender,
		params.CallbackUrl,
	)

	req.ID = params.ID
	req.ThreadID = params.ID

	if len(params.Scope) > 0 {
		req.Body.Scope = append(req.Body.Scope, params.Scope...)
	}

	timeout := DefaultTimeout

	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	err := s.reqCache.Set(ctx, params.ID, req, timeout)
	if err != nil {
		return protocol.AuthorizationRequestMessage{}, err
	}

	return req, nil
}

func (s *AuthService) AuthorizationCallback(ctx context.Context, id, token string, deleteOnSuccess bool) (*protocol.AuthorizationResponseMessage, error) {
	var request protocol.AuthorizationRequestMessage

	err := s.reqCache.Get(ctx, id, &request)
	if err != nil {
		return nil, err
	}

	var verificationKeyLoader = &loaders.FSKeyLoader{
		Dir: s.keyDir,
	}

	resolver := state.ETHResolver{
		RPCUrl:          s.rpcUrl,
		ContractAddress: common.HexToAddress(s.contractAddress),
	}

	resolvers := map[string]pubsignals.StateResolver{
		s.resolverPrefix: &resolver,
	}

	verifier, err := auth.NewVerifier(
		verificationKeyLoader,
		resolvers,
		auth.WithIPFSGateway(s.ipfsUrl),
	)
	if err != nil {
		return nil, err
	}

	response, err := verifier.FullVerify(
		ctx,
		token,
		request,
		pubsignals.WithAcceptedProofGenerationDelay(time.Minute*5),
	)
	if err != nil {
		return nil, err
	}

	if deleteOnSuccess {
		go func() {
			_ = s.reqCache.Delete(ctx, id)
		}()
	}

	return response, nil
}
