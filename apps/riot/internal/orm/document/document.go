package document

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.DocumentSelect
type Omit = riot.DocumentOmit
type QueryBuilder = riot.DocumentQueryBuilder
type CreateBuilder = riot.DocumentCreateBuilder
type Upsert = riot.DocumentUpsert
type ConflictBuilder[B any] = riot.DocumentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Document]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Document]) riot.PredicateOf[riot.Document] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Document]) riot.PredicateOf[riot.Document] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Document]) riot.PredicateOf[riot.Document] {
	return riot.Not(pred)
}

var Id = riot.StringField[riot.Document]{Column: "id"}

var Kind = riot.StringField[riot.Document]{Column: "kind"}

var Body = riot.Field[riot.Document, []byte]{Column: "body"}

var UpdatedAt = riot.Field[riot.Document, int32]{Column: "updated_at"}

type id_Kind struct {
	riot.CompositeUniqueConstraint[riot.Document]
	Column string
}

var Id_Kind = id_Kind{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.Document]{
		Name: "Id_Kind",
		Columns: []string{
			"id",
			"kind",
		},
	},
	Column: "Id_Kind",
}

func (f id_Kind) EQ(id string, kind string) riot.UniquePredicate[riot.Document] {
	return riot.UniquePredicate[riot.Document]{
		Data: riot.PredicateData{
			Column:   "Id_Kind",
			Operator: "AND",
			Value: map[string]any{
				"id":   id,
				"kind": kind,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "id", Operator: "=", Value: id},
				{Column: "kind", Operator: "=", Value: kind},
			},
		},
	}
}

type CreateInput = riot.DocumentCreate
type Create = riot.DocumentCreate

type CreateArgs = riot.DocumentCreateArgs
type CreateManyArgs = riot.DocumentCreateManyArgs
type CreateManyAndReturnArgs = riot.DocumentCreateManyAndReturnArgs

type CreateQuery = riot.DocumentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Document, error)

type CreateManyQuery = riot.DocumentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.DocumentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Document, error)

type FindUniqueArgs = riot.DocumentFindUniqueArgs
type FindUniqueQuery = riot.DocumentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Document, error)

type FindFirstArgs = riot.DocumentFindFirstArgs
type FindFirstQuery = riot.DocumentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Document, error)

type FindManyArgs = riot.DocumentFindManyArgs
type FindManyQuery = riot.DocumentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Document, error)

type CountArgs = riot.DocumentCountArgs
type CountQuery = riot.DocumentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.DocumentUpdate

type Extension = riot.DocumentExtension

// ConflictUpdate creates a custom ConflictAction for Document upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.DocumentUpsert)) *riot.ConflictAction {
	return riot.DocumentConflictUpdate(fn)
}
