package projectSpecialization

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.ProjectSpecializationSelect
type Omit = riot.ProjectSpecializationOmit
type QueryBuilder = riot.ProjectSpecializationQueryBuilder
type CreateBuilder = riot.ProjectSpecializationCreateBuilder
type Upsert = riot.ProjectSpecializationUpsert
type ConflictBuilder[B any] = riot.ProjectSpecializationConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.ProjectSpecialization]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.ProjectSpecialization]) riot.PredicateOf[riot.ProjectSpecialization] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.ProjectSpecialization]) riot.PredicateOf[riot.ProjectSpecialization] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.ProjectSpecialization]) riot.PredicateOf[riot.ProjectSpecialization] {
	return riot.Not(pred)
}

var GeneralizesName = riot.StringField[riot.ProjectSpecialization]{Column: "generalizesName"}

var SpecializesName = riot.StringField[riot.ProjectSpecialization]{Column: "specializesName"}

type generalizesName_SpecializesName struct {
	riot.CompositeUniqueConstraint[riot.ProjectSpecialization]
	Column string
}

var GeneralizesName_SpecializesName = generalizesName_SpecializesName{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.ProjectSpecialization]{
		Name: "GeneralizesName_SpecializesName",
		Columns: []string{
			"generalizesName",
			"specializesName",
		},
	},
	Column: "GeneralizesName_SpecializesName",
}

func (f generalizesName_SpecializesName) EQ(generalizesName string, specializesName string) riot.UniquePredicate[riot.ProjectSpecialization] {
	return riot.UniquePredicate[riot.ProjectSpecialization]{
		Data: riot.PredicateData{
			Column:   "GeneralizesName_SpecializesName",
			Operator: "AND",
			Value: map[string]any{
				"generalizesName": generalizesName,
				"specializesName": specializesName,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "generalizesName", Operator: "=", Value: generalizesName},
				{Column: "specializesName", Operator: "=", Value: specializesName},
			},
		},
	}
}

type CreateInput = riot.ProjectSpecializationCreate
type Create = riot.ProjectSpecializationCreate

type CreateArgs = riot.ProjectSpecializationCreateArgs
type CreateManyArgs = riot.ProjectSpecializationCreateManyArgs
type CreateManyAndReturnArgs = riot.ProjectSpecializationCreateManyAndReturnArgs

type CreateQuery = riot.ProjectSpecializationCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.ProjectSpecialization, error)

type CreateManyQuery = riot.ProjectSpecializationCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.ProjectSpecializationCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.ProjectSpecialization, error)

type FindUniqueArgs = riot.ProjectSpecializationFindUniqueArgs
type FindUniqueQuery = riot.ProjectSpecializationFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.ProjectSpecialization, error)

type FindFirstArgs = riot.ProjectSpecializationFindFirstArgs
type FindFirstQuery = riot.ProjectSpecializationFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.ProjectSpecialization, error)

type FindManyArgs = riot.ProjectSpecializationFindManyArgs
type FindManyQuery = riot.ProjectSpecializationFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.ProjectSpecialization, error)

type CountArgs = riot.ProjectSpecializationCountArgs
type CountQuery = riot.ProjectSpecializationCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.ProjectSpecializationUpdate

type Extension = riot.ProjectSpecializationExtension

// ConflictUpdate creates a custom ConflictAction for ProjectSpecialization upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.ProjectSpecializationUpsert)) *riot.ConflictAction {
	return riot.ProjectSpecializationConflictUpdate(fn)
}
