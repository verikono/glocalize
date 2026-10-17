package skillDependency

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.SkillDependencySelect
type Omit = riot.SkillDependencyOmit
type QueryBuilder = riot.SkillDependencyQueryBuilder
type CreateBuilder = riot.SkillDependencyCreateBuilder
type Upsert = riot.SkillDependencyUpsert
type ConflictBuilder[B any] = riot.SkillDependencyConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.SkillDependency]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.SkillDependency]) riot.PredicateOf[riot.SkillDependency] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.SkillDependency]) riot.PredicateOf[riot.SkillDependency] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.SkillDependency]) riot.PredicateOf[riot.SkillDependency] {
	return riot.Not(pred)
}

var DependencyId = riot.StringField[riot.SkillDependency]{Column: "dependencyId"}

var DependentId = riot.StringField[riot.SkillDependency]{Column: "dependentId"}

type dependencyId_DependentId struct {
	riot.CompositeUniqueConstraint[riot.SkillDependency]
	Column string
}

var DependencyId_DependentId = dependencyId_DependentId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.SkillDependency]{
		Name: "DependencyId_DependentId",
		Columns: []string{
			"dependencyId",
			"dependentId",
		},
	},
	Column: "DependencyId_DependentId",
}

func (f dependencyId_DependentId) EQ(dependencyId string, dependentId string) riot.UniquePredicate[riot.SkillDependency] {
	return riot.UniquePredicate[riot.SkillDependency]{
		Data: riot.PredicateData{
			Column:   "DependencyId_DependentId",
			Operator: "AND",
			Value: map[string]any{
				"dependencyId": dependencyId,
				"dependentId":  dependentId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "dependencyId", Operator: "=", Value: dependencyId},
				{Column: "dependentId", Operator: "=", Value: dependentId},
			},
		},
	}
}

type CreateInput = riot.SkillDependencyCreate
type Create = riot.SkillDependencyCreate

type CreateArgs = riot.SkillDependencyCreateArgs
type CreateManyArgs = riot.SkillDependencyCreateManyArgs
type CreateManyAndReturnArgs = riot.SkillDependencyCreateManyAndReturnArgs

type CreateQuery = riot.SkillDependencyCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.SkillDependency, error)

type CreateManyQuery = riot.SkillDependencyCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.SkillDependencyCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.SkillDependency, error)

type FindUniqueArgs = riot.SkillDependencyFindUniqueArgs
type FindUniqueQuery = riot.SkillDependencyFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.SkillDependency, error)

type FindFirstArgs = riot.SkillDependencyFindFirstArgs
type FindFirstQuery = riot.SkillDependencyFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.SkillDependency, error)

type FindManyArgs = riot.SkillDependencyFindManyArgs
type FindManyQuery = riot.SkillDependencyFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.SkillDependency, error)

type CountArgs = riot.SkillDependencyCountArgs
type CountQuery = riot.SkillDependencyCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.SkillDependencyUpdate

type Extension = riot.SkillDependencyExtension

// ConflictUpdate creates a custom ConflictAction for SkillDependency upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.SkillDependencyUpsert)) *riot.ConflictAction {
	return riot.SkillDependencyConflictUpdate(fn)
}
