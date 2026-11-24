package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// AppAgent represents the database model
type AppAgent struct {
	App_language AppLanguageType `db:"app_language" json:"app_language"`
	App_name     string          `db:"app_name" json:"app_name"`
	AgentId      string          `db:"agentId" json:"agentId"`
	Agent        *Agent          `json:"agent,omitempty"`
}

// AppAgentCreate contains model input fields for AppAgent creation operations.
//
// Fields for AppAgent:
//
//	app_language AppLanguage required
//	app_name     string      required
//	agentId      string      required
type AppAgentCreate struct {
	App_language AppLanguageType `json:"app_language"`
	App_name     string          `json:"app_name"`
	AgentId      string          `json:"agentId"`
}

// colMask returns a bit mask of columns that are set
func (s *AppAgentCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	return mask
}

// AppAgentUpdate contains model input fields for AppAgent update operations.
type AppAgentUpdate struct {
	App_language *AppLanguageType `json:"app_language"`
	App_name     *string          `json:"app_name"`
	AgentId      *string          `json:"agentId"`
}

func (u *AppAgentUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.App_language != nil {
		cols = append(cols, "app_language")
		vals = append(vals, u.App_language)
	}
	if u.App_name != nil {
		cols = append(cols, "app_name")
		vals = append(vals, u.App_name)
	}
	if u.AgentId != nil {
		cols = append(cols, "agentId")
		vals = append(vals, u.AgentId)
	}
	return cols, vals
}

func assignmentsToAppAgentUpdate(assignments []FieldAssignment) (AppAgentUpdate, error) {
	var input AppAgentUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "app_language":
			if v, ok := a.Val.(AppLanguageType); ok {
				input.App_language = &v
				if !v.IsValid() {
					errs.Add("app_language", v, "enum", fmt.Sprintf("invalid enum value %q for field app_language", v))
				}
			} else if v, ok := a.Val.(*AppLanguageType); ok {
				input.App_language = v
				if v != nil && !v.IsValid() {
					errs.Add("app_language", *v, "enum", fmt.Sprintf("invalid enum value %q for field app_language", *v))
				}
			} else {
				errs.Add("app_language", a.Val, "type", "field app_language must be of type AppLanguageType")
			}
		case "app_name":
			if v, ok := a.Val.(string); ok {
				input.App_name = &v
				errs.ValidateString("app_name", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.App_name = v
			} else {
				errs.Add("app_name", a.Val, "type", "field app_name must be of type string")
			}
		case "agentId":
			if v, ok := a.Val.(string); ok {
				input.AgentId = &v
				errs.ValidateString("agentId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.AgentId = v
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// AppAgentSelect specifies which scalar and relation fields to select for AppAgent.
//
// Selectable fields:
//
//	-- Scalars --
//	app_language (bool)
//	app_name     (bool)
//	agentId      (bool)
//	-- Relations --
//	agent        (Agent)
type AppAgentSelect struct {
	App_language bool         `json:"app_language"`
	App_name     bool         `json:"app_name"`
	AgentId      bool         `json:"agentId"`
	Agent        *AgentSelect `json:"agent,omitempty"`
}

var fullAppAgentSelectVal = &AppAgentSelect{
	App_language: true,
	App_name:     true,
	AgentId:      true,
}

func fullAppAgentSelect() *AppAgentSelect {
	return fullAppAgentSelectVal
}

func (s *AppAgentSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.App_language || s.App_name || s.AgentId
}

func (s *AppAgentSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type AppAgentOmit struct {
	App_language bool `json:"app_language"`
	App_name     bool `json:"app_name"`
	AgentId      bool `json:"agentId"`
}

type AppAgentSelectQuery interface {
	GetRelationParams() (*AppAgentSelect, *AppAgentOmit, QueryParams[AppAgent])
}

func (s *AppAgentSelect) GetRelationParams() (*AppAgentSelect, *AppAgentOmit, QueryParams[AppAgent]) {
	return s, nil, QueryParams[AppAgent]{}
}

type AppAgentQueryBuilder struct {
	selects *AppAgentSelect
	omits   *AppAgentOmit
	where   []PredicateOf[AppAgent]
	take    *int
	skip    *int
	orderBy []OrderBy[AppAgent]
	cursor  UniquePredicate[AppAgent]
}

func (b *AppAgentQueryBuilder) Where(preds ...PredicateOf[AppAgent]) *AppAgentQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AppAgentQueryBuilder) Take(limit int) *AppAgentQueryBuilder {
	b.take = &limit
	return b
}

func (b *AppAgentQueryBuilder) Skip(offset int) *AppAgentQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AppAgentQueryBuilder) OrderBy(orders ...OrderBy[AppAgent]) *AppAgentQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AppAgentQueryBuilder) Cursor(where UniquePredicate[AppAgent]) *AppAgentQueryBuilder {
	b.cursor = where
	return b
}

func (b *AppAgentQueryBuilder) Select(s AppAgentSelect) *AppAgentQueryBuilder {
	b.selects = &s
	return b
}

func (b *AppAgentQueryBuilder) Omit(o AppAgentOmit) *AppAgentQueryBuilder {
	b.omits = &o
	return b
}

func (b *AppAgentQueryBuilder) GetRelationParams() (*AppAgentSelect, *AppAgentOmit, QueryParams[AppAgent]) {
	if b == nil {
		return nil, nil, QueryParams[AppAgent]{}
	}
	return b.selects, b.omits, QueryParams[AppAgent]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// AppAgentCreateArgs is the input argument passed to AppAgent Create extension hooks.
//
// Fields for AppAgent:
//
//	app_language AppLanguage required
//	app_name     string      required
//	agentId      string      required
//
// Relations for AppAgent:
//
//	agent        (Agent)
type AppAgentCreateArgs struct {
	// Data contains the model fields to insert.
	Data *AppAgentCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *AppAgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// AppAgentCreateManyArgs is the input argument passed to AppAgent CreateMany extension hooks.
//
// Fields for AppAgent:
//
//	app_language AppLanguage required
//	app_name     string      required
//	agentId      string      required
type AppAgentCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppAgentCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppAgentCreateManyArgs) AppendData(builders ...*AppAgentCreateBuilder) *AppAgentCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToAppAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppAgentCreateManyAndReturnArgs is the input argument passed to AppAgent CreateManyAndReturn extension hooks.
//
// Fields for AppAgent:
//
//	app_language AppLanguage required
//	app_name     string      required
//	agentId      string      required
//
// Relations for AppAgent:
//
//	agent        (Agent)
type AppAgentCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppAgentCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *AppAgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppAgentCreateManyAndReturnArgs) AppendData(builders ...*AppAgentCreateBuilder) *AppAgentCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToAppAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppAgentFindUniqueArgs is the input argument passed to AppAgent FindUnique extension hooks.
type AppAgentFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppAgent]
	// Select specifies which scalar and relation fields to select and return.
	Select *AppAgentSelect
}

func (a *AppAgentFindUniqueArgs) SetWhere(unique UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *AppAgentFindUniqueArgs {
	a.Where = make([]PredicateOf[AppAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppAgentFindFirstArgs is the input argument passed to AppAgent FindFirst extension hooks.
type AppAgentFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AppAgent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AppAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppAgentSelect
}

func (a *AppAgentFindFirstArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentFindFirstArgs {
	a.Where = preds
	return a
}

func (a *AppAgentFindFirstArgs) SetOrderBy(orders ...OrderBy[AppAgent]) *AppAgentFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *AppAgentFindFirstArgs) SetCursor(cursor UniquePredicate[AppAgent]) *AppAgentFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *AppAgentFindFirstArgs) SetSkip(n int) *AppAgentFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *AppAgentFindFirstArgs) SetTake(n int) *AppAgentFindFirstArgs {
	a.Take = &n
	return a
}

// AppAgentFindManyArgs is the input argument passed to AppAgent FindMany extension hooks.
type AppAgentFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AppAgent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AppAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppAgentSelect
}

func (a *AppAgentFindManyArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentFindManyArgs {
	a.Where = preds
	return a
}

func (a *AppAgentFindManyArgs) SetOrderBy(orders ...OrderBy[AppAgent]) *AppAgentFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *AppAgentFindManyArgs) SetCursor(cursor UniquePredicate[AppAgent]) *AppAgentFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *AppAgentFindManyArgs) SetSkip(n int) *AppAgentFindManyArgs {
	a.Skip = &n
	return a
}

func (a *AppAgentFindManyArgs) SetTake(n int) *AppAgentFindManyArgs {
	a.Take = &n
	return a
}

// AppAgentCountArgs is the input argument passed to AppAgent Count extension hooks.
type AppAgentCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *AppAgentCountArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentCountArgs {
	a.Where = preds
	return a
}

func (a *AppAgentCountArgs) SetSkip(n int) *AppAgentCountArgs {
	a.Skip = &n
	return a
}

func (a *AppAgentCountArgs) SetTake(n int) *AppAgentCountArgs {
	a.Take = &n
	return a
}

// AppAgentDeleteArgs is the input argument passed to AppAgent Delete extension hooks.
type AppAgentDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppAgent]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *AppAgentSelect
}

func (a *AppAgentDeleteArgs) SetWhere(unique UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *AppAgentDeleteArgs {
	a.Where = make([]PredicateOf[AppAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppAgentDeleteManyArgs is the input argument passed to AppAgent DeleteMany extension hooks.
type AppAgentDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
}

func (a *AppAgentDeleteManyArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentDeleteManyArgs {
	a.Where = preds
	return a
}

// AppAgentUpdateArgs is the input argument passed to AppAgent Update extension hooks.
type AppAgentUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppAgent]
	// Data contains the model fields to update.
	Data *AppAgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppAgentSelect
}

func (a *AppAgentUpdateArgs) SetWhere(unique UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *AppAgentUpdateArgs {
	a.Where = make([]PredicateOf[AppAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppAgentUpdateManyArgs is the input argument passed to AppAgent UpdateMany extension hooks.
type AppAgentUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
	// Data contains the model fields to update.
	Data *AppAgentUpdate
}

func (a *AppAgentUpdateManyArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentUpdateManyArgs {
	a.Where = preds
	return a
}

// AppAgentUpdateManyAndReturnArgs is the input argument passed to AppAgent UpdateManyAndReturn extension hooks.
type AppAgentUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppAgent]
	// Data contains the model fields to update.
	Data *AppAgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppAgentSelect
}

func (a *AppAgentUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[AppAgent]) *AppAgentUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type AppAgentCreateQuery = func(ctx context.Context, args *AppAgentCreateArgs) (*AppAgent, error)
type AppAgentCreateManyQuery = func(ctx context.Context, args *AppAgentCreateManyArgs) (int64, error)
type AppAgentCreateManyAndReturnQuery = func(ctx context.Context, args *AppAgentCreateManyAndReturnArgs) ([]*AppAgent, error)
type AppAgentFindUniqueQuery = func(ctx context.Context, args *AppAgentFindUniqueArgs) (*AppAgent, error)
type AppAgentFindFirstQuery = func(ctx context.Context, args *AppAgentFindFirstArgs) (*AppAgent, error)
type AppAgentFindManyQuery = func(ctx context.Context, args *AppAgentFindManyArgs) ([]*AppAgent, error)
type AppAgentDeleteQuery = func(ctx context.Context, args *AppAgentDeleteArgs) (*AppAgent, error)
type AppAgentDeleteManyQuery = func(ctx context.Context, args *AppAgentDeleteManyArgs) (int64, error)
type AppAgentCountQuery = func(ctx context.Context, args *AppAgentCountArgs) (int64, error)
type AppAgentUpdateQuery = func(ctx context.Context, args *AppAgentUpdateArgs) (*AppAgent, error)
type AppAgentUpdateManyQuery = func(ctx context.Context, args *AppAgentUpdateManyArgs) (int64, error)
type AppAgentUpdateManyAndReturnQuery = func(ctx context.Context, args *AppAgentUpdateManyAndReturnArgs) ([]*AppAgent, error)

type AppAgentExtension struct {
	Create              func(ctx context.Context, args *AppAgentCreateArgs, next AppAgentCreateQuery) (*AppAgent, error)
	CreateMany          func(ctx context.Context, args *AppAgentCreateManyArgs, next AppAgentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *AppAgentCreateManyAndReturnArgs, next AppAgentCreateManyAndReturnQuery) ([]*AppAgent, error)
	FindUnique          func(ctx context.Context, args *AppAgentFindUniqueArgs, next AppAgentFindUniqueQuery) (*AppAgent, error)
	FindFirst           func(ctx context.Context, args *AppAgentFindFirstArgs, next AppAgentFindFirstQuery) (*AppAgent, error)
	FindMany            func(ctx context.Context, args *AppAgentFindManyArgs, next AppAgentFindManyQuery) ([]*AppAgent, error)
	Delete              func(ctx context.Context, args *AppAgentDeleteArgs, next AppAgentDeleteQuery) (*AppAgent, error)
	DeleteMany          func(ctx context.Context, args *AppAgentDeleteManyArgs, next AppAgentDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *AppAgentCountArgs, next AppAgentCountQuery) (int64, error)
	Update              func(ctx context.Context, args *AppAgentUpdateArgs, next AppAgentUpdateQuery) (*AppAgent, error)
	UpdateMany          func(ctx context.Context, args *AppAgentUpdateManyArgs, next AppAgentUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *AppAgentUpdateManyAndReturnArgs, next AppAgentUpdateManyAndReturnQuery) ([]*AppAgent, error)
}

type AppAgentDelegate struct {
	client     *Queries
	extensions []AppAgentExtension
}

func (d *AppAgentDelegate) Use(exts ...AppAgentExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *AppAgent) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "app_language":
			targets[i] = &m.App_language
		case "app_name":
			targets[i] = &m.App_name
		case "agentId":
			targets[i] = &m.AgentId
		}
	}
	return targets
}

var appAgentDefaultCols = []string{
	"app_language",
	"app_name",
	"agentId",
}

var appAgentPKCols = []string{
	"app_language",
	"app_name",
	"agentId",
}

var appAgentUniqueCols = []string{}

func selectAppAgentCols(selects *AppAgentSelect, omits *AppAgentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return appAgentDefaultCols
	}

	anySelected := selects != nil && (selects.App_language || selects.App_name || selects.AgentId || selects.Agent != nil)

	specs := []colSpec{
		{"app_language", selects != nil && selects.App_language, omits != nil && omits.App_language, false},
		{"app_name", selects != nil && selects.App_name, omits != nil && omits.App_name, false},
		{"agentId", selects != nil && selects.AgentId, omits != nil && omits.AgentId, selects != nil && selects.Agent != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *AppAgentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Agent != nil
}

type AppAgentUpsertBuilder struct {
	*CreateBuilder[AppAgent, AppAgentSelect, AppAgentOmit]
}

func (b *AppAgentUpsertBuilder) Select(s AppAgentSelect) *AppAgentUpsertBuilder {
	b.selects = &s
	return b
}

func (b *AppAgentUpsertBuilder) Omit(o AppAgentOmit) *AppAgentUpsertBuilder {
	b.omits = &o
	return b
}

type AppAgentCreateBuilder struct {
	*CreateBuilder[AppAgent, AppAgentSelect, AppAgentOmit]
}

func (b *AppAgentCreateBuilder) Select(s AppAgentSelect) *AppAgentCreateBuilder {
	b.selects = &s
	return b
}

func (b *AppAgentCreateBuilder) Omit(o AppAgentOmit) *AppAgentCreateBuilder {
	b.omits = &o
	return b
}

func (b *AppAgentCreateBuilder) OnConflict(target UniqueConstraintTarget) *AppAgentConflictBuilder[AppAgentUpsertBuilder] {
	upsertBuilder := &AppAgentUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &AppAgentConflictBuilder[AppAgentUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AppAgentCreateBuilder) SetApp_language(v AppLanguageType) *AppAgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_language", Val: v})
	return b
}
func (b *AppAgentCreateBuilder) SetApp_name(v string) *AppAgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_name", Val: v})
	return b
}
func (b *AppAgentCreateBuilder) SetAgentId(v string) *AppAgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AppAgentCreateBuilder) Assignments(assignments ...FieldAssignmentOf[AppAgent]) *AppAgentCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppAgentDelegate) Create() *AppAgentCreateBuilder {
	return &AppAgentCreateBuilder{
		CreateBuilder: &CreateBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedAppAgentApp_language uint64 = 1 << 0
	providedAppAgentApp_name     uint64 = 1 << 1
	providedAppAgentAgentId      uint64 = 1 << 2
)

func assignmentsToAppAgentCreate(assignments []FieldAssignment) (AppAgentCreate, error) {
	var input AppAgentCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "app_language":
			provided |= providedAppAgentApp_language
			if v, ok := a.Val.(AppLanguageType); ok {
				input.App_language = v
				if !v.IsValid() {
					errs.Add("app_language", v, "enum", fmt.Sprintf("invalid enum value %q for field app_language", v))
				}
			} else {
				errs.Add("app_language", a.Val, "type", "field app_language must be of type AppLanguageType")
			}
		case "app_name":
			provided |= providedAppAgentApp_name
			if v, ok := a.Val.(string); ok {
				input.App_name = v
				errs.ValidateString("app_name", v, true, 0, false, false)
			} else {
				errs.Add("app_name", a.Val, "type", "field app_name must be of type string")
			}
		case "agentId":
			provided |= providedAppAgentAgentId
			if v, ok := a.Val.(string); ok {
				input.AgentId = v
				errs.ValidateString("agentId", v, true, 0, false, false)
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
			}
		}
	}
	if provided&providedAppAgentApp_language == 0 {
		errs.Add("app_language", nil, "required", "field App_language is required")
	}
	if provided&providedAppAgentApp_name == 0 {
		errs.Add("app_name", "", "required", "field App_name is required")
	}
	if provided&providedAppAgentAgentId == 0 {
		errs.Add("agentId", "", "required", "field AgentId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AppAgentCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 3)
	vals = make([]any, 0, 3)
	cols = append(cols, "app_language")
	vals = append(vals, s.App_language)
	cols = append(cols, "app_name")
	vals = append(vals, s.App_name)
	cols = append(cols, "agentId")
	vals = append(vals, s.AgentId)
	return
}

func partitionAppAgentInputs(dialect Dialect, inputs []*AppAgentCreate) [][]*AppAgentCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AppAgentCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AppAgentCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AppAgentCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AppAgentCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AppAgentCreate{inputs}
}

func (d *AppAgentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AppAgentSelect, omits *AppAgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*AppAgent, error) {
	input, err := assignmentsToAppAgentCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAppAgentCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppAgentCreateArgs) (*AppAgent, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectAppAgentCols(a.Select, omits)
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
			curr = func(c context.Context, a *AppAgentCreateArgs) (*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type AppAgentCreateManyBuilder struct {
	*CreateManyBuilder[AppAgent]
}

func (b *AppAgentCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AppAgentConflictBuilder[AppAgentCreateManyBuilder] {
	return &AppAgentConflictBuilder[AppAgentCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AppAgentCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[AppAgent, AppAgentSelect, AppAgentOmit]
}

func (b *AppAgentCreateManyAndReturnBuilder) Select(s AppAgentSelect) *AppAgentCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *AppAgentCreateManyAndReturnBuilder) Omit(o AppAgentOmit) *AppAgentCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *AppAgentCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AppAgentConflictBuilder[AppAgentCreateManyAndReturnBuilder] {
	return &AppAgentConflictBuilder[AppAgentCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToAppAgentRecordInputs(builders []*AppAgentCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *AppAgentDelegate) CreateMany(builders ...*AppAgentCreateBuilder) *AppAgentCreateManyBuilder {
	return &AppAgentCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[AppAgent]{
			records:  createBuildersToAppAgentRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AppAgentDelegate) CreateManyAndReturn(builders ...*AppAgentCreateBuilder) *AppAgentCreateManyAndReturnBuilder {
	return &AppAgentCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
			records:  createBuildersToAppAgentRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToAppAgentCreateInputs(records []RecordInput) ([]*AppAgentCreate, error) {
	structs := make([]AppAgentCreate, len(records))
	inputs := make([]*AppAgentCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToAppAgentCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *AppAgentDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToAppAgentCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &AppAgentCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppAgentCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppAgentCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AppAgentSelect, omits *AppAgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*AppAgent, error) {
	inputs, err := recordsToAppAgentCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppAgentCreateManyAndReturnArgs) ([]*AppAgent, error) {
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
			curr = func(c context.Context, a *AppAgentCreateManyAndReturnArgs) ([]*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *AppAgentSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*AppAgent, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *AppAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AppAgent.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.AppAgent.loadRelations(ctx, []*AppAgent{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "AppAgent", cols, returningCols, appAgentPKCols, conflictTarget, conflictAction, len(vals))
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

		var res AppAgent
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, appAgentPKCols)
}

func (d *AppAgentDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*AppAgent, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("AppAgent") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "AppAgent")
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

	var res AppAgent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *AppAgentDelegate) buildBulkInsertSQL(q *Queries, batch []*AppAgentCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 3)
	for i, c := range appAgentDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "AppAgent")
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
			case "app_language":
				vals = append(vals, input.App_language)
			case "app_name":
				vals = append(vals, input.App_name)
			case "agentId":
				vals = append(vals, input.AgentId)
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

func applyAppAgentConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanAppAgentRows(rows *sql.Rows, returningCols []string) ([]*AppAgent, error) {
	var records []*AppAgent
	for rows.Next() {
		var res AppAgent
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

func (d *AppAgentDelegate) runCreateMany(ctx context.Context, inputs []*AppAgentCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAppAgentInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppAgentConflictClause(d.client.dialect, queryStr, vals, cols, appAgentPKCols, conflictTarget, conflictAction)

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

func (d *AppAgentDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AppAgentCreate,
	selects *AppAgentSelect,
	omits *AppAgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AppAgent, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAppAgentInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*AppAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.AppAgent.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.AppAgent.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectAppAgentCols(selects, omits, appAgentPKCols...)
	recordsOut := make([]*AppAgent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppAgentConflictClause(d.client.dialect, queryStr, vals, cols, appAgentPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanAppAgentRows(rows, returningCols)
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

func (d *AppAgentDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*AppAgentCreate,
	selects *AppAgentSelect,
	omits *AppAgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AppAgent, error) {
	batches := partitionAppAgentInputs(d.client.dialect, inputs)
	returningCols := selectAppAgentCols(selects, omits, appAgentPKCols...)
	recordsOut := make([]*AppAgent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppAgentConflictClause(d.client.dialect, queryStr, vals, cols, appAgentPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("AppAgent") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "AppAgent")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appAgentPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appAgentPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanAppAgentRows(rows, returningCols)
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

type AppAgentConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AppAgentConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppAgentConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppAgentConflictBuilder[B]) Update(fn func(u *AppAgentUpsert)) *B {
	cb.setAction(*AppAgentConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// AppAgentConflictUpdate creates a custom ConflictAction for AppAgent upsert conflicts.
func AppAgentConflictUpdate(fn func(u *AppAgentUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newAppAgentUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type AppAgentUpsert struct {
	App_language fieldUpsert[AppLanguageType]
	App_name     fieldUpsert[string]
	AgentId      fieldUpsert[string]
}

func newAppAgentUpsert(up *ConflictUpdate) *AppAgentUpsert {
	return &AppAgentUpsert{
		App_language: fieldUpsert[AppLanguageType]{column: "app_language", update: up},
		App_name:     fieldUpsert[string]{column: "app_name", update: up},
		AgentId:      fieldUpsert[string]{column: "agentId", update: up},
	}
}

type AppAgentUpdateBuilder struct {
	*UpdateBuilder[AppAgent, AppAgentSelect, AppAgentOmit]
}

type AppAgentUpdateManyBuilder struct {
	*UpdateManyBuilder[AppAgent]
}

type AppAgentUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[AppAgent, AppAgentSelect, AppAgentOmit]
}

func (b *AppAgentUpdateBuilder) SetApp_language(v AppLanguageType) *AppAgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_language", Val: v})
	return b
}

func (b *AppAgentUpdateManyBuilder) SetApp_language(v AppLanguageType) *AppAgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_language", Val: v})
	return b
}

func (b *AppAgentUpdateManyAndReturnBuilder) SetApp_language(v AppLanguageType) *AppAgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_language", Val: v})
	return b
}
func (b *AppAgentUpdateBuilder) SetApp_name(v string) *AppAgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_name", Val: v})
	return b
}

func (b *AppAgentUpdateManyBuilder) SetApp_name(v string) *AppAgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_name", Val: v})
	return b
}

func (b *AppAgentUpdateManyAndReturnBuilder) SetApp_name(v string) *AppAgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_name", Val: v})
	return b
}
func (b *AppAgentUpdateBuilder) SetAgentId(v string) *AppAgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AppAgentUpdateManyBuilder) SetAgentId(v string) *AppAgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AppAgentUpdateManyAndReturnBuilder) SetAgentId(v string) *AppAgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AppAgentUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[AppAgent]) *AppAgentUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppAgentUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[AppAgent]) *AppAgentUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppAgentUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[AppAgent]) *AppAgentUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppAgentDelegate) Update(where UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *AppAgentUpdateBuilder {
	return &AppAgentUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *AppAgentDelegate) UpdateMany(preds ...PredicateOf[AppAgent]) *AppAgentUpdateManyBuilder {
	return &AppAgentUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[AppAgent]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *AppAgentDelegate) UpdateManyAndReturn(preds ...PredicateOf[AppAgent]) *AppAgentUpdateManyAndReturnBuilder {
	return &AppAgentUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *AppAgentDelegate) buildUpdateSQL(preds []PredicateOf[AppAgent], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "AppAgent")
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

func (d *AppAgentDelegate) executeUpdate(ctx context.Context, where UniquePredicate[AppAgent], additional []PredicateOf[AppAgent], assignments []FieldAssignment, selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
	allWhere := make([]PredicateOf[AppAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToAppAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppAgentUpdateArgs) (*AppAgent, error) {
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
			curr = func(c context.Context, a *AppAgentUpdateArgs) (*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runUpdate(ctx context.Context, preds []PredicateOf[AppAgent], cols []string, vals []any, selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
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
		var res *AppAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AppAgent.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AppAgent.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppAgentCols(selects, omits, appAgentPKCols...)
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
		return nil, &NotFoundError{Model: "AppAgent"}
	}

	var res AppAgent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AppAgent{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *AppAgentDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[AppAgent], cols []string, vals []any) (int64, error) {
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

func (d *AppAgentDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[AppAgent], cols []string, vals []any, selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "AppAgent"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *AppAgentDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[AppAgent], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToAppAgentUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &AppAgentUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *AppAgentUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppAgentUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *AppAgentDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AppAgent], assignments []FieldAssignment, selects *AppAgentSelect, omits *AppAgentOmit) ([]*AppAgent, error) {
	input, err := assignmentsToAppAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppAgentUpdateManyAndReturnArgs) ([]*AppAgent, error) {
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
			curr = func(c context.Context, a *AppAgentUpdateManyAndReturnArgs) ([]*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AppAgent], cols []string, vals []any, selects *AppAgentSelect, omits *AppAgentOmit) ([]*AppAgent, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[AppAgent]{Where: preds}, selects, omits)
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
		var res []*AppAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AppAgent.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AppAgent.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppAgentCols(selects, omits, appAgentPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanAppAgentRows(rows, returningCols)
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

func (d *AppAgentDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[AppAgent], cols []string, vals []any, selects *AppAgentSelect, omits *AppAgentOmit) ([]*AppAgent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*AppAgent{}, nil
	}
	return d.runFindMany(ctx, QueryParams[AppAgent]{Where: preds}, selects, omits)
}
func (d *AppAgentDelegate) FindUnique(where UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *FindUniqueBuilder[AppAgent, AppAgentSelect, AppAgentOmit] {
	return &FindUniqueBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AppAgentDelegate) FindFirst(preds ...PredicateOf[AppAgent]) *FindFirstBuilder[AppAgent, AppAgentSelect, AppAgentOmit] {
	return &FindFirstBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AppAgentDelegate) FindMany(preds ...PredicateOf[AppAgent]) *FindManyBuilder[AppAgent, AppAgentSelect, AppAgentOmit] {
	return &FindManyBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AppAgentDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[AppAgent], additional []PredicateOf[AppAgent], selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
	allWhere := make([]PredicateOf[AppAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppAgentFindUniqueArgs) (*AppAgent, error) {
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
			curr = func(c context.Context, a *AppAgentFindUniqueArgs) (*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[AppAgent],
	selects *AppAgentSelect,
	omits *AppAgentOmit,
) (*AppAgent, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppAgentFindFirstArgs) (*AppAgent, error) {
		return d.runFindFirst(c, QueryParams[AppAgent]{
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
			curr = func(c context.Context, a *AppAgentFindFirstArgs) (*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[AppAgent],
	selects *AppAgentSelect,
	omits *AppAgentOmit,
) ([]*AppAgent, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppAgentFindManyArgs) ([]*AppAgent, error) {
		return d.runFindMany(c, QueryParams[AppAgent]{
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
			curr = func(c context.Context, a *AppAgentFindManyArgs) ([]*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runFindUnique(ctx context.Context, where []PredicateOf[AppAgent], selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
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
	returningCols := selectAppAgentCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AppAgent{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *AppAgentDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[AppAgent],
	selects *AppAgentSelect,
	omits *AppAgentOmit,
) (*AppAgent, error) {
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

func (d *AppAgentDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[AppAgent],
	selects *AppAgentSelect,
	omits *AppAgentOmit,
) ([]*AppAgent, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, appAgentPKCols, appAgentUniqueCols, "AppAgent", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, appAgentPKCols, appAgentUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAppAgentCols(selects, omits)

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

func (d *AppAgentDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*AppAgent, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "AppAgent", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res AppAgent
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *AppAgentDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*AppAgent, error) {
	query := buildSelectSQL(d.client, "AppAgent", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*AppAgent, 0)
	for rows.Next() {
		var res AppAgent
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
func (d *AppAgentDelegate) DeleteMany(preds ...PredicateOf[AppAgent]) *DeleteManyBuilder[AppAgent] {
	return &DeleteManyBuilder[AppAgent]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AppAgentDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[AppAgent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &AppAgentDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *AppAgentDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppAgentDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[AppAgent]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "AppAgent")
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

func (d *AppAgentDelegate) Delete(where UniquePredicate[AppAgent], additional ...PredicateOf[AppAgent]) *DeleteBuilder[AppAgent, AppAgentSelect, AppAgentOmit] {
	return &DeleteBuilder[AppAgent, AppAgentSelect, AppAgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *AppAgentDelegate) executeDelete(ctx context.Context, where UniquePredicate[AppAgent], additional []PredicateOf[AppAgent], selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
	allWhere := make([]PredicateOf[AppAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppAgentSelect()
	}

	args := &AppAgentDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppAgentDeleteArgs) (*AppAgent, error) {
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
			curr = func(c context.Context, a *AppAgentDeleteArgs) (*AppAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runDelete(ctx context.Context, where []PredicateOf[AppAgent], selects *AppAgentSelect, omits *AppAgentOmit) (*AppAgent, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectAppAgentCols(selects, omits, appAgentPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *AppAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AppAgent.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "AppAgent"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "AppAgent")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[AppAgent]
			pkPreds = append(pkPreds, Predicate[AppAgent]{
				Data: PredicateData{
					Column:   "app_language",
					Operator: "=",
					Value:    res.App_language,
				},
			})
			pkPreds = append(pkPreds, Predicate[AppAgent]{
				Data: PredicateData{
					Column:   "app_name",
					Operator: "=",
					Value:    res.App_name,
				},
			})
			pkPreds = append(pkPreds, Predicate[AppAgent]{
				Data: PredicateData{
					Column:   "agentId",
					Operator: "=",
					Value:    res.AgentId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "AppAgent")

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
		return nil, &NotFoundError{Model: "AppAgent"}
	}

	var row AppAgent
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *AppAgentDelegate) Count(preds ...PredicateOf[AppAgent]) *CountBuilder[AppAgent] {
	return &CountBuilder[AppAgent]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AppAgentDelegate) executeCount(ctx context.Context, params QueryParams[AppAgent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &AppAgentCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *AppAgentCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[AppAgent]{
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
			curr = func(c context.Context, a *AppAgentCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppAgentDelegate) runCount(ctx context.Context, params QueryParams[AppAgent]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "AppAgent")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "AppAgent")
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
func (d *AppAgentDelegate) loadRelations(ctx context.Context, records []*AppAgent, selects *AppAgentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Agent != nil {
		relationSelects, relationOmits, relationParams := selects.Agent.GetRelationParams()
		returningCols := selectAgentCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: AppAgent.agentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *AppAgent) string { return p.AgentId }),
			"Agent",
			"id",
			returningCols,
			scanInto(returningCols, (*Agent).ScanFields),
			directKey(func(c *Agent) string { return c.Id }),
			setOne(func(p *AppAgent, c *Agent) { p.Agent = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agent: %w", err)
		}
		if err := d.client.Agent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
