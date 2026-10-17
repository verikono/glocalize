package trash

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
	"time"
)

type Select = riot.TrashSelect
type Omit = riot.TrashOmit
type QueryBuilder = riot.TrashQueryBuilder
type CreateBuilder = riot.TrashCreateBuilder
type Upsert = riot.TrashUpsert
type ConflictBuilder[B any] = riot.TrashConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Trash]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Trash]) riot.PredicateOf[riot.Trash] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Trash]) riot.PredicateOf[riot.Trash] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Trash]) riot.PredicateOf[riot.Trash] {
	return riot.Not(pred)
}

var Id = riot.UniqueField[riot.Trash, int32]{Column: "id"}

var CreatedAt = riot.Field[riot.Trash, time.Time]{Column: "createdAt"}

var Content = riot.StringField[riot.Trash]{Column: "content"}

type CreateInput = riot.TrashCreate
type Create = riot.TrashCreate

type CreateArgs = riot.TrashCreateArgs
type CreateManyArgs = riot.TrashCreateManyArgs
type CreateManyAndReturnArgs = riot.TrashCreateManyAndReturnArgs

type CreateQuery = riot.TrashCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Trash, error)

type CreateManyQuery = riot.TrashCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.TrashCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Trash, error)

type FindUniqueArgs = riot.TrashFindUniqueArgs
type FindUniqueQuery = riot.TrashFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Trash, error)

type FindFirstArgs = riot.TrashFindFirstArgs
type FindFirstQuery = riot.TrashFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Trash, error)

type FindManyArgs = riot.TrashFindManyArgs
type FindManyQuery = riot.TrashFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Trash, error)

type CountArgs = riot.TrashCountArgs
type CountQuery = riot.TrashCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.TrashUpdate

type Extension = riot.TrashExtension

// ConflictUpdate creates a custom ConflictAction for Trash upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.TrashUpsert)) *riot.ConflictAction {
	return riot.TrashConflictUpdate(fn)
}
