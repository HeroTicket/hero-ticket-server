package did

import "github.com/iden3/iden3comm/v2/protocol"

type Service interface {
	Login(id, audience, callbackUrl string) (protocol.AuthorizationRequestMessage, error)
	LoginCallback(id, token string) (*protocol.AuthorizationResponseMessage, error)
	// TODO: add more methods
	// CreateCredential()
	// VerifyCredential()
	// VerifyCredentialCallback()
	// RevokeCredential()
}
