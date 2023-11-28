package ticket

import "context"

type Query interface {
	FindTicketByAddress(ctx context.Context, address string) (*Ticket, error)
	FindTicketByOwnerAddress(ctx context.Context, ownerAddress string) ([]*Ticket, error)
}

type Command interface{}

type Repository interface {
	Query
	Command
}
