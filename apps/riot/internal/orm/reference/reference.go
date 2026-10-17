package reference

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.ReferenceSelect
type Omit = riot.ReferenceOmit
type QueryBuilder = riot.ReferenceQueryBuilder
type CreateBuilder = riot.ReferenceCreateBuilder
type Upsert = riot.ReferenceUpsert
type ConflictBuilder[B any] = riot.ReferenceConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Reference]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Reference]) riot.PredicateOf[riot.Reference] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Reference]) riot.PredicateOf[riot.Reference] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Reference]) riot.PredicateOf[riot.Reference] {
	return riot.Not(pred)
}

var Name = riot.StringUniqueField[riot.Reference]{Column: "name"}

var Description = riot.StringField[riot.Reference]{Column: "description"}

var AccessType = riot.Field[riot.Reference, riot.URIType]{Column: "accessType"}

var ResourceType = riot.Field[riot.Reference, riot.ReferenceTypeType]{Column: "resourceType"}

var Content = riot.OptionalStringField[riot.Reference]{StringField: riot.StringField[riot.Reference]{Column: "content"}}

type CreateInput = riot.ReferenceCreate
type Create = riot.ReferenceCreate

type CreateArgs = riot.ReferenceCreateArgs
type CreateManyArgs = riot.ReferenceCreateManyArgs
type CreateManyAndReturnArgs = riot.ReferenceCreateManyAndReturnArgs

type CreateQuery = riot.ReferenceCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Reference, error)

type CreateManyQuery = riot.ReferenceCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.ReferenceCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Reference, error)

type FindUniqueArgs = riot.ReferenceFindUniqueArgs
type FindUniqueQuery = riot.ReferenceFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Reference, error)

type FindFirstArgs = riot.ReferenceFindFirstArgs
type FindFirstQuery = riot.ReferenceFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Reference, error)

type FindManyArgs = riot.ReferenceFindManyArgs
type FindManyQuery = riot.ReferenceFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Reference, error)

type CountArgs = riot.ReferenceCountArgs
type CountQuery = riot.ReferenceCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.ReferenceUpdate

type Extension = riot.ReferenceExtension

// ConflictUpdate creates a custom ConflictAction for Reference upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.ReferenceUpsert)) *riot.ConflictAction {
	return riot.ReferenceConflictUpdate(fn)
}
