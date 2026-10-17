package styleToken

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.StyleTokenSelect
type Omit = riot.StyleTokenOmit
type QueryBuilder = riot.StyleTokenQueryBuilder
type CreateBuilder = riot.StyleTokenCreateBuilder
type Upsert = riot.StyleTokenUpsert
type ConflictBuilder[B any] = riot.StyleTokenConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.StyleToken]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.StyleToken]) riot.PredicateOf[riot.StyleToken] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.StyleToken]) riot.PredicateOf[riot.StyleToken] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.StyleToken]) riot.PredicateOf[riot.StyleToken] {
	return riot.Not(pred)
}

var Style = riot.StringField[riot.StyleToken]{Column: "style"}

var Value = riot.StringField[riot.StyleToken]{Column: "value"}

var Mode = riot.Field[riot.StyleToken, riot.ThemeModeType]{Column: "mode"}

var Type = riot.Field[riot.StyleToken, riot.StyleTypeType]{Column: "type"}

var Token = riot.Field[riot.StyleToken, riot.TokenTypeType]{Column: "token"}

var App_id = riot.StringField[riot.StyleToken]{Column: "app_id"}

var Variable = riot.OptionalStringField[riot.StyleToken]{StringField: riot.StringField[riot.StyleToken]{Column: "variable"}}

type app_id_Mode_Token_Style struct {
	riot.CompositeUniqueConstraint[riot.StyleToken]
	Column string
}

var App_id_Mode_Token_Style = app_id_Mode_Token_Style{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.StyleToken]{
		Name: "App_id_Mode_Token_Style",
		Columns: []string{
			"app_id",
			"mode",
			"token",
			"style",
		},
	},
	Column: "App_id_Mode_Token_Style",
}

func (f app_id_Mode_Token_Style) EQ(app_id string, mode riot.ThemeModeType, token riot.TokenTypeType, style string) riot.UniquePredicate[riot.StyleToken] {
	return riot.UniquePredicate[riot.StyleToken]{
		Data: riot.PredicateData{
			Column:   "App_id_Mode_Token_Style",
			Operator: "AND",
			Value: map[string]any{
				"app_id": app_id,
				"mode":   mode,
				"token":  token,
				"style":  style,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "app_id", Operator: "=", Value: app_id},
				{Column: "mode", Operator: "=", Value: mode},
				{Column: "token", Operator: "=", Value: token},
				{Column: "style", Operator: "=", Value: style},
			},
		},
	}
}

type CreateInput = riot.StyleTokenCreate
type Create = riot.StyleTokenCreate

type CreateArgs = riot.StyleTokenCreateArgs
type CreateManyArgs = riot.StyleTokenCreateManyArgs
type CreateManyAndReturnArgs = riot.StyleTokenCreateManyAndReturnArgs

type CreateQuery = riot.StyleTokenCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.StyleToken, error)

type CreateManyQuery = riot.StyleTokenCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.StyleTokenCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.StyleToken, error)

type FindUniqueArgs = riot.StyleTokenFindUniqueArgs
type FindUniqueQuery = riot.StyleTokenFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.StyleToken, error)

type FindFirstArgs = riot.StyleTokenFindFirstArgs
type FindFirstQuery = riot.StyleTokenFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.StyleToken, error)

type FindManyArgs = riot.StyleTokenFindManyArgs
type FindManyQuery = riot.StyleTokenFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.StyleToken, error)

type CountArgs = riot.StyleTokenCountArgs
type CountQuery = riot.StyleTokenCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.StyleTokenUpdate

type Extension = riot.StyleTokenExtension

// ConflictUpdate creates a custom ConflictAction for StyleToken upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.StyleTokenUpsert)) *riot.ConflictAction {
	return riot.StyleTokenConflictUpdate(fn)
}
