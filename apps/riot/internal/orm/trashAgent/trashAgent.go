package trashAgent

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.TrashAgentSelect
type Omit = riot.TrashAgentOmit
type QueryBuilder = riot.TrashAgentQueryBuilder
type CreateBuilder = riot.TrashAgentCreateBuilder
type Upsert = riot.TrashAgentUpsert
type ConflictBuilder[B any] = riot.TrashAgentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.TrashAgent]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.TrashAgent]) riot.PredicateOf[riot.TrashAgent] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.TrashAgent]) riot.PredicateOf[riot.TrashAgent] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.TrashAgent]) riot.PredicateOf[riot.TrashAgent] {
	return riot.Not(pred)
}

var TrashId = riot.Field[riot.TrashAgent, int32]{Column: "trashId"}

var AgentId = riot.StringField[riot.TrashAgent]{Column: "agentId"}

type trashId_AgentId struct {
	riot.CompositeUniqueConstraint[riot.TrashAgent]
	Column string
}

var TrashId_AgentId = trashId_AgentId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.TrashAgent]{
		Name: "TrashId_AgentId",
		Columns: []string{
			"trashId",
			"agentId",
		},
	},
	Column: "TrashId_AgentId",
}

func (f trashId_AgentId) EQ(trashId int32, agentId string) riot.UniquePredicate[riot.TrashAgent] {
	return riot.UniquePredicate[riot.TrashAgent]{
		Data: riot.PredicateData{
			Column:   "TrashId_AgentId",
			Operator: "AND",
			Value: map[string]any{
				"trashId": trashId,
				"agentId": agentId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "trashId", Operator: "=", Value: trashId},
				{Column: "agentId", Operator: "=", Value: agentId},
			},
		},
	}
}

type CreateInput = riot.TrashAgentCreate
type Create = riot.TrashAgentCreate

type CreateArgs = riot.TrashAgentCreateArgs
type CreateManyArgs = riot.TrashAgentCreateManyArgs
type CreateManyAndReturnArgs = riot.TrashAgentCreateManyAndReturnArgs

type CreateQuery = riot.TrashAgentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.TrashAgent, error)

type CreateManyQuery = riot.TrashAgentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.TrashAgentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.TrashAgent, error)

type FindUniqueArgs = riot.TrashAgentFindUniqueArgs
type FindUniqueQuery = riot.TrashAgentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.TrashAgent, error)

type FindFirstArgs = riot.TrashAgentFindFirstArgs
type FindFirstQuery = riot.TrashAgentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.TrashAgent, error)

type FindManyArgs = riot.TrashAgentFindManyArgs
type FindManyQuery = riot.TrashAgentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.TrashAgent, error)

type CountArgs = riot.TrashAgentCountArgs
type CountQuery = riot.TrashAgentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.TrashAgentUpdate

type Extension = riot.TrashAgentExtension

// ConflictUpdate creates a custom ConflictAction for TrashAgent upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.TrashAgentUpsert)) *riot.ConflictAction {
	return riot.TrashAgentConflictUpdate(fn)
}
