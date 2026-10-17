package projectSkill

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.ProjectSkillSelect
type Omit = riot.ProjectSkillOmit
type QueryBuilder = riot.ProjectSkillQueryBuilder
type CreateBuilder = riot.ProjectSkillCreateBuilder
type Upsert = riot.ProjectSkillUpsert
type ConflictBuilder[B any] = riot.ProjectSkillConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.ProjectSkill]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.ProjectSkill]) riot.PredicateOf[riot.ProjectSkill] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.ProjectSkill]) riot.PredicateOf[riot.ProjectSkill] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.ProjectSkill]) riot.PredicateOf[riot.ProjectSkill] {
	return riot.Not(pred)
}

var ProjectName = riot.StringField[riot.ProjectSkill]{Column: "projectName"}

var SkillId = riot.StringField[riot.ProjectSkill]{Column: "skillId"}

type projectName_SkillId struct {
	riot.CompositeUniqueConstraint[riot.ProjectSkill]
	Column string
}

var ProjectName_SkillId = projectName_SkillId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.ProjectSkill]{
		Name: "ProjectName_SkillId",
		Columns: []string{
			"projectName",
			"skillId",
		},
	},
	Column: "ProjectName_SkillId",
}

func (f projectName_SkillId) EQ(projectName string, skillId string) riot.UniquePredicate[riot.ProjectSkill] {
	return riot.UniquePredicate[riot.ProjectSkill]{
		Data: riot.PredicateData{
			Column:   "ProjectName_SkillId",
			Operator: "AND",
			Value: map[string]any{
				"projectName": projectName,
				"skillId":     skillId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "projectName", Operator: "=", Value: projectName},
				{Column: "skillId", Operator: "=", Value: skillId},
			},
		},
	}
}

type CreateInput = riot.ProjectSkillCreate
type Create = riot.ProjectSkillCreate

type CreateArgs = riot.ProjectSkillCreateArgs
type CreateManyArgs = riot.ProjectSkillCreateManyArgs
type CreateManyAndReturnArgs = riot.ProjectSkillCreateManyAndReturnArgs

type CreateQuery = riot.ProjectSkillCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.ProjectSkill, error)

type CreateManyQuery = riot.ProjectSkillCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.ProjectSkillCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.ProjectSkill, error)

type FindUniqueArgs = riot.ProjectSkillFindUniqueArgs
type FindUniqueQuery = riot.ProjectSkillFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.ProjectSkill, error)

type FindFirstArgs = riot.ProjectSkillFindFirstArgs
type FindFirstQuery = riot.ProjectSkillFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.ProjectSkill, error)

type FindManyArgs = riot.ProjectSkillFindManyArgs
type FindManyQuery = riot.ProjectSkillFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.ProjectSkill, error)

type CountArgs = riot.ProjectSkillCountArgs
type CountQuery = riot.ProjectSkillCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.ProjectSkillUpdate

type Extension = riot.ProjectSkillExtension

// ConflictUpdate creates a custom ConflictAction for ProjectSkill upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.ProjectSkillUpsert)) *riot.ConflictAction {
	return riot.ProjectSkillConflictUpdate(fn)
}
