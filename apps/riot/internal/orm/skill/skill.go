package skill

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
	"time"
)

type Select = riot.SkillSelect
type Omit = riot.SkillOmit
type QueryBuilder = riot.SkillQueryBuilder
type CreateBuilder = riot.SkillCreateBuilder
type Upsert = riot.SkillUpsert
type ConflictBuilder[B any] = riot.SkillConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.Skill]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.Skill]) riot.PredicateOf[riot.Skill] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.Skill]) riot.PredicateOf[riot.Skill] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.Skill]) riot.PredicateOf[riot.Skill] {
	return riot.Not(pred)
}

var Id = riot.StringUniqueField[riot.Skill]{Column: "id"}

var Name = riot.StringField[riot.Skill]{Column: "name"}

var Description = riot.StringField[riot.Skill]{Column: "description"}

var Url = riot.StringUniqueField[riot.Skill]{Column: "url"}

var Content = riot.StringField[riot.Skill]{Column: "content"}

var CreatedAt = riot.Field[riot.Skill, time.Time]{Column: "createdAt"}

var UpdatedAt = riot.Field[riot.Skill, time.Time]{Column: "updatedAt"}

type CreateInput = riot.SkillCreate
type Create = riot.SkillCreate

type CreateArgs = riot.SkillCreateArgs
type CreateManyArgs = riot.SkillCreateManyArgs
type CreateManyAndReturnArgs = riot.SkillCreateManyAndReturnArgs

type CreateQuery = riot.SkillCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.Skill, error)

type CreateManyQuery = riot.SkillCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.SkillCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.Skill, error)

type FindUniqueArgs = riot.SkillFindUniqueArgs
type FindUniqueQuery = riot.SkillFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.Skill, error)

type FindFirstArgs = riot.SkillFindFirstArgs
type FindFirstQuery = riot.SkillFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.Skill, error)

type FindManyArgs = riot.SkillFindManyArgs
type FindManyQuery = riot.SkillFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.Skill, error)

type CountArgs = riot.SkillCountArgs
type CountQuery = riot.SkillCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.SkillUpdate

type Extension = riot.SkillExtension

// ConflictUpdate creates a custom ConflictAction for Skill upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.SkillUpsert)) *riot.ConflictAction {
	return riot.SkillConflictUpdate(fn)
}
