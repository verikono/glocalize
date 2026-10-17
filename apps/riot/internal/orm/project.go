package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Project represents the database model
type Project struct {
	Name        string                   `db:"name" json:"name"`
	SkillsPath  string                   `db:"skillsPath" json:"skillsPath"`
	AgentsPath  string                   `db:"agentsPath" json:"agentsPath"`
	Skills      []*ProjectSkill          `json:"skills,omitempty"`
	Generalizes []*ProjectSpecialization `json:"generalizes,omitempty"`
	Specializes []*ProjectSpecialization `json:"specializes,omitempty"`
	Apps        []*App                   `json:"apps,omitempty"`
	Agents      []*Agent                 `json:"agents,omitempty"`
}

// ProjectCreate contains model input fields for Project creation operations.
//
// Fields for Project:
//
//	name       string required
//	skillsPath string required
//	agentsPath string required
type ProjectCreate struct {
	Name       string `json:"name"`
	SkillsPath string `json:"skillsPath"`
	AgentsPath string `json:"agentsPath"`
}

// colMask returns a bit mask of columns that are set
func (s *ProjectCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	return mask
}

// ProjectUpdate contains model input fields for Project update operations.
type ProjectUpdate struct {
	Name       *string `json:"name"`
	SkillsPath *string `json:"skillsPath"`
	AgentsPath *string `json:"agentsPath"`
}

func (u *ProjectUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	if u.SkillsPath != nil {
		cols = append(cols, "skillsPath")
		vals = append(vals, u.SkillsPath)
	}
	if u.AgentsPath != nil {
		cols = append(cols, "agentsPath")
		vals = append(vals, u.AgentsPath)
	}
	return cols, vals
}

func assignmentsToProjectUpdate(assignments []FieldAssignment) (ProjectUpdate, error) {
	var input ProjectUpdate
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
		case "skillsPath":
			if v, ok := a.Val.(string); ok {
				input.SkillsPath = &v
				errs.ValidateString("skillsPath", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.SkillsPath = v
			} else {
				errs.Add("skillsPath", a.Val, "type", "field skillsPath must be of type string")
			}
		case "agentsPath":
			if v, ok := a.Val.(string); ok {
				input.AgentsPath = &v
				errs.ValidateString("agentsPath", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.AgentsPath = v
			} else {
				errs.Add("agentsPath", a.Val, "type", "field agentsPath must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// ProjectSelect specifies which scalar and relation fields to select for Project.
//
// Selectable fields:
//
//	-- Scalars --
//	name       (bool)
//	skillsPath (bool)
//	agentsPath (bool)
//	-- Relations --
//	skills     ([]ProjectSkill)
//	generalizes ([]ProjectSpecialization)
//	specializes ([]ProjectSpecialization)
//	apps       ([]App)
//	agents     ([]Agent)
type ProjectSelect struct {
	Name        bool                             `json:"name"`
	SkillsPath  bool                             `json:"skillsPath"`
	AgentsPath  bool                             `json:"agentsPath"`
	Skills      ProjectSkillSelectQuery          `json:"skills,omitempty"`
	Generalizes ProjectSpecializationSelectQuery `json:"generalizes,omitempty"`
	Specializes ProjectSpecializationSelectQuery `json:"specializes,omitempty"`
	Apps        AppSelectQuery                   `json:"apps,omitempty"`
	Agents      AgentSelectQuery                 `json:"agents,omitempty"`
}

var fullProjectSelectVal = &ProjectSelect{
	Name:       true,
	SkillsPath: true,
	AgentsPath: true,
}

func fullProjectSelect() *ProjectSelect {
	return fullProjectSelectVal
}

func (s *ProjectSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Name || s.SkillsPath || s.AgentsPath
}

func (s *ProjectSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type ProjectOmit struct {
	Name       bool `json:"name"`
	SkillsPath bool `json:"skillsPath"`
	AgentsPath bool `json:"agentsPath"`
}

type ProjectSelectQuery interface {
	GetRelationParams() (*ProjectSelect, *ProjectOmit, QueryParams[Project])
}

func (s *ProjectSelect) GetRelationParams() (*ProjectSelect, *ProjectOmit, QueryParams[Project]) {
	return s, nil, QueryParams[Project]{}
}

type ProjectQueryBuilder struct {
	selects *ProjectSelect
	omits   *ProjectOmit
	where   []PredicateOf[Project]
	take    *int
	skip    *int
	orderBy []OrderBy[Project]
	cursor  UniquePredicate[Project]
}

func (b *ProjectQueryBuilder) Where(preds ...PredicateOf[Project]) *ProjectQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *ProjectQueryBuilder) Take(limit int) *ProjectQueryBuilder {
	b.take = &limit
	return b
}

func (b *ProjectQueryBuilder) Skip(offset int) *ProjectQueryBuilder {
	b.skip = &offset
	return b
}

func (b *ProjectQueryBuilder) OrderBy(orders ...OrderBy[Project]) *ProjectQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *ProjectQueryBuilder) Cursor(where UniquePredicate[Project]) *ProjectQueryBuilder {
	b.cursor = where
	return b
}

func (b *ProjectQueryBuilder) Select(s ProjectSelect) *ProjectQueryBuilder {
	b.selects = &s
	return b
}

func (b *ProjectQueryBuilder) Omit(o ProjectOmit) *ProjectQueryBuilder {
	b.omits = &o
	return b
}

func (b *ProjectQueryBuilder) GetRelationParams() (*ProjectSelect, *ProjectOmit, QueryParams[Project]) {
	if b == nil {
		return nil, nil, QueryParams[Project]{}
	}
	return b.selects, b.omits, QueryParams[Project]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// ProjectCreateArgs is the input argument passed to Project Create extension hooks.
//
// Fields for Project:
//
//	name       string required
//	skillsPath string required
//	agentsPath string required
//
// Relations for Project:
//
//	skills     ([]ProjectSkill)
//	generalizes ([]ProjectSpecialization)
//	specializes ([]ProjectSpecialization)
//	apps       ([]App)
//	agents     ([]Agent)
type ProjectCreateArgs struct {
	// Data contains the model fields to insert.
	Data *ProjectCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *ProjectSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// ProjectCreateManyArgs is the input argument passed to Project CreateMany extension hooks.
//
// Fields for Project:
//
//	name       string required
//	skillsPath string required
//	agentsPath string required
type ProjectCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectCreateManyArgs) AppendData(builders ...*ProjectCreateBuilder) *ProjectCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectCreateManyAndReturnArgs is the input argument passed to Project CreateManyAndReturn extension hooks.
//
// Fields for Project:
//
//	name       string required
//	skillsPath string required
//	agentsPath string required
//
// Relations for Project:
//
//	skills     ([]ProjectSkill)
//	generalizes ([]ProjectSpecialization)
//	specializes ([]ProjectSpecialization)
//	apps       ([]App)
//	agents     ([]Agent)
type ProjectCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ProjectCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *ProjectSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ProjectCreateManyAndReturnArgs) AppendData(builders ...*ProjectCreateBuilder) *ProjectCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToProjectCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ProjectFindUniqueArgs is the input argument passed to Project FindUnique extension hooks.
type ProjectFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Project]
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSelect
}

func (a *ProjectFindUniqueArgs) SetWhere(unique UniquePredicate[Project], additional ...PredicateOf[Project]) *ProjectFindUniqueArgs {
	a.Where = make([]PredicateOf[Project], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectFindFirstArgs is the input argument passed to Project FindFirst extension hooks.
type ProjectFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Project]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Project]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSelect
}

func (a *ProjectFindFirstArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectFindFirstArgs {
	a.Where = preds
	return a
}

func (a *ProjectFindFirstArgs) SetOrderBy(orders ...OrderBy[Project]) *ProjectFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectFindFirstArgs) SetCursor(cursor UniquePredicate[Project]) *ProjectFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectFindFirstArgs) SetSkip(n int) *ProjectFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *ProjectFindFirstArgs) SetTake(n int) *ProjectFindFirstArgs {
	a.Take = &n
	return a
}

// ProjectFindManyArgs is the input argument passed to Project FindMany extension hooks.
type ProjectFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Project]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Project]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ProjectSelect
}

func (a *ProjectFindManyArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectFindManyArgs {
	a.Where = preds
	return a
}

func (a *ProjectFindManyArgs) SetOrderBy(orders ...OrderBy[Project]) *ProjectFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *ProjectFindManyArgs) SetCursor(cursor UniquePredicate[Project]) *ProjectFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *ProjectFindManyArgs) SetSkip(n int) *ProjectFindManyArgs {
	a.Skip = &n
	return a
}

func (a *ProjectFindManyArgs) SetTake(n int) *ProjectFindManyArgs {
	a.Take = &n
	return a
}

// ProjectCountArgs is the input argument passed to Project Count extension hooks.
type ProjectCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *ProjectCountArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectCountArgs {
	a.Where = preds
	return a
}

func (a *ProjectCountArgs) SetSkip(n int) *ProjectCountArgs {
	a.Skip = &n
	return a
}

func (a *ProjectCountArgs) SetTake(n int) *ProjectCountArgs {
	a.Take = &n
	return a
}

// ProjectDeleteArgs is the input argument passed to Project Delete extension hooks.
type ProjectDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Project]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *ProjectSelect
}

func (a *ProjectDeleteArgs) SetWhere(unique UniquePredicate[Project], additional ...PredicateOf[Project]) *ProjectDeleteArgs {
	a.Where = make([]PredicateOf[Project], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectDeleteManyArgs is the input argument passed to Project DeleteMany extension hooks.
type ProjectDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
}

func (a *ProjectDeleteManyArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectDeleteManyArgs {
	a.Where = preds
	return a
}

// ProjectUpdateArgs is the input argument passed to Project Update extension hooks.
type ProjectUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Project]
	// Data contains the model fields to update.
	Data *ProjectUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSelect
}

func (a *ProjectUpdateArgs) SetWhere(unique UniquePredicate[Project], additional ...PredicateOf[Project]) *ProjectUpdateArgs {
	a.Where = make([]PredicateOf[Project], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ProjectUpdateManyArgs is the input argument passed to Project UpdateMany extension hooks.
type ProjectUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
	// Data contains the model fields to update.
	Data *ProjectUpdate
}

func (a *ProjectUpdateManyArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectUpdateManyArgs {
	a.Where = preds
	return a
}

// ProjectUpdateManyAndReturnArgs is the input argument passed to Project UpdateManyAndReturn extension hooks.
type ProjectUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Project]
	// Data contains the model fields to update.
	Data *ProjectUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ProjectSelect
}

func (a *ProjectUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Project]) *ProjectUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type ProjectCreateQuery = func(ctx context.Context, args *ProjectCreateArgs) (*Project, error)
type ProjectCreateManyQuery = func(ctx context.Context, args *ProjectCreateManyArgs) (int64, error)
type ProjectCreateManyAndReturnQuery = func(ctx context.Context, args *ProjectCreateManyAndReturnArgs) ([]*Project, error)
type ProjectFindUniqueQuery = func(ctx context.Context, args *ProjectFindUniqueArgs) (*Project, error)
type ProjectFindFirstQuery = func(ctx context.Context, args *ProjectFindFirstArgs) (*Project, error)
type ProjectFindManyQuery = func(ctx context.Context, args *ProjectFindManyArgs) ([]*Project, error)
type ProjectDeleteQuery = func(ctx context.Context, args *ProjectDeleteArgs) (*Project, error)
type ProjectDeleteManyQuery = func(ctx context.Context, args *ProjectDeleteManyArgs) (int64, error)
type ProjectCountQuery = func(ctx context.Context, args *ProjectCountArgs) (int64, error)
type ProjectUpdateQuery = func(ctx context.Context, args *ProjectUpdateArgs) (*Project, error)
type ProjectUpdateManyQuery = func(ctx context.Context, args *ProjectUpdateManyArgs) (int64, error)
type ProjectUpdateManyAndReturnQuery = func(ctx context.Context, args *ProjectUpdateManyAndReturnArgs) ([]*Project, error)

type ProjectExtension struct {
	Create              func(ctx context.Context, args *ProjectCreateArgs, next ProjectCreateQuery) (*Project, error)
	CreateMany          func(ctx context.Context, args *ProjectCreateManyArgs, next ProjectCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *ProjectCreateManyAndReturnArgs, next ProjectCreateManyAndReturnQuery) ([]*Project, error)
	FindUnique          func(ctx context.Context, args *ProjectFindUniqueArgs, next ProjectFindUniqueQuery) (*Project, error)
	FindFirst           func(ctx context.Context, args *ProjectFindFirstArgs, next ProjectFindFirstQuery) (*Project, error)
	FindMany            func(ctx context.Context, args *ProjectFindManyArgs, next ProjectFindManyQuery) ([]*Project, error)
	Delete              func(ctx context.Context, args *ProjectDeleteArgs, next ProjectDeleteQuery) (*Project, error)
	DeleteMany          func(ctx context.Context, args *ProjectDeleteManyArgs, next ProjectDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *ProjectCountArgs, next ProjectCountQuery) (int64, error)
	Update              func(ctx context.Context, args *ProjectUpdateArgs, next ProjectUpdateQuery) (*Project, error)
	UpdateMany          func(ctx context.Context, args *ProjectUpdateManyArgs, next ProjectUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *ProjectUpdateManyAndReturnArgs, next ProjectUpdateManyAndReturnQuery) ([]*Project, error)
}

type ProjectDelegate struct {
	client     *Queries
	extensions []ProjectExtension
}

func (d *ProjectDelegate) Use(exts ...ProjectExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Project) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "name":
			targets[i] = &m.Name
		case "skillsPath":
			targets[i] = &m.SkillsPath
		case "agentsPath":
			targets[i] = &m.AgentsPath
		}
	}
	return targets
}

var projectDefaultCols = []string{
	"name",
	"skillsPath",
	"agentsPath",
}

var projectPKCols = []string{
	"name",
}

var projectUniqueCols = []string{
	"name",
}

func selectProjectCols(selects *ProjectSelect, omits *ProjectOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return projectDefaultCols
	}

	anySelected := selects != nil && (selects.Name || selects.SkillsPath || selects.AgentsPath || selects.Skills != nil || selects.Generalizes != nil || selects.Specializes != nil || selects.Apps != nil || selects.Agents != nil)

	specs := []colSpec{
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, selects != nil && selects.hasAnyRelation()},
		{"skillsPath", selects != nil && selects.SkillsPath, omits != nil && omits.SkillsPath, false},
		{"agentsPath", selects != nil && selects.AgentsPath, omits != nil && omits.AgentsPath, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *ProjectSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Skills != nil || s.Generalizes != nil || s.Specializes != nil || s.Apps != nil || s.Agents != nil
}

type ProjectUpsertBuilder struct {
	*CreateBuilder[Project, ProjectSelect, ProjectOmit]
}

func (b *ProjectUpsertBuilder) Select(s ProjectSelect) *ProjectUpsertBuilder {
	b.selects = &s
	return b
}

func (b *ProjectUpsertBuilder) Omit(o ProjectOmit) *ProjectUpsertBuilder {
	b.omits = &o
	return b
}

type ProjectCreateBuilder struct {
	*CreateBuilder[Project, ProjectSelect, ProjectOmit]
}

func (b *ProjectCreateBuilder) Select(s ProjectSelect) *ProjectCreateBuilder {
	b.selects = &s
	return b
}

func (b *ProjectCreateBuilder) Omit(o ProjectOmit) *ProjectCreateBuilder {
	b.omits = &o
	return b
}

func (b *ProjectCreateBuilder) OnConflict(target UniqueConstraintTarget) *ProjectConflictBuilder[ProjectUpsertBuilder] {
	upsertBuilder := &ProjectUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &ProjectConflictBuilder[ProjectUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *ProjectCreateBuilder) SetName(v string) *ProjectCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *ProjectCreateBuilder) SetSkillsPath(v string) *ProjectCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillsPath", Val: v})
	return b
}
func (b *ProjectCreateBuilder) SetAgentsPath(v string) *ProjectCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentsPath", Val: v})
	return b
}

func (b *ProjectCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Project]) *ProjectCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectDelegate) Create() *ProjectCreateBuilder {
	return &ProjectCreateBuilder{
		CreateBuilder: &CreateBuilder[Project, ProjectSelect, ProjectOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedProjectName       uint64 = 1 << 0
	providedProjectSkillsPath uint64 = 1 << 1
	providedProjectAgentsPath uint64 = 1 << 2
)

func assignmentsToProjectCreate(assignments []FieldAssignment) (ProjectCreate, error) {
	var input ProjectCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "name":
			provided |= providedProjectName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "skillsPath":
			provided |= providedProjectSkillsPath
			if v, ok := a.Val.(string); ok {
				input.SkillsPath = v
				errs.ValidateString("skillsPath", v, true, 0, false, false)
			} else {
				errs.Add("skillsPath", a.Val, "type", "field skillsPath must be of type string")
			}
		case "agentsPath":
			provided |= providedProjectAgentsPath
			if v, ok := a.Val.(string); ok {
				input.AgentsPath = v
				errs.ValidateString("agentsPath", v, true, 0, false, false)
			} else {
				errs.Add("agentsPath", a.Val, "type", "field agentsPath must be of type string")
			}
		}
	}
	if provided&providedProjectName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}
	if provided&providedProjectSkillsPath == 0 {
		errs.Add("skillsPath", "", "required", "field SkillsPath is required")
	}
	if provided&providedProjectAgentsPath == 0 {
		errs.Add("agentsPath", "", "required", "field AgentsPath is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *ProjectCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 3)
	vals = make([]any, 0, 3)
	cols = append(cols, "name")
	if s.Name != "" {
		vals = append(vals, s.Name)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "skillsPath")
	vals = append(vals, s.SkillsPath)
	cols = append(cols, "agentsPath")
	vals = append(vals, s.AgentsPath)
	return
}

func partitionProjectInputs(dialect Dialect, inputs []*ProjectCreate) [][]*ProjectCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*ProjectCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*ProjectCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*ProjectCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*ProjectCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*ProjectCreate{inputs}
}

func (d *ProjectDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *ProjectSelect, omits *ProjectOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Project, error) {
	input, err := assignmentsToProjectCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectProjectCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectCreateArgs) (*Project, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectProjectCols(a.Select, omits)
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
			curr = func(c context.Context, a *ProjectCreateArgs) (*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type ProjectCreateManyBuilder struct {
	*CreateManyBuilder[Project]
}

func (b *ProjectCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *ProjectConflictBuilder[ProjectCreateManyBuilder] {
	return &ProjectConflictBuilder[ProjectCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type ProjectCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Project, ProjectSelect, ProjectOmit]
}

func (b *ProjectCreateManyAndReturnBuilder) Select(s ProjectSelect) *ProjectCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *ProjectCreateManyAndReturnBuilder) Omit(o ProjectOmit) *ProjectCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *ProjectCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *ProjectConflictBuilder[ProjectCreateManyAndReturnBuilder] {
	return &ProjectConflictBuilder[ProjectCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToProjectRecordInputs(builders []*ProjectCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *ProjectDelegate) CreateMany(builders ...*ProjectCreateBuilder) *ProjectCreateManyBuilder {
	return &ProjectCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Project]{
			records:  createBuildersToProjectRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *ProjectDelegate) CreateManyAndReturn(builders ...*ProjectCreateBuilder) *ProjectCreateManyAndReturnBuilder {
	return &ProjectCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Project, ProjectSelect, ProjectOmit]{
			records:  createBuildersToProjectRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToProjectCreateInputs(records []RecordInput) ([]*ProjectCreate, error) {
	structs := make([]ProjectCreate, len(records))
	inputs := make([]*ProjectCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToProjectCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *ProjectDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToProjectCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &ProjectCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProjectSelect, omits *ProjectOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Project, error) {
	inputs, err := recordsToProjectCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ProjectCreateManyAndReturnArgs) ([]*Project, error) {
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
			curr = func(c context.Context, a *ProjectCreateManyAndReturnArgs) ([]*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *ProjectSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Project, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Project
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Project.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Project.loadRelations(ctx, []*Project{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Project", cols, returningCols, projectPKCols, conflictTarget, conflictAction, len(vals))
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

		var res Project
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, projectPKCols)
}

func (d *ProjectDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Project, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Project") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Project")
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

	var res Project
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *ProjectDelegate) buildBulkInsertSQL(q *Queries, batch []*ProjectCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 3)
	for i, c := range projectDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Project")
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
			case "skillsPath":
				vals = append(vals, input.SkillsPath)
			case "agentsPath":
				vals = append(vals, input.AgentsPath)
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

func applyProjectConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanProjectRows(rows *sql.Rows, returningCols []string) ([]*Project, error) {
	var records []*Project
	for rows.Next() {
		var res Project
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

func (d *ProjectDelegate) runCreateMany(ctx context.Context, inputs []*ProjectCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionProjectInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectConflictClause(d.client.dialect, queryStr, vals, cols, projectPKCols, conflictTarget, conflictAction)

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

func (d *ProjectDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*ProjectCreate,
	selects *ProjectSelect,
	omits *ProjectOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Project, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionProjectInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Project
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Project.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Project.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectProjectCols(selects, omits, projectPKCols...)
	recordsOut := make([]*Project, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectConflictClause(d.client.dialect, queryStr, vals, cols, projectPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanProjectRows(rows, returningCols)
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

func (d *ProjectDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*ProjectCreate,
	selects *ProjectSelect,
	omits *ProjectOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Project, error) {
	batches := partitionProjectInputs(d.client.dialect, inputs)
	returningCols := selectProjectCols(selects, omits, projectPKCols...)
	recordsOut := make([]*Project, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyProjectConflictClause(d.client.dialect, queryStr, vals, cols, projectPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Project") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Project")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, projectPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanProjectRows(rows, returningCols)
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

type ProjectConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *ProjectConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProjectConflictBuilder[B]) Update(fn func(u *ProjectUpsert)) *B {
	cb.setAction(*ProjectConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// ProjectConflictUpdate creates a custom ConflictAction for Project upsert conflicts.
func ProjectConflictUpdate(fn func(u *ProjectUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newProjectUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type ProjectUpsert struct {
	Name       fieldUpsert[string]
	SkillsPath fieldUpsert[string]
	AgentsPath fieldUpsert[string]
}

func newProjectUpsert(up *ConflictUpdate) *ProjectUpsert {
	return &ProjectUpsert{
		Name:       fieldUpsert[string]{column: "name", update: up},
		SkillsPath: fieldUpsert[string]{column: "skillsPath", update: up},
		AgentsPath: fieldUpsert[string]{column: "agentsPath", update: up},
	}
}

type ProjectUpdateBuilder struct {
	*UpdateBuilder[Project, ProjectSelect, ProjectOmit]
}

type ProjectUpdateManyBuilder struct {
	*UpdateManyBuilder[Project]
}

type ProjectUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Project, ProjectSelect, ProjectOmit]
}

func (b *ProjectUpdateBuilder) SetName(v string) *ProjectUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *ProjectUpdateManyBuilder) SetName(v string) *ProjectUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *ProjectUpdateManyAndReturnBuilder) SetName(v string) *ProjectUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *ProjectUpdateBuilder) SetSkillsPath(v string) *ProjectUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillsPath", Val: v})
	return b
}

func (b *ProjectUpdateManyBuilder) SetSkillsPath(v string) *ProjectUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillsPath", Val: v})
	return b
}

func (b *ProjectUpdateManyAndReturnBuilder) SetSkillsPath(v string) *ProjectUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillsPath", Val: v})
	return b
}
func (b *ProjectUpdateBuilder) SetAgentsPath(v string) *ProjectUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentsPath", Val: v})
	return b
}

func (b *ProjectUpdateManyBuilder) SetAgentsPath(v string) *ProjectUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentsPath", Val: v})
	return b
}

func (b *ProjectUpdateManyAndReturnBuilder) SetAgentsPath(v string) *ProjectUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentsPath", Val: v})
	return b
}

func (b *ProjectUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Project]) *ProjectUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Project]) *ProjectUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ProjectUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Project]) *ProjectUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ProjectDelegate) Update(where UniquePredicate[Project], additional ...PredicateOf[Project]) *ProjectUpdateBuilder {
	return &ProjectUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Project, ProjectSelect, ProjectOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *ProjectDelegate) UpdateMany(preds ...PredicateOf[Project]) *ProjectUpdateManyBuilder {
	return &ProjectUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Project]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *ProjectDelegate) UpdateManyAndReturn(preds ...PredicateOf[Project]) *ProjectUpdateManyAndReturnBuilder {
	return &ProjectUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Project, ProjectSelect, ProjectOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *ProjectDelegate) buildUpdateSQL(preds []PredicateOf[Project], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Project")
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

func (d *ProjectDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Project], additional []PredicateOf[Project], assignments []FieldAssignment, selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
	allWhere := make([]PredicateOf[Project], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToProjectUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectUpdateArgs) (*Project, error) {
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
			curr = func(c context.Context, a *ProjectUpdateArgs) (*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Project], cols []string, vals []any, selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
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
		var res *Project
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Project.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Project.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectCols(selects, omits, projectPKCols...)
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
		return nil, &NotFoundError{Model: "Project"}
	}

	var res Project
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Project{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *ProjectDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Project], cols []string, vals []any) (int64, error) {
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

func (d *ProjectDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Project], cols []string, vals []any, selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Project"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *ProjectDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Project], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToProjectUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &ProjectUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *ProjectUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *ProjectDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Project], assignments []FieldAssignment, selects *ProjectSelect, omits *ProjectOmit) ([]*Project, error) {
	input, err := assignmentsToProjectUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectUpdateManyAndReturnArgs) ([]*Project, error) {
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
			curr = func(c context.Context, a *ProjectUpdateManyAndReturnArgs) ([]*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Project], cols []string, vals []any, selects *ProjectSelect, omits *ProjectOmit) ([]*Project, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Project]{Where: preds}, selects, omits)
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
		var res []*Project
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Project.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Project.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectProjectCols(selects, omits, projectPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanProjectRows(rows, returningCols)
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

func (d *ProjectDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Project], cols []string, vals []any, selects *ProjectSelect, omits *ProjectOmit) ([]*Project, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Project{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Project]{Where: preds}, selects, omits)
}
func (d *ProjectDelegate) FindUnique(where UniquePredicate[Project], additional ...PredicateOf[Project]) *FindUniqueBuilder[Project, ProjectSelect, ProjectOmit] {
	return &FindUniqueBuilder[Project, ProjectSelect, ProjectOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *ProjectDelegate) FindFirst(preds ...PredicateOf[Project]) *FindFirstBuilder[Project, ProjectSelect, ProjectOmit] {
	return &FindFirstBuilder[Project, ProjectSelect, ProjectOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *ProjectDelegate) FindMany(preds ...PredicateOf[Project]) *FindManyBuilder[Project, ProjectSelect, ProjectOmit] {
	return &FindManyBuilder[Project, ProjectSelect, ProjectOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *ProjectDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Project], additional []PredicateOf[Project], selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
	allWhere := make([]PredicateOf[Project], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectFindUniqueArgs) (*Project, error) {
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
			curr = func(c context.Context, a *ProjectFindUniqueArgs) (*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Project],
	selects *ProjectSelect,
	omits *ProjectOmit,
) (*Project, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectFindFirstArgs) (*Project, error) {
		return d.runFindFirst(c, QueryParams[Project]{
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
			curr = func(c context.Context, a *ProjectFindFirstArgs) (*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Project],
	selects *ProjectSelect,
	omits *ProjectOmit,
) ([]*Project, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ProjectFindManyArgs) ([]*Project, error) {
		return d.runFindMany(c, QueryParams[Project]{
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
			curr = func(c context.Context, a *ProjectFindManyArgs) ([]*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Project], selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
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
	returningCols := selectProjectCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Project{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *ProjectDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Project],
	selects *ProjectSelect,
	omits *ProjectOmit,
) (*Project, error) {
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

func (d *ProjectDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Project],
	selects *ProjectSelect,
	omits *ProjectOmit,
) ([]*Project, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, projectPKCols, projectUniqueCols, "Project", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, projectPKCols, projectUniqueCols, isCursorQuery, params.Take)
	returningCols := selectProjectCols(selects, omits)

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

func (d *ProjectDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Project, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Project", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Project
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *ProjectDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Project, error) {
	query := buildSelectSQL(d.client, "Project", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Project, 0)
	for rows.Next() {
		var res Project
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
func (d *ProjectDelegate) DeleteMany(preds ...PredicateOf[Project]) *DeleteManyBuilder[Project] {
	return &DeleteManyBuilder[Project]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *ProjectDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Project]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &ProjectDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *ProjectDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ProjectDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Project]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Project")
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

func (d *ProjectDelegate) Delete(where UniquePredicate[Project], additional ...PredicateOf[Project]) *DeleteBuilder[Project, ProjectSelect, ProjectOmit] {
	return &DeleteBuilder[Project, ProjectSelect, ProjectOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *ProjectDelegate) executeDelete(ctx context.Context, where UniquePredicate[Project], additional []PredicateOf[Project], selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
	allWhere := make([]PredicateOf[Project], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullProjectSelect()
	}

	args := &ProjectDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ProjectDeleteArgs) (*Project, error) {
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
			curr = func(c context.Context, a *ProjectDeleteArgs) (*Project, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runDelete(ctx context.Context, where []PredicateOf[Project], selects *ProjectSelect, omits *ProjectOmit) (*Project, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectProjectCols(selects, omits, projectPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Project
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Project.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Project"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Project")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Project]
			pkPreds = append(pkPreds, Predicate[Project]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "Project")

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
		return nil, &NotFoundError{Model: "Project"}
	}

	var row Project
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *ProjectDelegate) Count(preds ...PredicateOf[Project]) *CountBuilder[Project] {
	return &CountBuilder[Project]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *ProjectDelegate) executeCount(ctx context.Context, params QueryParams[Project]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &ProjectCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *ProjectCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Project]{
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
			curr = func(c context.Context, a *ProjectCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ProjectDelegate) runCount(ctx context.Context, params QueryParams[Project]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Project")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Project")
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
func (d *ProjectDelegate) loadRelations(ctx context.Context, records []*Project, selects *ProjectSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Skills != nil {
		relationSelects, relationOmits, relationParams := selects.Skills.GetRelationParams()
		returningCols := selectProjectSkillCols(relationSelects, relationOmits, "projectName")
		// Inverse holds the FK: ProjectSkill.projectName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Project) string { return p.Name }),
			"ProjectSkill",
			"projectName",
			returningCols,
			scanInto(returningCols, (*ProjectSkill).ScanFields),
			directKey(func(c *ProjectSkill) string { return c.ProjectName }),
			appendMany(func(p *Project) *[]*ProjectSkill { return &p.Skills }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skills: %w", err)
		}
		if err := d.client.ProjectSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Generalizes != nil {
		relationSelects, relationOmits, relationParams := selects.Generalizes.GetRelationParams()
		returningCols := selectProjectSpecializationCols(relationSelects, relationOmits, "generalizesName")
		// Inverse holds the FK: ProjectSpecialization.generalizesName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Project) string { return p.Name }),
			"ProjectSpecialization",
			"generalizesName",
			returningCols,
			scanInto(returningCols, (*ProjectSpecialization).ScanFields),
			directKey(func(c *ProjectSpecialization) string { return c.GeneralizesName }),
			appendMany(func(p *Project) *[]*ProjectSpecialization { return &p.Generalizes }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading generalizes: %w", err)
		}
		if err := d.client.ProjectSpecialization.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Specializes != nil {
		relationSelects, relationOmits, relationParams := selects.Specializes.GetRelationParams()
		returningCols := selectProjectSpecializationCols(relationSelects, relationOmits, "specializesName")
		// Inverse holds the FK: ProjectSpecialization.specializesName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Project) string { return p.Name }),
			"ProjectSpecialization",
			"specializesName",
			returningCols,
			scanInto(returningCols, (*ProjectSpecialization).ScanFields),
			directKey(func(c *ProjectSpecialization) string { return c.SpecializesName }),
			appendMany(func(p *Project) *[]*ProjectSpecialization { return &p.Specializes }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading specializes: %w", err)
		}
		if err := d.client.ProjectSpecialization.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Apps != nil {
		relationSelects, relationOmits, relationParams := selects.Apps.GetRelationParams()
		returningCols := selectAppCols(relationSelects, relationOmits, "project_id")
		// Inverse holds the FK: App.project_id
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Project) string { return p.Name }),
			"App",
			"project_id",
			returningCols,
			scanInto(returningCols, (*App).ScanFields),
			directKey(func(c *App) string { return c.Project_id }),
			appendMany(func(p *Project) *[]*App { return &p.Apps }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading apps: %w", err)
		}
		if err := d.client.App.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Agents != nil {
		// Many-to-many through join table _ProjectAgents
		relationSelects, relationOmits, relationParams := selects.Agents.GetRelationParams()
		returningCols := selectAgentCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *Project) string { return p.Name }),
			"Agent",
			"id",
			"_ProjectAgents",
			"projectName",
			"agentId",
			returningCols,
			(*Agent).ScanFields,
			appendMany(func(p *Project) *[]*Agent { return &p.Agents }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agents: %w", err)
		}
		if err := d.client.Agent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
