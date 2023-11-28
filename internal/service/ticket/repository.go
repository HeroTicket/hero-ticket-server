package ticket

import (
	"context"
)

type Query interface {
	FindTicketByAddress(ctx context.Context, address string) (*Ticket, error)
	FindTicketByOwnerAddress(ctx context.Context, ownerAddress string) ([]*Ticket, error)
}

type Command interface {
	CreateTicket(ctx context.Context, params Ticket) (*Ticket, error)
	DeleteTicket(ctx context.Context, id string) error
	CreateTicketCollection(ctx context.Context, params TicketCollection) (*TicketCollection, error)
	CreateTBA(ctx context.Context, params TbaAddresses) (*TbaAddresses, error)
}

type Repository interface {
	Query
	Command
}
