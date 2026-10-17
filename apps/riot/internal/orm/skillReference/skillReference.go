package skillReference

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.SkillReferenceSelect
type Omit = riot.SkillReferenceOmit
type QueryBuilder = riot.SkillReferenceQueryBuilder
type CreateBuilder = riot.SkillReferenceCreateBuilder
type Upsert = riot.SkillReferenceUpsert
type ConflictBuilder[B any] = riot.SkillReferenceConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.SkillReference]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.SkillReference]) riot.PredicateOf[riot.SkillReference] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.SkillReference]) riot.PredicateOf[riot.SkillReference] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.SkillReference]) riot.PredicateOf[riot.SkillReference] {
	return riot.Not(pred)
}

var SkillId = riot.StringField[riot.SkillReference]{Column: "skillId"}

var ReferenceName = riot.StringField[riot.SkillReference]{Column: "referenceName"}

type skillId_ReferenceName struct {
	riot.CompositeUniqueConstraint[riot.SkillReference]
	Column string
}

var SkillId_ReferenceName = skillId_ReferenceName{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.SkillReference]{
		Name: "SkillId_ReferenceName",
		Columns: []string{
			"skillId",
			"referenceName",
		},
	},
	Column: "SkillId_ReferenceName",
}

func (f skillId_ReferenceName) EQ(skillId string, referenceName string) riot.UniquePredicate[riot.SkillReference] {
	return riot.UniquePredicate[riot.SkillReference]{
		Data: riot.PredicateData{
			Column:   "SkillId_ReferenceName",
			Operator: "AND",
			Value: map[string]any{
				"skillId":       skillId,
				"referenceName": referenceName,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "skillId", Operator: "=", Value: skillId},
				{Column: "referenceName", Operator: "=", Value: referenceName},
			},
		},
	}
}

type CreateInput = riot.SkillReferenceCreate
type Create = riot.SkillReferenceCreate

type CreateArgs = riot.SkillReferenceCreateArgs
type CreateManyArgs = riot.SkillReferenceCreateManyArgs
type CreateManyAndReturnArgs = riot.SkillReferenceCreateManyAndReturnArgs

type CreateQuery = riot.SkillReferenceCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.SkillReference, error)

type CreateManyQuery = riot.SkillReferenceCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.SkillReferenceCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.SkillReference, error)

type FindUniqueArgs = riot.SkillReferenceFindUniqueArgs
type FindUniqueQuery = riot.SkillReferenceFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.SkillReference, error)

type FindFirstArgs = riot.SkillReferenceFindFirstArgs
type FindFirstQuery = riot.SkillReferenceFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.SkillReference, error)

type FindManyArgs = riot.SkillReferenceFindManyArgs
type FindManyQuery = riot.SkillReferenceFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.SkillReference, error)

type CountArgs = riot.SkillReferenceCountArgs
type CountQuery = riot.SkillReferenceCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.SkillReferenceUpdate

type Extension = riot.SkillReferenceExtension

// ConflictUpdate creates a custom ConflictAction for SkillReference upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.SkillReferenceUpsert)) *riot.ConflictAction {
	return riot.SkillReferenceConflictUpdate(fn)
}
