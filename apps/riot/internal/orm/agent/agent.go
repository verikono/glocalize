package agent

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
	"time"
)

type Select = riot.AgentSelect
type Omit = riot.AgentOmit
type QueryBuilder = riot.AgentQueryBuilder
type CreateBuilder = riot.AgentCreateBuilder
type Upsert = riot.AgentUpsert
type ConflictBuilder[B any] = riot.AgentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Agent]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Agent]) riot.PredicateOf[riot.Agent] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Agent]) riot.PredicateOf[riot.Agent] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Agent]) riot.PredicateOf[riot.Agent] {
	return riot.Not(pred)
}

var Id = riot.StringUniqueField[riot.Agent]{Column: "id"}

var Name = riot.StringUniqueField[riot.Agent]{Column: "name"}

var Description = riot.StringField[riot.Agent]{Column: "description"}

var Content = riot.StringField[riot.Agent]{Column: "content"}

var AllowedTools = riot.StringField[riot.Agent]{Column: "allowed-tools"}

var Metadata = riot.StringField[riot.Agent]{Column: "metadata"}

var CreatedAt = riot.Field[riot.Agent, time.Time]{Column: "createdAt"}

var UpdatedAt = riot.Field[riot.Agent, time.Time]{Column: "updatedAt"}

type CreateInput = riot.AgentCreate
type Create = riot.AgentCreate

type CreateArgs = riot.AgentCreateArgs
type CreateManyArgs = riot.AgentCreateManyArgs
type CreateManyAndReturnArgs = riot.AgentCreateManyAndReturnArgs

type CreateQuery = riot.AgentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Agent, error)

type CreateManyQuery = riot.AgentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.AgentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Agent, error)

type FindUniqueArgs = riot.AgentFindUniqueArgs
type FindUniqueQuery = riot.AgentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Agent, error)

type FindFirstArgs = riot.AgentFindFirstArgs
type FindFirstQuery = riot.AgentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Agent, error)

type FindManyArgs = riot.AgentFindManyArgs
type FindManyQuery = riot.AgentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Agent, error)

type CountArgs = riot.AgentCountArgs
type CountQuery = riot.AgentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.AgentUpdate

type Extension = riot.AgentExtension

// ConflictUpdate creates a custom ConflictAction for Agent upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.AgentUpsert)) *riot.ConflictAction {
	return riot.AgentConflictUpdate(fn)
}
