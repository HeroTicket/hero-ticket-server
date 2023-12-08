package rest

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
)

type ClaimCtrl struct {
	did    did.Service
	jwt    jwt.Service
	ticket ticket.Service
	user   user.Service
}

func NewClaimCtrl(did did.Service, jwt jwt.Service, ticket ticket.Service, user user.Service) *ClaimCtrl {
	return &ClaimCtrl{
		did:    did,
		jwt:    jwt,
		ticket: ticket,
		user:   user,
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
// @Success			202	{object}	CommonResponse
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
	rawContractAddress := strings.ToLower(chi.URLParam(r, "contractAddress"))

	// 3. find user by id
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user", "error", err)
		ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		return
	}

	// 4. check if user has claim
	_, err = c.did.FindClaim(r.Context(), u.ID, rawContractAddress)
	if err != nil {
		if err != did.ErrClaimNotFound {
			logger.Error("failed to find claim", "error", err)
			ErrorJSON(w, "failed to find claim", http.StatusInternalServerError)
			return
		}
	} else {
		resp := CommonResponse{
			Status: http.StatusAccepted,
			Data:   "claim already exists",
		}

		_ = WriteJSON(w, http.StatusAccepted, resp)
		return
	}

	// 5. check if user's tba has ticket with contract address
	contractAddress := web3.HexToAddress(rawContractAddress)
	tbaAddress := web3.HexToAddress(u.TbaAddress)

	ok, err := c.ticket.HasTicket(r.Context(), contractAddress, tbaAddress)
	if err != nil {
		logger.Error("failed to check if user has ticket", "error", err)
		ErrorJSON(w, "failed to check if user has ticket", http.StatusInternalServerError)
		return
	}

	if !ok {
		ErrorJSON(w, "user does not have ticket", http.StatusBadRequest)
		return
	}

	// 6. get admin id
	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// 7. request claim
	claimResp, err := c.did.CreateClaim(r.Context(), admin.ID, did.CreateClaimRequest{
		CredentialSchema: "ipfs://QmeoSVXtH3sjpD5ccsRnteaBn8ft1wVC7uRD6uGX6pzbKR",
		CredentialSubject: map[string]interface{}{
			"id":             u.ID,
			"dapp_name":      "Hero Ticket",
			"ticket_address": contractAddress,
		},
		Type: "Ownership",
	})
	if err != nil {
		logger.Error("failed to create claim", "error", err)
		ErrorJSON(w, "failed to create claim", http.StatusInternalServerError)
		return
	}

	// 8. save claim in db
	claim, err := c.did.SaveClaim(r.Context(), did.SaveClaimParams{
		ID:              claimResp.ID,
		UserID:          jwtUser.ID,
		ContractAddress: rawContractAddress,
	})
	if err != nil {
		logger.Error("failed to save claim", "error", err)
		ErrorJSON(w, "failed to save claim", http.StatusInternalServerError)
		return
	}

	// 9. return success response
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
	rawContractAddress := strings.ToLower(chi.URLParam(r, "contractAddress"))

	// 3. find claim from db by user id and contract address
	claim, err := c.did.FindClaim(r.Context(), jwtUser.ID, rawContractAddress)
	if err != nil {
		if err == did.ErrClaimNotFound {
			ErrorJSON(w, "claim not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to find claim", "error", err)
		ErrorJSON(w, "failed to find claim", http.StatusInternalServerError)
		return
	}

	// 4. get admin id
	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// 5. request qr code from did service
	qrResp, err := c.did.GetClaimQrCode(r.Context(), admin.ID, claim.ID)
	if err != nil {
		logger.Error("failed to get claim qr code", "error", err)
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
