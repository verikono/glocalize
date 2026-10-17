package project

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.ProjectSelect
type Omit = riot.ProjectOmit
type QueryBuilder = riot.ProjectQueryBuilder
type CreateBuilder = riot.ProjectCreateBuilder
type Upsert = riot.ProjectUpsert
type ConflictBuilder[B any] = riot.ProjectConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Project]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Project]) riot.PredicateOf[riot.Project] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Project]) riot.PredicateOf[riot.Project] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Project]) riot.PredicateOf[riot.Project] {
	return riot.Not(pred)
}

var Name = riot.StringUniqueField[riot.Project]{Column: "name"}

var SkillsPath = riot.StringField[riot.Project]{Column: "skillsPath"}

var AgentsPath = riot.StringField[riot.Project]{Column: "agentsPath"}

type CreateInput = riot.ProjectCreate
type Create = riot.ProjectCreate

type CreateArgs = riot.ProjectCreateArgs
type CreateManyArgs = riot.ProjectCreateManyArgs
type CreateManyAndReturnArgs = riot.ProjectCreateManyAndReturnArgs

type CreateQuery = riot.ProjectCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Project, error)

type CreateManyQuery = riot.ProjectCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.ProjectCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Project, error)

type FindUniqueArgs = riot.ProjectFindUniqueArgs
type FindUniqueQuery = riot.ProjectFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Project, error)

type FindFirstArgs = riot.ProjectFindFirstArgs
type FindFirstQuery = riot.ProjectFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Project, error)

type FindManyArgs = riot.ProjectFindManyArgs
type FindManyQuery = riot.ProjectFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Project, error)

type CountArgs = riot.ProjectCountArgs
type CountQuery = riot.ProjectCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.ProjectUpdate

type Extension = riot.ProjectExtension

// ConflictUpdate creates a custom ConflictAction for Project upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.ProjectUpsert)) *riot.ConflictAction {
	return riot.ProjectConflictUpdate(fn)
}
