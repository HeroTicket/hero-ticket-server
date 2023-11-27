package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/user"
	"go.uber.org/zap"
)

type ClaimCtrl struct {
	adminID string

	did  did.Service
	jwt  jwt.Service
	user user.Service
	// ticket ticket.Service
}

func NewClaimCtrl(did did.Service, jwt jwt.Service, user user.Service, adminID string) *ClaimCtrl {
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
// @Success			201	{object}	CommonResponse{data=did.CreateClaimResponse}
// @Failure			400	{object}	CommonResponse
// @Failure			401	{object}	CommonResponse
// @Failure			404	{object}	CommonResponse
// @Failure			500	{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/claims/{contractAddress}	[post]
func (c *ClaimCtrl) requestClaim(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get contract address from path
	contractAddress := chi.URLParam(r, "contractAddress")

	// 3. check if user has claim
	_, err = c.did.FindClaim(r.Context(), jwtUser.ID, contractAddress)
	if err != nil {
		if err != did.ErrClaimNotFound {
			zap.L().Error("failed to find claim", zap.Error(err))
			ErrorJSON(w, "failed to find claim", http.StatusInternalServerError)
			return
		}
	} else {
		ErrorJSON(w, "claim already exists", http.StatusBadRequest)
		return
	}

	// TODO: check if user has ticket

	// TODO: create claim

	// 4. if not, create claim

	// 4.1 generate hash

	// 4.2 request claim
	claimResp, err := c.did.CreateClaim(r.Context(), c.adminID, did.CreateClaimRequest{
		CredentialSchema: "",
		CredentialSubject: map[string]interface{}{
			"id":     jwtUser.ID,
			"ticket": contractAddress,
			"hash":   "",
		},
		Type: "",
	})
	if err != nil {
		zap.L().Error("failed to create claim", zap.Error(err))
		ErrorJSON(w, "failed to create claim", http.StatusInternalServerError)
		return
	}

	// 5. save claim in db
	claim, err := c.did.SaveClaim(r.Context(), did.SaveClaimParams{
		ID:              claimResp.ID,
		UserID:          jwtUser.ID,
		ContractAddress: contractAddress,
	})
	if err != nil {
		zap.L().Error("failed to save claim", zap.Error(err))
		ErrorJSON(w, "failed to save claim", http.StatusInternalServerError)
		return
	}

	// 6.update ticket status

	// 7. return success response
	resp := CommonResponse{
		Status:  http.StatusCreated,
		Message: "Successfully requested claim",
		Data:    claim,
	}

	_ = WriteJSON(w, http.StatusCreated, resp)
}

// ClaimQR godoc
// @Tags			claims
// @Summary			returns claim qr
// @Description		returns claim qr
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Success			200	{object}	CommonResponse{did.GetClaimQrCodeResponse}
// @Failure			400	{object}	CommonResponse
// @Failure			401	{object}	CommonResponse
// @Failure			404	{object}	CommonResponse
// @Failure			500	{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/claims/{contractAddress}	[get]
func (c *ClaimCtrl) claimQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get contract address from path
	contractAddress := chi.URLParam(r, "contractAddress")

	// 3. get claim from db
	claim, err := c.did.FindClaim(r.Context(), jwtUser.ID, contractAddress)
	if err != nil {
		zap.L().Error("failed to find claim", zap.Error(err))
		ErrorJSON(w, "failed to find claim", http.StatusInternalServerError)
		return
	}

	// 4. request qr code from did service
	qrResp, err := c.did.GetClaimQrCode(r.Context(), c.adminID, claim.ID)
	if err != nil {
		zap.L().Error("failed to get claim qr code", zap.Error(err))
		ErrorJSON(w, "failed to get claim qr code", http.StatusInternalServerError)
		return
	}

	// 5. return qr code
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully retrieved claim qr code",
		Data:    qrResp,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}
