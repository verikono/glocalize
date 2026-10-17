package app

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.AppSelect
type Omit = riot.AppOmit
type QueryBuilder = riot.AppQueryBuilder
type CreateBuilder = riot.AppCreateBuilder
type Upsert = riot.AppUpsert
type ConflictBuilder[B any] = riot.AppConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.App]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.App]) riot.PredicateOf[riot.App] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.App]) riot.PredicateOf[riot.App] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.App]) riot.PredicateOf[riot.App] {
	return riot.Not(pred)
}

var Name = riot.StringUniqueField[riot.App]{Column: "name"}

var Project_id = riot.StringField[riot.App]{Column: "project_id"}

type CreateInput = riot.AppCreate
type Create = riot.AppCreate

type CreateArgs = riot.AppCreateArgs
type CreateManyArgs = riot.AppCreateManyArgs
type CreateManyAndReturnArgs = riot.AppCreateManyAndReturnArgs

type CreateQuery = riot.AppCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.App, error)

type CreateManyQuery = riot.AppCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.AppCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.App, error)

type FindUniqueArgs = riot.AppFindUniqueArgs
type FindUniqueQuery = riot.AppFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.App, error)

type FindFirstArgs = riot.AppFindFirstArgs
type FindFirstQuery = riot.AppFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.App, error)

type FindManyArgs = riot.AppFindManyArgs
type FindManyQuery = riot.AppFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.App, error)

type CountArgs = riot.AppCountArgs
type CountQuery = riot.AppCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.AppUpdate

type Extension = riot.AppExtension

// ConflictUpdate creates a custom ConflictAction for App upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.AppUpsert)) *riot.ConflictAction {
	return riot.AppConflictUpdate(fn)
}
