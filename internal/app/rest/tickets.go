package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/ipfs"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/user"
)

type TicketCtrl struct {
	serverUrl string

	auth auth.Service
	ipfs ipfs.Service
	jwt  jwt.Service
	user user.Service
	// ticket ticket.Service
}

func NewTicketCtrl(auth auth.Service, ipfs ipfs.Service, jwt jwt.Service, user user.Service, serverUrl string) *TicketCtrl {
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

	r.Get("/", c.Tickets)
	r.Get("/{contractAddress}", c.ticket)

	r.Post("/purchase-callback", c.purchaseCallback)
	r.Post("/verify-callback", c.verifyCallback)

	r.Group(func(r chi.Router) {
		r.Use(TokenRequired(c.jwt))
		r.Get("/{contractAddress}/purchase-qr", c.purchaseQR)
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
func (c *TicketCtrl) Tickets(w http.ResponseWriter, r *http.Request) {
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
func (c *TicketCtrl) ticket(w http.ResponseWriter, r *http.Request) {
	// 1. get contract address from path

	// 2. get ticket from db, if user is logged in, also get user ticket ownership

	// 3. return ticket
}

// PurchaseQR godoc
//
// @Tags			tickets
// @Summary		returns purchase authorization qr code
// @Description	returns purchase authorization qr code
// @Accept			json
// @Produce		json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			sessionId		query	string	true	"session id"
// @Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Security 		BearerAuth
// @Router			/v1/tickets/{contractAddress}/purchase-qr [get]
func (c *TicketCtrl) purchaseQR(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt user from context

	// 2. get contract address from path

	// 3. get session id from query

	// 4. get user from db

	// 5. check if ticket collection exists

	// 6. check if user has ticket

	// 7. create qr code

	// 8. return qr code
}

// PurchaseCallback godoc
//
// @Tags			tickets
// @Summary		purchase callback
// @Description	purchase callback
// @Accept			json
// @Produce		json
// @Param			contractAddress	path	string	true	"contract address"
// @Param			accountAddress	query	string	true	"account address"
// @Param			sessionId		query	string	true	"session id"
// @Param			token			body	string	true	"token"
// @Success		200			{object}	CommonResponse
// @Failure		400			{object}	CommonResponse
// @Failure		500			{object}	CommonResponse
// @Router			/v1/tickets/purchase-callback [post]
func (c *TicketCtrl) purchaseCallback(w http.ResponseWriter, r *http.Request) {
	// 1. get contract address from path

	// 2. get account address from query

	// 3. get session id from query

	// 4. get token from body

	// 5. get user from db

	// 6. check if ticket collection exists

	// 7. check if user has ticket

	// 8. verify token

	// 9. get user id from verification response

	// 10. check if user id matches user id from db

	// 11. call contract to set user address on whitelist

	// 12. return success response
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
