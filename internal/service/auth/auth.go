package auth

import (
	"time"

	"github.com/iden3/iden3comm/v2/protocol"
)

var DefaultTimeout = 10 * time.Minute

type AuthorizationRequestParams struct {
	ID          string
	Reason      string
	Message     string
	Sender      string
	CallbackUrl string
	Scope       []protocol.ZeroKnowledgeProofRequest
	Timeout     time.Duration
}
