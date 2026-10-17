package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ProjectSpecialization represents the database model
type ProjectSpecialization struct {
	GeneralizesName string   `db:"generalizesName" json:"generalizesName"`
	SpecializesName string   `db:"specializesName" json:"specializesName"`
	Generalizes     *Project `json:"generalizes,omitempty"`
	Specializes     *Project `json:"specializes,omitempty"`
}

// ProjectSpecializationCreate contains model input fields for ProjectSpecialization creation operations.
//
// Fields for ProjectSpecialization:
//
//	generalizesName string required
//	specializesName string required
type ProjectSpecializationCreate struct {
	GeneralizesName string `json:"generalizesName"`
	SpecializesName string `json:"specializesName"`
}

// colMask returns a bit mask of columns that are set
func (s *ProjectSpecializationCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// ProjectSpecializationUpdate contains model input fields for ProjectSpecialization update operations.
type ProjectSpecializationUpdate struct {
	GeneralizesName *string `json:"generalizesName"`
	SpecializesName *string `json:"specializesName"`
}

func (u *ProjectSpecializationUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.GeneralizesName != nil {
		cols = append(cols, "generalizesName")
		vals = append(vals, u.GeneralizesName)
	}
	if u.SpecializesName != nil {
		cols = append(cols, "specializesName")
		vals = append(vals, u.SpecializesName)
	}
	return cols, vals
}

func assignmentsToProjectSpecializationUpdate(assignments []FieldAssignment) (ProjectSpecializationUpdate, error) {
	var input ProjectSpecializationUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "generalizesName":
			if v, ok := a.Val.(string); ok {
				input.GeneralizesName = &v
				errs.ValidateString("generalizesName", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.GeneralizesName = v
			} else {
				errs.Add("generalizesName", a.Val, "type", "field generalizesName must be of type string")
			}
		case "specializesName":
			if v, ok := a.Val.(string); ok {
				input.SpecializesName = &v
				errs.ValidateString("specializesName", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.SpecializesName = v
			} else {
				errs.Add("specializesName", a.Val, "type", "field specializesName must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// ProjectSpecializationSelect specifies which scalar and relation fields to select for ProjectSpecialization.
//
// Selectable fields:
//
//	-- Scalars --
//	generalizesName (bool)
//	specializesName (bool)
//	-- Relations --
//	generalizes     (Project)
//	specializes     (Project)
type ProjectSpecializationSelect struct {
	GeneralizesName bool           `json:"generalizesName"`
	SpecializesName bool           `json:"specializesName"`
	Generalizes     *ProjectSelect `json:"generalizes,omitempty"`
	Specializes     *ProjectSelect `json:"specializes,omitempty"`
}

var fullProjectSpecializationSelectVal = &ProjectSpecializationSelect{
	GeneralizesName: true,
	SpecializesName: true,
}

func fullProjectSpecializationSelect() *ProjectSpecializationSelect {
	return fullProjectSpecializationSelectVal
}

func (s *ProjectSpecializationSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.GeneralizesName || s.SpecializesName
}

func (s *ProjectSpecializationSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type ProjectSpecializationOmit struct {
	GeneralizesName bool `json:"generalizesName"`
	SpecializesName bool `json:"specializesName"`
}

type ProjectSpecializationSelectQuery interface {
	GetRelationParams() (*ProjectSpecializationSelect, *ProjectSpecializationOmit, QueryParams[ProjectSpecialization])
}

func (s *ProjectSpecializationSelect) GetRelationParams() (*ProjectSpecializationSelect, *ProjectSpecializationOmit, QueryParams[ProjectSpecialization]) {
	return s, nil, QueryParams[ProjectSpecialization]{}
}

type ProjectSpecializationQueryBuilder struct {
	selects *ProjectSpecializationSelect
	omits   *ProjectSpecializationOmit
	where   []PredicateOf[ProjectSpecialization]
	take    *int
	skip    *int
	orderBy []OrderBy[ProjectSpecialization]
	cursor  UniquePredicate[ProjectSpecialization]
}

func (b *ProjectSpecializationQueryBuilder) Where(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *ProjectSpecializationQueryBuilder) Take(limit int) *ProjectSpecializationQueryBuilder {
	b.take = &limit
	return b
}

func (b *ProjectSpecializationQueryBuilder) Skip(offset int) *ProjectSpecializationQueryBuilder {
	b.skip = &offset
	return b
}

func (b *ProjectSpecializationQueryBuilder) OrderBy(orders ...OrderBy[ProjectSpecialization]) *ProjectSpecializationQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *ProjectSpecializationQueryBuilder) Cursor(where UniquePredicate[ProjectSpecialization]) *ProjectSpecializationQueryBuilder {
	b.cursor = where
	return b
}

func (b *ProjectSpecializationQueryBuilder) Select(s ProjectSpecializationSelect) *ProjectSpecializationQueryBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSpecializationQueryBuilder) Omit(o ProjectSpecializationOmit) *ProjectSpecializationQueryBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSpecializationQueryBuilder) GetRelationParams() (*ProjectSpecializationSelect, *ProjectSpecializationOmit, QueryParams[ProjectSpecialization]) {
	if b == nil {
		return nil, nil, QueryParams[ProjectSpecialization]{}
	}
	return b.selects, b.omits, QueryParams[ProjectSpecialization]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// ProjectSpecializationCreateArgs is the input argument passed to ProjectSpecialization Create extension hooks.
//
// Fields for ProjectSpecialization:
//
//	generalizesName string required
//	specializesName string required
//
// Relations for ProjectSpecialization:
//
//	generalizes     (Project)
//	specializes     (Project)
type ProjectSpecializationCreateArgs struct {
	// Data contains the model fields to insert.
	Data *ProjectSpecializationCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *ProjectSpecializationSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// ProjectSpecializationCreateManyArgs is the input argument passed to ProjectSpecialization CreateMany extension hooks.
//
// Fields for ProjectSpecialization:
//
//	generalizesName string required
//	specializesName string required
type ProjectSpecializationCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectSpecializationCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectSpecializationCreateManyArgs) AppendData(builders ...*ProjectSpecializationCreateBuilder) *ProjectSpecializationCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectSpecializationCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectSpecializationCreateManyAndReturnArgs is the input argument passed to ProjectSpecialization CreateManyAndReturn extension hooks.
//
// Fields for ProjectSpecialization:
//
//	generalizesName string required
//	specializesName string required
//
// Relations for ProjectSpecialization:
//
//	generalizes     (Project)
//	specializes     (Project)
type ProjectSpecializationCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectSpecializationCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *ProjectSpecializationSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectSpecializationCreateManyAndReturnArgs) AppendData(builders ...*ProjectSpecializationCreateBuilder) *ProjectSpecializationCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectSpecializationCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectSpecializationFindUniqueArgs is the input argument passed to ProjectSpecialization FindUnique extension hooks.
type ProjectSpecializationFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSpecialization]
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationFindUniqueArgs) SetWhere(unique UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationFindUniqueArgs {
	a.Where = make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSpecializationFindFirstArgs is the input argument passed to ProjectSpecialization FindFirst extension hooks.
type ProjectSpecializationFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[ProjectSpecialization]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[ProjectSpecialization]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationFindFirstArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationFindFirstArgs {
	a.Where = preds
	return a
}

func (a *ProjectSpecializationFindFirstArgs) SetOrderBy(orders ...OrderBy[ProjectSpecialization]) *ProjectSpecializationFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectSpecializationFindFirstArgs) SetCursor(cursor UniquePredicate[ProjectSpecialization]) *ProjectSpecializationFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectSpecializationFindFirstArgs) SetSkip(n int) *ProjectSpecializationFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSpecializationFindFirstArgs) SetTake(n int) *ProjectSpecializationFindFirstArgs {
	a.Take = &n
	return a
}

// ProjectSpecializationFindManyArgs is the input argument passed to ProjectSpecialization FindMany extension hooks.
type ProjectSpecializationFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[ProjectSpecialization]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[ProjectSpecialization]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationFindManyArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationFindManyArgs {
	a.Where = preds
	return a
}

func (a *ProjectSpecializationFindManyArgs) SetOrderBy(orders ...OrderBy[ProjectSpecialization]) *ProjectSpecializationFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectSpecializationFindManyArgs) SetCursor(cursor UniquePredicate[ProjectSpecialization]) *ProjectSpecializationFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectSpecializationFindManyArgs) SetSkip(n int) *ProjectSpecializationFindManyArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSpecializationFindManyArgs) SetTake(n int) *ProjectSpecializationFindManyArgs {
	a.Take = &n
	return a
}

// ProjectSpecializationCountArgs is the input argument passed to ProjectSpecialization Count extension hooks.
type ProjectSpecializationCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *ProjectSpecializationCountArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationCountArgs {
	a.Where = preds
	return a
}

func (a *ProjectSpecializationCountArgs) SetSkip(n int) *ProjectSpecializationCountArgs {
	a.Skip = &n
	return a
}

func (a *ProjectSpecializationCountArgs) SetTake(n int) *ProjectSpecializationCountArgs {
	a.Take = &n
	return a
}

// ProjectSpecializationDeleteArgs is the input argument passed to ProjectSpecialization Delete extension hooks.
type ProjectSpecializationDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSpecialization]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationDeleteArgs) SetWhere(unique UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationDeleteArgs {
	a.Where = make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSpecializationDeleteManyArgs is the input argument passed to ProjectSpecialization DeleteMany extension hooks.
type ProjectSpecializationDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
}

func (a *ProjectSpecializationDeleteManyArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationDeleteManyArgs {
	a.Where = preds
	return a
}

// ProjectSpecializationUpdateArgs is the input argument passed to ProjectSpecialization Update extension hooks.
type ProjectSpecializationUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[ProjectSpecialization]
	// Data contains the model fields to update.
	Data *ProjectSpecializationUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationUpdateArgs) SetWhere(unique UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateArgs {
	a.Where = make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectSpecializationUpdateManyArgs is the input argument passed to ProjectSpecialization UpdateMany extension hooks.
type ProjectSpecializationUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
	// Data contains the model fields to update.
	Data *ProjectSpecializationUpdate
}

func (a *ProjectSpecializationUpdateManyArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyArgs {
	a.Where = preds
	return a
}

// ProjectSpecializationUpdateManyAndReturnArgs is the input argument passed to ProjectSpecialization UpdateManyAndReturn extension hooks.
type ProjectSpecializationUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[ProjectSpecialization]
	// Data contains the model fields to update.
	Data *ProjectSpecializationUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSpecializationSelect
}

func (a *ProjectSpecializationUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type ProjectSpecializationCreateQuery = func(ctx context.Context, args *ProjectSpecializationCreateArgs) (*ProjectSpecialization, error)
type ProjectSpecializationCreateManyQuery = func(ctx context.Context, args *ProjectSpecializationCreateManyArgs) (int64, error)
type ProjectSpecializationCreateManyAndReturnQuery = func(ctx context.Context, args *ProjectSpecializationCreateManyAndReturnArgs) ([]*ProjectSpecialization, error)
type ProjectSpecializationFindUniqueQuery = func(ctx context.Context, args *ProjectSpecializationFindUniqueArgs) (*ProjectSpecialization, error)
type ProjectSpecializationFindFirstQuery = func(ctx context.Context, args *ProjectSpecializationFindFirstArgs) (*ProjectSpecialization, error)
type ProjectSpecializationFindManyQuery = func(ctx context.Context, args *ProjectSpecializationFindManyArgs) ([]*ProjectSpecialization, error)
type ProjectSpecializationDeleteQuery = func(ctx context.Context, args *ProjectSpecializationDeleteArgs) (*ProjectSpecialization, error)
type ProjectSpecializationDeleteManyQuery = func(ctx context.Context, args *ProjectSpecializationDeleteManyArgs) (int64, error)
type ProjectSpecializationCountQuery = func(ctx context.Context, args *ProjectSpecializationCountArgs) (int64, error)
type ProjectSpecializationUpdateQuery = func(ctx context.Context, args *ProjectSpecializationUpdateArgs) (*ProjectSpecialization, error)
type ProjectSpecializationUpdateManyQuery = func(ctx context.Context, args *ProjectSpecializationUpdateManyArgs) (int64, error)
type ProjectSpecializationUpdateManyAndReturnQuery = func(ctx context.Context, args *ProjectSpecializationUpdateManyAndReturnArgs) ([]*ProjectSpecialization, error)

type ProjectSpecializationExtension struct {
	Create              func(ctx context.Context, args *ProjectSpecializationCreateArgs, next ProjectSpecializationCreateQuery) (*ProjectSpecialization, error)
	CreateMany          func(ctx context.Context, args *ProjectSpecializationCreateManyArgs, next ProjectSpecializationCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *ProjectSpecializationCreateManyAndReturnArgs, next ProjectSpecializationCreateManyAndReturnQuery) ([]*ProjectSpecialization, error)
	FindUnique          func(ctx context.Context, args *ProjectSpecializationFindUniqueArgs, next ProjectSpecializationFindUniqueQuery) (*ProjectSpecialization, error)
	FindFirst           func(ctx context.Context, args *ProjectSpecializationFindFirstArgs, next ProjectSpecializationFindFirstQuery) (*ProjectSpecialization, error)
	FindMany            func(ctx context.Context, args *ProjectSpecializationFindManyArgs, next ProjectSpecializationFindManyQuery) ([]*ProjectSpecialization, error)
	Delete              func(ctx context.Context, args *ProjectSpecializationDeleteArgs, next ProjectSpecializationDeleteQuery) (*ProjectSpecialization, error)
	DeleteMany          func(ctx context.Context, args *ProjectSpecializationDeleteManyArgs, next ProjectSpecializationDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *ProjectSpecializationCountArgs, next ProjectSpecializationCountQuery) (int64, error)
	Update              func(ctx context.Context, args *ProjectSpecializationUpdateArgs, next ProjectSpecializationUpdateQuery) (*ProjectSpecialization, error)
	UpdateMany          func(ctx context.Context, args *ProjectSpecializationUpdateManyArgs, next ProjectSpecializationUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *ProjectSpecializationUpdateManyAndReturnArgs, next ProjectSpecializationUpdateManyAndReturnQuery) ([]*ProjectSpecialization, error)
}

type ProjectSpecializationDelegate struct {
	client     *Queries
	extensions []ProjectSpecializationExtension
}

func (d *ProjectSpecializationDelegate) Use(exts ...ProjectSpecializationExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *ProjectSpecialization) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "generalizesName":
			targets[i] = &m.GeneralizesName
		case "specializesName":
			targets[i] = &m.SpecializesName
		}
	}
	return targets
}

var projectSpecializationDefaultCols = []string{
	"generalizesName",
	"specializesName",
}

var projectSpecializationPKCols = []string{
	"generalizesName",
	"specializesName",
}

var projectSpecializationUniqueCols = []string{}

func selectProjectSpecializationCols(selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return projectSpecializationDefaultCols
	}

	anySelected := selects != nil && (selects.GeneralizesName || selects.SpecializesName || selects.Generalizes != nil || selects.Specializes != nil)

	specs := []colSpec{
		{"generalizesName", selects != nil && selects.GeneralizesName, omits != nil && omits.GeneralizesName, selects != nil && selects.Generalizes != nil},
		{"specializesName", selects != nil && selects.SpecializesName, omits != nil && omits.SpecializesName, selects != nil && selects.Specializes != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *ProjectSpecializationSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Generalizes != nil || s.Specializes != nil
}

type ProjectSpecializationUpsertBuilder struct {
	*CreateBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]
}

func (b *ProjectSpecializationUpsertBuilder) Select(s ProjectSpecializationSelect) *ProjectSpecializationUpsertBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSpecializationUpsertBuilder) Omit(o ProjectSpecializationOmit) *ProjectSpecializationUpsertBuilder {
	b.omits = &o
	return b
}

type ProjectSpecializationCreateBuilder struct {
	*CreateBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]
}

func (b *ProjectSpecializationCreateBuilder) Select(s ProjectSpecializationSelect) *ProjectSpecializationCreateBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSpecializationCreateBuilder) Omit(o ProjectSpecializationOmit) *ProjectSpecializationCreateBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSpecializationCreateBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSpecializationConflictBuilder[ProjectSpecializationUpsertBuilder] {
	upsertBuilder := &ProjectSpecializationUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &ProjectSpecializationConflictBuilder[ProjectSpecializationUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *ProjectSpecializationCreateBuilder) SetGeneralizesName(v string) *ProjectSpecializationCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "generalizesName", Val: v})
	return b
}
func (b *ProjectSpecializationCreateBuilder) SetSpecializesName(v string) *ProjectSpecializationCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "specializesName", Val: v})
	return b
}

func (b *ProjectSpecializationCreateBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSpecialization]) *ProjectSpecializationCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectSpecializationDelegate) Create() *ProjectSpecializationCreateBuilder {
	return &ProjectSpecializationCreateBuilder{
		CreateBuilder: &CreateBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedProjectSpecializationGeneralizesName uint64 = 1 << 0
	providedProjectSpecializationSpecializesName uint64 = 1 << 1
)

func assignmentsToProjectSpecializationCreate(assignments []FieldAssignment) (ProjectSpecializationCreate, error) {
	var input ProjectSpecializationCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "generalizesName":
			provided |= providedProjectSpecializationGeneralizesName
			if v, ok := a.Val.(string); ok {
				input.GeneralizesName = v
				errs.ValidateString("generalizesName", v, true, 0, false, false)
			} else {
				errs.Add("generalizesName", a.Val, "type", "field generalizesName must be of type string")
			}
		case "specializesName":
			provided |= providedProjectSpecializationSpecializesName
			if v, ok := a.Val.(string); ok {
				input.SpecializesName = v
				errs.ValidateString("specializesName", v, true, 0, false, false)
			} else {
				errs.Add("specializesName", a.Val, "type", "field specializesName must be of type string")
			}
		}
	}
	if provided&providedProjectSpecializationGeneralizesName == 0 {
		errs.Add("generalizesName", "", "required", "field GeneralizesName is required")
	}
	if provided&providedProjectSpecializationSpecializesName == 0 {
		errs.Add("specializesName", "", "required", "field SpecializesName is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *ProjectSpecializationCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "generalizesName")
	vals = append(vals, s.GeneralizesName)
	cols = append(cols, "specializesName")
	vals = append(vals, s.SpecializesName)
	return
}

func partitionProjectSpecializationInputs(dialect Dialect, inputs []*ProjectSpecializationCreate) [][]*ProjectSpecializationCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*ProjectSpecializationCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*ProjectSpecializationCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*ProjectSpecializationCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*ProjectSpecializationCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*ProjectSpecializationCreate{inputs}
}

func (d *ProjectSpecializationDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*ProjectSpecialization, error) {
	input, err := assignmentsToProjectSpecializationCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectProjectSpecializationCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSpecializationCreateArgs) (*ProjectSpecialization, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectProjectSpecializationCols(a.Select, omits)
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
			curr = func(c context.Context, a *ProjectSpecializationCreateArgs) (*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type ProjectSpecializationCreateManyBuilder struct {
	*CreateManyBuilder[ProjectSpecialization]
}

func (b *ProjectSpecializationCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSpecializationConflictBuilder[ProjectSpecializationCreateManyBuilder] {
	return &ProjectSpecializationConflictBuilder[ProjectSpecializationCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type ProjectSpecializationCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]
}

func (b *ProjectSpecializationCreateManyAndReturnBuilder) Select(s ProjectSpecializationSelect) *ProjectSpecializationCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *ProjectSpecializationCreateManyAndReturnBuilder) Omit(o ProjectSpecializationOmit) *ProjectSpecializationCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *ProjectSpecializationCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *ProjectSpecializationConflictBuilder[ProjectSpecializationCreateManyAndReturnBuilder] {
	return &ProjectSpecializationConflictBuilder[ProjectSpecializationCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToProjectSpecializationRecordInputs(builders []*ProjectSpecializationCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *ProjectSpecializationDelegate) CreateMany(builders ...*ProjectSpecializationCreateBuilder) *ProjectSpecializationCreateManyBuilder {
	return &ProjectSpecializationCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[ProjectSpecialization]{
			records:  createBuildersToProjectSpecializationRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *ProjectSpecializationDelegate) CreateManyAndReturn(builders ...*ProjectSpecializationCreateBuilder) *ProjectSpecializationCreateManyAndReturnBuilder {
	return &ProjectSpecializationCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
			records:  createBuildersToProjectSpecializationRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToProjectSpecializationCreateInputs(records []RecordInput) ([]*ProjectSpecializationCreate, error) {
	structs := make([]ProjectSpecializationCreate, len(records))
	inputs := make([]*ProjectSpecializationCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToProjectSpecializationCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *ProjectSpecializationDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToProjectSpecializationCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &ProjectSpecializationCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSpecializationCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*ProjectSpecialization, error) {
	inputs, err := recordsToProjectSpecializationCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectSpecializationCreateManyAndReturnArgs) ([]*ProjectSpecialization, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationCreateManyAndReturnArgs) ([]*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *ProjectSpecializationSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*ProjectSpecialization, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *ProjectSpecialization
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.ProjectSpecialization.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.ProjectSpecialization.loadRelations(ctx, []*ProjectSpecialization{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "ProjectSpecialization", cols, returningCols, projectSpecializationPKCols, conflictTarget, conflictAction, len(vals))
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

		var res ProjectSpecialization
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, projectSpecializationPKCols)
}

func (d *ProjectSpecializationDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*ProjectSpecialization, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("ProjectSpecialization") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "ProjectSpecialization")
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

	var res ProjectSpecialization
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *ProjectSpecializationDelegate) buildBulkInsertSQL(q *Queries, batch []*ProjectSpecializationCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range projectSpecializationDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "ProjectSpecialization")
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
			case "generalizesName":
				vals = append(vals, input.GeneralizesName)
			case "specializesName":
				vals = append(vals, input.SpecializesName)
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

func applyProjectSpecializationConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanProjectSpecializationRows(rows *sql.Rows, returningCols []string) ([]*ProjectSpecialization, error) {
	var records []*ProjectSpecialization
	for rows.Next() {
		var res ProjectSpecialization
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

func (d *ProjectSpecializationDelegate) runCreateMany(ctx context.Context, inputs []*ProjectSpecializationCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionProjectSpecializationInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSpecializationConflictClause(d.client.dialect, queryStr, vals, cols, projectSpecializationPKCols, conflictTarget, conflictAction)

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

func (d *ProjectSpecializationDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*ProjectSpecializationCreate,
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*ProjectSpecialization, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionProjectSpecializationInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*ProjectSpecialization
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.ProjectSpecialization.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.ProjectSpecialization.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectProjectSpecializationCols(selects, omits, projectSpecializationPKCols...)
	recordsOut := make([]*ProjectSpecialization, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSpecializationConflictClause(d.client.dialect, queryStr, vals, cols, projectSpecializationPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanProjectSpecializationRows(rows, returningCols)
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

func (d *ProjectSpecializationDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*ProjectSpecializationCreate,
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*ProjectSpecialization, error) {
	batches := partitionProjectSpecializationInputs(d.client.dialect, inputs)
	returningCols := selectProjectSpecializationCols(selects, omits, projectSpecializationPKCols...)
	recordsOut := make([]*ProjectSpecialization, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectSpecializationConflictClause(d.client.dialect, queryStr, vals, cols, projectSpecializationPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("ProjectSpecialization") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "ProjectSpecialization")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectSpecializationPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectSpecializationPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanProjectSpecializationRows(rows, returningCols)
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

type ProjectSpecializationConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *ProjectSpecializationConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectSpecializationConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectSpecializationConflictBuilder[B]) Update(fn func(u *ProjectSpecializationUpsert)) *B {
	cb.setAction(*ProjectSpecializationConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// ProjectSpecializationConflictUpdate creates a custom ConflictAction for ProjectSpecialization upsert conflicts.
func ProjectSpecializationConflictUpdate(fn func(u *ProjectSpecializationUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newProjectSpecializationUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type ProjectSpecializationUpsert struct {
	GeneralizesName fieldUpsert[string]
	SpecializesName fieldUpsert[string]
}

func newProjectSpecializationUpsert(up *ConflictUpdate) *ProjectSpecializationUpsert {
	return &ProjectSpecializationUpsert{
		GeneralizesName: fieldUpsert[string]{column: "generalizesName", update: up},
		SpecializesName: fieldUpsert[string]{column: "specializesName", update: up},
	}
}

type ProjectSpecializationUpdateBuilder struct {
	*UpdateBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]
}

type ProjectSpecializationUpdateManyBuilder struct {
	*UpdateManyBuilder[ProjectSpecialization]
}

type ProjectSpecializationUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]
}

func (b *ProjectSpecializationUpdateBuilder) SetGeneralizesName(v string) *ProjectSpecializationUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "generalizesName", Val: v})
	return b
}

func (b *ProjectSpecializationUpdateManyBuilder) SetGeneralizesName(v string) *ProjectSpecializationUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "generalizesName", Val: v})
	return b
}

func (b *ProjectSpecializationUpdateManyAndReturnBuilder) SetGeneralizesName(v string) *ProjectSpecializationUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "generalizesName", Val: v})
	return b
}
func (b *ProjectSpecializationUpdateBuilder) SetSpecializesName(v string) *ProjectSpecializationUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "specializesName", Val: v})
	return b
}

func (b *ProjectSpecializationUpdateManyBuilder) SetSpecializesName(v string) *ProjectSpecializationUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "specializesName", Val: v})
	return b
}

func (b *ProjectSpecializationUpdateManyAndReturnBuilder) SetSpecializesName(v string) *ProjectSpecializationUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "specializesName", Val: v})
	return b
}

func (b *ProjectSpecializationUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSpecialization]) *ProjectSpecializationUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectSpecializationUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectSpecializationUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectSpecializationDelegate) Update(where UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateBuilder {
	return &ProjectSpecializationUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *ProjectSpecializationDelegate) UpdateMany(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyBuilder {
	return &ProjectSpecializationUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[ProjectSpecialization]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *ProjectSpecializationDelegate) UpdateManyAndReturn(preds ...PredicateOf[ProjectSpecialization]) *ProjectSpecializationUpdateManyAndReturnBuilder {
	return &ProjectSpecializationUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *ProjectSpecializationDelegate) buildUpdateSQL(preds []PredicateOf[ProjectSpecialization], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSpecialization")
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

func (d *ProjectSpecializationDelegate) executeUpdate(ctx context.Context, where UniquePredicate[ProjectSpecialization], additional []PredicateOf[ProjectSpecialization], assignments []FieldAssignment, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
	allWhere := make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToProjectSpecializationUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationUpdateArgs) (*ProjectSpecialization, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationUpdateArgs) (*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runUpdate(ctx context.Context, preds []PredicateOf[ProjectSpecialization], cols []string, vals []any, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
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
		var res *ProjectSpecialization
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.ProjectSpecialization.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.ProjectSpecialization.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectSpecializationCols(selects, omits, projectSpecializationPKCols...)
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
		return nil, &NotFoundError{Model: "ProjectSpecialization"}
	}

	var res ProjectSpecialization
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*ProjectSpecialization{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *ProjectSpecializationDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[ProjectSpecialization], cols []string, vals []any) (int64, error) {
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

func (d *ProjectSpecializationDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[ProjectSpecialization], cols []string, vals []any, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "ProjectSpecialization"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *ProjectSpecializationDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[ProjectSpecialization], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToProjectSpecializationUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &ProjectSpecializationUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *ProjectSpecializationUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *ProjectSpecializationDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[ProjectSpecialization], assignments []FieldAssignment, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) ([]*ProjectSpecialization, error) {
	input, err := assignmentsToProjectSpecializationUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationUpdateManyAndReturnArgs) ([]*ProjectSpecialization, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationUpdateManyAndReturnArgs) ([]*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[ProjectSpecialization], cols []string, vals []any, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) ([]*ProjectSpecialization, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[ProjectSpecialization]{Where: preds}, selects, omits)
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
		var res []*ProjectSpecialization
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.ProjectSpecialization.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.ProjectSpecialization.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectSpecializationCols(selects, omits, projectSpecializationPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanProjectSpecializationRows(rows, returningCols)
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

func (d *ProjectSpecializationDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[ProjectSpecialization], cols []string, vals []any, selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) ([]*ProjectSpecialization, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*ProjectSpecialization{}, nil
	}
	return d.runFindMany(ctx, QueryParams[ProjectSpecialization]{Where: preds}, selects, omits)
}
func (d *ProjectSpecializationDelegate) FindUnique(where UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *FindUniqueBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit] {
	return &FindUniqueBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *ProjectSpecializationDelegate) FindFirst(preds ...PredicateOf[ProjectSpecialization]) *FindFirstBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit] {
	return &FindFirstBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *ProjectSpecializationDelegate) FindMany(preds ...PredicateOf[ProjectSpecialization]) *FindManyBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit] {
	return &FindManyBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *ProjectSpecializationDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[ProjectSpecialization], additional []PredicateOf[ProjectSpecialization], selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
	allWhere := make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationFindUniqueArgs) (*ProjectSpecialization, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationFindUniqueArgs) (*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[ProjectSpecialization],
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
) (*ProjectSpecialization, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationFindFirstArgs) (*ProjectSpecialization, error) {
		return d.runFindFirst(c, QueryParams[ProjectSpecialization]{
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
			curr = func(c context.Context, a *ProjectSpecializationFindFirstArgs) (*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[ProjectSpecialization],
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
) ([]*ProjectSpecialization, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationFindManyArgs) ([]*ProjectSpecialization, error) {
		return d.runFindMany(c, QueryParams[ProjectSpecialization]{
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
			curr = func(c context.Context, a *ProjectSpecializationFindManyArgs) ([]*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runFindUnique(ctx context.Context, where []PredicateOf[ProjectSpecialization], selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
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
	returningCols := selectProjectSpecializationCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*ProjectSpecialization{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *ProjectSpecializationDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[ProjectSpecialization],
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
) (*ProjectSpecialization, error) {
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

func (d *ProjectSpecializationDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[ProjectSpecialization],
	selects *ProjectSpecializationSelect,
	omits *ProjectSpecializationOmit,
) ([]*ProjectSpecialization, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, projectSpecializationPKCols, projectSpecializationUniqueCols, "ProjectSpecialization", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, projectSpecializationPKCols, projectSpecializationUniqueCols, isCursorQuery, params.Take)
	returningCols := selectProjectSpecializationCols(selects, omits)

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

func (d *ProjectSpecializationDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*ProjectSpecialization, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "ProjectSpecialization", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res ProjectSpecialization
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *ProjectSpecializationDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*ProjectSpecialization, error) {
	query := buildSelectSQL(d.client, "ProjectSpecialization", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*ProjectSpecialization, 0)
	for rows.Next() {
		var res ProjectSpecialization
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
func (d *ProjectSpecializationDelegate) DeleteMany(preds ...PredicateOf[ProjectSpecialization]) *DeleteManyBuilder[ProjectSpecialization] {
	return &DeleteManyBuilder[ProjectSpecialization]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *ProjectSpecializationDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[ProjectSpecialization]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &ProjectSpecializationDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *ProjectSpecializationDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[ProjectSpecialization]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSpecialization")
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

func (d *ProjectSpecializationDelegate) Delete(where UniquePredicate[ProjectSpecialization], additional ...PredicateOf[ProjectSpecialization]) *DeleteBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit] {
	return &DeleteBuilder[ProjectSpecialization, ProjectSpecializationSelect, ProjectSpecializationOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *ProjectSpecializationDelegate) executeDelete(ctx context.Context, where UniquePredicate[ProjectSpecialization], additional []PredicateOf[ProjectSpecialization], selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
	allWhere := make([]PredicateOf[ProjectSpecialization], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSpecializationSelect()
	}

	args := &ProjectSpecializationDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectSpecializationDeleteArgs) (*ProjectSpecialization, error) {
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
			curr = func(c context.Context, a *ProjectSpecializationDeleteArgs) (*ProjectSpecialization, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runDelete(ctx context.Context, where []PredicateOf[ProjectSpecialization], selects *ProjectSpecializationSelect, omits *ProjectSpecializationOmit) (*ProjectSpecialization, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectProjectSpecializationCols(selects, omits, projectSpecializationPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *ProjectSpecialization
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.ProjectSpecialization.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "ProjectSpecialization"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "ProjectSpecialization")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[ProjectSpecialization]
			pkPreds = append(pkPreds, Predicate[ProjectSpecialization]{
				Data: PredicateData{
					Column:   "generalizesName",
					Operator: "=",
					Value:    res.GeneralizesName,
				},
			})
			pkPreds = append(pkPreds, Predicate[ProjectSpecialization]{
				Data: PredicateData{
					Column:   "specializesName",
					Operator: "=",
					Value:    res.SpecializesName,
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
	d.client.dialect.WriteQuotedIdent(&sb, "ProjectSpecialization")

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
		return nil, &NotFoundError{Model: "ProjectSpecialization"}
	}

	var row ProjectSpecialization
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *ProjectSpecializationDelegate) Count(preds ...PredicateOf[ProjectSpecialization]) *CountBuilder[ProjectSpecialization] {
	return &CountBuilder[ProjectSpecialization]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *ProjectSpecializationDelegate) executeCount(ctx context.Context, params QueryParams[ProjectSpecialization]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &ProjectSpecializationCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *ProjectSpecializationCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[ProjectSpecialization]{
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
			curr = func(c context.Context, a *ProjectSpecializationCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectSpecializationDelegate) runCount(ctx context.Context, params QueryParams[ProjectSpecialization]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "ProjectSpecialization")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "ProjectSpecialization")
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
func (d *ProjectSpecializationDelegate) loadRelations(ctx context.Context, records []*ProjectSpecialization, selects *ProjectSpecializationSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Generalizes != nil {
		relationSelects, relationOmits, relationParams := selects.Generalizes.GetRelationParams()
		returningCols := selectProjectCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: ProjectSpecialization.generalizesName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *ProjectSpecialization) string { return p.GeneralizesName }),
			"Project",
			"name",
			returningCols,
			scanInto(returningCols, (*Project).ScanFields),
			directKey(func(c *Project) string { return c.Name }),
			setOne(func(p *ProjectSpecialization, c *Project) { p.Generalizes = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading generalizes: %w", err)
		}
		if err := d.client.Project.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Specializes != nil {
		relationSelects, relationOmits, relationParams := selects.Specializes.GetRelationParams()
		returningCols := selectProjectCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: ProjectSpecialization.specializesName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *ProjectSpecialization) string { return p.SpecializesName }),
			"Project",
			"name",
			returningCols,
			scanInto(returningCols, (*Project).ScanFields),
			directKey(func(c *Project) string { return c.Name }),
			setOne(func(p *ProjectSpecialization, c *Project) { p.Specializes = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading specializes: %w", err)
		}
		if err := d.client.Project.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
