package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/user"
)

type ClaimCtrl struct {
	did  did.Service
	jwt  jwt.Service
	user user.Service
	// ticket ticket.Service
}

func NewClaimCtrl(did did.Service, jwt jwt.Service, user user.Service) *ClaimCtrl {
	return &ClaimCtrl{
		did:  did,
		jwt:  jwt,
		user: user,
	}
}

func (c *ClaimCtrl) Pattern() string {
	return "/claims"
}

func (c *ClaimCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(TokenRequired(c.jwt))
	r.Post("/{contractAddress}", c.requestClaim)
	r.Get("/{contractAddress}", c.claimQR)

	return r
}

// RequestClaim godoc
// @Tags			claims
// @Summary			requests claim
// @Description		requests claim
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Success			201	{object}	CommonResponse
// @Failure			400	{object}	CommonResponse
// @Failure			401	{object}	CommonResponse
// @Failure			404	{object}	CommonResponse
// @Failure			500	{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/claims/{contractAddress}	[post]
func (c *ClaimCtrl) requestClaim(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get contract address from path

	// 3. check if user has claim

	// 4. if not, create claim

	// 5. save claim in db

	// 6. return success response
}

// ClaimQR godoc
// @Tags			claims
// @Summary			returns claim qr
// @Description		returns claim qr
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Success			200	{object}	CommonResponse
// @Failure			400	{object}	CommonResponse
// @Failure			401	{object}	CommonResponse
// @Failure			404	{object}	CommonResponse
// @Failure			500	{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/claims/{contractAddress}	[get]
func (c *ClaimCtrl) claimQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get contract address from path

	// 3. get claim from db

	// 4. if not, return error

	// 5. request qr code from did service

	// 6. return qr code
}
