package appSkill

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.AppSkillSelect
type Omit = riot.AppSkillOmit
type QueryBuilder = riot.AppSkillQueryBuilder
type CreateBuilder = riot.AppSkillCreateBuilder
type Upsert = riot.AppSkillUpsert
type ConflictBuilder[B any] = riot.AppSkillConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.AppSkill]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.AppSkill]) riot.PredicateOf[riot.AppSkill] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.AppSkill]) riot.PredicateOf[riot.AppSkill] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.AppSkill]) riot.PredicateOf[riot.AppSkill] {
	return riot.Not(pred)
}

var App_language = riot.Field[riot.AppSkill, riot.AppLanguageType]{Column: "app_language"}

var App_name = riot.StringField[riot.AppSkill]{Column: "app_name"}

var SkillId = riot.StringField[riot.AppSkill]{Column: "skillId"}

var Priority = riot.Field[riot.AppSkill, riot.PriorityType]{Column: "priority"}

type app_language_App_name_SkillId struct {
	riot.CompositeUniqueConstraint[riot.AppSkill]
	Column string
}

var App_language_App_name_SkillId = app_language_App_name_SkillId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.AppSkill]{
		Name: "App_language_App_name_SkillId",
		Columns: []string{
			"app_language",
			"app_name",
			"skillId",
		},
	},
	Column: "App_language_App_name_SkillId",
}

func (f app_language_App_name_SkillId) EQ(app_language riot.AppLanguageType, app_name string, skillId string) riot.UniquePredicate[riot.AppSkill] {
	return riot.UniquePredicate[riot.AppSkill]{
		Data: riot.PredicateData{
			Column:   "App_language_App_name_SkillId",
			Operator: "AND",
			Value: map[string]any{
				"app_language": app_language,
				"app_name":     app_name,
				"skillId":      skillId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "app_language", Operator: "=", Value: app_language},
				{Column: "app_name", Operator: "=", Value: app_name},
				{Column: "skillId", Operator: "=", Value: skillId},
			},
		},
	}
}

type CreateInput = riot.AppSkillCreate
type Create = riot.AppSkillCreate

type CreateArgs = riot.AppSkillCreateArgs
type CreateManyArgs = riot.AppSkillCreateManyArgs
type CreateManyAndReturnArgs = riot.AppSkillCreateManyAndReturnArgs

type CreateQuery = riot.AppSkillCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.AppSkill, error)

type CreateManyQuery = riot.AppSkillCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.AppSkillCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.AppSkill, error)

type FindUniqueArgs = riot.AppSkillFindUniqueArgs
type FindUniqueQuery = riot.AppSkillFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.AppSkill, error)

type FindFirstArgs = riot.AppSkillFindFirstArgs
type FindFirstQuery = riot.AppSkillFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.AppSkill, error)

type FindManyArgs = riot.AppSkillFindManyArgs
type FindManyQuery = riot.AppSkillFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.AppSkill, error)

type CountArgs = riot.AppSkillCountArgs
type CountQuery = riot.AppSkillCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.AppSkillUpdate

type Extension = riot.AppSkillExtension

// ConflictUpdate creates a custom ConflictAction for AppSkill upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.AppSkillUpsert)) *riot.ConflictAction {
	return riot.AppSkillConflictUpdate(fn)
}
