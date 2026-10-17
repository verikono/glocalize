package agentSkill

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.AgentSkillSelect
type Omit = riot.AgentSkillOmit
type QueryBuilder = riot.AgentSkillQueryBuilder
type CreateBuilder = riot.AgentSkillCreateBuilder
type Upsert = riot.AgentSkillUpsert
type ConflictBuilder[B any] = riot.AgentSkillConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.AgentSkill]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.AgentSkill]) riot.PredicateOf[riot.AgentSkill] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.AgentSkill]) riot.PredicateOf[riot.AgentSkill] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.AgentSkill]) riot.PredicateOf[riot.AgentSkill] {
	return riot.Not(pred)
}

var AgentId = riot.StringField[riot.AgentSkill]{Column: "agentId"}

var SkillId = riot.StringField[riot.AgentSkill]{Column: "skillId"}

var Priority = riot.Field[riot.AgentSkill, riot.PriorityType]{Column: "priority"}

type agentId_SkillId struct {
	riot.CompositeUniqueConstraint[riot.AgentSkill]
	Column string
}

var AgentId_SkillId = agentId_SkillId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.AgentSkill]{
		Name: "AgentId_SkillId",
		Columns: []string{
			"agentId",
			"skillId",
		},
	},
	Column: "AgentId_SkillId",
}

func (f agentId_SkillId) EQ(agentId string, skillId string) riot.UniquePredicate[riot.AgentSkill] {
	return riot.UniquePredicate[riot.AgentSkill]{
		Data: riot.PredicateData{
			Column:   "AgentId_SkillId",
			Operator: "AND",
			Value: map[string]any{
				"agentId": agentId,
				"skillId": skillId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "agentId", Operator: "=", Value: agentId},
				{Column: "skillId", Operator: "=", Value: skillId},
			},
		},
	}
}

type CreateInput = riot.AgentSkillCreate
type Create = riot.AgentSkillCreate

type CreateArgs = riot.AgentSkillCreateArgs
type CreateManyArgs = riot.AgentSkillCreateManyArgs
type CreateManyAndReturnArgs = riot.AgentSkillCreateManyAndReturnArgs

type CreateQuery = riot.AgentSkillCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.AgentSkill, error)

type CreateManyQuery = riot.AgentSkillCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.AgentSkillCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.AgentSkill, error)

type FindUniqueArgs = riot.AgentSkillFindUniqueArgs
type FindUniqueQuery = riot.AgentSkillFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.AgentSkill, error)

type FindFirstArgs = riot.AgentSkillFindFirstArgs
type FindFirstQuery = riot.AgentSkillFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.AgentSkill, error)

type FindManyArgs = riot.AgentSkillFindManyArgs
type FindManyQuery = riot.AgentSkillFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.AgentSkill, error)

type CountArgs = riot.AgentSkillCountArgs
type CountQuery = riot.AgentSkillCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.AgentSkillUpdate

type Extension = riot.AgentSkillExtension

// ConflictUpdate creates a custom ConflictAction for AgentSkill upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.AgentSkillUpsert)) *riot.ConflictAction {
	return riot.AgentSkillConflictUpdate(fn)
}
