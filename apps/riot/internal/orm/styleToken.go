package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// StyleToken represents the database model
type StyleToken struct {
	Style    string        `db:"style" json:"style"`
	Value    string        `db:"value" json:"value"`
	Mode     ThemeModeType `db:"mode" json:"mode"`
	Type     StyleTypeType `db:"type" json:"type"`
	Token    TokenTypeType `db:"token" json:"token"`
	App_id   string        `db:"app_id" json:"app_id"`
	Variable *string       `db:"variable" json:"variable,omitempty"`
	App      *App          `json:"app,omitempty"`
}

// StyleTokenCreate contains model input fields for StyleToken creation operations.
//
// Fields for StyleToken:
//
//	style    string    required
//	value    string    required
//	mode     ThemeMode required
//	type     StyleType required
//	token    TokenType required
//	app_id   string    required
//	variable string    optional
type StyleTokenCreate struct {
	Style    string        `json:"style"`
	Value    string        `json:"value"`
	Mode     ThemeModeType `json:"mode"`
	Type     StyleTypeType `json:"type"`
	Token    TokenTypeType `json:"token"`
	App_id   string        `json:"app_id"`
	Variable *string       `json:"variable"`
}

// colMask returns a bit mask of columns that are set
func (s *StyleTokenCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	mask |= 1 << 3
	mask |= 1 << 4
	mask |= 1 << 5
	if s.Variable != nil {
		mask |= 1 << 6
	}
	return mask
}

// StyleTokenUpdate contains model input fields for StyleToken update operations.
type StyleTokenUpdate struct {
	Style    *string        `json:"style"`
	Value    *string        `json:"value"`
	Mode     *ThemeModeType `json:"mode"`
	Type     *StyleTypeType `json:"type"`
	Token    *TokenTypeType `json:"token"`
	App_id   *string        `json:"app_id"`
	Variable *string        `json:"variable"`
}

func (u *StyleTokenUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Style != nil {
		cols = append(cols, "style")
		vals = append(vals, u.Style)
	}
	if u.Value != nil {
		cols = append(cols, "value")
		vals = append(vals, u.Value)
	}
	if u.Mode != nil {
		cols = append(cols, "mode")
		vals = append(vals, u.Mode)
	}
	if u.Type != nil {
		cols = append(cols, "type")
		vals = append(vals, u.Type)
	}
	if u.Token != nil {
		cols = append(cols, "token")
		vals = append(vals, u.Token)
	}
	if u.App_id != nil {
		cols = append(cols, "app_id")
		vals = append(vals, u.App_id)
	}
	if u.Variable != nil {
		cols = append(cols, "variable")
		vals = append(vals, u.Variable)
	}
	return cols, vals
}

func assignmentsToStyleTokenUpdate(assignments []FieldAssignment) (StyleTokenUpdate, error) {
	var input StyleTokenUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "style":
			if v, ok := a.Val.(string); ok {
				input.Style = &v
				errs.ValidateString("style", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Style = v
			} else {
				errs.Add("style", a.Val, "type", "field style must be of type string")
			}
		case "value":
			if v, ok := a.Val.(string); ok {
				input.Value = &v
				errs.ValidateString("value", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Value = v
			} else {
				errs.Add("value", a.Val, "type", "field value must be of type string")
			}
		case "mode":
			if v, ok := a.Val.(ThemeModeType); ok {
				input.Mode = &v
				if !v.IsValid() {
					errs.Add("mode", v, "enum", fmt.Sprintf("invalid enum value %q for field mode", v))
				}
			} else if v, ok := a.Val.(*ThemeModeType); ok {
				input.Mode = v
				if v != nil && !v.IsValid() {
					errs.Add("mode", *v, "enum", fmt.Sprintf("invalid enum value %q for field mode", *v))
				}
			} else {
				errs.Add("mode", a.Val, "type", "field mode must be of type ThemeModeType")
			}
		case "type":
			if v, ok := a.Val.(StyleTypeType); ok {
				input.Type = &v
				if !v.IsValid() {
					errs.Add("type", v, "enum", fmt.Sprintf("invalid enum value %q for field type", v))
				}
			} else if v, ok := a.Val.(*StyleTypeType); ok {
				input.Type = v
				if v != nil && !v.IsValid() {
					errs.Add("type", *v, "enum", fmt.Sprintf("invalid enum value %q for field type", *v))
				}
			} else {
				errs.Add("type", a.Val, "type", "field type must be of type StyleTypeType")
			}
		case "token":
			if v, ok := a.Val.(TokenTypeType); ok {
				input.Token = &v
				if !v.IsValid() {
					errs.Add("token", v, "enum", fmt.Sprintf("invalid enum value %q for field token", v))
				}
			} else if v, ok := a.Val.(*TokenTypeType); ok {
				input.Token = v
				if v != nil && !v.IsValid() {
					errs.Add("token", *v, "enum", fmt.Sprintf("invalid enum value %q for field token", *v))
				}
			} else {
				errs.Add("token", a.Val, "type", "field token must be of type TokenTypeType")
			}
		case "app_id":
			if v, ok := a.Val.(string); ok {
				input.App_id = &v
				errs.ValidateString("app_id", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.App_id = v
			} else {
				errs.Add("app_id", a.Val, "type", "field app_id must be of type string")
			}
		case "variable":
			if v, ok := a.Val.(string); ok {
				input.Variable = &v
				errs.ValidateString("variable", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Variable = v
			} else {
				errs.Add("variable", a.Val, "type", "field variable must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// StyleTokenSelect specifies which scalar and relation fields to select for StyleToken.
//
// Selectable fields:
//
//	-- Scalars --
//	style    (bool)
//	value    (bool)
//	mode     (bool)
//	type     (bool)
//	token    (bool)
//	app_id   (bool)
//	variable (bool)
//	-- Relations --
//	app      (App)
type StyleTokenSelect struct {
	Style    bool       `json:"style"`
	Value    bool       `json:"value"`
	Mode     bool       `json:"mode"`
	Type     bool       `json:"type"`
	Token    bool       `json:"token"`
	App_id   bool       `json:"app_id"`
	Variable bool       `json:"variable"`
	App      *AppSelect `json:"app,omitempty"`
}

var fullStyleTokenSelectVal = &StyleTokenSelect{
	Style:    true,
	Value:    true,
	Mode:     true,
	Type:     true,
	Token:    true,
	App_id:   true,
	Variable: true,
}

func fullStyleTokenSelect() *StyleTokenSelect {
	return fullStyleTokenSelectVal
}

func (s *StyleTokenSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Style || s.Value || s.Mode || s.Type || s.Token || s.App_id || s.Variable
}

func (s *StyleTokenSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type StyleTokenOmit struct {
	Style    bool `json:"style"`
	Value    bool `json:"value"`
	Mode     bool `json:"mode"`
	Type     bool `json:"type"`
	Token    bool `json:"token"`
	App_id   bool `json:"app_id"`
	Variable bool `json:"variable"`
}

type StyleTokenSelectQuery interface {
	GetRelationParams() (*StyleTokenSelect, *StyleTokenOmit, QueryParams[StyleToken])
}

func (s *StyleTokenSelect) GetRelationParams() (*StyleTokenSelect, *StyleTokenOmit, QueryParams[StyleToken]) {
	return s, nil, QueryParams[StyleToken]{}
}

type StyleTokenQueryBuilder struct {
	selects *StyleTokenSelect
	omits   *StyleTokenOmit
	where   []PredicateOf[StyleToken]
	take    *int
	skip    *int
	orderBy []OrderBy[StyleToken]
	cursor  UniquePredicate[StyleToken]
}

func (b *StyleTokenQueryBuilder) Where(preds ...PredicateOf[StyleToken]) *StyleTokenQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *StyleTokenQueryBuilder) Take(limit int) *StyleTokenQueryBuilder {
	b.take = &limit
	return b
}

func (b *StyleTokenQueryBuilder) Skip(offset int) *StyleTokenQueryBuilder {
	b.skip = &offset
	return b
}

func (b *StyleTokenQueryBuilder) OrderBy(orders ...OrderBy[StyleToken]) *StyleTokenQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *StyleTokenQueryBuilder) Cursor(where UniquePredicate[StyleToken]) *StyleTokenQueryBuilder {
	b.cursor = where
	return b
}

func (b *StyleTokenQueryBuilder) Select(s StyleTokenSelect) *StyleTokenQueryBuilder {
	b.selects = &s
	return b
}

func (b *StyleTokenQueryBuilder) Omit(o StyleTokenOmit) *StyleTokenQueryBuilder {
	b.omits = &o
	return b
}

func (b *StyleTokenQueryBuilder) GetRelationParams() (*StyleTokenSelect, *StyleTokenOmit, QueryParams[StyleToken]) {
	if b == nil {
		return nil, nil, QueryParams[StyleToken]{}
	}
	return b.selects, b.omits, QueryParams[StyleToken]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// StyleTokenCreateArgs is the input argument passed to StyleToken Create extension hooks.
//
// Fields for StyleToken:
//
//	style    string    required
//	value    string    required
//	mode     ThemeMode required
//	type     StyleType required
//	token    TokenType required
//	app_id   string    required
//	variable string    optional
//
// Relations for StyleToken:
//
//	app      (App)
type StyleTokenCreateArgs struct {
	// Data contains the model fields to insert.
	Data *StyleTokenCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *StyleTokenSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// StyleTokenCreateManyArgs is the input argument passed to StyleToken CreateMany extension hooks.
//
// Fields for StyleToken:
//
//	style    string    required
//	value    string    required
//	mode     ThemeMode required
//	type     StyleType required
//	token    TokenType required
//	app_id   string    required
//	variable string    optional
type StyleTokenCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*StyleTokenCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *StyleTokenCreateManyArgs) AppendData(builders ...*StyleTokenCreateBuilder) *StyleTokenCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToStyleTokenCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// StyleTokenCreateManyAndReturnArgs is the input argument passed to StyleToken CreateManyAndReturn extension hooks.
//
// Fields for StyleToken:
//
//	style    string    required
//	value    string    required
//	mode     ThemeMode required
//	type     StyleType required
//	token    TokenType required
//	app_id   string    required
//	variable string    optional
//
// Relations for StyleToken:
//
//	app      (App)
type StyleTokenCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*StyleTokenCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *StyleTokenSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *StyleTokenCreateManyAndReturnArgs) AppendData(builders ...*StyleTokenCreateBuilder) *StyleTokenCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToStyleTokenCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// StyleTokenFindUniqueArgs is the input argument passed to StyleToken FindUnique extension hooks.
type StyleTokenFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[StyleToken]
	// Select specifies which scalar and relation fields to select and return.
	Select *StyleTokenSelect
}

func (a *StyleTokenFindUniqueArgs) SetWhere(unique UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *StyleTokenFindUniqueArgs {
	a.Where = make([]PredicateOf[StyleToken], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// StyleTokenFindFirstArgs is the input argument passed to StyleToken FindFirst extension hooks.
type StyleTokenFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[StyleToken]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[StyleToken]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *StyleTokenSelect
}

func (a *StyleTokenFindFirstArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenFindFirstArgs {
	a.Where = preds
	return a
}

func (a *StyleTokenFindFirstArgs) SetOrderBy(orders ...OrderBy[StyleToken]) *StyleTokenFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *StyleTokenFindFirstArgs) SetCursor(cursor UniquePredicate[StyleToken]) *StyleTokenFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *StyleTokenFindFirstArgs) SetSkip(n int) *StyleTokenFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *StyleTokenFindFirstArgs) SetTake(n int) *StyleTokenFindFirstArgs {
	a.Take = &n
	return a
}

// StyleTokenFindManyArgs is the input argument passed to StyleToken FindMany extension hooks.
type StyleTokenFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[StyleToken]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[StyleToken]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *StyleTokenSelect
}

func (a *StyleTokenFindManyArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenFindManyArgs {
	a.Where = preds
	return a
}

func (a *StyleTokenFindManyArgs) SetOrderBy(orders ...OrderBy[StyleToken]) *StyleTokenFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *StyleTokenFindManyArgs) SetCursor(cursor UniquePredicate[StyleToken]) *StyleTokenFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *StyleTokenFindManyArgs) SetSkip(n int) *StyleTokenFindManyArgs {
	a.Skip = &n
	return a
}

func (a *StyleTokenFindManyArgs) SetTake(n int) *StyleTokenFindManyArgs {
	a.Take = &n
	return a
}

// StyleTokenCountArgs is the input argument passed to StyleToken Count extension hooks.
type StyleTokenCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *StyleTokenCountArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenCountArgs {
	a.Where = preds
	return a
}

func (a *StyleTokenCountArgs) SetSkip(n int) *StyleTokenCountArgs {
	a.Skip = &n
	return a
}

func (a *StyleTokenCountArgs) SetTake(n int) *StyleTokenCountArgs {
	a.Take = &n
	return a
}

// StyleTokenDeleteArgs is the input argument passed to StyleToken Delete extension hooks.
type StyleTokenDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[StyleToken]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *StyleTokenSelect
}

func (a *StyleTokenDeleteArgs) SetWhere(unique UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *StyleTokenDeleteArgs {
	a.Where = make([]PredicateOf[StyleToken], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// StyleTokenDeleteManyArgs is the input argument passed to StyleToken DeleteMany extension hooks.
type StyleTokenDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
}

func (a *StyleTokenDeleteManyArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenDeleteManyArgs {
	a.Where = preds
	return a
}

// StyleTokenUpdateArgs is the input argument passed to StyleToken Update extension hooks.
type StyleTokenUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[StyleToken]
	// Data contains the model fields to update.
	Data *StyleTokenUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *StyleTokenSelect
}

func (a *StyleTokenUpdateArgs) SetWhere(unique UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *StyleTokenUpdateArgs {
	a.Where = make([]PredicateOf[StyleToken], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// StyleTokenUpdateManyArgs is the input argument passed to StyleToken UpdateMany extension hooks.
type StyleTokenUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
	// Data contains the model fields to update.
	Data *StyleTokenUpdate
}

func (a *StyleTokenUpdateManyArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenUpdateManyArgs {
	a.Where = preds
	return a
}

// StyleTokenUpdateManyAndReturnArgs is the input argument passed to StyleToken UpdateManyAndReturn extension hooks.
type StyleTokenUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[StyleToken]
	// Data contains the model fields to update.
	Data *StyleTokenUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *StyleTokenSelect
}

func (a *StyleTokenUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[StyleToken]) *StyleTokenUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type StyleTokenCreateQuery = func(ctx context.Context, args *StyleTokenCreateArgs) (*StyleToken, error)
type StyleTokenCreateManyQuery = func(ctx context.Context, args *StyleTokenCreateManyArgs) (int64, error)
type StyleTokenCreateManyAndReturnQuery = func(ctx context.Context, args *StyleTokenCreateManyAndReturnArgs) ([]*StyleToken, error)
type StyleTokenFindUniqueQuery = func(ctx context.Context, args *StyleTokenFindUniqueArgs) (*StyleToken, error)
type StyleTokenFindFirstQuery = func(ctx context.Context, args *StyleTokenFindFirstArgs) (*StyleToken, error)
type StyleTokenFindManyQuery = func(ctx context.Context, args *StyleTokenFindManyArgs) ([]*StyleToken, error)
type StyleTokenDeleteQuery = func(ctx context.Context, args *StyleTokenDeleteArgs) (*StyleToken, error)
type StyleTokenDeleteManyQuery = func(ctx context.Context, args *StyleTokenDeleteManyArgs) (int64, error)
type StyleTokenCountQuery = func(ctx context.Context, args *StyleTokenCountArgs) (int64, error)
type StyleTokenUpdateQuery = func(ctx context.Context, args *StyleTokenUpdateArgs) (*StyleToken, error)
type StyleTokenUpdateManyQuery = func(ctx context.Context, args *StyleTokenUpdateManyArgs) (int64, error)
type StyleTokenUpdateManyAndReturnQuery = func(ctx context.Context, args *StyleTokenUpdateManyAndReturnArgs) ([]*StyleToken, error)

type StyleTokenExtension struct {
	Create              func(ctx context.Context, args *StyleTokenCreateArgs, next StyleTokenCreateQuery) (*StyleToken, error)
	CreateMany          func(ctx context.Context, args *StyleTokenCreateManyArgs, next StyleTokenCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *StyleTokenCreateManyAndReturnArgs, next StyleTokenCreateManyAndReturnQuery) ([]*StyleToken, error)
	FindUnique          func(ctx context.Context, args *StyleTokenFindUniqueArgs, next StyleTokenFindUniqueQuery) (*StyleToken, error)
	FindFirst           func(ctx context.Context, args *StyleTokenFindFirstArgs, next StyleTokenFindFirstQuery) (*StyleToken, error)
	FindMany            func(ctx context.Context, args *StyleTokenFindManyArgs, next StyleTokenFindManyQuery) ([]*StyleToken, error)
	Delete              func(ctx context.Context, args *StyleTokenDeleteArgs, next StyleTokenDeleteQuery) (*StyleToken, error)
	DeleteMany          func(ctx context.Context, args *StyleTokenDeleteManyArgs, next StyleTokenDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *StyleTokenCountArgs, next StyleTokenCountQuery) (int64, error)
	Update              func(ctx context.Context, args *StyleTokenUpdateArgs, next StyleTokenUpdateQuery) (*StyleToken, error)
	UpdateMany          func(ctx context.Context, args *StyleTokenUpdateManyArgs, next StyleTokenUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *StyleTokenUpdateManyAndReturnArgs, next StyleTokenUpdateManyAndReturnQuery) ([]*StyleToken, error)
}

type StyleTokenDelegate struct {
	client     *Queries
	extensions []StyleTokenExtension
}

func (d *StyleTokenDelegate) Use(exts ...StyleTokenExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *StyleToken) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "style":
			targets[i] = &m.Style
		case "value":
			targets[i] = &m.Value
		case "mode":
			targets[i] = &m.Mode
		case "type":
			targets[i] = &m.Type
		case "token":
			targets[i] = &m.Token
		case "app_id":
			targets[i] = &m.App_id
		case "variable":
			targets[i] = &m.Variable
		}
	}
	return targets
}

var styleTokenDefaultCols = []string{
	"style",
	"value",
	"mode",
	"type",
	"token",
	"app_id",
	"variable",
}

var styleTokenPKCols = []string{
	"app_id",
	"mode",
	"token",
	"style",
}

var styleTokenUniqueCols = []string{}

func selectStyleTokenCols(selects *StyleTokenSelect, omits *StyleTokenOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return styleTokenDefaultCols
	}

	anySelected := selects != nil && (selects.Style || selects.Value || selects.Mode || selects.Type || selects.Token || selects.App_id || selects.Variable || selects.App != nil)

	specs := []colSpec{
		{"style", selects != nil && selects.Style, omits != nil && omits.Style, false},
		{"value", selects != nil && selects.Value, omits != nil && omits.Value, false},
		{"mode", selects != nil && selects.Mode, omits != nil && omits.Mode, false},
		{"type", selects != nil && selects.Type, omits != nil && omits.Type, false},
		{"token", selects != nil && selects.Token, omits != nil && omits.Token, false},
		{"app_id", selects != nil && selects.App_id, omits != nil && omits.App_id, selects != nil && selects.App != nil},
		{"variable", selects != nil && selects.Variable, omits != nil && omits.Variable, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *StyleTokenSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.App != nil
}

type StyleTokenUpsertBuilder struct {
	*CreateBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]
}

func (b *StyleTokenUpsertBuilder) Select(s StyleTokenSelect) *StyleTokenUpsertBuilder {
	b.selects = &s
	return b
}

func (b *StyleTokenUpsertBuilder) Omit(o StyleTokenOmit) *StyleTokenUpsertBuilder {
	b.omits = &o
	return b
}

type StyleTokenCreateBuilder struct {
	*CreateBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]
}

func (b *StyleTokenCreateBuilder) Select(s StyleTokenSelect) *StyleTokenCreateBuilder {
	b.selects = &s
	return b
}

func (b *StyleTokenCreateBuilder) Omit(o StyleTokenOmit) *StyleTokenCreateBuilder {
	b.omits = &o
	return b
}

func (b *StyleTokenCreateBuilder) OnConflict(target UniqueConstraintTarget) *StyleTokenConflictBuilder[StyleTokenUpsertBuilder] {
	upsertBuilder := &StyleTokenUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &StyleTokenConflictBuilder[StyleTokenUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *StyleTokenCreateBuilder) SetStyle(v string) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "style", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetValue(v string) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "value", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetMode(v ThemeModeType) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "mode", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetType(v StyleTypeType) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "type", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetToken(v TokenTypeType) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "token", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetApp_id(v string) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}
func (b *StyleTokenCreateBuilder) SetVariable(v string) *StyleTokenCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "variable", Val: v})
	return b
}

func (b *StyleTokenCreateBuilder) Assignments(assignments ...FieldAssignmentOf[StyleToken]) *StyleTokenCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *StyleTokenDelegate) Create() *StyleTokenCreateBuilder {
	return &StyleTokenCreateBuilder{
		CreateBuilder: &CreateBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedStyleTokenStyle    uint64 = 1 << 0
	providedStyleTokenValue    uint64 = 1 << 1
	providedStyleTokenMode     uint64 = 1 << 2
	providedStyleTokenType     uint64 = 1 << 3
	providedStyleTokenToken    uint64 = 1 << 4
	providedStyleTokenApp_id   uint64 = 1 << 5
	providedStyleTokenVariable uint64 = 1 << 6
)

func assignmentsToStyleTokenCreate(assignments []FieldAssignment) (StyleTokenCreate, error) {
	var input StyleTokenCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "style":
			provided |= providedStyleTokenStyle
			if v, ok := a.Val.(string); ok {
				input.Style = v
				errs.ValidateString("style", v, true, 0, false, false)
			} else {
				errs.Add("style", a.Val, "type", "field style must be of type string")
			}
		case "value":
			provided |= providedStyleTokenValue
			if v, ok := a.Val.(string); ok {
				input.Value = v
				errs.ValidateString("value", v, true, 0, false, false)
			} else {
				errs.Add("value", a.Val, "type", "field value must be of type string")
			}
		case "mode":
			provided |= providedStyleTokenMode
			if v, ok := a.Val.(ThemeModeType); ok {
				input.Mode = v
				if !v.IsValid() {
					errs.Add("mode", v, "enum", fmt.Sprintf("invalid enum value %q for field mode", v))
				}
			} else {
				errs.Add("mode", a.Val, "type", "field mode must be of type ThemeModeType")
			}
		case "type":
			provided |= providedStyleTokenType
			if v, ok := a.Val.(StyleTypeType); ok {
				input.Type = v
				if !v.IsValid() {
					errs.Add("type", v, "enum", fmt.Sprintf("invalid enum value %q for field type", v))
				}
			} else {
				errs.Add("type", a.Val, "type", "field type must be of type StyleTypeType")
			}
		case "token":
			provided |= providedStyleTokenToken
			if v, ok := a.Val.(TokenTypeType); ok {
				input.Token = v
				if !v.IsValid() {
					errs.Add("token", v, "enum", fmt.Sprintf("invalid enum value %q for field token", v))
				}
			} else {
				errs.Add("token", a.Val, "type", "field token must be of type TokenTypeType")
			}
		case "app_id":
			provided |= providedStyleTokenApp_id
			if v, ok := a.Val.(string); ok {
				input.App_id = v
				errs.ValidateString("app_id", v, true, 0, false, false)
			} else {
				errs.Add("app_id", a.Val, "type", "field app_id must be of type string")
			}
		case "variable":
			provided |= providedStyleTokenVariable
			if v, ok := a.Val.(string); ok {
				input.Variable = &v
				errs.ValidateString("variable", v, false, 0, false, false)
			} else {
				errs.Add("variable", a.Val, "type", "field variable must be of type string")
			}
		}
	}
	if provided&providedStyleTokenStyle == 0 {
		errs.Add("style", "", "required", "field Style is required")
	}
	if provided&providedStyleTokenValue == 0 {
		errs.Add("value", "", "required", "field Value is required")
	}
	if provided&providedStyleTokenMode == 0 {
		errs.Add("mode", nil, "required", "field Mode is required")
	}
	if provided&providedStyleTokenType == 0 {
		errs.Add("type", nil, "required", "field Type is required")
	}
	if provided&providedStyleTokenToken == 0 {
		errs.Add("token", nil, "required", "field Token is required")
	}
	if provided&providedStyleTokenApp_id == 0 {
		errs.Add("app_id", "", "required", "field App_id is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *StyleTokenCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 7)
	vals = make([]any, 0, 7)
	cols = append(cols, "style")
	vals = append(vals, s.Style)
	cols = append(cols, "value")
	vals = append(vals, s.Value)
	cols = append(cols, "mode")
	vals = append(vals, s.Mode)
	cols = append(cols, "type")
	vals = append(vals, s.Type)
	cols = append(cols, "token")
	vals = append(vals, s.Token)
	cols = append(cols, "app_id")
	vals = append(vals, s.App_id)
	if s.Variable != nil {
		cols = append(cols, "variable")
		vals = append(vals, *s.Variable)
	}
	return
}

func partitionStyleTokenInputs(dialect Dialect, inputs []*StyleTokenCreate) [][]*StyleTokenCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*StyleTokenCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*StyleTokenCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*StyleTokenCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*StyleTokenCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*StyleTokenCreate{inputs}
}

func (d *StyleTokenDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *StyleTokenSelect, omits *StyleTokenOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*StyleToken, error) {
	input, err := assignmentsToStyleTokenCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectStyleTokenCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *StyleTokenCreateArgs) (*StyleToken, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectStyleTokenCols(a.Select, omits)
		return d.runCreate(c, cCols, cVals, cReturningCols, a.Select, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Create != nil {
			return ext.Create(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, a *StyleTokenCreateArgs) (*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type StyleTokenCreateManyBuilder struct {
	*CreateManyBuilder[StyleToken]
}

func (b *StyleTokenCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *StyleTokenConflictBuilder[StyleTokenCreateManyBuilder] {
	return &StyleTokenConflictBuilder[StyleTokenCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type StyleTokenCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]
}

func (b *StyleTokenCreateManyAndReturnBuilder) Select(s StyleTokenSelect) *StyleTokenCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *StyleTokenCreateManyAndReturnBuilder) Omit(o StyleTokenOmit) *StyleTokenCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *StyleTokenCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *StyleTokenConflictBuilder[StyleTokenCreateManyAndReturnBuilder] {
	return &StyleTokenConflictBuilder[StyleTokenCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToStyleTokenRecordInputs(builders []*StyleTokenCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *StyleTokenDelegate) CreateMany(builders ...*StyleTokenCreateBuilder) *StyleTokenCreateManyBuilder {
	return &StyleTokenCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[StyleToken]{
			records:  createBuildersToStyleTokenRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *StyleTokenDelegate) CreateManyAndReturn(builders ...*StyleTokenCreateBuilder) *StyleTokenCreateManyAndReturnBuilder {
	return &StyleTokenCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
			records:  createBuildersToStyleTokenRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToStyleTokenCreateInputs(records []RecordInput) ([]*StyleTokenCreate, error) {
	structs := make([]StyleTokenCreate, len(records))
	inputs := make([]*StyleTokenCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToStyleTokenCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *StyleTokenDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToStyleTokenCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &StyleTokenCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *StyleTokenCreateManyArgs) (int64, error) {
		return d.runCreateMany(c, a.Data, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateMany != nil {
			return ext.CreateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, a *StyleTokenCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *StyleTokenSelect, omits *StyleTokenOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*StyleToken, error) {
	inputs, err := recordsToStyleTokenCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *StyleTokenCreateManyAndReturnArgs) ([]*StyleToken, error) {
		return d.runCreateManyAndReturn(c, a.Data, a.Select, omits, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateManyAndReturn != nil {
			return ext.CreateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, a *StyleTokenCreateManyAndReturnArgs) ([]*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *StyleTokenSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*StyleToken, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *StyleToken
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.StyleToken.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.StyleToken.loadRelations(ctx, []*StyleToken{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "StyleToken", cols, returningCols, styleTokenPKCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}

		if !rows.Next() {
			err := rows.Err()
			rows.Close()
			if err != nil {
				return nil, TranslateDBError(err)
			}
			return nil, nil
		}

		var res StyleToken
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, styleTokenPKCols)
}

func (d *StyleTokenDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*StyleToken, error) {
	result, err := d.client.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var pkVals []any
	for _, pkCol := range pkCols {
		var val any
		for i, c := range cols {
			if c == pkCol {
				val = vals[i]
				break
			}
		}
		if val == nil && len(pkCols) == 1 {
			lastID, err := result.LastInsertId()
			if err != nil {
				return nil, TranslateDBError(err)
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("StyleToken") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "StyleToken")
	selectSb.WriteString(" WHERE ")
	for i, pkCol := range pkCols {
		if i > 0 {
			selectSb.WriteString(" AND ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, pkCol)
		selectSb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&selectSb, i+1)
	}

	rows, err := d.client.query(ctx, selectSb.String(), pkVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res StyleToken
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *StyleTokenDelegate) buildBulkInsertSQL(q *Queries, batch []*StyleTokenCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 7)
	for i, c := range styleTokenDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "StyleToken")
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES ")

	paramIdx := paramStartIdx
	for ri, input := range batch {
		if ri > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, col := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			writeDefault := false
			switch col {
			case "style":
				vals = append(vals, input.Style)
			case "value":
				vals = append(vals, input.Value)
			case "mode":
				vals = append(vals, input.Mode)
			case "type":
				vals = append(vals, input.Type)
			case "token":
				vals = append(vals, input.Token)
			case "app_id":
				vals = append(vals, input.App_id)
			case "variable":
				if input.Variable != nil {
					vals = append(vals, *input.Variable)
				} else {
					writeDefault = true
				}
			}
			if writeDefault {
				sb.WriteString("DEFAULT")
			} else {
				q.dialect.WritePlaceholder(&sb, paramIdx)
				paramIdx++
			}
		}
		sb.WriteString(")")
	}
	queryStr = sb.String()
	return cols, vals, queryStr
}

func applyStyleTokenConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
	var conflictCols []string
	if conflictTarget != nil {
		conflictCols = conflictTarget.UniqueColumns()
	}
	var nonConflictCols []string
	if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
		nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
	}
	clause, clauseArgs := dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
	queryStr += clause
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}
	return queryStr, vals
}

func scanStyleTokenRows(rows *sql.Rows, returningCols []string) ([]*StyleToken, error) {
	var records []*StyleToken
	for rows.Next() {
		var res StyleToken
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			rows.Close()
			return nil, TranslateDBError(err)
		}
		records = append(records, &res)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, TranslateDBError(err)
	}
	rows.Close()
	return records, nil
}

func (d *StyleTokenDelegate) runCreateMany(ctx context.Context, inputs []*StyleTokenCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionStyleTokenInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyStyleTokenConflictClause(d.client.dialect, queryStr, vals, cols, styleTokenPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += affected
	}
	return count, nil
}

func (d *StyleTokenDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*StyleTokenCreate,
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*StyleToken, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionStyleTokenInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*StyleToken
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.StyleToken.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.StyleToken.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectStyleTokenCols(selects, omits, styleTokenPKCols...)
	recordsOut := make([]*StyleToken, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyStyleTokenConflictClause(d.client.dialect, queryStr, vals, cols, styleTokenPKCols, conflictTarget, conflictAction)

		if len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				d.client.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
		}

		rows, err := d.client.query(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		scanned, err := scanStyleTokenRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, recordsOut, selects); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

func (d *StyleTokenDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*StyleTokenCreate,
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*StyleToken, error) {
	batches := partitionStyleTokenInputs(d.client.dialect, inputs)
	returningCols := selectStyleTokenCols(selects, omits, styleTokenPKCols...)
	recordsOut := make([]*StyleToken, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyStyleTokenConflictClause(d.client.dialect, queryStr, vals, cols, styleTokenPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("StyleToken") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "StyleToken")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, styleTokenPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, styleTokenPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanStyleTokenRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, recordsOut, selects); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

type StyleTokenConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *StyleTokenConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *StyleTokenConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *StyleTokenConflictBuilder[B]) Update(fn func(u *StyleTokenUpsert)) *B {
	cb.setAction(*StyleTokenConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// StyleTokenConflictUpdate creates a custom ConflictAction for StyleToken upsert conflicts.
func StyleTokenConflictUpdate(fn func(u *StyleTokenUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newStyleTokenUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type StyleTokenUpsert struct {
	Style    fieldUpsert[string]
	Value    fieldUpsert[string]
	Mode     fieldUpsert[ThemeModeType]
	Type     fieldUpsert[StyleTypeType]
	Token    fieldUpsert[TokenTypeType]
	App_id   fieldUpsert[string]
	Variable fieldUpsert[*string]
}

func newStyleTokenUpsert(up *ConflictUpdate) *StyleTokenUpsert {
	return &StyleTokenUpsert{
		Style:    fieldUpsert[string]{column: "style", update: up},
		Value:    fieldUpsert[string]{column: "value", update: up},
		Mode:     fieldUpsert[ThemeModeType]{column: "mode", update: up},
		Type:     fieldUpsert[StyleTypeType]{column: "type", update: up},
		Token:    fieldUpsert[TokenTypeType]{column: "token", update: up},
		App_id:   fieldUpsert[string]{column: "app_id", update: up},
		Variable: fieldUpsert[*string]{column: "variable", update: up},
	}
}

type StyleTokenUpdateBuilder struct {
	*UpdateBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]
}

type StyleTokenUpdateManyBuilder struct {
	*UpdateManyBuilder[StyleToken]
}

type StyleTokenUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]
}

func (b *StyleTokenUpdateBuilder) SetStyle(v string) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "style", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetStyle(v string) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "style", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetStyle(v string) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "style", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetValue(v string) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "value", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetValue(v string) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "value", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetValue(v string) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "value", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetMode(v ThemeModeType) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "mode", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetMode(v ThemeModeType) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "mode", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetMode(v ThemeModeType) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "mode", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetType(v StyleTypeType) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "type", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetType(v StyleTypeType) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "type", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetType(v StyleTypeType) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "type", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetToken(v TokenTypeType) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "token", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetToken(v TokenTypeType) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "token", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetToken(v TokenTypeType) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "token", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetApp_id(v string) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetApp_id(v string) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetApp_id(v string) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}
func (b *StyleTokenUpdateBuilder) SetVariable(v string) *StyleTokenUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "variable", Val: v})
	return b
}

func (b *StyleTokenUpdateManyBuilder) SetVariable(v string) *StyleTokenUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "variable", Val: v})
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) SetVariable(v string) *StyleTokenUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "variable", Val: v})
	return b
}

func (b *StyleTokenUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[StyleToken]) *StyleTokenUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *StyleTokenUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[StyleToken]) *StyleTokenUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *StyleTokenUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[StyleToken]) *StyleTokenUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *StyleTokenDelegate) Update(where UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *StyleTokenUpdateBuilder {
	return &StyleTokenUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *StyleTokenDelegate) UpdateMany(preds ...PredicateOf[StyleToken]) *StyleTokenUpdateManyBuilder {
	return &StyleTokenUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[StyleToken]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *StyleTokenDelegate) UpdateManyAndReturn(preds ...PredicateOf[StyleToken]) *StyleTokenUpdateManyAndReturnBuilder {
	return &StyleTokenUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *StyleTokenDelegate) buildUpdateSQL(preds []PredicateOf[StyleToken], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "StyleToken")
	sb.WriteString(" SET ")

	setVals := make([]any, 0, len(cols)+len(predVals))
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
		sb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&sb, i+1)
		setVals = append(setVals, vals[i])
	}

	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
		setVals = append(setVals, predVals...)
	}

	if len(returningCols) > 0 && d.client.dialect.SupportsUpdateReturning {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&sb, col)
		}
	}

	return sb.String(), setVals
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func (d *StyleTokenDelegate) executeUpdate(ctx context.Context, where UniquePredicate[StyleToken], additional []PredicateOf[StyleToken], assignments []FieldAssignment, selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	allWhere := make([]PredicateOf[StyleToken], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToStyleTokenUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *StyleTokenUpdateArgs) (*StyleToken, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdate(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Update != nil {
			return ext.Update(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Update != nil {
			next, hook := curr, ext.Update
			curr = func(c context.Context, a *StyleTokenUpdateArgs) (*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runUpdate(ctx context.Context, preds []PredicateOf[StyleToken], cols []string, vals []any, selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	if len(cols) == 0 {
		return d.runFindUnique(ctx, preds, selects, omits)
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res *StyleToken
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.StyleToken.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.StyleToken.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectStyleTokenCols(selects, omits, styleTokenPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "StyleToken"}
	}

	var res StyleToken
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*StyleToken{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *StyleTokenDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[StyleToken], cols []string, vals []any) (int64, error) {
	if len(cols) == 0 {
		return 0, nil
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	query, setVals := d.buildUpdateSQL(preds, cols, vals, nil)
	result, err := d.client.exec(ctx, query, setVals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *StyleTokenDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[StyleToken], cols []string, vals []any, selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "StyleToken"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *StyleTokenDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[StyleToken], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToStyleTokenUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &StyleTokenUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *StyleTokenUpdateManyArgs) (int64, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.execUpdateStmt(c, a.Where, extCols, extVals)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateMany != nil {
			return ext.UpdateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateMany != nil {
			next, hook := curr, ext.UpdateMany
			curr = func(c context.Context, a *StyleTokenUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *StyleTokenDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[StyleToken], assignments []FieldAssignment, selects *StyleTokenSelect, omits *StyleTokenOmit) ([]*StyleToken, error) {
	input, err := assignmentsToStyleTokenUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *StyleTokenUpdateManyAndReturnArgs) ([]*StyleToken, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdateManyAndReturn(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateManyAndReturn != nil {
			return ext.UpdateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateManyAndReturn != nil {
			next, hook := curr, ext.UpdateManyAndReturn
			curr = func(c context.Context, a *StyleTokenUpdateManyAndReturnArgs) ([]*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[StyleToken], cols []string, vals []any, selects *StyleTokenSelect, omits *StyleTokenOmit) ([]*StyleToken, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[StyleToken]{Where: preds}, selects, omits)
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res []*StyleToken
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.StyleToken.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.StyleToken.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectStyleTokenCols(selects, omits, styleTokenPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanStyleTokenRows(rows, returningCols)
	if err != nil {
		return nil, err
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, scanned, selects); err != nil {
			return nil, err
		}
	}

	return scanned, nil
}

func (d *StyleTokenDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[StyleToken], cols []string, vals []any, selects *StyleTokenSelect, omits *StyleTokenOmit) ([]*StyleToken, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*StyleToken{}, nil
	}
	return d.runFindMany(ctx, QueryParams[StyleToken]{Where: preds}, selects, omits)
}
func (d *StyleTokenDelegate) FindUnique(where UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *FindUniqueBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit] {
	return &FindUniqueBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *StyleTokenDelegate) FindFirst(preds ...PredicateOf[StyleToken]) *FindFirstBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit] {
	return &FindFirstBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *StyleTokenDelegate) FindMany(preds ...PredicateOf[StyleToken]) *FindManyBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit] {
	return &FindManyBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *StyleTokenDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[StyleToken], additional []PredicateOf[StyleToken], selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	allWhere := make([]PredicateOf[StyleToken], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *StyleTokenFindUniqueArgs) (*StyleToken, error) {
		return d.runFindUnique(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindUnique != nil {
			return ext.FindUnique(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, a *StyleTokenFindUniqueArgs) (*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[StyleToken],
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
) (*StyleToken, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *StyleTokenFindFirstArgs) (*StyleToken, error) {
		return d.runFindFirst(c, QueryParams[StyleToken]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindFirst != nil {
			return ext.FindFirst(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, a *StyleTokenFindFirstArgs) (*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[StyleToken],
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
) ([]*StyleToken, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *StyleTokenFindManyArgs) ([]*StyleToken, error) {
		return d.runFindMany(c, QueryParams[StyleToken]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindMany != nil {
			return ext.FindMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, a *StyleTokenFindManyArgs) ([]*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runFindUnique(ctx context.Context, where []PredicateOf[StyleToken], selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectStyleTokenCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*StyleToken{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *StyleTokenDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[StyleToken],
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
) (*StyleToken, error) {
	one := 1
	negOne := -1
	if params.Take != nil && *params.Take < 0 {
		params.Take = &negOne
	} else {
		params.Take = &one
	}

	results, err := d.runFindMany(ctx, params, selects, omits)
	if err != nil || len(results) == 0 {
		return nil, err
	}
	return results[0], nil
}

func (d *StyleTokenDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[StyleToken],
	selects *StyleTokenSelect,
	omits *StyleTokenOmit,
) ([]*StyleToken, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, styleTokenPKCols, styleTokenUniqueCols, "StyleToken", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, styleTokenPKCols, styleTokenUniqueCols, isCursorQuery, params.Take)
	returningCols := selectStyleTokenCols(selects, omits)

	results, err := d.queryMany(ctx, whereClause, orderByClause, vals, returningCols, params.Take, params.Skip)
	if err != nil || len(results) == 0 {
		return results, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, results, selects); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *StyleTokenDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*StyleToken, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "StyleToken", returningCols, whereClause, orderByClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
				return nil, nil
			}
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res StyleToken
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *StyleTokenDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*StyleToken, error) {
	query := buildSelectSQL(d.client, "StyleToken", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*StyleToken, 0)
	for rows.Next() {
		var res StyleToken
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, TranslateDBError(err)
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, TranslateDBError(err)
	}
	if take != nil && *take < 0 {
		reverseSlice(results)
	}
	return results, nil
}
func (d *StyleTokenDelegate) DeleteMany(preds ...PredicateOf[StyleToken]) *DeleteManyBuilder[StyleToken] {
	return &DeleteManyBuilder[StyleToken]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *StyleTokenDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[StyleToken]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &StyleTokenDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *StyleTokenDeleteManyArgs) (int64, error) {
		return d.runDeleteMany(c, a.Where)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.DeleteMany != nil {
			return ext.DeleteMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, a *StyleTokenDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[StyleToken]) (int64, error) {
	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, preds)

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "StyleToken")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	result, err := d.client.exec(ctx, sb.String(), vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *StyleTokenDelegate) Delete(where UniquePredicate[StyleToken], additional ...PredicateOf[StyleToken]) *DeleteBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit] {
	return &DeleteBuilder[StyleToken, StyleTokenSelect, StyleTokenOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *StyleTokenDelegate) executeDelete(ctx context.Context, where UniquePredicate[StyleToken], additional []PredicateOf[StyleToken], selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	allWhere := make([]PredicateOf[StyleToken], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullStyleTokenSelect()
	}

	args := &StyleTokenDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *StyleTokenDeleteArgs) (*StyleToken, error) {
		return d.runDelete(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Delete != nil {
			return ext.Delete(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, a *StyleTokenDeleteArgs) (*StyleToken, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runDelete(ctx context.Context, where []PredicateOf[StyleToken], selects *StyleTokenSelect, omits *StyleTokenOmit) (*StyleToken, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectStyleTokenCols(selects, omits, styleTokenPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *StyleToken
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.StyleToken.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "StyleToken"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "StyleToken")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[StyleToken]
			pkPreds = append(pkPreds, Predicate[StyleToken]{
				Data: PredicateData{
					Column:   "app_id",
					Operator: "=",
					Value:    res.App_id,
				},
			})
			pkPreds = append(pkPreds, Predicate[StyleToken]{
				Data: PredicateData{
					Column:   "mode",
					Operator: "=",
					Value:    res.Mode,
				},
			})
			pkPreds = append(pkPreds, Predicate[StyleToken]{
				Data: PredicateData{
					Column:   "token",
					Operator: "=",
					Value:    res.Token,
				},
			})
			pkPreds = append(pkPreds, Predicate[StyleToken]{
				Data: PredicateData{
					Column:   "style",
					Operator: "=",
					Value:    res.Style,
				},
			})

			whereClause, vals, _ := CompilePredicates(txQ.dialect, pkPreds)
			deleteSb.WriteString(whereClause)

			_, err = txQ.exec(ctx, deleteSb.String(), vals...)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Dialect supports RETURNING, and no relations need loading: run direct DELETE ... RETURNING
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "StyleToken")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	sb.WriteString(" RETURNING ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
	}

	rows, err := d.client.query(ctx, sb.String(), vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "StyleToken"}
	}

	var row StyleToken
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *StyleTokenDelegate) Count(preds ...PredicateOf[StyleToken]) *CountBuilder[StyleToken] {
	return &CountBuilder[StyleToken]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *StyleTokenDelegate) executeCount(ctx context.Context, params QueryParams[StyleToken]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &StyleTokenCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *StyleTokenCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[StyleToken]{
			Where: a.Where,
			Skip:  a.Skip,
			Take:  a.Take,
		})
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Count != nil {
			return ext.Count(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, a *StyleTokenCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *StyleTokenDelegate) runCount(ctx context.Context, params QueryParams[StyleToken]) (int64, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	var query string
	if params.Take != nil || params.Skip != nil {
		var subQuery strings.Builder
		subQuery.WriteString("SELECT 1 FROM ")
		d.client.dialect.WriteQuotedIdent(&subQuery, "StyleToken")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "StyleToken")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, TranslateDBError(err)
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, TranslateDBError(err)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, TranslateDBError(err)
	}
	return count, nil
}
func (d *StyleTokenDelegate) loadRelations(ctx context.Context, records []*StyleToken, selects *StyleTokenSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.App != nil {
		relationSelects, relationOmits, relationParams := selects.App.GetRelationParams()
		returningCols := selectAppCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: StyleToken.app_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *StyleToken) string { return p.App_id }),
			"App",
			"name",
			returningCols,
			scanInto(returningCols, (*App).ScanFields),
			directKey(func(c *App) string { return c.Name }),
			setOne(func(p *StyleToken, c *App) { p.App = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading app: %w", err)
		}
		if err := d.client.App.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
