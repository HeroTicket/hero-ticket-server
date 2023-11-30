package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/app/ws"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/ipfs"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
)

type TicketCtrl struct {
	serverUrl string

	auth   auth.Service
	ipfs   ipfs.Service
	jwt    jwt.Service
	user   user.Service
	ticket ticket.Service
}

func NewTicketCtrl(auth auth.Service, ipfs ipfs.Service, jwt jwt.Service, user user.Service, ticket ticket.Service, serverUrl string) *TicketCtrl {
	return &TicketCtrl{
		auth:      auth,
		ipfs:      ipfs,
		jwt:       jwt,
		user:      user,
		serverUrl: serverUrl,
	}
}

func (c *TicketCtrl) Pattern() string {
	return "/tickets"
}

func (c *TicketCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.tickets)
	r.With(TokenCheck(c.jwt)).Get("/{contractAddress}", c.ticketByContractAddress)
	r.Post("/{contractAddress}/whitelist-callback", c.whitelistCallback)
	r.Post("/{contractAddress}/direct-purchase-callback", c.directPurchaseCallback)
	r.Post("/verify-callback", c.verifyCallback)

	r.Group(func(r chi.Router) {
		r.Use(TokenRequired(c.jwt))
		r.Get("/{contractAddress}/whitelist-qr", c.whitelistQR)
		r.Get("/{contractAddress}/direct-purchase-qr", c.directPurchaseQR)
		r.Get("/{contractAddress}/verify-qr", c.verifyQR)
		r.Post("/create", c.createTicket)
	})

	return r
}

// Tickets godoc
//
// @Tags			tickets
// @Summary		returns tickets
// @Description	returns tickets
// @Accept			json
// @Produce		json
// @Param			page	query	int	false	"page number"
// @Param			limit	query	int	false	"page size"
// @Success		200			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets [get]
func (c *TicketCtrl) tickets(w http.ResponseWriter, r *http.Request) {
	// 1. get page and limit from query

	// 2. get tickets from db

	// 3. return tickets
}

// Ticket godoc
//
// @Tags			tickets
// @Summary		returns ticket
// @Description	returns ticket
// @Accept		json
// @Produce		json
// @Param			contractAddress	path	string	true	"contract address"
// @Success		200			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets/{contractAddress} [get]
func (c *TicketCtrl) ticketByContractAddress(w http.ResponseWriter, r *http.Request) {
	// // 1. get contract address from path
	// contractAddress := chi.URLParam(r, "contractAddress")

	// // 2. get ticket from db, if user is logged in, also get user ticket ownership

	// jwtUser, err := c.jwt.FromContext(r.Context())
	// if err != nil {
	// 	// user is not logged in
	// 	ticketCollection, err = c.ticket.TicketByContractAddress(r.Context(), contractAddress)
	// } else {
	// 	// user is logged in
	// 	ticketCollection, err = c.ticket.TicketByContractAddress(r.Context(), contractAddress, jwtUser.ID)
	// }

	// // 3. return ticket
	// resp := CommonResponse{
	// 	Status:  http.StatusOK,
	// 	Message: "",
	// 	Data:    ticketCollection,
	// }

	// _ = WriteJSON(w, http.StatusOK, resp)
}

// WhitelistQR godoc
//
// @Tags			tickets
// @Summary		returns purchase authorization qr code
// @Description	returns purchase authorization qr code
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			sessionId		query	string	true	"session id"
// @Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/{contractAddress}/whitelist-qr [get]
func (c *TicketCtrl) whitelistQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		logger.Error("failed to get jwt user from context", "error", err)
		ErrorJSON(w, "failed to get jwt user from context", http.StatusInternalServerError)
		return
	}

	// 2. get contract address from path
	rawContractAddress := chi.URLParam(r, "contractAddress")

	// 3. get session id from query
	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "whitelist-qr",
			Status: ws.InProgress,
		},
	})

	// 4. get user from db
	user, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user by id", "error", err)
		ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
		return
	}

	// 5. check if ticket collection exists
	contractAddress := common.HexToAddress(rawContractAddress)

	ok, err := c.ticket.IsIssuedTicket(r.Context(), contractAddress)
	if err != nil {
		logger.Error("failed to check if ticket collection exists", "error", err)
		ErrorJSON(w, "failed to check if ticket collection exists", http.StatusInternalServerError)
		return
	}

	if !ok {
		ErrorJSON(w, "ticket collection does not exist", http.StatusBadRequest)
		return
	}

	// 6. check if user has ticket
	tbaAddress := common.HexToAddress(user.TbaAddress)

	ok, err = c.ticket.HasTicket(r.Context(), contractAddress, tbaAddress)
	if err != nil {
		logger.Error("failed to check if user has ticket", "error", err)
		ErrorJSON(w, "failed to check if user has ticket", http.StatusInternalServerError)
		return
	}

	if ok {
		ErrorJSON(w, "user already has ticket", http.StatusBadRequest)
		return
	}

	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	callbackUrl := fmt.Sprintf("%s/v1/tickets/%s/whitelist-callback?sessionId=%s&accountAddress=%s", c.serverUrl, rawContractAddress, sessionId, user.AccountAddress)

	// 7. create qr code
	qrCode, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Update whitelist for purchase authorization",
		Message:     "Scan the QR code to update whitelist for purchase authorization",
		Sender:      admin.ID,
		CallbackUrl: callbackUrl,
	})
	if err != nil {
		logger.Error("failed to create authorization request", "error", err)
		ErrorJSON(w, "failed to create authorization request", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "whitelist-qr",
			Status: ws.Done,
		},
	})

	// 8. return qr code
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully created authorization request",
		Data:    qrCode,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}

// DirectPurchaseQR godoc
//
// @Tags			tickets
// @Summary		returns direct purchase authorization qr code
// @Description	returns direct purchase authorization qr code
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			sessionId		query	string	true	"session id"
// @Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/{contractAddress}/direct-purchase-qr [get]
func (c *TicketCtrl) directPurchaseQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get contract address from path

	// 3. get session id from query

	// 4. get user from db

	// 5. check if ticket collection exists

	// 6. check if user has ticket

	// 7. create qr code

	// 8. return qr code
}

// WhitelistCallback godoc
//
// @Tags			tickets
// @Summary			whitelist callback
// @Description		whitelist callback
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			accountAddress	query	string	true	"account address"
// @Param			sessionId		query	string	true	"session id"
// @Param			token			body	string	true	"token"
// @Success		200			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets/{contractAddress}/whitelist-callback [post]
func (c *TicketCtrl) whitelistCallback(w http.ResponseWriter, r *http.Request) {
	// 1. get contract address from path
	rawContractAddress := strings.ToLower(chi.URLParam(r, "contractAddress"))

	// 2. get account address from query
	rawAccountAddress := strings.ToLower(r.URL.Query().Get("accountAddress"))

	// 3. get session id from query
	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)
	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
		return
	}

	// 4. get token from body
	token, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("failed to read token from body", "error", err)
		ErrorJSON(w, "failed to read token from body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "whitelist-callback",
			Status: ws.InProgress,
		},
	})

	// 5. get user from db
	user, err := c.user.FindUserByAccountAddress(r.Context(), rawAccountAddress)
	if err != nil {
		logger.Error("failed to find user by account address", "error", err)
		ErrorJSON(w, "failed to find user by account address", http.StatusInternalServerError)
		return
	}

	// 6. check if ticket collection exists
	contractAddress := common.HexToAddress(rawContractAddress)

	ok, err := c.ticket.IsIssuedTicket(r.Context(), contractAddress)
	if err != nil {
		logger.Error("failed to check if ticket collection exists", "error", err)
		ErrorJSON(w, "failed to check if ticket collection exists", http.StatusInternalServerError)
		return
	}

	if !ok {
		ErrorJSON(w, "ticket collection does not exist", http.StatusBadRequest)
		return
	}

	// 7. check if user has ticket
	tbaAddress := common.HexToAddress(user.TbaAddress)

	ok, err = c.ticket.HasTicket(r.Context(), contractAddress, tbaAddress)
	if err != nil {
		logger.Error("failed to check if user has ticket", "error", err)
		ErrorJSON(w, "failed to check if user has ticket", http.StatusInternalServerError)
		return
	}

	if ok {
		ErrorJSON(w, "user already has ticket", http.StatusBadRequest)
		return
	}

	// 8. verify token
	resp, err := c.auth.AuthorizationCallback(r.Context(), sessionId, string(token))
	if err != nil {
		logger.Error("failed to handle whitelist callback", "error", err)
		ErrorJSON(w, "failed to handle whitelist callback", http.StatusInternalServerError)
		return
	}

	// 9. get user id from verification response
	userID := resp.From

	// 10. check if user id matches user id from db
	if userID != user.ID {
		ErrorJSON(w, "user id does not match", http.StatusBadRequest)
		return
	}

	// 11. call contract to set user address on whitelist
	accountAddress := common.HexToAddress(rawAccountAddress)

	err = c.ticket.UpdateWhitelist(r.Context(), contractAddress, accountAddress)
	if err != nil {
		logger.Error("failed to update whitelist", "error", err)
		ErrorJSON(w, "failed to update whitelist", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "whitelist-callback",
			Status: ws.Done,
			Data:   "Successfully updated whitelist",
		},
	})

	// 12. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Successfully updated whitelist for user with ID %s", userID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// DirectPurchaseCallback godoc
//
// @Tags			tickets
// @Summary			direct purchase callback
// @Description		direct purchase callback
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			accountAddress	query	string	true	"account address"
// @Param			sessionId		query	string	true	"session id"
// @Param			token			body	string	true	"token"
// @Success			200			{object}	CommonResponse
// @Failure			400			{object}	CommonResponse
// @Failure			500			{object}	CommonResponse
// @Router			/v1/tickets/direct-purchase-callback [post]
func (c *TicketCtrl) directPurchaseCallback(w http.ResponseWriter, r *http.Request) {

}

// VerifyQR godoc
//
// @Tags			tickets
// @Summary		returns verify authorization qr code
// @Description	returns verify authorization qr code
// @Accept			json
// @Produce		json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			sessionId		query	string	true	"session id"
// @Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/{contractAddress}/verify-qr [get]
func (c *TicketCtrl) verifyQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get contract address from path

	// 3. get session id from query

	// 4. get user from db

	// 5. check if ticket collection exists

	// 6. check if user is owner of ticket collection

	// 7. create qr code

	// 8. return qr code
}

// VerifyCallback godoc
//
// @Tags			tickets
// @Summary		verify callback
// @Description	verify callback
// @Accept			json
// @Produce		json
// @Param			sessionId		query	string	true	"session id"
// @Param			token			body	string	true	"token"
// @Success		200			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets/verify-callback [post]
func (c *TicketCtrl) verifyCallback(w http.ResponseWriter, r *http.Request) {
	// 1. get session id from query

	// 2. get token from body

	// 3. verify token

	// 4. return success response
}

// CreateTicket godoc
//
// @Tags			tickets
// @Summary		creates ticket
// @Description	creates ticket
// @Accept			json
// @Produce		json
// @Param			name		formData	string	true	"ticket name"
// @Param			symbol		formData	string	true	"ticket symbol"
// @Param			description	formData	string	true	"ticket description"
// @Param			organizer	formData	string	true	"ticket organizer"
// @Param			location	formData	string	true	"ticket location"
// @Param			date		formData	string	true	"ticket date"
// @Param			bannerImage	formData	file	true	"ticket banner image"
// @Param			ticketImage	formData	file	true	"ticket image"
// @Param			price		formData	int64	true	"ticket price"
// @Param			totalSupply	formData	int64	true	"ticket total supply"
// @Success		201			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/create [post]
func (c *TicketCtrl) createTicket(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get params from form data

	// 3. call contract to create new ticket collection

	// 4. save ticket collection to db

	// 5. return success response
}
