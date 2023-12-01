package rest

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
)

type ProfileCtrl struct {
	ticket ticket.Service
	user   user.Service
}

type ProfileResponse struct {
	UserInfo      user.User                  `json:"userInfo"`
	OwnedTickets  []ticket.NFT               `json:"ownedTickets"`
	IssuedTickets []*ticket.TicketCollection `json:"issuedTickets"`
}

func NewProfileCtrl(ticket ticket.Service, user user.Service) *ProfileCtrl {
	return &ProfileCtrl{
		ticket: ticket,
		user:   user,
	}
}

func (c *ProfileCtrl) Pattern() string {
	return "/profile"
}

func (c *ProfileCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/{accountAddress}", c.profile)

	return r
}

// Profile godoc
//
//	@Tags			profile
//	@Summary		returns user profile
//	@Description	returns user profile
//	@Accept			json
//	@Produce		json
//	@Param 			accountAddress	path	string	true	"account address"
//	@Success		200			{object}	CommonResponse{data=ProfileResponse}
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/v1/profile/{accountAddress} [get]
func (c *ProfileCtrl) profile(w http.ResponseWriter, r *http.Request) {
	// 1. check params
	accountAddress := strings.ToLower(chi.URLParam(r, "accountAddress"))

	// 2. get user
	u, err := c.user.FindUserByAccountAddress(r.Context(), accountAddress)
	if err != nil {
		logger.Error("failed to find user", "error", err)
		if err == user.ErrUserNotFound {
			ErrorJSON(w, "user not found", http.StatusBadRequest)
		} else {
			ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		}
		return
	}

	// TODO: get user profile
	// 3. get purchased tickets and issued tickets
	tickets, err := c.ticket.GetOwnedNFT(r.Context(), web3.HexToAddress(u.AccountAddress))
	if err != nil {
		logger.Error("failed to get owned nft", "error", err)
		ErrorJSON(w, "failed to get owned nft", http.StatusInternalServerError)
		return
	}
	//  4. get issued ticket by user
	ticketCollections, err := c.ticket.FindTicketCollections(r.Context(), ticket.TicketCollectionFilter{
		IssuerAddress: u.AccountAddress,
	})
	if err != nil {
		logger.Error("failed to get ticket collections", "error", err)
		ErrorJSON(w, "failed to get ticket collections", http.StatusInternalServerError)
		return
	}

	// 5. return user profile
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully get user profile",
		Data:    ProfileResponse{UserInfo: *u, OwnedTickets: tickets.NFTs, IssuedTickets: ticketCollections},
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}
