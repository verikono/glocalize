package trashSkill

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.TrashSkillSelect
type Omit = riot.TrashSkillOmit
type QueryBuilder = riot.TrashSkillQueryBuilder
type CreateBuilder = riot.TrashSkillCreateBuilder
type Upsert = riot.TrashSkillUpsert
type ConflictBuilder[B any] = riot.TrashSkillConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.TrashSkill]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.TrashSkill]) riot.PredicateOf[riot.TrashSkill] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.TrashSkill]) riot.PredicateOf[riot.TrashSkill] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.TrashSkill]) riot.PredicateOf[riot.TrashSkill] {
	return riot.Not(pred)
}

var TrashId = riot.Field[riot.TrashSkill, int32]{Column: "trashId"}

var SkillId = riot.StringField[riot.TrashSkill]{Column: "skillId"}

type trashId_SkillId struct {
	riot.CompositeUniqueConstraint[riot.TrashSkill]
	Column string
}

var TrashId_SkillId = trashId_SkillId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.TrashSkill]{
		Name: "TrashId_SkillId",
		Columns: []string{
			"trashId",
			"skillId",
		},
	},
	Column: "TrashId_SkillId",
}

func (f trashId_SkillId) EQ(trashId int32, skillId string) riot.UniquePredicate[riot.TrashSkill] {
	return riot.UniquePredicate[riot.TrashSkill]{
		Data: riot.PredicateData{
			Column:   "TrashId_SkillId",
			Operator: "AND",
			Value: map[string]any{
				"trashId": trashId,
				"skillId": skillId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "trashId", Operator: "=", Value: trashId},
				{Column: "skillId", Operator: "=", Value: skillId},
			},
		},
	}
}

type CreateInput = riot.TrashSkillCreate
type Create = riot.TrashSkillCreate

type CreateArgs = riot.TrashSkillCreateArgs
type CreateManyArgs = riot.TrashSkillCreateManyArgs
type CreateManyAndReturnArgs = riot.TrashSkillCreateManyAndReturnArgs

type CreateQuery = riot.TrashSkillCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.TrashSkill, error)

type CreateManyQuery = riot.TrashSkillCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.TrashSkillCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.TrashSkill, error)

type FindUniqueArgs = riot.TrashSkillFindUniqueArgs
type FindUniqueQuery = riot.TrashSkillFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.TrashSkill, error)

type FindFirstArgs = riot.TrashSkillFindFirstArgs
type FindFirstQuery = riot.TrashSkillFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.TrashSkill, error)

type FindManyArgs = riot.TrashSkillFindManyArgs
type FindManyQuery = riot.TrashSkillFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.TrashSkill, error)

type CountArgs = riot.TrashSkillCountArgs
type CountQuery = riot.TrashSkillCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.TrashSkillUpdate

type Extension = riot.TrashSkillExtension

// ConflictUpdate creates a custom ConflictAction for TrashSkill upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.TrashSkillUpsert)) *riot.ConflictAction {
	return riot.TrashSkillConflictUpdate(fn)
}
