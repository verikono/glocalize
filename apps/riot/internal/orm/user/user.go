package user

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.UserSelect
type Omit = riot.UserOmit
type QueryBuilder = riot.UserQueryBuilder
type CreateBuilder = riot.UserCreateBuilder
type Upsert = riot.UserUpsert
type ConflictBuilder[B any] = riot.UserConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.User]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.User]) riot.PredicateOf[riot.User] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.User]) riot.PredicateOf[riot.User] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.User]) riot.PredicateOf[riot.User] {
	return riot.Not(pred)
}

var Id = riot.StringUniqueField[riot.User]{Column: "id"}

var Name = riot.StringField[riot.User]{Column: "name"}

type CreateInput = riot.UserCreate
type Create = riot.UserCreate

type CreateArgs = riot.UserCreateArgs
type CreateManyArgs = riot.UserCreateManyArgs
type CreateManyAndReturnArgs = riot.UserCreateManyAndReturnArgs

type CreateQuery = riot.UserCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.User, error)

type CreateManyQuery = riot.UserCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.UserCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.User, error)

type FindUniqueArgs = riot.UserFindUniqueArgs
type FindUniqueQuery = riot.UserFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.User, error)

type FindFirstArgs = riot.UserFindFirstArgs
type FindFirstQuery = riot.UserFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.User, error)

type FindManyArgs = riot.UserFindManyArgs
type FindManyQuery = riot.UserFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.User, error)

type CountArgs = riot.UserCountArgs
type CountQuery = riot.UserCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.UserUpdate

type Extension = riot.UserExtension

// ConflictUpdate creates a custom ConflictAction for User upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.UserUpsert)) *riot.ConflictAction {
	return riot.UserConflictUpdate(fn)
}
