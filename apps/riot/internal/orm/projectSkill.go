package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ProjectSkill represents the database model
type ProjectSkill struct {
	ProjectName string   `db:"projectName" json:"projectName"`
	SkillId     string   `db:"skillId" json:"skillId"`
	Project     *Project `json:"project,omitempty"`
	Skill       *Skill   `json:"skill,omitempty"`
}

// ProjectSkillCreate contains model input fields for ProjectSkill creation operations.
//
// Fields for ProjectSkill:
//
//	projectName string required
//	skillId     string required
type ProjectSkillCreate struct {
	ProjectName string `json:"projectName"`
	SkillId     string `json:"skillId"`
}

// colMask returns a bit mask of columns that are set
func (s *ProjectSkillCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// ProjectSkillUpdate contains model input fields for ProjectSkill update operations.
type ProjectSkillUpdate struct {
	ProjectName *string `json:"projectName"`
	SkillId     *string `json:"skillId"`
}

func (u *ProjectSkillUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.ProjectName != nil {
		cols = append(cols, "projectName")
		vals = append(vals, u.ProjectName)
	}
	if u.SkillId != nil {
		cols = append(cols, "skillId")
		vals = append(vals, u.SkillId)
	}
	return cols, vals
}

func assignmentsToProjectSkillUpdate(assignments []FieldAssignment) (ProjectSkillUpdate, error) {
	var input ProjectSkillUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "projectName":
			if v, ok := a.Val.(string); ok {
				input.ProjectName = &v
				errs.ValidateString("projectName", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.ProjectName = v
			} else {
				errs.Add("projectName", a.Val, "type", "field projectName must be of type string")
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
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// ProjectSkillSelect specifies which scalar and relation fields to select for ProjectSkill.
//
// Selectable fields:
//
//	-- Scalars --
//	projectName (bool)
//	skillId     (bool)
//	-- Relations --
//	project     (Project)
//	skill       (Skill)
type ProjectSkillSelect struct {
	ProjectName bool           `json:"projectName"`
	SkillId     bool           `json:"skillId"`
	Project     *ProjectSelect `json:"project,omitempty"`
	Skill       *SkillSelect   `json:"skill,omitempty"`
}

var fullProjectSkillSelectVal = &ProjectSkillSelect{
	ProjectName: true,
	SkillId:     true,
}

func fullProjectSkillSelect() *ProjectSkillSelect {
	return fullProjectSkillSelectVal
}

func (s *ProjectSkillSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.ProjectName || s.SkillId
}

func (s *ProjectSkillSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type ProjectSkillOmit struct {
	ProjectName bool `json:"projectName"`
	SkillId     bool `json:"skillId"`
}

type ProjectSkillSelectQuery interface {
	GetRelationParams() (*ProjectSkillSelect, *ProjectSkillOmit, QueryParams[ProjectSkill])
}

func (s *ProjectSkillSelect) GetRelationParams() (*ProjectSkillSelect, *ProjectSkillOmit, QueryParams[ProjectSkill]) {
	return s, nil, QueryParams[ProjectSkill]{}
}

type ProjectSkillQueryBuilder struct {
	selects *ProjectSkillSelect
	omits   *ProjectSkillOmit
	where   []PredicateOf[ProjectSkill]
	take    *int
	skip    *int
	orderBy []OrderBy[ProjectSkill]
	cursor  UniquePredicate[ProjectSkill]
}

func (b *ProjectSkillQueryBuilder) Where(preds ...PredicateOf[ProjectSkill]) *ProjectSkillQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *ProjectSkillQueryBuilder) Take(limit int) *ProjectSkillQueryBuilder {
	b.take = &limit
	return b
}

func (b *ProjectSkillQueryBuilder) Skip(offset int) *ProjectSkillQueryBuilder {
	b.skip = &offset
	return b
}

func (b *ProjectSkillQueryBuilder) OrderBy(orders ...OrderBy[ProjectSkill]) *ProjectSkillQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *ProjectSkillQueryBuilder) Cursor(where UniquePredicate[ProjectSkill]) *ProjectSkillQueryBuilder {
	b.cursor = where
	return b
}

func (b *ProjectSkillQueryBuilder) Select(s ProjectSkillSelect) *ProjectSkillQueryBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSkillQueryBuilder) Omit(o ProjectSkillOmit) *ProjectSkillQueryBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSkillQueryBuilder) GetRelationParams() (*ProjectSkillSelect, *ProjectSkillOmit, QueryParams[ProjectSkill]) {
	if b == nil {
		return nil, nil, QueryParams[ProjectSkill]{}
	}
	return b.selects, b.omits, QueryParams[ProjectSkill]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// ProjectSkillCreateArgs is the input argument passed to ProjectSkill Create extension hooks.
//
// Fields for ProjectSkill:
//
//	projectName string required
//	skillId     string required
//
// Relations for ProjectSkill:
//
//	project     (Project)
//	skill       (Skill)
type ProjectSkillCreateArgs struct {
	// Data contains the model fields to insert.
	Data *ProjectSkillCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *ProjectSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// ProjectSkillCreateManyArgs is the input argument passed to ProjectSkill CreateMany extension hooks.
//
// Fields for ProjectSkill:
//
//	projectName string required
//	skillId     string required
type ProjectSkillCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectSkillCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectSkillCreateManyArgs) AppendData(builders ...*ProjectSkillCreateBuilder) *ProjectSkillCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectSkillCreateManyAndReturnArgs is the input argument passed to ProjectSkill CreateManyAndReturn extension hooks.
//
// Fields for ProjectSkill:
//
//	projectName string required
//	skillId     string required
//
// Relations for ProjectSkill:
//
//	project     (Project)
//	skill       (Skill)
type ProjectSkillCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectSkillCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *ProjectSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectSkillCreateManyAndReturnArgs) AppendData(builders ...*ProjectSkillCreateBuilder) *ProjectSkillCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectSkillFindUniqueArgs is the input argument passed to ProjectSkill FindUnique extension hooks.
type ProjectSkillFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSkill]
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillFindUniqueArgs) SetWhere(unique UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *ProjectSkillFindUniqueArgs {
	a.Where = make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSkillFindFirstArgs is the input argument passed to ProjectSkill FindFirst extension hooks.
type ProjectSkillFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[ProjectSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[ProjectSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillFindFirstArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillFindFirstArgs {
	a.Where = preds
	return a
}

func (a *ProjectSkillFindFirstArgs) SetOrderBy(orders ...OrderBy[ProjectSkill]) *ProjectSkillFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectSkillFindFirstArgs) SetCursor(cursor UniquePredicate[ProjectSkill]) *ProjectSkillFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectSkillFindFirstArgs) SetSkip(n int) *ProjectSkillFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSkillFindFirstArgs) SetTake(n int) *ProjectSkillFindFirstArgs {
	a.Take = &n
	return a
}

// ProjectSkillFindManyArgs is the input argument passed to ProjectSkill FindMany extension hooks.
type ProjectSkillFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[ProjectSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[ProjectSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillFindManyArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillFindManyArgs {
	a.Where = preds
	return a
}

func (a *ProjectSkillFindManyArgs) SetOrderBy(orders ...OrderBy[ProjectSkill]) *ProjectSkillFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectSkillFindManyArgs) SetCursor(cursor UniquePredicate[ProjectSkill]) *ProjectSkillFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectSkillFindManyArgs) SetSkip(n int) *ProjectSkillFindManyArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSkillFindManyArgs) SetTake(n int) *ProjectSkillFindManyArgs {
	a.Take = &n
	return a
}

// ProjectSkillCountArgs is the input argument passed to ProjectSkill Count extension hooks.
type ProjectSkillCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *ProjectSkillCountArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillCountArgs {
	a.Where = preds
	return a
}

func (a *ProjectSkillCountArgs) SetSkip(n int) *ProjectSkillCountArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSkillCountArgs) SetTake(n int) *ProjectSkillCountArgs {
	a.Take = &n
	return a
}

// ProjectSkillDeleteArgs is the input argument passed to ProjectSkill Delete extension hooks.
type ProjectSkillDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSkill]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillDeleteArgs) SetWhere(unique UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *ProjectSkillDeleteArgs {
	a.Where = make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSkillDeleteManyArgs is the input argument passed to ProjectSkill DeleteMany extension hooks.
type ProjectSkillDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
}

func (a *ProjectSkillDeleteManyArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillDeleteManyArgs {
	a.Where = preds
	return a
}

// ProjectSkillUpdateArgs is the input argument passed to ProjectSkill Update extension hooks.
type ProjectSkillUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSkill]
	// Data contains the model fields to update.
	Data *ProjectSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillUpdateArgs) SetWhere(unique UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateArgs {
	a.Where = make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSkillUpdateManyArgs is the input argument passed to ProjectSkill UpdateMany extension hooks.
type ProjectSkillUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
	// Data contains the model fields to update.
	Data *ProjectSkillUpdate
}

func (a *ProjectSkillUpdateManyArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateManyArgs {
	a.Where = preds
	return a
}

// ProjectSkillUpdateManyAndReturnArgs is the input argument passed to ProjectSkill UpdateManyAndReturn extension hooks.
type ProjectSkillUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSkill]
	// Data contains the model fields to update.
	Data *ProjectSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSkillSelect
}

func (a *ProjectSkillUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type ProjectSkillCreateQuery = func(ctx context.Context, args *ProjectSkillCreateArgs) (*ProjectSkill, error)
type ProjectSkillCreateManyQuery = func(ctx context.Context, args *ProjectSkillCreateManyArgs) (int64, error)
type ProjectSkillCreateManyAndReturnQuery = func(ctx context.Context, args *ProjectSkillCreateManyAndReturnArgs) ([]*ProjectSkill, error)
type ProjectSkillFindUniqueQuery = func(ctx context.Context, args *ProjectSkillFindUniqueArgs) (*ProjectSkill, error)
type ProjectSkillFindFirstQuery = func(ctx context.Context, args *ProjectSkillFindFirstArgs) (*ProjectSkill, error)
type ProjectSkillFindManyQuery = func(ctx context.Context, args *ProjectSkillFindManyArgs) ([]*ProjectSkill, error)
type ProjectSkillDeleteQuery = func(ctx context.Context, args *ProjectSkillDeleteArgs) (*ProjectSkill, error)
type ProjectSkillDeleteManyQuery = func(ctx context.Context, args *ProjectSkillDeleteManyArgs) (int64, error)
type ProjectSkillCountQuery = func(ctx context.Context, args *ProjectSkillCountArgs) (int64, error)
type ProjectSkillUpdateQuery = func(ctx context.Context, args *ProjectSkillUpdateArgs) (*ProjectSkill, error)
type ProjectSkillUpdateManyQuery = func(ctx context.Context, args *ProjectSkillUpdateManyArgs) (int64, error)
type ProjectSkillUpdateManyAndReturnQuery = func(ctx context.Context, args *ProjectSkillUpdateManyAndReturnArgs) ([]*ProjectSkill, error)

type ProjectSkillExtension struct {
	Create              func(ctx context.Context, args *ProjectSkillCreateArgs, next ProjectSkillCreateQuery) (*ProjectSkill, error)
	CreateMany          func(ctx context.Context, args *ProjectSkillCreateManyArgs, next ProjectSkillCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *ProjectSkillCreateManyAndReturnArgs, next ProjectSkillCreateManyAndReturnQuery) ([]*ProjectSkill, error)
	FindUnique          func(ctx context.Context, args *ProjectSkillFindUniqueArgs, next ProjectSkillFindUniqueQuery) (*ProjectSkill, error)
	FindFirst           func(ctx context.Context, args *ProjectSkillFindFirstArgs, next ProjectSkillFindFirstQuery) (*ProjectSkill, error)
	FindMany            func(ctx context.Context, args *ProjectSkillFindManyArgs, next ProjectSkillFindManyQuery) ([]*ProjectSkill, error)
	Delete              func(ctx context.Context, args *ProjectSkillDeleteArgs, next ProjectSkillDeleteQuery) (*ProjectSkill, error)
	DeleteMany          func(ctx context.Context, args *ProjectSkillDeleteManyArgs, next ProjectSkillDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *ProjectSkillCountArgs, next ProjectSkillCountQuery) (int64, error)
	Update              func(ctx context.Context, args *ProjectSkillUpdateArgs, next ProjectSkillUpdateQuery) (*ProjectSkill, error)
	UpdateMany          func(ctx context.Context, args *ProjectSkillUpdateManyArgs, next ProjectSkillUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *ProjectSkillUpdateManyAndReturnArgs, next ProjectSkillUpdateManyAndReturnQuery) ([]*ProjectSkill, error)
}

type ProjectSkillDelegate struct {
	client     *Queries
	extensions []ProjectSkillExtension
}

func (d *ProjectSkillDelegate) Use(exts ...ProjectSkillExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *ProjectSkill) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "projectName":
			targets[i] = &m.ProjectName
		case "skillId":
			targets[i] = &m.SkillId
		}
	}
	return targets
}

var projectSkillDefaultCols = []string{
	"projectName",
	"skillId",
}

var projectSkillPKCols = []string{
	"projectName",
	"skillId",
}

var projectSkillUniqueCols = []string{}

func selectProjectSkillCols(selects *ProjectSkillSelect, omits *ProjectSkillOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return projectSkillDefaultCols
	}

	anySelected := selects != nil && (selects.ProjectName || selects.SkillId || selects.Project != nil || selects.Skill != nil)

	specs := []colSpec{
		{"projectName", selects != nil && selects.ProjectName, omits != nil && omits.ProjectName, selects != nil && selects.Project != nil},
		{"skillId", selects != nil && selects.SkillId, omits != nil && omits.SkillId, selects != nil && selects.Skill != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *ProjectSkillSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Project != nil || s.Skill != nil
}

type ProjectSkillUpsertBuilder struct {
	*CreateBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]
}

func (b *ProjectSkillUpsertBuilder) Select(s ProjectSkillSelect) *ProjectSkillUpsertBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSkillUpsertBuilder) Omit(o ProjectSkillOmit) *ProjectSkillUpsertBuilder {
	b.omits = &o
	return b
}

type ProjectSkillCreateBuilder struct {
	*CreateBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]
}

func (b *ProjectSkillCreateBuilder) Select(s ProjectSkillSelect) *ProjectSkillCreateBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSkillCreateBuilder) Omit(o ProjectSkillOmit) *ProjectSkillCreateBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSkillCreateBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSkillConflictBuilder[ProjectSkillUpsertBuilder] {
	upsertBuilder := &ProjectSkillUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &ProjectSkillConflictBuilder[ProjectSkillUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *ProjectSkillCreateBuilder) SetProjectName(v string) *ProjectSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "projectName", Val: v})
	return b
}
func (b *ProjectSkillCreateBuilder) SetSkillId(v string) *ProjectSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *ProjectSkillCreateBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSkill]) *ProjectSkillCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectSkillDelegate) Create() *ProjectSkillCreateBuilder {
	return &ProjectSkillCreateBuilder{
		CreateBuilder: &CreateBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedProjectSkillProjectName uint64 = 1 << 0
	providedProjectSkillSkillId     uint64 = 1 << 1
)

func assignmentsToProjectSkillCreate(assignments []FieldAssignment) (ProjectSkillCreate, error) {
	var input ProjectSkillCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "projectName":
			provided |= providedProjectSkillProjectName
			if v, ok := a.Val.(string); ok {
				input.ProjectName = v
				errs.ValidateString("projectName", v, true, 0, false, false)
			} else {
				errs.Add("projectName", a.Val, "type", "field projectName must be of type string")
			}
		case "skillId":
			provided |= providedProjectSkillSkillId
			if v, ok := a.Val.(string); ok {
				input.SkillId = v
				errs.ValidateString("skillId", v, true, 0, false, false)
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		}
	}
	if provided&providedProjectSkillProjectName == 0 {
		errs.Add("projectName", "", "required", "field ProjectName is required")
	}
	if provided&providedProjectSkillSkillId == 0 {
		errs.Add("skillId", "", "required", "field SkillId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *ProjectSkillCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "projectName")
	vals = append(vals, s.ProjectName)
	cols = append(cols, "skillId")
	vals = append(vals, s.SkillId)
	return
}

func partitionProjectSkillInputs(dialect Dialect, inputs []*ProjectSkillCreate) [][]*ProjectSkillCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*ProjectSkillCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*ProjectSkillCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*ProjectSkillCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*ProjectSkillCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*ProjectSkillCreate{inputs}
}

func (d *ProjectSkillDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *ProjectSkillSelect, omits *ProjectSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*ProjectSkill, error) {
	input, err := assignmentsToProjectSkillCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectProjectSkillCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSkillCreateArgs) (*ProjectSkill, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectProjectSkillCols(a.Select, omits)
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
			curr = func(c context.Context, a *ProjectSkillCreateArgs) (*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type ProjectSkillCreateManyBuilder struct {
	*CreateManyBuilder[ProjectSkill]
}

func (b *ProjectSkillCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSkillConflictBuilder[ProjectSkillCreateManyBuilder] {
	return &ProjectSkillConflictBuilder[ProjectSkillCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type ProjectSkillCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]
}

func (b *ProjectSkillCreateManyAndReturnBuilder) Select(s ProjectSkillSelect) *ProjectSkillCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSkillCreateManyAndReturnBuilder) Omit(o ProjectSkillOmit) *ProjectSkillCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSkillCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSkillConflictBuilder[ProjectSkillCreateManyAndReturnBuilder] {
	return &ProjectSkillConflictBuilder[ProjectSkillCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToProjectSkillRecordInputs(builders []*ProjectSkillCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *ProjectSkillDelegate) CreateMany(builders ...*ProjectSkillCreateBuilder) *ProjectSkillCreateManyBuilder {
	return &ProjectSkillCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[ProjectSkill]{
			records:  createBuildersToProjectSkillRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *ProjectSkillDelegate) CreateManyAndReturn(builders ...*ProjectSkillCreateBuilder) *ProjectSkillCreateManyAndReturnBuilder {
	return &ProjectSkillCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
			records:  createBuildersToProjectSkillRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToProjectSkillCreateInputs(records []RecordInput) ([]*ProjectSkillCreate, error) {
	structs := make([]ProjectSkillCreate, len(records))
	inputs := make([]*ProjectSkillCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToProjectSkillCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *ProjectSkillDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToProjectSkillCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &ProjectSkillCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSkillCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSkillCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProjectSkillSelect, omits *ProjectSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*ProjectSkill, error) {
	inputs, err := recordsToProjectSkillCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSkillCreateManyAndReturnArgs) ([]*ProjectSkill, error) {
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
			curr = func(c context.Context, a *ProjectSkillCreateManyAndReturnArgs) ([]*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *ProjectSkillSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*ProjectSkill, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *ProjectSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.ProjectSkill.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.ProjectSkill.loadRelations(ctx, []*ProjectSkill{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "ProjectSkill", cols, returningCols, projectSkillPKCols, conflictTarget, conflictAction, len(vals))
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

		var res ProjectSkill
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, projectSkillPKCols)
}

func (d *ProjectSkillDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*ProjectSkill, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("ProjectSkill") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "ProjectSkill")
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

	var res ProjectSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *ProjectSkillDelegate) buildBulkInsertSQL(q *Queries, batch []*ProjectSkillCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range projectSkillDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "ProjectSkill")
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
			case "projectName":
				vals = append(vals, input.ProjectName)
			case "skillId":
				vals = append(vals, input.SkillId)
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

func applyProjectSkillConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanProjectSkillRows(rows *sql.Rows, returningCols []string) ([]*ProjectSkill, error) {
	var records []*ProjectSkill
	for rows.Next() {
		var res ProjectSkill
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

func (d *ProjectSkillDelegate) runCreateMany(ctx context.Context, inputs []*ProjectSkillCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionProjectSkillInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSkillConflictClause(d.client.dialect, queryStr, vals, cols, projectSkillPKCols, conflictTarget, conflictAction)

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

func (d *ProjectSkillDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*ProjectSkillCreate,
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*ProjectSkill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionProjectSkillInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*ProjectSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.ProjectSkill.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.ProjectSkill.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectProjectSkillCols(selects, omits, projectSkillPKCols...)
	recordsOut := make([]*ProjectSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSkillConflictClause(d.client.dialect, queryStr, vals, cols, projectSkillPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanProjectSkillRows(rows, returningCols)
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

func (d *ProjectSkillDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*ProjectSkillCreate,
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*ProjectSkill, error) {
	batches := partitionProjectSkillInputs(d.client.dialect, inputs)
	returningCols := selectProjectSkillCols(selects, omits, projectSkillPKCols...)
	recordsOut := make([]*ProjectSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSkillConflictClause(d.client.dialect, queryStr, vals, cols, projectSkillPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("ProjectSkill") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "ProjectSkill")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectSkillPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectSkillPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanProjectSkillRows(rows, returningCols)
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

type ProjectSkillConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *ProjectSkillConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectSkillConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectSkillConflictBuilder[B]) Update(fn func(u *ProjectSkillUpsert)) *B {
	cb.setAction(*ProjectSkillConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// ProjectSkillConflictUpdate creates a custom ConflictAction for ProjectSkill upsert conflicts.
func ProjectSkillConflictUpdate(fn func(u *ProjectSkillUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newProjectSkillUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type ProjectSkillUpsert struct {
	ProjectName fieldUpsert[string]
	SkillId     fieldUpsert[string]
}

func newProjectSkillUpsert(up *ConflictUpdate) *ProjectSkillUpsert {
	return &ProjectSkillUpsert{
		ProjectName: fieldUpsert[string]{column: "projectName", update: up},
		SkillId:     fieldUpsert[string]{column: "skillId", update: up},
	}
}

type ProjectSkillUpdateBuilder struct {
	*UpdateBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]
}

type ProjectSkillUpdateManyBuilder struct {
	*UpdateManyBuilder[ProjectSkill]
}

type ProjectSkillUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]
}

func (b *ProjectSkillUpdateBuilder) SetProjectName(v string) *ProjectSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "projectName", Val: v})
	return b
}

func (b *ProjectSkillUpdateManyBuilder) SetProjectName(v string) *ProjectSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "projectName", Val: v})
	return b
}

func (b *ProjectSkillUpdateManyAndReturnBuilder) SetProjectName(v string) *ProjectSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "projectName", Val: v})
	return b
}
func (b *ProjectSkillUpdateBuilder) SetSkillId(v string) *ProjectSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *ProjectSkillUpdateManyBuilder) SetSkillId(v string) *ProjectSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *ProjectSkillUpdateManyAndReturnBuilder) SetSkillId(v string) *ProjectSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *ProjectSkillUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSkill]) *ProjectSkillUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectSkillUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSkill]) *ProjectSkillUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectSkillUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSkill]) *ProjectSkillUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectSkillDelegate) Update(where UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateBuilder {
	return &ProjectSkillUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *ProjectSkillDelegate) UpdateMany(preds ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateManyBuilder {
	return &ProjectSkillUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[ProjectSkill]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *ProjectSkillDelegate) UpdateManyAndReturn(preds ...PredicateOf[ProjectSkill]) *ProjectSkillUpdateManyAndReturnBuilder {
	return &ProjectSkillUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *ProjectSkillDelegate) buildUpdateSQL(preds []PredicateOf[ProjectSkill], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSkill")
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

func (d *ProjectSkillDelegate) executeUpdate(ctx context.Context, where UniquePredicate[ProjectSkill], additional []PredicateOf[ProjectSkill], assignments []FieldAssignment, selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
	allWhere := make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToProjectSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSkillUpdateArgs) (*ProjectSkill, error) {
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
			curr = func(c context.Context, a *ProjectSkillUpdateArgs) (*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runUpdate(ctx context.Context, preds []PredicateOf[ProjectSkill], cols []string, vals []any, selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
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
		var res *ProjectSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.ProjectSkill.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.ProjectSkill.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectSkillCols(selects, omits, projectSkillPKCols...)
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
		return nil, &NotFoundError{Model: "ProjectSkill"}
	}

	var res ProjectSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*ProjectSkill{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *ProjectSkillDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[ProjectSkill], cols []string, vals []any) (int64, error) {
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

func (d *ProjectSkillDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[ProjectSkill], cols []string, vals []any, selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "ProjectSkill"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *ProjectSkillDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[ProjectSkill], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToProjectSkillUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &ProjectSkillUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *ProjectSkillUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSkillUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *ProjectSkillDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[ProjectSkill], assignments []FieldAssignment, selects *ProjectSkillSelect, omits *ProjectSkillOmit) ([]*ProjectSkill, error) {
	input, err := assignmentsToProjectSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSkillUpdateManyAndReturnArgs) ([]*ProjectSkill, error) {
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
			curr = func(c context.Context, a *ProjectSkillUpdateManyAndReturnArgs) ([]*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[ProjectSkill], cols []string, vals []any, selects *ProjectSkillSelect, omits *ProjectSkillOmit) ([]*ProjectSkill, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[ProjectSkill]{Where: preds}, selects, omits)
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
		var res []*ProjectSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.ProjectSkill.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.ProjectSkill.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectSkillCols(selects, omits, projectSkillPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanProjectSkillRows(rows, returningCols)
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

func (d *ProjectSkillDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[ProjectSkill], cols []string, vals []any, selects *ProjectSkillSelect, omits *ProjectSkillOmit) ([]*ProjectSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*ProjectSkill{}, nil
	}
	return d.runFindMany(ctx, QueryParams[ProjectSkill]{Where: preds}, selects, omits)
}
func (d *ProjectSkillDelegate) FindUnique(where UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *FindUniqueBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit] {
	return &FindUniqueBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *ProjectSkillDelegate) FindFirst(preds ...PredicateOf[ProjectSkill]) *FindFirstBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit] {
	return &FindFirstBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *ProjectSkillDelegate) FindMany(preds ...PredicateOf[ProjectSkill]) *FindManyBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit] {
	return &FindManyBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *ProjectSkillDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[ProjectSkill], additional []PredicateOf[ProjectSkill], selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
	allWhere := make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSkillFindUniqueArgs) (*ProjectSkill, error) {
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
			curr = func(c context.Context, a *ProjectSkillFindUniqueArgs) (*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[ProjectSkill],
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
) (*ProjectSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectSkillFindFirstArgs) (*ProjectSkill, error) {
		return d.runFindFirst(c, QueryParams[ProjectSkill]{
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
			curr = func(c context.Context, a *ProjectSkillFindFirstArgs) (*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[ProjectSkill],
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
) ([]*ProjectSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectSkillFindManyArgs) ([]*ProjectSkill, error) {
		return d.runFindMany(c, QueryParams[ProjectSkill]{
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
			curr = func(c context.Context, a *ProjectSkillFindManyArgs) ([]*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runFindUnique(ctx context.Context, where []PredicateOf[ProjectSkill], selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
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
	returningCols := selectProjectSkillCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*ProjectSkill{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *ProjectSkillDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[ProjectSkill],
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
) (*ProjectSkill, error) {
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

func (d *ProjectSkillDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[ProjectSkill],
	selects *ProjectSkillSelect,
	omits *ProjectSkillOmit,
) ([]*ProjectSkill, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, projectSkillPKCols, projectSkillUniqueCols, "ProjectSkill", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, projectSkillPKCols, projectSkillUniqueCols, isCursorQuery, params.Take)
	returningCols := selectProjectSkillCols(selects, omits)

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

func (d *ProjectSkillDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*ProjectSkill, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "ProjectSkill", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res ProjectSkill
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *ProjectSkillDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*ProjectSkill, error) {
	query := buildSelectSQL(d.client, "ProjectSkill", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*ProjectSkill, 0)
	for rows.Next() {
		var res ProjectSkill
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
func (d *ProjectSkillDelegate) DeleteMany(preds ...PredicateOf[ProjectSkill]) *DeleteManyBuilder[ProjectSkill] {
	return &DeleteManyBuilder[ProjectSkill]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *ProjectSkillDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[ProjectSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &ProjectSkillDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *ProjectSkillDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSkillDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[ProjectSkill]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSkill")
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

func (d *ProjectSkillDelegate) Delete(where UniquePredicate[ProjectSkill], additional ...PredicateOf[ProjectSkill]) *DeleteBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit] {
	return &DeleteBuilder[ProjectSkill, ProjectSkillSelect, ProjectSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *ProjectSkillDelegate) executeDelete(ctx context.Context, where UniquePredicate[ProjectSkill], additional []PredicateOf[ProjectSkill], selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
	allWhere := make([]PredicateOf[ProjectSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSkillSelect()
	}

	args := &ProjectSkillDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSkillDeleteArgs) (*ProjectSkill, error) {
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
			curr = func(c context.Context, a *ProjectSkillDeleteArgs) (*ProjectSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runDelete(ctx context.Context, where []PredicateOf[ProjectSkill], selects *ProjectSkillSelect, omits *ProjectSkillOmit) (*ProjectSkill, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectProjectSkillCols(selects, omits, projectSkillPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *ProjectSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.ProjectSkill.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "ProjectSkill"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "ProjectSkill")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[ProjectSkill]
			pkPreds = append(pkPreds, Predicate[ProjectSkill]{
				Data: PredicateData{
					Column:   "projectName",
					Operator: "=",
					Value:    res.ProjectName,
				},
			})
			pkPreds = append(pkPreds, Predicate[ProjectSkill]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSkill")

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
		return nil, &NotFoundError{Model: "ProjectSkill"}
	}

	var row ProjectSkill
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *ProjectSkillDelegate) Count(preds ...PredicateOf[ProjectSkill]) *CountBuilder[ProjectSkill] {
	return &CountBuilder[ProjectSkill]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *ProjectSkillDelegate) executeCount(ctx context.Context, params QueryParams[ProjectSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &ProjectSkillCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *ProjectSkillCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[ProjectSkill]{
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
			curr = func(c context.Context, a *ProjectSkillCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSkillDelegate) runCount(ctx context.Context, params QueryParams[ProjectSkill]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "ProjectSkill")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "ProjectSkill")
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
func (d *ProjectSkillDelegate) loadRelations(ctx context.Context, records []*ProjectSkill, selects *ProjectSkillSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Project != nil {
		relationSelects, relationOmits, relationParams := selects.Project.GetRelationParams()
		returningCols := selectProjectCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: ProjectSkill.projectName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *ProjectSkill) string { return p.ProjectName }),
			"Project",
			"name",
			returningCols,
			scanInto(returningCols, (*Project).ScanFields),
			directKey(func(c *Project) string { return c.Name }),
			setOne(func(p *ProjectSkill, c *Project) { p.Project = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading project: %w", err)
		}
		if err := d.client.Project.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skill != nil {
		relationSelects, relationOmits, relationParams := selects.Skill.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: ProjectSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *ProjectSkill) string { return p.SkillId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *ProjectSkill, c *Skill) { p.Skill = c }),
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
