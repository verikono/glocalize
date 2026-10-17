package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// AppSkill represents the database model
type AppSkill struct {
	App_id   string       `db:"app_id" json:"app_id"`
	SkillId  string       `db:"skillId" json:"skillId"`
	Priority PriorityType `db:"priority" json:"priority"`
	App      *App         `json:"app,omitempty"`
	Skill    *Skill       `json:"skill,omitempty"`
}

// AppSkillCreate contains model input fields for AppSkill creation operations.
//
// Fields for AppSkill:
//
//	app_id   string   required
//	skillId  string   required
//	priority Priority default: NORMAL
type AppSkillCreate struct {
	App_id   string        `json:"app_id"`
	SkillId  string        `json:"skillId"`
	Priority *PriorityType `json:"priority"`
}

// colMask returns a bit mask of columns that are set
func (s *AppSkillCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	if s.Priority != nil {
		mask |= 1 << 2
	}
	return mask
}

// AppSkillUpdate contains model input fields for AppSkill update operations.
type AppSkillUpdate struct {
	App_id   *string       `json:"app_id"`
	SkillId  *string       `json:"skillId"`
	Priority *PriorityType `json:"priority"`
}

func (u *AppSkillUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.App_id != nil {
		cols = append(cols, "app_id")
		vals = append(vals, u.App_id)
	}
	if u.SkillId != nil {
		cols = append(cols, "skillId")
		vals = append(vals, u.SkillId)
	}
	if u.Priority != nil {
		cols = append(cols, "priority")
		vals = append(vals, u.Priority)
	}
	return cols, vals
}

func assignmentsToAppSkillUpdate(assignments []FieldAssignment) (AppSkillUpdate, error) {
	var input AppSkillUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "app_id":
			if v, ok := a.Val.(string); ok {
				input.App_id = &v
				errs.ValidateString("app_id", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.App_id = v
			} else {
				errs.Add("app_id", a.Val, "type", "field app_id must be of type string")
			}
		case "skillId":
			if v, ok := a.Val.(string); ok {
				input.SkillId = &v
				errs.ValidateString("skillId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.SkillId = v
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		case "priority":
			if v, ok := a.Val.(PriorityType); ok {
				input.Priority = &v
				if !v.IsValid() {
					errs.Add("priority", v, "enum", fmt.Sprintf("invalid enum value %q for field priority", v))
				}
			} else if v, ok := a.Val.(*PriorityType); ok {
				input.Priority = v
				if v != nil && !v.IsValid() {
					errs.Add("priority", *v, "enum", fmt.Sprintf("invalid enum value %q for field priority", *v))
				}
			} else {
				errs.Add("priority", a.Val, "type", "field priority must be of type PriorityType")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// AppSkillSelect specifies which scalar and relation fields to select for AppSkill.
//
// Selectable fields:
//
//	-- Scalars --
//	app_id   (bool)
//	skillId  (bool)
//	priority (bool)
//	-- Relations --
//	app      (App)
//	skill    (Skill)
type AppSkillSelect struct {
	App_id   bool         `json:"app_id"`
	SkillId  bool         `json:"skillId"`
	Priority bool         `json:"priority"`
	App      *AppSelect   `json:"app,omitempty"`
	Skill    *SkillSelect `json:"skill,omitempty"`
}

var fullAppSkillSelectVal = &AppSkillSelect{
	App_id:   true,
	SkillId:  true,
	Priority: true,
}

func fullAppSkillSelect() *AppSkillSelect {
	return fullAppSkillSelectVal
}

func (s *AppSkillSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.App_id || s.SkillId || s.Priority
}

func (s *AppSkillSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type AppSkillOmit struct {
	App_id   bool `json:"app_id"`
	SkillId  bool `json:"skillId"`
	Priority bool `json:"priority"`
}

type AppSkillSelectQuery interface {
	GetRelationParams() (*AppSkillSelect, *AppSkillOmit, QueryParams[AppSkill])
}

func (s *AppSkillSelect) GetRelationParams() (*AppSkillSelect, *AppSkillOmit, QueryParams[AppSkill]) {
	return s, nil, QueryParams[AppSkill]{}
}

type AppSkillQueryBuilder struct {
	selects *AppSkillSelect
	omits   *AppSkillOmit
	where   []PredicateOf[AppSkill]
	take    *int
	skip    *int
	orderBy []OrderBy[AppSkill]
	cursor  UniquePredicate[AppSkill]
}

func (b *AppSkillQueryBuilder) Where(preds ...PredicateOf[AppSkill]) *AppSkillQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AppSkillQueryBuilder) Take(limit int) *AppSkillQueryBuilder {
	b.take = &limit
	return b
}

func (b *AppSkillQueryBuilder) Skip(offset int) *AppSkillQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AppSkillQueryBuilder) OrderBy(orders ...OrderBy[AppSkill]) *AppSkillQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AppSkillQueryBuilder) Cursor(where UniquePredicate[AppSkill]) *AppSkillQueryBuilder {
	b.cursor = where
	return b
}

func (b *AppSkillQueryBuilder) Select(s AppSkillSelect) *AppSkillQueryBuilder {
	b.selects = &s
	return b
}

func (b *AppSkillQueryBuilder) Omit(o AppSkillOmit) *AppSkillQueryBuilder {
	b.omits = &o
	return b
}

func (b *AppSkillQueryBuilder) GetRelationParams() (*AppSkillSelect, *AppSkillOmit, QueryParams[AppSkill]) {
	if b == nil {
		return nil, nil, QueryParams[AppSkill]{}
	}
	return b.selects, b.omits, QueryParams[AppSkill]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// AppSkillCreateArgs is the input argument passed to AppSkill Create extension hooks.
//
// Fields for AppSkill:
//
//	app_id   string   required
//	skillId  string   required
//	priority Priority default: NORMAL
//
// Relations for AppSkill:
//
//	app      (App)
//	skill    (Skill)
type AppSkillCreateArgs struct {
	// Data contains the model fields to insert.
	Data *AppSkillCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *AppSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// AppSkillCreateManyArgs is the input argument passed to AppSkill CreateMany extension hooks.
//
// Fields for AppSkill:
//
//	app_id   string   required
//	skillId  string   required
//	priority Priority default: NORMAL
type AppSkillCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppSkillCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppSkillCreateManyArgs) AppendData(builders ...*AppSkillCreateBuilder) *AppSkillCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToAppSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppSkillCreateManyAndReturnArgs is the input argument passed to AppSkill CreateManyAndReturn extension hooks.
//
// Fields for AppSkill:
//
//	app_id   string   required
//	skillId  string   required
//	priority Priority default: NORMAL
//
// Relations for AppSkill:
//
//	app      (App)
//	skill    (Skill)
type AppSkillCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppSkillCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *AppSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppSkillCreateManyAndReturnArgs) AppendData(builders ...*AppSkillCreateBuilder) *AppSkillCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToAppSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppSkillFindUniqueArgs is the input argument passed to AppSkill FindUnique extension hooks.
type AppSkillFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppSkill]
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSkillSelect
}

func (a *AppSkillFindUniqueArgs) SetWhere(unique UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *AppSkillFindUniqueArgs {
	a.Where = make([]PredicateOf[AppSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppSkillFindFirstArgs is the input argument passed to AppSkill FindFirst extension hooks.
type AppSkillFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AppSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AppSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSkillSelect
}

func (a *AppSkillFindFirstArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillFindFirstArgs {
	a.Where = preds
	return a
}

func (a *AppSkillFindFirstArgs) SetOrderBy(orders ...OrderBy[AppSkill]) *AppSkillFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *AppSkillFindFirstArgs) SetCursor(cursor UniquePredicate[AppSkill]) *AppSkillFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *AppSkillFindFirstArgs) SetSkip(n int) *AppSkillFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *AppSkillFindFirstArgs) SetTake(n int) *AppSkillFindFirstArgs {
	a.Take = &n
	return a
}

// AppSkillFindManyArgs is the input argument passed to AppSkill FindMany extension hooks.
type AppSkillFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AppSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AppSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSkillSelect
}

func (a *AppSkillFindManyArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillFindManyArgs {
	a.Where = preds
	return a
}

func (a *AppSkillFindManyArgs) SetOrderBy(orders ...OrderBy[AppSkill]) *AppSkillFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *AppSkillFindManyArgs) SetCursor(cursor UniquePredicate[AppSkill]) *AppSkillFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *AppSkillFindManyArgs) SetSkip(n int) *AppSkillFindManyArgs {
	a.Skip = &n
	return a
}

func (a *AppSkillFindManyArgs) SetTake(n int) *AppSkillFindManyArgs {
	a.Take = &n
	return a
}

// AppSkillCountArgs is the input argument passed to AppSkill Count extension hooks.
type AppSkillCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *AppSkillCountArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillCountArgs {
	a.Where = preds
	return a
}

func (a *AppSkillCountArgs) SetSkip(n int) *AppSkillCountArgs {
	a.Skip = &n
	return a
}

func (a *AppSkillCountArgs) SetTake(n int) *AppSkillCountArgs {
	a.Take = &n
	return a
}

// AppSkillDeleteArgs is the input argument passed to AppSkill Delete extension hooks.
type AppSkillDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppSkill]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *AppSkillSelect
}

func (a *AppSkillDeleteArgs) SetWhere(unique UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *AppSkillDeleteArgs {
	a.Where = make([]PredicateOf[AppSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppSkillDeleteManyArgs is the input argument passed to AppSkill DeleteMany extension hooks.
type AppSkillDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
}

func (a *AppSkillDeleteManyArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillDeleteManyArgs {
	a.Where = preds
	return a
}

// AppSkillUpdateArgs is the input argument passed to AppSkill Update extension hooks.
type AppSkillUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AppSkill]
	// Data contains the model fields to update.
	Data *AppSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppSkillSelect
}

func (a *AppSkillUpdateArgs) SetWhere(unique UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *AppSkillUpdateArgs {
	a.Where = make([]PredicateOf[AppSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppSkillUpdateManyArgs is the input argument passed to AppSkill UpdateMany extension hooks.
type AppSkillUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
	// Data contains the model fields to update.
	Data *AppSkillUpdate
}

func (a *AppSkillUpdateManyArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillUpdateManyArgs {
	a.Where = preds
	return a
}

// AppSkillUpdateManyAndReturnArgs is the input argument passed to AppSkill UpdateManyAndReturn extension hooks.
type AppSkillUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AppSkill]
	// Data contains the model fields to update.
	Data *AppSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppSkillSelect
}

func (a *AppSkillUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[AppSkill]) *AppSkillUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type AppSkillCreateQuery = func(ctx context.Context, args *AppSkillCreateArgs) (*AppSkill, error)
type AppSkillCreateManyQuery = func(ctx context.Context, args *AppSkillCreateManyArgs) (int64, error)
type AppSkillCreateManyAndReturnQuery = func(ctx context.Context, args *AppSkillCreateManyAndReturnArgs) ([]*AppSkill, error)
type AppSkillFindUniqueQuery = func(ctx context.Context, args *AppSkillFindUniqueArgs) (*AppSkill, error)
type AppSkillFindFirstQuery = func(ctx context.Context, args *AppSkillFindFirstArgs) (*AppSkill, error)
type AppSkillFindManyQuery = func(ctx context.Context, args *AppSkillFindManyArgs) ([]*AppSkill, error)
type AppSkillDeleteQuery = func(ctx context.Context, args *AppSkillDeleteArgs) (*AppSkill, error)
type AppSkillDeleteManyQuery = func(ctx context.Context, args *AppSkillDeleteManyArgs) (int64, error)
type AppSkillCountQuery = func(ctx context.Context, args *AppSkillCountArgs) (int64, error)
type AppSkillUpdateQuery = func(ctx context.Context, args *AppSkillUpdateArgs) (*AppSkill, error)
type AppSkillUpdateManyQuery = func(ctx context.Context, args *AppSkillUpdateManyArgs) (int64, error)
type AppSkillUpdateManyAndReturnQuery = func(ctx context.Context, args *AppSkillUpdateManyAndReturnArgs) ([]*AppSkill, error)

type AppSkillExtension struct {
	Create              func(ctx context.Context, args *AppSkillCreateArgs, next AppSkillCreateQuery) (*AppSkill, error)
	CreateMany          func(ctx context.Context, args *AppSkillCreateManyArgs, next AppSkillCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *AppSkillCreateManyAndReturnArgs, next AppSkillCreateManyAndReturnQuery) ([]*AppSkill, error)
	FindUnique          func(ctx context.Context, args *AppSkillFindUniqueArgs, next AppSkillFindUniqueQuery) (*AppSkill, error)
	FindFirst           func(ctx context.Context, args *AppSkillFindFirstArgs, next AppSkillFindFirstQuery) (*AppSkill, error)
	FindMany            func(ctx context.Context, args *AppSkillFindManyArgs, next AppSkillFindManyQuery) ([]*AppSkill, error)
	Delete              func(ctx context.Context, args *AppSkillDeleteArgs, next AppSkillDeleteQuery) (*AppSkill, error)
	DeleteMany          func(ctx context.Context, args *AppSkillDeleteManyArgs, next AppSkillDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *AppSkillCountArgs, next AppSkillCountQuery) (int64, error)
	Update              func(ctx context.Context, args *AppSkillUpdateArgs, next AppSkillUpdateQuery) (*AppSkill, error)
	UpdateMany          func(ctx context.Context, args *AppSkillUpdateManyArgs, next AppSkillUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *AppSkillUpdateManyAndReturnArgs, next AppSkillUpdateManyAndReturnQuery) ([]*AppSkill, error)
}

type AppSkillDelegate struct {
	client     *Queries
	extensions []AppSkillExtension
}

func (d *AppSkillDelegate) Use(exts ...AppSkillExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *AppSkill) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "app_id":
			targets[i] = &m.App_id
		case "skillId":
			targets[i] = &m.SkillId
		case "priority":
			targets[i] = &m.Priority
		}
	}
	return targets
}

var appSkillDefaultCols = []string{
	"app_id",
	"skillId",
	"priority",
}

var appSkillPKCols = []string{
	"app_id",
	"skillId",
}

var appSkillUniqueCols = []string{}

func selectAppSkillCols(selects *AppSkillSelect, omits *AppSkillOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return appSkillDefaultCols
	}

	anySelected := selects != nil && (selects.App_id || selects.SkillId || selects.Priority || selects.App != nil || selects.Skill != nil)

	specs := []colSpec{
		{"app_id", selects != nil && selects.App_id, omits != nil && omits.App_id, selects != nil && selects.App != nil},
		{"skillId", selects != nil && selects.SkillId, omits != nil && omits.SkillId, selects != nil && selects.Skill != nil},
		{"priority", selects != nil && selects.Priority, omits != nil && omits.Priority, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *AppSkillSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.App != nil || s.Skill != nil
}

type AppSkillUpsertBuilder struct {
	*CreateBuilder[AppSkill, AppSkillSelect, AppSkillOmit]
}

func (b *AppSkillUpsertBuilder) Select(s AppSkillSelect) *AppSkillUpsertBuilder {
	b.selects = &s
	return b
}

func (b *AppSkillUpsertBuilder) Omit(o AppSkillOmit) *AppSkillUpsertBuilder {
	b.omits = &o
	return b
}

type AppSkillCreateBuilder struct {
	*CreateBuilder[AppSkill, AppSkillSelect, AppSkillOmit]
}

func (b *AppSkillCreateBuilder) Select(s AppSkillSelect) *AppSkillCreateBuilder {
	b.selects = &s
	return b
}

func (b *AppSkillCreateBuilder) Omit(o AppSkillOmit) *AppSkillCreateBuilder {
	b.omits = &o
	return b
}

func (b *AppSkillCreateBuilder) OnConflict(target UniqueConstraintTarget) *AppSkillConflictBuilder[AppSkillUpsertBuilder] {
	upsertBuilder := &AppSkillUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &AppSkillConflictBuilder[AppSkillUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AppSkillCreateBuilder) SetApp_id(v string) *AppSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}
func (b *AppSkillCreateBuilder) SetSkillId(v string) *AppSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *AppSkillCreateBuilder) SetPriority(v PriorityType) *AppSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AppSkillCreateBuilder) Assignments(assignments ...FieldAssignmentOf[AppSkill]) *AppSkillCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppSkillDelegate) Create() *AppSkillCreateBuilder {
	return &AppSkillCreateBuilder{
		CreateBuilder: &CreateBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedAppSkillApp_id   uint64 = 1 << 0
	providedAppSkillSkillId  uint64 = 1 << 1
	providedAppSkillPriority uint64 = 1 << 2
)

func assignmentsToAppSkillCreate(assignments []FieldAssignment) (AppSkillCreate, error) {
	var input AppSkillCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "app_id":
			provided |= providedAppSkillApp_id
			if v, ok := a.Val.(string); ok {
				input.App_id = v
				errs.ValidateString("app_id", v, true, 0, false, false)
			} else {
				errs.Add("app_id", a.Val, "type", "field app_id must be of type string")
			}
		case "skillId":
			provided |= providedAppSkillSkillId
			if v, ok := a.Val.(string); ok {
				input.SkillId = v
				errs.ValidateString("skillId", v, true, 0, false, false)
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		case "priority":
			provided |= providedAppSkillPriority
			if v, ok := a.Val.(PriorityType); ok {
				input.Priority = &v
				if !v.IsValid() {
					errs.Add("priority", v, "enum", fmt.Sprintf("invalid enum value %q for field priority", v))
				}
			} else {
				errs.Add("priority", a.Val, "type", "field priority must be of type PriorityType")
			}
		}
	}
	if provided&providedAppSkillApp_id == 0 {
		errs.Add("app_id", "", "required", "field App_id is required")
	}
	if provided&providedAppSkillSkillId == 0 {
		errs.Add("skillId", "", "required", "field SkillId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AppSkillCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 3)
	vals = make([]any, 0, 3)
	cols = append(cols, "app_id")
	vals = append(vals, s.App_id)
	cols = append(cols, "skillId")
	vals = append(vals, s.SkillId)
	if s.Priority != nil {
		cols = append(cols, "priority")
		vals = append(vals, *s.Priority)
	}
	return
}

func partitionAppSkillInputs(dialect Dialect, inputs []*AppSkillCreate) [][]*AppSkillCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AppSkillCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AppSkillCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AppSkillCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AppSkillCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AppSkillCreate{inputs}
}

func (d *AppSkillDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AppSkillSelect, omits *AppSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*AppSkill, error) {
	input, err := assignmentsToAppSkillCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAppSkillCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppSkillCreateArgs) (*AppSkill, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectAppSkillCols(a.Select, omits)
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
			curr = func(c context.Context, a *AppSkillCreateArgs) (*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type AppSkillCreateManyBuilder struct {
	*CreateManyBuilder[AppSkill]
}

func (b *AppSkillCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AppSkillConflictBuilder[AppSkillCreateManyBuilder] {
	return &AppSkillConflictBuilder[AppSkillCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AppSkillCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[AppSkill, AppSkillSelect, AppSkillOmit]
}

func (b *AppSkillCreateManyAndReturnBuilder) Select(s AppSkillSelect) *AppSkillCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *AppSkillCreateManyAndReturnBuilder) Omit(o AppSkillOmit) *AppSkillCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *AppSkillCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AppSkillConflictBuilder[AppSkillCreateManyAndReturnBuilder] {
	return &AppSkillConflictBuilder[AppSkillCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToAppSkillRecordInputs(builders []*AppSkillCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *AppSkillDelegate) CreateMany(builders ...*AppSkillCreateBuilder) *AppSkillCreateManyBuilder {
	return &AppSkillCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[AppSkill]{
			records:  createBuildersToAppSkillRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AppSkillDelegate) CreateManyAndReturn(builders ...*AppSkillCreateBuilder) *AppSkillCreateManyAndReturnBuilder {
	return &AppSkillCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
			records:  createBuildersToAppSkillRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToAppSkillCreateInputs(records []RecordInput) ([]*AppSkillCreate, error) {
	structs := make([]AppSkillCreate, len(records))
	inputs := make([]*AppSkillCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToAppSkillCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *AppSkillDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToAppSkillCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &AppSkillCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppSkillCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppSkillCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AppSkillSelect, omits *AppSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*AppSkill, error) {
	inputs, err := recordsToAppSkillCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppSkillCreateManyAndReturnArgs) ([]*AppSkill, error) {
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
			curr = func(c context.Context, a *AppSkillCreateManyAndReturnArgs) ([]*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *AppSkillSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*AppSkill, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *AppSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AppSkill.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.AppSkill.loadRelations(ctx, []*AppSkill{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "AppSkill", cols, returningCols, appSkillPKCols, conflictTarget, conflictAction, len(vals))
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

		var res AppSkill
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, appSkillPKCols)
}

func (d *AppSkillDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*AppSkill, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("AppSkill") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "AppSkill")
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

	var res AppSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *AppSkillDelegate) buildBulkInsertSQL(q *Queries, batch []*AppSkillCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 3)
	for i, c := range appSkillDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "AppSkill")
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
			case "app_id":
				vals = append(vals, input.App_id)
			case "skillId":
				vals = append(vals, input.SkillId)
			case "priority":
				if input.Priority != nil {
					vals = append(vals, *input.Priority)
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

func applyAppSkillConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanAppSkillRows(rows *sql.Rows, returningCols []string) ([]*AppSkill, error) {
	var records []*AppSkill
	for rows.Next() {
		var res AppSkill
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

func (d *AppSkillDelegate) runCreateMany(ctx context.Context, inputs []*AppSkillCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAppSkillInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppSkillConflictClause(d.client.dialect, queryStr, vals, cols, appSkillPKCols, conflictTarget, conflictAction)

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

func (d *AppSkillDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AppSkillCreate,
	selects *AppSkillSelect,
	omits *AppSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AppSkill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAppSkillInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*AppSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.AppSkill.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.AppSkill.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectAppSkillCols(selects, omits, appSkillPKCols...)
	recordsOut := make([]*AppSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppSkillConflictClause(d.client.dialect, queryStr, vals, cols, appSkillPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanAppSkillRows(rows, returningCols)
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

func (d *AppSkillDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*AppSkillCreate,
	selects *AppSkillSelect,
	omits *AppSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AppSkill, error) {
	batches := partitionAppSkillInputs(d.client.dialect, inputs)
	returningCols := selectAppSkillCols(selects, omits, appSkillPKCols...)
	recordsOut := make([]*AppSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppSkillConflictClause(d.client.dialect, queryStr, vals, cols, appSkillPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("AppSkill") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "AppSkill")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appSkillPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appSkillPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanAppSkillRows(rows, returningCols)
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

type AppSkillConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AppSkillConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppSkillConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppSkillConflictBuilder[B]) Update(fn func(u *AppSkillUpsert)) *B {
	cb.setAction(*AppSkillConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// AppSkillConflictUpdate creates a custom ConflictAction for AppSkill upsert conflicts.
func AppSkillConflictUpdate(fn func(u *AppSkillUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newAppSkillUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type AppSkillUpsert struct {
	App_id   fieldUpsert[string]
	SkillId  fieldUpsert[string]
	Priority fieldUpsert[PriorityType]
}

func newAppSkillUpsert(up *ConflictUpdate) *AppSkillUpsert {
	return &AppSkillUpsert{
		App_id:   fieldUpsert[string]{column: "app_id", update: up},
		SkillId:  fieldUpsert[string]{column: "skillId", update: up},
		Priority: fieldUpsert[PriorityType]{column: "priority", update: up},
	}
}

type AppSkillUpdateBuilder struct {
	*UpdateBuilder[AppSkill, AppSkillSelect, AppSkillOmit]
}

type AppSkillUpdateManyBuilder struct {
	*UpdateManyBuilder[AppSkill]
}

type AppSkillUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[AppSkill, AppSkillSelect, AppSkillOmit]
}

func (b *AppSkillUpdateBuilder) SetApp_id(v string) *AppSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}

func (b *AppSkillUpdateManyBuilder) SetApp_id(v string) *AppSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}

func (b *AppSkillUpdateManyAndReturnBuilder) SetApp_id(v string) *AppSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "app_id", Val: v})
	return b
}
func (b *AppSkillUpdateBuilder) SetSkillId(v string) *AppSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *AppSkillUpdateManyBuilder) SetSkillId(v string) *AppSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *AppSkillUpdateManyAndReturnBuilder) SetSkillId(v string) *AppSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *AppSkillUpdateBuilder) SetPriority(v PriorityType) *AppSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AppSkillUpdateManyBuilder) SetPriority(v PriorityType) *AppSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AppSkillUpdateManyAndReturnBuilder) SetPriority(v PriorityType) *AppSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AppSkillUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[AppSkill]) *AppSkillUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppSkillUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[AppSkill]) *AppSkillUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppSkillUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[AppSkill]) *AppSkillUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppSkillDelegate) Update(where UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *AppSkillUpdateBuilder {
	return &AppSkillUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *AppSkillDelegate) UpdateMany(preds ...PredicateOf[AppSkill]) *AppSkillUpdateManyBuilder {
	return &AppSkillUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[AppSkill]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *AppSkillDelegate) UpdateManyAndReturn(preds ...PredicateOf[AppSkill]) *AppSkillUpdateManyAndReturnBuilder {
	return &AppSkillUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *AppSkillDelegate) buildUpdateSQL(preds []PredicateOf[AppSkill], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "AppSkill")
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

func (d *AppSkillDelegate) executeUpdate(ctx context.Context, where UniquePredicate[AppSkill], additional []PredicateOf[AppSkill], assignments []FieldAssignment, selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
	allWhere := make([]PredicateOf[AppSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToAppSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppSkillUpdateArgs) (*AppSkill, error) {
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
			curr = func(c context.Context, a *AppSkillUpdateArgs) (*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runUpdate(ctx context.Context, preds []PredicateOf[AppSkill], cols []string, vals []any, selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
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
		var res *AppSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AppSkill.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AppSkill.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppSkillCols(selects, omits, appSkillPKCols...)
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
		return nil, &NotFoundError{Model: "AppSkill"}
	}

	var res AppSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AppSkill{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *AppSkillDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[AppSkill], cols []string, vals []any) (int64, error) {
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

func (d *AppSkillDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[AppSkill], cols []string, vals []any, selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "AppSkill"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *AppSkillDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[AppSkill], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToAppSkillUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &AppSkillUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *AppSkillUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppSkillUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *AppSkillDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AppSkill], assignments []FieldAssignment, selects *AppSkillSelect, omits *AppSkillOmit) ([]*AppSkill, error) {
	input, err := assignmentsToAppSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppSkillUpdateManyAndReturnArgs) ([]*AppSkill, error) {
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
			curr = func(c context.Context, a *AppSkillUpdateManyAndReturnArgs) ([]*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AppSkill], cols []string, vals []any, selects *AppSkillSelect, omits *AppSkillOmit) ([]*AppSkill, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[AppSkill]{Where: preds}, selects, omits)
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
		var res []*AppSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AppSkill.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AppSkill.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppSkillCols(selects, omits, appSkillPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanAppSkillRows(rows, returningCols)
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

func (d *AppSkillDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[AppSkill], cols []string, vals []any, selects *AppSkillSelect, omits *AppSkillOmit) ([]*AppSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*AppSkill{}, nil
	}
	return d.runFindMany(ctx, QueryParams[AppSkill]{Where: preds}, selects, omits)
}
func (d *AppSkillDelegate) FindUnique(where UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *FindUniqueBuilder[AppSkill, AppSkillSelect, AppSkillOmit] {
	return &FindUniqueBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AppSkillDelegate) FindFirst(preds ...PredicateOf[AppSkill]) *FindFirstBuilder[AppSkill, AppSkillSelect, AppSkillOmit] {
	return &FindFirstBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AppSkillDelegate) FindMany(preds ...PredicateOf[AppSkill]) *FindManyBuilder[AppSkill, AppSkillSelect, AppSkillOmit] {
	return &FindManyBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AppSkillDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[AppSkill], additional []PredicateOf[AppSkill], selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
	allWhere := make([]PredicateOf[AppSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppSkillFindUniqueArgs) (*AppSkill, error) {
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
			curr = func(c context.Context, a *AppSkillFindUniqueArgs) (*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[AppSkill],
	selects *AppSkillSelect,
	omits *AppSkillOmit,
) (*AppSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppSkillFindFirstArgs) (*AppSkill, error) {
		return d.runFindFirst(c, QueryParams[AppSkill]{
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
			curr = func(c context.Context, a *AppSkillFindFirstArgs) (*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[AppSkill],
	selects *AppSkillSelect,
	omits *AppSkillOmit,
) ([]*AppSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppSkillFindManyArgs) ([]*AppSkill, error) {
		return d.runFindMany(c, QueryParams[AppSkill]{
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
			curr = func(c context.Context, a *AppSkillFindManyArgs) ([]*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runFindUnique(ctx context.Context, where []PredicateOf[AppSkill], selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
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
	returningCols := selectAppSkillCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AppSkill{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *AppSkillDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[AppSkill],
	selects *AppSkillSelect,
	omits *AppSkillOmit,
) (*AppSkill, error) {
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

func (d *AppSkillDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[AppSkill],
	selects *AppSkillSelect,
	omits *AppSkillOmit,
) ([]*AppSkill, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, appSkillPKCols, appSkillUniqueCols, "AppSkill", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, appSkillPKCols, appSkillUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAppSkillCols(selects, omits)

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

func (d *AppSkillDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*AppSkill, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "AppSkill", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res AppSkill
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *AppSkillDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*AppSkill, error) {
	query := buildSelectSQL(d.client, "AppSkill", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*AppSkill, 0)
	for rows.Next() {
		var res AppSkill
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
func (d *AppSkillDelegate) DeleteMany(preds ...PredicateOf[AppSkill]) *DeleteManyBuilder[AppSkill] {
	return &DeleteManyBuilder[AppSkill]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AppSkillDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[AppSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &AppSkillDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *AppSkillDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppSkillDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[AppSkill]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "AppSkill")
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

func (d *AppSkillDelegate) Delete(where UniquePredicate[AppSkill], additional ...PredicateOf[AppSkill]) *DeleteBuilder[AppSkill, AppSkillSelect, AppSkillOmit] {
	return &DeleteBuilder[AppSkill, AppSkillSelect, AppSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *AppSkillDelegate) executeDelete(ctx context.Context, where UniquePredicate[AppSkill], additional []PredicateOf[AppSkill], selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
	allWhere := make([]PredicateOf[AppSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSkillSelect()
	}

	args := &AppSkillDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppSkillDeleteArgs) (*AppSkill, error) {
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
			curr = func(c context.Context, a *AppSkillDeleteArgs) (*AppSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runDelete(ctx context.Context, where []PredicateOf[AppSkill], selects *AppSkillSelect, omits *AppSkillOmit) (*AppSkill, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectAppSkillCols(selects, omits, appSkillPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *AppSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AppSkill.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "AppSkill"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "AppSkill")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[AppSkill]
			pkPreds = append(pkPreds, Predicate[AppSkill]{
				Data: PredicateData{
					Column:   "app_id",
					Operator: "=",
					Value:    res.App_id,
				},
			})
			pkPreds = append(pkPreds, Predicate[AppSkill]{
				Data: PredicateData{
					Column:   "skillId",
					Operator: "=",
					Value:    res.SkillId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "AppSkill")

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
		return nil, &NotFoundError{Model: "AppSkill"}
	}

	var row AppSkill
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *AppSkillDelegate) Count(preds ...PredicateOf[AppSkill]) *CountBuilder[AppSkill] {
	return &CountBuilder[AppSkill]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AppSkillDelegate) executeCount(ctx context.Context, params QueryParams[AppSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &AppSkillCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *AppSkillCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[AppSkill]{
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
			curr = func(c context.Context, a *AppSkillCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppSkillDelegate) runCount(ctx context.Context, params QueryParams[AppSkill]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "AppSkill")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "AppSkill")
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
func (d *AppSkillDelegate) loadRelations(ctx context.Context, records []*AppSkill, selects *AppSkillSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.App != nil {
		relationSelects, relationOmits, relationParams := selects.App.GetRelationParams()
		returningCols := selectAppCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: AppSkill.app_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *AppSkill) string { return p.App_id }),
			"App",
			"name",
			returningCols,
			scanInto(returningCols, (*App).ScanFields),
			directKey(func(c *App) string { return c.Name }),
			setOne(func(p *AppSkill, c *App) { p.App = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading app: %w", err)
		}
		if err := d.client.App.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skill != nil {
		relationSelects, relationOmits, relationParams := selects.Skill.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: AppSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *AppSkill) string { return p.SkillId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *AppSkill, c *Skill) { p.Skill = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skill: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
