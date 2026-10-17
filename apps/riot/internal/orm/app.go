package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// App represents the database model
type App struct {
	Name       string        `db:"name" json:"name"`
	Project_id string        `db:"project_id" json:"project_id"`
	Project    *Project      `json:"project,omitempty"`
	Styles     []*StyleToken `json:"styles,omitempty"`
	Skills     []*AppSkill   `json:"skills,omitempty"`
}

// AppCreate contains model input fields for App creation operations.
//
// Fields for App:
//
//	name       string required
//	project_id string required
type AppCreate struct {
	Name       string `json:"name"`
	Project_id string `json:"project_id"`
}

// colMask returns a bit mask of columns that are set
func (s *AppCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// AppUpdate contains model input fields for App update operations.
type AppUpdate struct {
	Name       *string `json:"name"`
	Project_id *string `json:"project_id"`
}

func (u *AppUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	if u.Project_id != nil {
		cols = append(cols, "project_id")
		vals = append(vals, u.Project_id)
	}
	return cols, vals
}

func assignmentsToAppUpdate(assignments []FieldAssignment) (AppUpdate, error) {
	var input AppUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "name":
			if v, ok := a.Val.(string); ok {
				input.Name = &v
				errs.ValidateString("name", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Name = v
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "project_id":
			if v, ok := a.Val.(string); ok {
				input.Project_id = &v
				errs.ValidateString("project_id", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Project_id = v
			} else {
				errs.Add("project_id", a.Val, "type", "field project_id must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// AppSelect specifies which scalar and relation fields to select for App.
//
// Selectable fields:
//
//	-- Scalars --
//	name       (bool)
//	project_id (bool)
//	-- Relations --
//	project    (Project)
//	styles     ([]StyleToken)
//	skills     ([]AppSkill)
type AppSelect struct {
	Name       bool                  `json:"name"`
	Project_id bool                  `json:"project_id"`
	Project    *ProjectSelect        `json:"project,omitempty"`
	Styles     StyleTokenSelectQuery `json:"styles,omitempty"`
	Skills     AppSkillSelectQuery   `json:"skills,omitempty"`
}

var fullAppSelectVal = &AppSelect{
	Name:       true,
	Project_id: true,
}

func fullAppSelect() *AppSelect {
	return fullAppSelectVal
}

func (s *AppSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Name || s.Project_id
}

func (s *AppSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type AppOmit struct {
	Name       bool `json:"name"`
	Project_id bool `json:"project_id"`
}

type AppSelectQuery interface {
	GetRelationParams() (*AppSelect, *AppOmit, QueryParams[App])
}

func (s *AppSelect) GetRelationParams() (*AppSelect, *AppOmit, QueryParams[App]) {
	return s, nil, QueryParams[App]{}
}

type AppQueryBuilder struct {
	selects *AppSelect
	omits   *AppOmit
	where   []PredicateOf[App]
	take    *int
	skip    *int
	orderBy []OrderBy[App]
	cursor  UniquePredicate[App]
}

func (b *AppQueryBuilder) Where(preds ...PredicateOf[App]) *AppQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AppQueryBuilder) Take(limit int) *AppQueryBuilder {
	b.take = &limit
	return b
}

func (b *AppQueryBuilder) Skip(offset int) *AppQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AppQueryBuilder) OrderBy(orders ...OrderBy[App]) *AppQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AppQueryBuilder) Cursor(where UniquePredicate[App]) *AppQueryBuilder {
	b.cursor = where
	return b
}

func (b *AppQueryBuilder) Select(s AppSelect) *AppQueryBuilder {
	b.selects = &s
	return b
}

func (b *AppQueryBuilder) Omit(o AppOmit) *AppQueryBuilder {
	b.omits = &o
	return b
}

func (b *AppQueryBuilder) GetRelationParams() (*AppSelect, *AppOmit, QueryParams[App]) {
	if b == nil {
		return nil, nil, QueryParams[App]{}
	}
	return b.selects, b.omits, QueryParams[App]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// AppCreateArgs is the input argument passed to App Create extension hooks.
//
// Fields for App:
//
//	name       string required
//	project_id string required
//
// Relations for App:
//
//	project    (Project)
//	styles     ([]StyleToken)
//	skills     ([]AppSkill)
type AppCreateArgs struct {
	// Data contains the model fields to insert.
	Data *AppCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *AppSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// AppCreateManyArgs is the input argument passed to App CreateMany extension hooks.
//
// Fields for App:
//
//	name       string required
//	project_id string required
type AppCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppCreateManyArgs) AppendData(builders ...*AppCreateBuilder) *AppCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToAppCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppCreateManyAndReturnArgs is the input argument passed to App CreateManyAndReturn extension hooks.
//
// Fields for App:
//
//	name       string required
//	project_id string required
//
// Relations for App:
//
//	project    (Project)
//	styles     ([]StyleToken)
//	skills     ([]AppSkill)
type AppCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AppCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *AppSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AppCreateManyAndReturnArgs) AppendData(builders ...*AppCreateBuilder) *AppCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToAppCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AppFindUniqueArgs is the input argument passed to App FindUnique extension hooks.
type AppFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[App]
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSelect
}

func (a *AppFindUniqueArgs) SetWhere(unique UniquePredicate[App], additional ...PredicateOf[App]) *AppFindUniqueArgs {
	a.Where = make([]PredicateOf[App], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppFindFirstArgs is the input argument passed to App FindFirst extension hooks.
type AppFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[App]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[App]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSelect
}

func (a *AppFindFirstArgs) SetWhere(preds ...PredicateOf[App]) *AppFindFirstArgs {
	a.Where = preds
	return a
}

func (a *AppFindFirstArgs) SetOrderBy(orders ...OrderBy[App]) *AppFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *AppFindFirstArgs) SetCursor(cursor UniquePredicate[App]) *AppFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *AppFindFirstArgs) SetSkip(n int) *AppFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *AppFindFirstArgs) SetTake(n int) *AppFindFirstArgs {
	a.Take = &n
	return a
}

// AppFindManyArgs is the input argument passed to App FindMany extension hooks.
type AppFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[App]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[App]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AppSelect
}

func (a *AppFindManyArgs) SetWhere(preds ...PredicateOf[App]) *AppFindManyArgs {
	a.Where = preds
	return a
}

func (a *AppFindManyArgs) SetOrderBy(orders ...OrderBy[App]) *AppFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *AppFindManyArgs) SetCursor(cursor UniquePredicate[App]) *AppFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *AppFindManyArgs) SetSkip(n int) *AppFindManyArgs {
	a.Skip = &n
	return a
}

func (a *AppFindManyArgs) SetTake(n int) *AppFindManyArgs {
	a.Take = &n
	return a
}

// AppCountArgs is the input argument passed to App Count extension hooks.
type AppCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *AppCountArgs) SetWhere(preds ...PredicateOf[App]) *AppCountArgs {
	a.Where = preds
	return a
}

func (a *AppCountArgs) SetSkip(n int) *AppCountArgs {
	a.Skip = &n
	return a
}

func (a *AppCountArgs) SetTake(n int) *AppCountArgs {
	a.Take = &n
	return a
}

// AppDeleteArgs is the input argument passed to App Delete extension hooks.
type AppDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[App]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *AppSelect
}

func (a *AppDeleteArgs) SetWhere(unique UniquePredicate[App], additional ...PredicateOf[App]) *AppDeleteArgs {
	a.Where = make([]PredicateOf[App], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppDeleteManyArgs is the input argument passed to App DeleteMany extension hooks.
type AppDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
}

func (a *AppDeleteManyArgs) SetWhere(preds ...PredicateOf[App]) *AppDeleteManyArgs {
	a.Where = preds
	return a
}

// AppUpdateArgs is the input argument passed to App Update extension hooks.
type AppUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[App]
	// Data contains the model fields to update.
	Data *AppUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppSelect
}

func (a *AppUpdateArgs) SetWhere(unique UniquePredicate[App], additional ...PredicateOf[App]) *AppUpdateArgs {
	a.Where = make([]PredicateOf[App], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AppUpdateManyArgs is the input argument passed to App UpdateMany extension hooks.
type AppUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
	// Data contains the model fields to update.
	Data *AppUpdate
}

func (a *AppUpdateManyArgs) SetWhere(preds ...PredicateOf[App]) *AppUpdateManyArgs {
	a.Where = preds
	return a
}

// AppUpdateManyAndReturnArgs is the input argument passed to App UpdateManyAndReturn extension hooks.
type AppUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[App]
	// Data contains the model fields to update.
	Data *AppUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AppSelect
}

func (a *AppUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[App]) *AppUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type AppCreateQuery = func(ctx context.Context, args *AppCreateArgs) (*App, error)
type AppCreateManyQuery = func(ctx context.Context, args *AppCreateManyArgs) (int64, error)
type AppCreateManyAndReturnQuery = func(ctx context.Context, args *AppCreateManyAndReturnArgs) ([]*App, error)
type AppFindUniqueQuery = func(ctx context.Context, args *AppFindUniqueArgs) (*App, error)
type AppFindFirstQuery = func(ctx context.Context, args *AppFindFirstArgs) (*App, error)
type AppFindManyQuery = func(ctx context.Context, args *AppFindManyArgs) ([]*App, error)
type AppDeleteQuery = func(ctx context.Context, args *AppDeleteArgs) (*App, error)
type AppDeleteManyQuery = func(ctx context.Context, args *AppDeleteManyArgs) (int64, error)
type AppCountQuery = func(ctx context.Context, args *AppCountArgs) (int64, error)
type AppUpdateQuery = func(ctx context.Context, args *AppUpdateArgs) (*App, error)
type AppUpdateManyQuery = func(ctx context.Context, args *AppUpdateManyArgs) (int64, error)
type AppUpdateManyAndReturnQuery = func(ctx context.Context, args *AppUpdateManyAndReturnArgs) ([]*App, error)

type AppExtension struct {
	Create              func(ctx context.Context, args *AppCreateArgs, next AppCreateQuery) (*App, error)
	CreateMany          func(ctx context.Context, args *AppCreateManyArgs, next AppCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *AppCreateManyAndReturnArgs, next AppCreateManyAndReturnQuery) ([]*App, error)
	FindUnique          func(ctx context.Context, args *AppFindUniqueArgs, next AppFindUniqueQuery) (*App, error)
	FindFirst           func(ctx context.Context, args *AppFindFirstArgs, next AppFindFirstQuery) (*App, error)
	FindMany            func(ctx context.Context, args *AppFindManyArgs, next AppFindManyQuery) ([]*App, error)
	Delete              func(ctx context.Context, args *AppDeleteArgs, next AppDeleteQuery) (*App, error)
	DeleteMany          func(ctx context.Context, args *AppDeleteManyArgs, next AppDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *AppCountArgs, next AppCountQuery) (int64, error)
	Update              func(ctx context.Context, args *AppUpdateArgs, next AppUpdateQuery) (*App, error)
	UpdateMany          func(ctx context.Context, args *AppUpdateManyArgs, next AppUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *AppUpdateManyAndReturnArgs, next AppUpdateManyAndReturnQuery) ([]*App, error)
}

type AppDelegate struct {
	client     *Queries
	extensions []AppExtension
}

func (d *AppDelegate) Use(exts ...AppExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *App) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "name":
			targets[i] = &m.Name
		case "project_id":
			targets[i] = &m.Project_id
		}
	}
	return targets
}

var appDefaultCols = []string{
	"name",
	"project_id",
}

var appPKCols = []string{
	"name",
}

var appUniqueCols = []string{
	"name",
}

func selectAppCols(selects *AppSelect, omits *AppOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return appDefaultCols
	}

	anySelected := selects != nil && (selects.Name || selects.Project_id || selects.Project != nil || selects.Styles != nil || selects.Skills != nil)

	specs := []colSpec{
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, selects != nil && selects.hasAnyRelation()},
		{"project_id", selects != nil && selects.Project_id, omits != nil && omits.Project_id, selects != nil && selects.Project != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *AppSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Project != nil || s.Styles != nil || s.Skills != nil
}

type AppUpsertBuilder struct {
	*CreateBuilder[App, AppSelect, AppOmit]
}

func (b *AppUpsertBuilder) Select(s AppSelect) *AppUpsertBuilder {
	b.selects = &s
	return b
}

func (b *AppUpsertBuilder) Omit(o AppOmit) *AppUpsertBuilder {
	b.omits = &o
	return b
}

type AppCreateBuilder struct {
	*CreateBuilder[App, AppSelect, AppOmit]
}

func (b *AppCreateBuilder) Select(s AppSelect) *AppCreateBuilder {
	b.selects = &s
	return b
}

func (b *AppCreateBuilder) Omit(o AppOmit) *AppCreateBuilder {
	b.omits = &o
	return b
}

func (b *AppCreateBuilder) OnConflict(target UniqueConstraintTarget) *AppConflictBuilder[AppUpsertBuilder] {
	upsertBuilder := &AppUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &AppConflictBuilder[AppUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AppCreateBuilder) SetName(v string) *AppCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *AppCreateBuilder) SetProject_id(v string) *AppCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "project_id", Val: v})
	return b
}

func (b *AppCreateBuilder) Assignments(assignments ...FieldAssignmentOf[App]) *AppCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppDelegate) Create() *AppCreateBuilder {
	return &AppCreateBuilder{
		CreateBuilder: &CreateBuilder[App, AppSelect, AppOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedAppName       uint64 = 1 << 0
	providedAppProject_id uint64 = 1 << 1
)

func assignmentsToAppCreate(assignments []FieldAssignment) (AppCreate, error) {
	var input AppCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "name":
			provided |= providedAppName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "project_id":
			provided |= providedAppProject_id
			if v, ok := a.Val.(string); ok {
				input.Project_id = v
				errs.ValidateString("project_id", v, true, 0, false, false)
			} else {
				errs.Add("project_id", a.Val, "type", "field project_id must be of type string")
			}
		}
	}
	if provided&providedAppName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}
	if provided&providedAppProject_id == 0 {
		errs.Add("project_id", "", "required", "field Project_id is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AppCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "name")
	if s.Name != "" {
		vals = append(vals, s.Name)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "project_id")
	vals = append(vals, s.Project_id)
	return
}

func partitionAppInputs(dialect Dialect, inputs []*AppCreate) [][]*AppCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AppCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AppCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AppCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AppCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AppCreate{inputs}
}

func (d *AppDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AppSelect, omits *AppOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*App, error) {
	input, err := assignmentsToAppCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAppCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppCreateArgs) (*App, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectAppCols(a.Select, omits)
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
			curr = func(c context.Context, a *AppCreateArgs) (*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type AppCreateManyBuilder struct {
	*CreateManyBuilder[App]
}

func (b *AppCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AppConflictBuilder[AppCreateManyBuilder] {
	return &AppConflictBuilder[AppCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AppCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[App, AppSelect, AppOmit]
}

func (b *AppCreateManyAndReturnBuilder) Select(s AppSelect) *AppCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *AppCreateManyAndReturnBuilder) Omit(o AppOmit) *AppCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *AppCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AppConflictBuilder[AppCreateManyAndReturnBuilder] {
	return &AppConflictBuilder[AppCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToAppRecordInputs(builders []*AppCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *AppDelegate) CreateMany(builders ...*AppCreateBuilder) *AppCreateManyBuilder {
	return &AppCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[App]{
			records:  createBuildersToAppRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AppDelegate) CreateManyAndReturn(builders ...*AppCreateBuilder) *AppCreateManyAndReturnBuilder {
	return &AppCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[App, AppSelect, AppOmit]{
			records:  createBuildersToAppRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToAppCreateInputs(records []RecordInput) ([]*AppCreate, error) {
	structs := make([]AppCreate, len(records))
	inputs := make([]*AppCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToAppCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *AppDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToAppCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &AppCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AppSelect, omits *AppOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*App, error) {
	inputs, err := recordsToAppCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AppCreateManyAndReturnArgs) ([]*App, error) {
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
			curr = func(c context.Context, a *AppCreateManyAndReturnArgs) ([]*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *AppSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*App, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *App
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.App.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.App.loadRelations(ctx, []*App{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "App", cols, returningCols, appPKCols, conflictTarget, conflictAction, len(vals))
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

		var res App
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, appPKCols)
}

func (d *AppDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*App, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("App") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "App")
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

	var res App
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *AppDelegate) buildBulkInsertSQL(q *Queries, batch []*AppCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range appDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "App")
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
			case "name":
				if input.Name != "" {
					vals = append(vals, input.Name)
				} else {
					vals = append(vals, generateCUID())
				}
			case "project_id":
				vals = append(vals, input.Project_id)
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

func applyAppConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanAppRows(rows *sql.Rows, returningCols []string) ([]*App, error) {
	var records []*App
	for rows.Next() {
		var res App
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

func (d *AppDelegate) runCreateMany(ctx context.Context, inputs []*AppCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAppInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppConflictClause(d.client.dialect, queryStr, vals, cols, appPKCols, conflictTarget, conflictAction)

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

func (d *AppDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AppCreate,
	selects *AppSelect,
	omits *AppOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*App, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAppInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*App
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.App.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.App.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectAppCols(selects, omits, appPKCols...)
	recordsOut := make([]*App, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppConflictClause(d.client.dialect, queryStr, vals, cols, appPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanAppRows(rows, returningCols)
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

func (d *AppDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*AppCreate,
	selects *AppSelect,
	omits *AppOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*App, error) {
	batches := partitionAppInputs(d.client.dialect, inputs)
	returningCols := selectAppCols(selects, omits, appPKCols...)
	recordsOut := make([]*App, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAppConflictClause(d.client.dialect, queryStr, vals, cols, appPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("App") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "App")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, appPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanAppRows(rows, returningCols)
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

type AppConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AppConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AppConflictBuilder[B]) Update(fn func(u *AppUpsert)) *B {
	cb.setAction(*AppConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// AppConflictUpdate creates a custom ConflictAction for App upsert conflicts.
func AppConflictUpdate(fn func(u *AppUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newAppUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type AppUpsert struct {
	Name       fieldUpsert[string]
	Project_id fieldUpsert[string]
}

func newAppUpsert(up *ConflictUpdate) *AppUpsert {
	return &AppUpsert{
		Name:       fieldUpsert[string]{column: "name", update: up},
		Project_id: fieldUpsert[string]{column: "project_id", update: up},
	}
}

type AppUpdateBuilder struct {
	*UpdateBuilder[App, AppSelect, AppOmit]
}

type AppUpdateManyBuilder struct {
	*UpdateManyBuilder[App]
}

type AppUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[App, AppSelect, AppOmit]
}

func (b *AppUpdateBuilder) SetName(v string) *AppUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *AppUpdateManyBuilder) SetName(v string) *AppUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *AppUpdateManyAndReturnBuilder) SetName(v string) *AppUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *AppUpdateBuilder) SetProject_id(v string) *AppUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "project_id", Val: v})
	return b
}

func (b *AppUpdateManyBuilder) SetProject_id(v string) *AppUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "project_id", Val: v})
	return b
}

func (b *AppUpdateManyAndReturnBuilder) SetProject_id(v string) *AppUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "project_id", Val: v})
	return b
}

func (b *AppUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[App]) *AppUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[App]) *AppUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AppUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[App]) *AppUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AppDelegate) Update(where UniquePredicate[App], additional ...PredicateOf[App]) *AppUpdateBuilder {
	return &AppUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[App, AppSelect, AppOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *AppDelegate) UpdateMany(preds ...PredicateOf[App]) *AppUpdateManyBuilder {
	return &AppUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[App]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *AppDelegate) UpdateManyAndReturn(preds ...PredicateOf[App]) *AppUpdateManyAndReturnBuilder {
	return &AppUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[App, AppSelect, AppOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *AppDelegate) buildUpdateSQL(preds []PredicateOf[App], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "App")
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

func (d *AppDelegate) executeUpdate(ctx context.Context, where UniquePredicate[App], additional []PredicateOf[App], assignments []FieldAssignment, selects *AppSelect, omits *AppOmit) (*App, error) {
	allWhere := make([]PredicateOf[App], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToAppUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppUpdateArgs) (*App, error) {
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
			curr = func(c context.Context, a *AppUpdateArgs) (*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runUpdate(ctx context.Context, preds []PredicateOf[App], cols []string, vals []any, selects *AppSelect, omits *AppOmit) (*App, error) {
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
		var res *App
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.App.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.App.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppCols(selects, omits, appPKCols...)
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
		return nil, &NotFoundError{Model: "App"}
	}

	var res App
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*App{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *AppDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[App], cols []string, vals []any) (int64, error) {
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

func (d *AppDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[App], cols []string, vals []any, selects *AppSelect, omits *AppOmit) (*App, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "App"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *AppDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[App], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToAppUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &AppUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *AppUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *AppDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[App], assignments []FieldAssignment, selects *AppSelect, omits *AppOmit) ([]*App, error) {
	input, err := assignmentsToAppUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AppUpdateManyAndReturnArgs) ([]*App, error) {
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
			curr = func(c context.Context, a *AppUpdateManyAndReturnArgs) ([]*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[App], cols []string, vals []any, selects *AppSelect, omits *AppOmit) ([]*App, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[App]{Where: preds}, selects, omits)
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
		var res []*App
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.App.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.App.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAppCols(selects, omits, appPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanAppRows(rows, returningCols)
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

func (d *AppDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[App], cols []string, vals []any, selects *AppSelect, omits *AppOmit) ([]*App, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*App{}, nil
	}
	return d.runFindMany(ctx, QueryParams[App]{Where: preds}, selects, omits)
}
func (d *AppDelegate) FindUnique(where UniquePredicate[App], additional ...PredicateOf[App]) *FindUniqueBuilder[App, AppSelect, AppOmit] {
	return &FindUniqueBuilder[App, AppSelect, AppOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AppDelegate) FindFirst(preds ...PredicateOf[App]) *FindFirstBuilder[App, AppSelect, AppOmit] {
	return &FindFirstBuilder[App, AppSelect, AppOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AppDelegate) FindMany(preds ...PredicateOf[App]) *FindManyBuilder[App, AppSelect, AppOmit] {
	return &FindManyBuilder[App, AppSelect, AppOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AppDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[App], additional []PredicateOf[App], selects *AppSelect, omits *AppOmit) (*App, error) {
	allWhere := make([]PredicateOf[App], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppFindUniqueArgs) (*App, error) {
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
			curr = func(c context.Context, a *AppFindUniqueArgs) (*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[App],
	selects *AppSelect,
	omits *AppOmit,
) (*App, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppFindFirstArgs) (*App, error) {
		return d.runFindFirst(c, QueryParams[App]{
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
			curr = func(c context.Context, a *AppFindFirstArgs) (*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[App],
	selects *AppSelect,
	omits *AppOmit,
) ([]*App, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AppFindManyArgs) ([]*App, error) {
		return d.runFindMany(c, QueryParams[App]{
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
			curr = func(c context.Context, a *AppFindManyArgs) ([]*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runFindUnique(ctx context.Context, where []PredicateOf[App], selects *AppSelect, omits *AppOmit) (*App, error) {
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
	returningCols := selectAppCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*App{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *AppDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[App],
	selects *AppSelect,
	omits *AppOmit,
) (*App, error) {
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

func (d *AppDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[App],
	selects *AppSelect,
	omits *AppOmit,
) ([]*App, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, appPKCols, appUniqueCols, "App", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, appPKCols, appUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAppCols(selects, omits)

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

func (d *AppDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*App, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "App", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res App
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *AppDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*App, error) {
	query := buildSelectSQL(d.client, "App", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*App, 0)
	for rows.Next() {
		var res App
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
func (d *AppDelegate) DeleteMany(preds ...PredicateOf[App]) *DeleteManyBuilder[App] {
	return &DeleteManyBuilder[App]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AppDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[App]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &AppDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *AppDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AppDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[App]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "App")
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

func (d *AppDelegate) Delete(where UniquePredicate[App], additional ...PredicateOf[App]) *DeleteBuilder[App, AppSelect, AppOmit] {
	return &DeleteBuilder[App, AppSelect, AppOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *AppDelegate) executeDelete(ctx context.Context, where UniquePredicate[App], additional []PredicateOf[App], selects *AppSelect, omits *AppOmit) (*App, error) {
	allWhere := make([]PredicateOf[App], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAppSelect()
	}

	args := &AppDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AppDeleteArgs) (*App, error) {
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
			curr = func(c context.Context, a *AppDeleteArgs) (*App, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runDelete(ctx context.Context, where []PredicateOf[App], selects *AppSelect, omits *AppOmit) (*App, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectAppCols(selects, omits, appPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *App
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.App.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "App"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "App")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[App]
			pkPreds = append(pkPreds, Predicate[App]{
				Data: PredicateData{
					Column:   "name",
					Operator: "=",
					Value:    res.Name,
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
	d.client.dialect.WriteQuotedIdent(&sb, "App")

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
		return nil, &NotFoundError{Model: "App"}
	}

	var row App
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *AppDelegate) Count(preds ...PredicateOf[App]) *CountBuilder[App] {
	return &CountBuilder[App]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AppDelegate) executeCount(ctx context.Context, params QueryParams[App]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &AppCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *AppCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[App]{
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
			curr = func(c context.Context, a *AppCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AppDelegate) runCount(ctx context.Context, params QueryParams[App]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "App")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "App")
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
func (d *AppDelegate) loadRelations(ctx context.Context, records []*App, selects *AppSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Project != nil {
		relationSelects, relationOmits, relationParams := selects.Project.GetRelationParams()
		returningCols := selectProjectCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: App.project_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *App) string { return p.Project_id }),
			"Project",
			"name",
			returningCols,
			scanInto(returningCols, (*Project).ScanFields),
			directKey(func(c *Project) string { return c.Name }),
			setOne(func(p *App, c *Project) { p.Project = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading project: %w", err)
		}
		if err := d.client.Project.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Styles != nil {
		relationSelects, relationOmits, relationParams := selects.Styles.GetRelationParams()
		returningCols := selectStyleTokenCols(relationSelects, relationOmits, "app_id")
		// Inverse holds the FK: StyleToken.app_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *App) string { return p.Name }),
			"StyleToken",
			"app_id",
			returningCols,
			scanInto(returningCols, (*StyleToken).ScanFields),
			directKey(func(c *StyleToken) string { return c.App_id }),
			appendMany(func(p *App) *[]*StyleToken { return &p.Styles }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading styles: %w", err)
		}
		if err := d.client.StyleToken.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skills != nil {
		relationSelects, relationOmits, relationParams := selects.Skills.GetRelationParams()
		returningCols := selectAppSkillCols(relationSelects, relationOmits, "app_id")
		// Inverse holds the FK: AppSkill.app_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *App) string { return p.Name }),
			"AppSkill",
			"app_id",
			returningCols,
			scanInto(returningCols, (*AppSkill).ScanFields),
			directKey(func(c *AppSkill) string { return c.App_id }),
			appendMany(func(p *App) *[]*AppSkill { return &p.Skills }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skills: %w", err)
		}
		if err := d.client.AppSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
