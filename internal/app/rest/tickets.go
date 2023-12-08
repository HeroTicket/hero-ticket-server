package rest

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/heroticket/internal/app/ws"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/ipfs"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
	"github.com/iden3/go-circuits/v2"
	"github.com/iden3/iden3comm/v2/protocol"
)

type TicketCtrl struct {
	serverUrl string

	auth   auth.Service
	ipfs   ipfs.Service
	jwt    jwt.Service
	ticket ticket.Service
	user   user.Service
}

func NewTicketCtrl(auth auth.Service, ipfs ipfs.Service, jwt jwt.Service, ticket ticket.Service, user user.Service, serverUrl string) *TicketCtrl {
	return &TicketCtrl{
		auth:      auth,
		ipfs:      ipfs,
		jwt:       jwt,
		ticket:    ticket,
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
	r.Post("/{contractAddress}/token-purchase-callback", c.tokenPurchaseCallback)
	r.Post("/verify-callback", c.verifyCallback)

	r.Group(func(r chi.Router) {
		r.Use(TokenRequired(c.jwt))
		r.Get("/{contractAddress}/whitelist-qr", c.whitelistQR)
		r.Get("/{contractAddress}/token-purchase-qr", c.tokenPurchaseQR)
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
// @Success		200			{object}	CommonResponse{data=[]ticket.TicketCollection}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets [get]
func (c *TicketCtrl) tickets(w http.ResponseWriter, r *http.Request) {
	// 1. get page and limit from query

	// TODO: add pagination

	// 2. get ticket collections from db
	collections, err := c.ticket.FindTicketCollections(r.Context(), ticket.TicketCollectionFilter{})
	if err != nil {
		logger.Error("failed to find ticket collections", "error", err)
		ErrorJSON(w, "failed to find ticket collections", http.StatusInternalServerError)
		return
	}

	// 3. return tickets
	resp := CommonResponse{
		Status: http.StatusOK,
	}

	if len(collections) > 0 {
		resp.Message = "Successfully retrieved ticket collections"
		resp.Data = collections
	} else {
		resp.Message = "No ticket collections found"
		resp.Data = []ticket.TicketCollection{}
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}

// Ticket godoc
//
// @Tags			tickets
// @Summary		returns ticket
// @Description	returns ticket
// @Accept		json
// @Produce		json
// @Param			contractAddress	path	string	true	"contract address"
// @Success		200			{object}	CommonResponse{data=ticket.TicketCollectionDetail}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets/{contractAddress} [get]
func (c *TicketCtrl) ticketByContractAddress(w http.ResponseWriter, r *http.Request) {
	// 1. get contract address from path
	rawContractAddress := strings.ToLower(chi.URLParam(r, "contractAddress"))

	// 2. get ticket from db, if user is logged in, also get user ticket ownership

	var userHasTicket bool

	jwtUser, err := c.jwt.FromContext(r.Context())
	if err == nil {
		u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
		if err != nil {
			logger.Error("failed to find user by id", "error", err)
			ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
			return
		}

		contractAddress := web3.HexToAddress(rawContractAddress)

		tbaAddress := web3.HexToAddress(u.TbaAddress)

		// user is not logged in
		ok, err := c.ticket.HasTicket(r.Context(), contractAddress, tbaAddress)
		if err != nil {
			logger.Error("failed to check if user has ticket", "error", err)
			ErrorJSON(w, "failed to check if user has ticket", http.StatusInternalServerError)
			return
		}

		userHasTicket = ok
	}

	// 3. find ticket collection by contract address
	collection, err := c.ticket.FindTicketCollectionByContractAddress(r.Context(), rawContractAddress)
	if err != nil {
		if err == ticket.ErrTicketCollectionNotFound {
			ErrorJSON(w, "ticket collection not found", http.StatusBadRequest)
			return
		}

		logger.Error("failed to find ticket collection by contract address", "error", err)
		ErrorJSON(w, "failed to find ticket collection by contract address", http.StatusInternalServerError)
		return
	}

	var detail ticket.TicketCollectionDetail

	// 4. get onchain ticket collection info
	detail.TicketCollection = *collection
	detail.UserHasTicket = userHasTicket

	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully retrieved ticket collection",
		Data:    detail,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
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
// @Success		202			{object}	CommonResponse
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

	// 4. check if ticket collection exists
	contractAddress := web3.HexToAddress(rawContractAddress)

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

	// 5. get onchain ticket collection info
	onchainTicket, err := c.ticket.OnChainTicketInfo(r.Context(), contractAddress)
	if err != nil {
		logger.Error("failed to get onchain ticket collection", "error", err)
		ErrorJSON(w, "failed to get onchain ticket collection", http.StatusInternalServerError)
		return
	}

	// 6. check if ticket is on sale
	if onchainTicket.Remaining.Cmp(big.NewInt(0)) == 0 || onchainTicket.SaleEndAt.Int64() < time.Now().Unix() {
		ErrorJSON(w, "ticket is not on sale", http.StatusBadRequest)
		return
	}

	// 7. get user from db
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user by id", "error", err)
		ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
		return
	}

	// 8. check if user has ticket
	tbaAddress := web3.HexToAddress(u.TbaAddress)

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

	// 9. check if tba address is already on whitelist
	ok, err = c.ticket.IsWhitelisted(r.Context(), contractAddress, tbaAddress)
	if err != nil {
		logger.Error("failed to check if user is already on whitelist", "error", err)
		ErrorJSON(w, "failed to check if user is already on whitelist", http.StatusInternalServerError)
		return
	}

	if ok {
		go ws.Send(ws.Message{
			ID:   id,
			Type: ws.EventMessage,
			Event: ws.Event{
				Name:   "whitelist-qr",
				Status: ws.Done,
				Data:   "User is already on whitelist",
			},
		})

		resp := CommonResponse{
			Status:  http.StatusAccepted,
			Message: "User is already on whitelist",
		}

		_ = WriteJSON(w, http.StatusAccepted, resp)
		return
	}

	// 10. get admin from db
	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	callbackUrl := fmt.Sprintf("%s/v1/tickets/%s/whitelist-callback?sessionId=%s&accountAddress=%s", c.serverUrl, rawContractAddress, sessionId, u.AccountAddress)

	// 11. create qr code
	qrCode, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Update whitelist for purchase authentication",
		Message:     "Scan the QR code to update whitelist for purchase authentication",
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

	// 12. return qr code
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully created authorization request",
		Data:    qrCode,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}

// TokenPurchaseCallback godoc
//
// @Tags			tickets
// @Summary		returns token purchase authorization qr code
// @Description	returns token purchase authorization qr code
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			sessionId		query	string	true	"session id"
// @Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/{contractAddress}/token-purchase-qr [get]
func (c *TicketCtrl) tokenPurchaseQR(w http.ResponseWriter, r *http.Request) {
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
			Name:   "token-purchase-qr",
			Status: ws.InProgress,
		},
	})

	// 4. check if ticket collection exists
	contractAddress := web3.HexToAddress(rawContractAddress)

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

	// 5. get onchain ticket collection info
	onchainTicket, err := c.ticket.OnChainTicketInfo(r.Context(), contractAddress)
	if err != nil {
		logger.Error("failed to get onchain ticket collection", "error", err)
		ErrorJSON(w, "failed to get onchain ticket collection", http.StatusInternalServerError)
		return
	}

	// 6. check if ticket is on sale
	if onchainTicket.Remaining.Cmp(big.NewInt(0)) == 0 || onchainTicket.SaleEndAt.Int64() < time.Now().Unix() {
		ErrorJSON(w, "ticket is not on sale", http.StatusBadRequest)
		return
	}

	// 7. get user from db
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user by id", "error", err)
		ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
		return
	}

	// 8. check if user has ticket
	tbaAddress := web3.HexToAddress(u.TbaAddress)

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

	callbackUrl := fmt.Sprintf("%s/v1/tickets/%s/token-purchase-callback?sessionId=%s&accountAddress=%s", c.serverUrl, rawContractAddress, sessionId, u.AccountAddress)

	// 9. create qr code
	qrCode, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Ticket purchase authorization",
		Message:     fmt.Sprintf("Scan the QR code to authenticate ticket purchase for %s", rawContractAddress),
		Sender:      admin.ID,
		CallbackUrl: callbackUrl,
	})
	if err != nil {
		logger.Error("failed to create authorization request", "error", err)
		ErrorJSON(w, "failed to create authorization request", http.StatusInternalServerError)
		return
	}

	// 10. return qr code
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully created authorization request",
		Data:    qrCode,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
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
	tokenBytes, err := io.ReadAll(r.Body)
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
	contractAddress := web3.HexToAddress(rawContractAddress)

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
	tbaAddress := web3.HexToAddress(user.TbaAddress)

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
	resp, err := c.auth.AuthorizationCallback(r.Context(), sessionId, string(tokenBytes))
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
	accountAddress := web3.HexToAddress(rawAccountAddress)

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

// TokenPurchaseCallback godoc
//
// @Tags			tickets
// @Summary			token purchase callback
// @Description		token purchase callback
// @Accept			json
// @Produce			json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			accountAddress	query	string	true	"account address"
// @Param			sessionId		query	string	true	"session id"
// @Param			token			body	string	true	"token"
// @Success			200			{object}	CommonResponse
// @Failure			400			{object}	CommonResponse
// @Failure			500			{object}	CommonResponse
// @Router			/v1/tickets/{contractAddress}/token-purchase-callback [post]
func (c *TicketCtrl) tokenPurchaseCallback(w http.ResponseWriter, r *http.Request) {
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
	tokenBytes, err := io.ReadAll(r.Body)
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
			Name:   "token-purchase-callback",
			Status: ws.InProgress,
		},
	})

	// 5. get user from db
	u, err := c.user.FindUserByAccountAddress(r.Context(), rawAccountAddress)
	if err != nil {
		logger.Error("failed to find user by account address", "error", err)
		ErrorJSON(w, "failed to find user by account address", http.StatusInternalServerError)
		return
	}

	// 6. check if ticket collection exists
	contractAddress := web3.HexToAddress(rawContractAddress)

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
	tbaAddress := web3.HexToAddress(u.TbaAddress)

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
	resp, err := c.auth.AuthorizationCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		logger.Error("failed to handle token purchase callback", "error", err)
		ErrorJSON(w, "failed to handle token purchase callback", http.StatusInternalServerError)
		return
	}

	// 9. get user id from verification response
	userID := resp.From

	// 10. check if user id matches user id from db
	if userID != u.ID {
		ErrorJSON(w, "user id does not match", http.StatusBadRequest)
		return
	}

	// 11. call contract to mint token
	accountAddress := web3.HexToAddress(rawAccountAddress)

	_, err = c.ticket.BuyTicketByToken(r.Context(), contractAddress, accountAddress)
	if err != nil {
		logger.Error("failed to buy ticket by token", "error", err)
		ErrorJSON(w, "failed to buy ticket by token", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "token-purchase-callback",
			Status: ws.Done,
			Data:   "Successfully purchased ticket",
		},
	})

	// 12. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Successfully purchased ticket for user with ID %s", userID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
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
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		logger.Error("failed to get jwt user from context", "error", err)
		ErrorJSON(w, "failed to get jwt user from context", http.StatusInternalServerError)
		return
	}

	// 2. get contract address from path
	rawContractAddress := strings.ToLower(chi.URLParam(r, "contractAddress"))

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
			Name:   "verify-qr",
			Status: ws.InProgress,
		},
	})

	// 4. check if ticket collection exists
	contractAddress := web3.HexToAddress(rawContractAddress)

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

	// 5. get onchain ticket collection
	onchainTicket, err := c.ticket.OnChainTicketInfo(r.Context(), contractAddress)
	if err != nil {
		logger.Error("failed to get onchain ticket collection", "error", err)
		ErrorJSON(w, "failed to get onchain ticket collection", http.StatusInternalServerError)
		return
	}

	// 6. get user from db
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user by id", "error", err)
		ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
		return
	}

	// 7. check if user is owner of ticket collection
	accountAddress := web3.HexToAddress(u.AccountAddress)

	if onchainTicket.Issuer.Big().Cmp(accountAddress.Big()) != 0 {
		ErrorJSON(w, "user is not owner of ticket collection", http.StatusBadRequest)
		return
	}

	// 8. get admin id
	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// 8. create qr code
	var mtpProofRequest protocol.ZeroKnowledgeProofRequest

	// random number
	mtpProofRequest.ID = id.UUID().ID()
	mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigV2CircuitID)
	mtpProofRequest.Query = map[string]interface{}{
		"allowedIssuers": []string{admin.ID},
		"credentialSubject": map[string]interface{}{
			"ticket_address": map[string]interface{}{
				"$eq": rawContractAddress,
			},
		},
		"context": "ipfs://QmfNkUAwq73r1HmMmzYDZ9REBqLrqdXmQm8xBdq7QbQvHz",
		"type":    "Ownership",
	}

	qrCode, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Verify ticket ownership",
		Message:     fmt.Sprintf("Scan the QR code to verify ticket ownership for %s", rawContractAddress),
		Sender:      u.ID,
		CallbackUrl: fmt.Sprintf("%s/v1/tickets/verify-callback?sessionId=%s", c.serverUrl, sessionId),
		Scope: []protocol.ZeroKnowledgeProofRequest{
			mtpProofRequest,
		},
		Timeout: 1 * time.Hour,
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
			Name:   "verify-qr",
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
	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
		return
	}

	// 2. get token from body
	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("failed to read token from body", "error", err)
		ErrorJSON(w, "failed to read token from body", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "verify-callback",
			Status: ws.InProgress,
		},
	})

	// 3. verify token
	resp, err := c.auth.AuthorizationCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		logger.Error("failed to handle verify callback", "error", err)
		ErrorJSON(w, "failed to handle verify callback", http.StatusInternalServerError)
		return
	}

	// 4. get user id from verification response
	userID := resp.From

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "verify-callback",
			Status: ws.Done,
			Data:   "Successfully verified ticket ownership",
		},
	})

	// 5. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Successfully verified ticket ownership for user with ID %s", userID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// CreateTicket godoc
//
// @Tags			tickets
// @Summary		creates ticket
// @Description	creates ticket
// @Accept			json
// @Produce		json
// @Param			name			formData	string	true	"ticket name"
// @Param			symbol			formData	string	true	"ticket symbol"
// @Param			description		formData	string	true	"ticket description"
// @Param			organizer		formData	string	true	"ticket organizer"
// @Param			location		formData	string	true	"ticket location"
// @Param			date			formData	string	true	"ticket usage date "
// @Param			bannerImage		formData	file	true	"ticket banner image file"
// @Param			ticketUri		formData	string	true	"ticket uri (ipfs hash)"
// @Param			ethPrice		formData	int64	true	"ticket eth price (min 1 gwei = 1e9)"
// @Param			tokenPrice		formData	int64	true	"ticket token price (min 1 token)"
// @Param			totalSupply		formData	int64	true	"ticket total supply (min 1 ticket)"
// @Param			saleDuration	formData	int64	true	"ticket sale duration in days (min 1 day)"
// @Success		201			{object}	CommonResponse{data=ticket.TicketCollection}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/create [post]
func (c *TicketCtrl) createTicket(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		logger.Error("failed to get jwt user from context", "error", err)
		ErrorJSON(w, "failed to get jwt user from context", http.StatusInternalServerError)
		return
	}

	// 2. get user from db
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		logger.Error("failed to find user by id", "error", err)
		ErrorJSON(w, "failed to find user by id", http.StatusInternalServerError)
		return
	}

	// 3. get params from form data
	name := r.FormValue("name")
	symbol := r.FormValue("symbol")
	description := r.FormValue("description")
	organizer := r.FormValue("organizer")
	location := r.FormValue("location")
	date := r.FormValue("date")

	bannerImage, imgHeader, err := r.FormFile("bannerImage")
	if err != nil {
		logger.Error("failed to get banner image from form data", "error", err)
		ErrorJSON(w, "failed to get banner image from form data", http.StatusInternalServerError)
		return
	}
	defer bannerImage.Close()

	ticketUri := r.FormValue("ticketUri")
	ethPrice := r.FormValue("ethPrice")
	tokenPrice := r.FormValue("tokenPrice")
	totalSupply := r.FormValue("totalSupply")
	saleDuration := r.FormValue("saleDuration")

	// TODO: validate params
	if name == "" {
		ErrorJSON(w, "name is required", http.StatusBadRequest)
		return
	}

	if symbol == "" {
		ErrorJSON(w, "symbol is required", http.StatusBadRequest)
		return
	}

	if description == "" {
		ErrorJSON(w, "description is required", http.StatusBadRequest)
		return
	}

	if organizer == "" {
		ErrorJSON(w, "organizer is required", http.StatusBadRequest)
		return
	}

	if location == "" {
		ErrorJSON(w, "location is required", http.StatusBadRequest)
		return
	}

	if date == "" {
		ErrorJSON(w, "date is required", http.StatusBadRequest)
		return
	}

	ethPriceBigInt, ok := big.NewInt(0).SetString(ethPrice, 10)
	if !ok {
		ErrorJSON(w, "failed to parse eth price", http.StatusInternalServerError)
		return
	}

	if ethPriceBigInt.Cmp(big.NewInt(1e9)) == -1 {
		ErrorJSON(w, "eth price must be greater than 1 gwei", http.StatusBadRequest)
		return
	}

	tokenPriceBigInt, ok := big.NewInt(0).SetString(tokenPrice, 10)
	if !ok {
		ErrorJSON(w, "failed to parse token price", http.StatusInternalServerError)
		return
	}

	if tokenPriceBigInt.Cmp(big.NewInt(1)) == -1 {
		ErrorJSON(w, "token price must be greater than 1 token", http.StatusBadRequest)
		return
	}

	totalSupplyBigInt, ok := big.NewInt(0).SetString(totalSupply, 10)
	if !ok {
		ErrorJSON(w, "failed to parse total supply", http.StatusInternalServerError)
		return
	}

	if totalSupplyBigInt.Cmp(big.NewInt(1)) == -1 {
		ErrorJSON(w, "total supply must be greater than 1 ticket", http.StatusBadRequest)
		return
	}

	saleDurationInt, err := strconv.ParseUint(saleDuration, 10, 64)
	if err != nil {
		logger.Error("failed to parse sale duration", "error", err)
		ErrorJSON(w, "failed to parse sale duration", http.StatusInternalServerError)
		return
	}

	if saleDurationInt < 1 {
		ErrorJSON(w, "sale duration must be greater than 1 day", http.StatusBadRequest)
		return
	}

	// 4. upload banner image to ipfs
	pinResp, err := c.ipfs.PinFile(r.Context(), bannerImage, fmt.Sprintf("%s_%s", uuid.New().String(), imgHeader.Filename))
	if err != nil {
		logger.Error("failed to pin file to ipfs", "error", err)
		ErrorJSON(w, "failed to pin file to ipfs", http.StatusInternalServerError)
		return
	}

	bannerUrl := fmt.Sprintf("https://ipfs.io/ipfs/%s", pinResp.IpfsHash)

	if !strings.HasPrefix(ticketUri, "https://ipfs.io/ipfs/") {
		ticketUri = fmt.Sprintf("https://ipfs.io/ipfs/%s", ticketUri)
	}

	// 5. call contract to create new ticket collection
	ticketIssued, err := c.ticket.IssueTicket(r.Context(), ticket.IssueTicketParams{
		TicketName:       name,
		TicketSymbol:     symbol,
		TicketUri:        ticketUri,
		Issuer:           web3.HexToAddress(u.AccountAddress),
		TicketAmount:     totalSupplyBigInt,
		TicketEthPrice:   ethPriceBigInt,
		TicketTokenPrice: tokenPriceBigInt,
		SaleDuration:     big.NewInt(0).Mul(big.NewInt(int64(saleDurationInt)), big.NewInt(86400)),
	})
	if err != nil {
		logger.Error("failed to issue ticket", "error", err)
		ErrorJSON(w, "failed to issue ticket", http.StatusInternalServerError)
		return
	}

	// 4. get onchain ticket collection data
	onchainTicket, err := c.ticket.OnChainTicketInfo(r.Context(), ticketIssued.TicketAddress)
	if err != nil {
		logger.Error("failed to get onchain ticket collection data", "error", err)
		ErrorJSON(w, "failed to get onchain ticket collection data", http.StatusInternalServerError)
		return
	}

	// 5. save ticket collection to db
	params := ticket.CreateTicketCollectionParams{
		ContractAddress: strings.ToLower(ticketIssued.TicketAddress.Hex()),
		IssuerAddress:   strings.ToLower(u.AccountAddress),
		Name:            name,
		Symbol:          symbol,
		Description:     description,
		Organizer:       organizer,
		Location:        location,
		Date:            date,
		BannerUrl:       bannerUrl,
		TicketUrl:       ticketUri,
		EthPrice:        ethPriceBigInt.String(),
		TokenPrice:      tokenPriceBigInt.String(),
		TotalSupply:     totalSupplyBigInt.String(),
		Remaining:       totalSupplyBigInt.String(),
		SaleStartAt:     onchainTicket.SaleStartAt.Int64(),
		SaleEndAt:       onchainTicket.SaleEndAt.Int64(),
	}

	ticketCollection, err := c.ticket.CreateTicketCollection(r.Context(), params)
	if err != nil {
		logger.Error("failed to save ticket collection", "error", err)
		ErrorJSON(w, "failed to save ticket collection", http.StatusInternalServerError)
		return
	}

	// 6. return success response
	resp := CommonResponse{
		Status:  http.StatusCreated,
		Message: "Successfully created ticket collection",
		Data:    ticketCollection,
	}

	_ = WriteJSON(w, http.StatusCreated, resp)
}
