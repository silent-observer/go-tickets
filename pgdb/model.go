package pgdb

import (
	"time"

	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel `bun:"table:projects,alias:p"`

	Id		uint32	`bun:",pk,autoincrement"`
	Name	string	`bun:",varchar(20),unique,notnull"`
}

type Ticket struct {
	bun.BaseModel `bun:"table:tickets,alias:t"`

	Id			uint32		`bun:",pk,autoincrement"`
	ProjectId 	uint32		`bun:",notnull"`
	TicketNum	uint32		`bun:",notnull"`
	Title 		string		`bun:",notnull"`
	Content		string		`bun:",notnull"`
	StatusId	uint32		`bun:",notnull"`
	CreatedDate	time.Time	`bun:",type:timestamptz,nullzero,notnull,default:current_timestamp"`
	EditedDate	time.Time	`bun:",type:timestamptz,nullzero,notnull,default:current_timestamp"`
	CreatedBy	uint32		`bun:",notnull"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id 			uint32		`bun:",pk,autoincrement"`
	Username	string		`bun:",type:varchar(120),notnull,unique"`
}

type Status struct {
	bun.BaseModel `bun:"table:statuses,alias:s"`

	Id			uint32		`bun:",pk,autoincrement"`
	ProjectId 	uint32		`bun:",notnull"`
	Name		string		`bun:",type:varchar(50),notnull"`
	Icon		StatusIcon	`bun:",type:varchar(20),notnull"`
}

type StatusIcon string

const (
    StatusIconOpen			StatusIcon = "open"
    StatusIconInProgress	StatusIcon = "progress"
    StatusIconDone			StatusIcon = "done"
)

type TicketLink struct {
	bun.BaseModel `bun:"table:ticket_links,alias:tl"`

	TicketFrom	uint32		`bun:",pk"`
	TicketTo	uint32		`bun:",pk"`
	Type		LinkType	`bun:",type:varchar(20),notnull"`
}

type LinkType string

const (
    LinkTypeRelated		LinkType = "related"
    LinkTypeDuplicate	LinkType = "duplicate"
    LinkTypeChild		LinkType = "child"
    LinkTypeBlocks		LinkType = "blocks"
)

type Tag struct {
	bun.BaseModel `bun:"table:tags"`

	Id			uint32		`bun:",pk,autoincrement"`
	ProjectId 	uint32		`bun:",notnull"`
	Name		string		`bun:",type:varchar(50),notnull"`
}

type TicketTag struct {
	bun.BaseModel `bun:"table:ticket_tags,alias:tt"`

	TicketId	uint32		`bun:",pk"`
	TagId		uint32		`bun:",pk"`
	Value		*string		`bun:"val,type:varchar(50)"`
}

type TicketAssignedUser struct {
	bun.BaseModel `bun:"table:ticket_assigned,alias:ta"`

	TicketId	uint32		`bun:",pk"`
	UserId		uint32		`bun:",pk"`
}

type Board struct {
	bun.BaseModel `bun:"table:boards,alias:b"`

	Id			uint32		`bun:",pk,autoincrement"`
	ProjectId 	uint32		`bun:",notnull"`
	Name 		string		`bun:",type:varchar(50),notnull"`
}

type BoardTickets struct {
	bun.BaseModel `bun:"table:board_tickets,alias:bt"`

	BoardId		uint32		`bun:",pk"`
	TicketId	uint32		`bun:",pk"`
}