package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Skill represents the database model
type Skill struct {
	Id            string             `db:"id" json:"id"`
	Name          string             `db:"name" json:"name"`
	Description   string             `db:"description" json:"description"`
	Url           string             `db:"url" json:"url"`
	Content       string             `db:"content" json:"content"`
	Keywords      string             `db:"keywords" json:"keywords"`
	CreatedAt     time.Time          `db:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `db:"updatedAt" json:"updatedAt"`
	Dependencies  []*SkillDependency `json:"dependencies,omitempty"`
	Dependents    []*SkillDependency `json:"dependents,omitempty"`
	Projects      []*ProjectSkill    `json:"projects,omitempty"`
	References    []*SkillReference  `json:"references,omitempty"`
	Trashes       []*TrashSkill      `json:"trashes,omitempty"`
	Agents        []*AgentSkill      `json:"agents,omitempty"`
	Apps          []*AppSkill        `json:"apps,omitempty"`
	Users         []*User            `json:"users,omitempty"`
	Userblacklist []*User            `json:"userblacklist,omitempty"`
}

// SkillCreate contains model input fields for Skill creation operations.
//
// Fields for Skill:
//
//	id          string    default: cuid()
//	name        string    required
//	description string    required
//	url         string    required
//	content     string    default:
//	keywords    string    default:
//	createdAt   time.Time default: now()
//	updatedAt   time.Time required
type SkillCreate struct {
	Id          *string    `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Url         string     `json:"url"`
	Content     *string    `json:"content"`
	Keywords    *string    `json:"keywords"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// colMask returns a bit mask of columns that are set
func (s *SkillCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	mask |= 1 << 3
	if s.Content != nil {
		mask |= 1 << 4
	}
	if s.Keywords != nil {
		mask |= 1 << 5
	}
	mask |= 1 << 6
	mask |= 1 << 7
	return mask
}

// SkillUpdate contains model input fields for Skill update operations.
type SkillUpdate struct {
	Id          *string    `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Url         *string    `json:"url"`
	Content     *string    `json:"content"`
	Keywords    *string    `json:"keywords"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (u *SkillUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, u.Id)
	}
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	if u.Description != nil {
		cols = append(cols, "description")
		vals = append(vals, u.Description)
	}
	if u.Url != nil {
		cols = append(cols, "url")
		vals = append(vals, u.Url)
	}
	if u.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, u.Content)
	}
	if u.Keywords != nil {
		cols = append(cols, "keywords")
		vals = append(vals, u.Keywords)
	}
	if u.CreatedAt != nil {
		cols = append(cols, "createdAt")
		vals = append(vals, u.CreatedAt)
	}
	if u.UpdatedAt != nil {
		cols = append(cols, "updatedAt")
		vals = append(vals, u.UpdatedAt)
	}
	return cols, vals
}

func assignmentsToSkillUpdate(assignments []FieldAssignment) (SkillUpdate, error) {
	var input SkillUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				errs.ValidateString("id", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Id = v
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "name":
			if v, ok := a.Val.(string); ok {
				input.Name = &v
				errs.ValidateString("name", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Name = v
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "description":
			if v, ok := a.Val.(string); ok {
				input.Description = &v
				errs.ValidateString("description", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Description = v
			} else {
				errs.Add("description", a.Val, "type", "field description must be of type string")
			}
		case "url":
			if v, ok := a.Val.(string); ok {
				input.Url = &v
				errs.ValidateString("url", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Url = v
			} else {
				errs.Add("url", a.Val, "type", "field url must be of type string")
			}
		case "content":
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Content = v
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "keywords":
			if v, ok := a.Val.(string); ok {
				input.Keywords = &v
				errs.ValidateString("keywords", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Keywords = v
			} else {
				errs.Add("keywords", a.Val, "type", "field keywords must be of type string")
			}
		case "createdAt":
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else if v, ok := a.Val.(*time.Time); ok {
				input.CreatedAt = v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		case "updatedAt":
			if v, ok := a.Val.(time.Time); ok {
				input.UpdatedAt = &v
			} else if v, ok := a.Val.(*time.Time); ok {
				input.UpdatedAt = v
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type time.Time")
			}
		}
	}
	if input.UpdatedAt == nil {
		now := time.Now().Truncate(time.Microsecond)
		input.UpdatedAt = &now
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// SkillSelect specifies which scalar and relation fields to select for Skill.
//
// Selectable fields:
//
//	-- Scalars --
//	id          (bool)
//	name        (bool)
//	description (bool)
//	url         (bool)
//	content     (bool)
//	keywords    (bool)
//	createdAt   (bool)
//	updatedAt   (bool)
//	-- Relations --
//	dependencies ([]SkillDependency)
//	dependents  ([]SkillDependency)
//	projects    ([]ProjectSkill)
//	references  ([]SkillReference)
//	trashes     ([]TrashSkill)
//	agents      ([]AgentSkill)
//	apps        ([]AppSkill)
//	users       ([]User)
//	userblacklist ([]User)
type SkillSelect struct {
	Id            bool                       `json:"id"`
	Name          bool                       `json:"name"`
	Description   bool                       `json:"description"`
	Url           bool                       `json:"url"`
	Content       bool                       `json:"content"`
	Keywords      bool                       `json:"keywords"`
	CreatedAt     bool                       `json:"createdAt"`
	UpdatedAt     bool                       `json:"updatedAt"`
	Dependencies  SkillDependencySelectQuery `json:"dependencies,omitempty"`
	Dependents    SkillDependencySelectQuery `json:"dependents,omitempty"`
	Projects      ProjectSkillSelectQuery    `json:"projects,omitempty"`
	References    SkillReferenceSelectQuery  `json:"references,omitempty"`
	Trashes       TrashSkillSelectQuery      `json:"trashes,omitempty"`
	Agents        AgentSkillSelectQuery      `json:"agents,omitempty"`
	Apps          AppSkillSelectQuery        `json:"apps,omitempty"`
	Users         UserSelectQuery            `json:"users,omitempty"`
	Userblacklist UserSelectQuery            `json:"userblacklist,omitempty"`
}

var fullSkillSelectVal = &SkillSelect{
	Id:          true,
	Name:        true,
	Description: true,
	Url:         true,
	Content:     true,
	Keywords:    true,
	CreatedAt:   true,
	UpdatedAt:   true,
}

func fullSkillSelect() *SkillSelect {
	return fullSkillSelectVal
}

func (s *SkillSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.Name || s.Description || s.Url || s.Content || s.Keywords || s.CreatedAt || s.UpdatedAt
}

func (s *SkillSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type SkillOmit struct {
	Id          bool `json:"id"`
	Name        bool `json:"name"`
	Description bool `json:"description"`
	Url         bool `json:"url"`
	Content     bool `json:"content"`
	Keywords    bool `json:"keywords"`
	CreatedAt   bool `json:"createdAt"`
	UpdatedAt   bool `json:"updatedAt"`
}

type SkillSelectQuery interface {
	GetRelationParams() (*SkillSelect, *SkillOmit, QueryParams[Skill])
}

func (s *SkillSelect) GetRelationParams() (*SkillSelect, *SkillOmit, QueryParams[Skill]) {
	return s, nil, QueryParams[Skill]{}
}

type SkillQueryBuilder struct {
	selects *SkillSelect
	omits   *SkillOmit
	where   []PredicateOf[Skill]
	take    *int
	skip    *int
	orderBy []OrderBy[Skill]
	cursor  UniquePredicate[Skill]
}

func (b *SkillQueryBuilder) Where(preds ...PredicateOf[Skill]) *SkillQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *SkillQueryBuilder) Take(limit int) *SkillQueryBuilder {
	b.take = &limit
	return b
}

func (b *SkillQueryBuilder) Skip(offset int) *SkillQueryBuilder {
	b.skip = &offset
	return b
}

func (b *SkillQueryBuilder) OrderBy(orders ...OrderBy[Skill]) *SkillQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *SkillQueryBuilder) Cursor(where UniquePredicate[Skill]) *SkillQueryBuilder {
	b.cursor = where
	return b
}

func (b *SkillQueryBuilder) Select(s SkillSelect) *SkillQueryBuilder {
	b.selects = &s
	return b
}

func (b *SkillQueryBuilder) Omit(o SkillOmit) *SkillQueryBuilder {
	b.omits = &o
	return b
}

func (b *SkillQueryBuilder) GetRelationParams() (*SkillSelect, *SkillOmit, QueryParams[Skill]) {
	if b == nil {
		return nil, nil, QueryParams[Skill]{}
	}
	return b.selects, b.omits, QueryParams[Skill]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// SkillCreateArgs is the input argument passed to Skill Create extension hooks.
//
// Fields for Skill:
//
//	id          string    default: cuid()
//	name        string    required
//	description string    required
//	url         string    required
//	content     string    default:
//	keywords    string    default:
//	createdAt   time.Time default: now()
//	updatedAt   time.Time required
//
// Relations for Skill:
//
//	dependencies ([]SkillDependency)
//	dependents  ([]SkillDependency)
//	projects    ([]ProjectSkill)
//	references  ([]SkillReference)
//	trashes     ([]TrashSkill)
//	agents      ([]AgentSkill)
//	apps        ([]AppSkill)
//	users       ([]User)
//	userblacklist ([]User)
type SkillCreateArgs struct {
	// Data contains the model fields to insert.
	Data *SkillCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *SkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// SkillCreateManyArgs is the input argument passed to Skill CreateMany extension hooks.
//
// Fields for Skill:
//
//	id          string    default: cuid()
//	name        string    required
//	description string    required
//	url         string    required
//	content     string    default:
//	keywords    string    default:
//	createdAt   time.Time default: now()
//	updatedAt   time.Time required
type SkillCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillCreateManyArgs) AppendData(builders ...*SkillCreateBuilder) *SkillCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillCreateManyAndReturnArgs is the input argument passed to Skill CreateManyAndReturn extension hooks.
//
// Fields for Skill:
//
//	id          string    default: cuid()
//	name        string    required
//	description string    required
//	url         string    required
//	content     string    default:
//	keywords    string    default:
//	createdAt   time.Time default: now()
//	updatedAt   time.Time required
//
// Relations for Skill:
//
//	dependencies ([]SkillDependency)
//	dependents  ([]SkillDependency)
//	projects    ([]ProjectSkill)
//	references  ([]SkillReference)
//	trashes     ([]TrashSkill)
//	agents      ([]AgentSkill)
//	apps        ([]AppSkill)
//	users       ([]User)
//	userblacklist ([]User)
type SkillCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *SkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillCreateManyAndReturnArgs) AppendData(builders ...*SkillCreateBuilder) *SkillCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillFindUniqueArgs is the input argument passed to Skill FindUnique extension hooks.
type SkillFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Skill]
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillSelect
}

func (a *SkillFindUniqueArgs) SetWhere(unique UniquePredicate[Skill], additional ...PredicateOf[Skill]) *SkillFindUniqueArgs {
	a.Where = make([]PredicateOf[Skill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillFindFirstArgs is the input argument passed to Skill FindFirst extension hooks.
type SkillFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Skill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Skill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillSelect
}

func (a *SkillFindFirstArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillFindFirstArgs {
	a.Where = preds
	return a
}

func (a *SkillFindFirstArgs) SetOrderBy(orders ...OrderBy[Skill]) *SkillFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillFindFirstArgs) SetCursor(cursor UniquePredicate[Skill]) *SkillFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillFindFirstArgs) SetSkip(n int) *SkillFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *SkillFindFirstArgs) SetTake(n int) *SkillFindFirstArgs {
	a.Take = &n
	return a
}

// SkillFindManyArgs is the input argument passed to Skill FindMany extension hooks.
type SkillFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Skill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Skill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillSelect
}

func (a *SkillFindManyArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillFindManyArgs {
	a.Where = preds
	return a
}

func (a *SkillFindManyArgs) SetOrderBy(orders ...OrderBy[Skill]) *SkillFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillFindManyArgs) SetCursor(cursor UniquePredicate[Skill]) *SkillFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillFindManyArgs) SetSkip(n int) *SkillFindManyArgs {
	a.Skip = &n
	return a
}

func (a *SkillFindManyArgs) SetTake(n int) *SkillFindManyArgs {
	a.Take = &n
	return a
}

// SkillCountArgs is the input argument passed to Skill Count extension hooks.
type SkillCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *SkillCountArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillCountArgs {
	a.Where = preds
	return a
}

func (a *SkillCountArgs) SetSkip(n int) *SkillCountArgs {
	a.Skip = &n
	return a
}

func (a *SkillCountArgs) SetTake(n int) *SkillCountArgs {
	a.Take = &n
	return a
}

// SkillDeleteArgs is the input argument passed to Skill Delete extension hooks.
type SkillDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Skill]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *SkillSelect
}

func (a *SkillDeleteArgs) SetWhere(unique UniquePredicate[Skill], additional ...PredicateOf[Skill]) *SkillDeleteArgs {
	a.Where = make([]PredicateOf[Skill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillDeleteManyArgs is the input argument passed to Skill DeleteMany extension hooks.
type SkillDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
}

func (a *SkillDeleteManyArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillDeleteManyArgs {
	a.Where = preds
	return a
}

// SkillUpdateArgs is the input argument passed to Skill Update extension hooks.
type SkillUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Skill]
	// Data contains the model fields to update.
	Data *SkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillSelect
}

func (a *SkillUpdateArgs) SetWhere(unique UniquePredicate[Skill], additional ...PredicateOf[Skill]) *SkillUpdateArgs {
	a.Where = make([]PredicateOf[Skill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillUpdateManyArgs is the input argument passed to Skill UpdateMany extension hooks.
type SkillUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
	// Data contains the model fields to update.
	Data *SkillUpdate
}

func (a *SkillUpdateManyArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillUpdateManyArgs {
	a.Where = preds
	return a
}

// SkillUpdateManyAndReturnArgs is the input argument passed to Skill UpdateManyAndReturn extension hooks.
type SkillUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Skill]
	// Data contains the model fields to update.
	Data *SkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillSelect
}

func (a *SkillUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Skill]) *SkillUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type SkillCreateQuery = func(ctx context.Context, args *SkillCreateArgs) (*Skill, error)
type SkillCreateManyQuery = func(ctx context.Context, args *SkillCreateManyArgs) (int64, error)
type SkillCreateManyAndReturnQuery = func(ctx context.Context, args *SkillCreateManyAndReturnArgs) ([]*Skill, error)
type SkillFindUniqueQuery = func(ctx context.Context, args *SkillFindUniqueArgs) (*Skill, error)
type SkillFindFirstQuery = func(ctx context.Context, args *SkillFindFirstArgs) (*Skill, error)
type SkillFindManyQuery = func(ctx context.Context, args *SkillFindManyArgs) ([]*Skill, error)
type SkillDeleteQuery = func(ctx context.Context, args *SkillDeleteArgs) (*Skill, error)
type SkillDeleteManyQuery = func(ctx context.Context, args *SkillDeleteManyArgs) (int64, error)
type SkillCountQuery = func(ctx context.Context, args *SkillCountArgs) (int64, error)
type SkillUpdateQuery = func(ctx context.Context, args *SkillUpdateArgs) (*Skill, error)
type SkillUpdateManyQuery = func(ctx context.Context, args *SkillUpdateManyArgs) (int64, error)
type SkillUpdateManyAndReturnQuery = func(ctx context.Context, args *SkillUpdateManyAndReturnArgs) ([]*Skill, error)

type SkillExtension struct {
	Create              func(ctx context.Context, args *SkillCreateArgs, next SkillCreateQuery) (*Skill, error)
	CreateMany          func(ctx context.Context, args *SkillCreateManyArgs, next SkillCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *SkillCreateManyAndReturnArgs, next SkillCreateManyAndReturnQuery) ([]*Skill, error)
	FindUnique          func(ctx context.Context, args *SkillFindUniqueArgs, next SkillFindUniqueQuery) (*Skill, error)
	FindFirst           func(ctx context.Context, args *SkillFindFirstArgs, next SkillFindFirstQuery) (*Skill, error)
	FindMany            func(ctx context.Context, args *SkillFindManyArgs, next SkillFindManyQuery) ([]*Skill, error)
	Delete              func(ctx context.Context, args *SkillDeleteArgs, next SkillDeleteQuery) (*Skill, error)
	DeleteMany          func(ctx context.Context, args *SkillDeleteManyArgs, next SkillDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *SkillCountArgs, next SkillCountQuery) (int64, error)
	Update              func(ctx context.Context, args *SkillUpdateArgs, next SkillUpdateQuery) (*Skill, error)
	UpdateMany          func(ctx context.Context, args *SkillUpdateManyArgs, next SkillUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *SkillUpdateManyAndReturnArgs, next SkillUpdateManyAndReturnQuery) ([]*Skill, error)
}

type SkillDelegate struct {
	client     *Queries
	extensions []SkillExtension
}

func (d *SkillDelegate) Use(exts ...SkillExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Skill) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "name":
			targets[i] = &m.Name
		case "description":
			targets[i] = &m.Description
		case "url":
			targets[i] = &m.Url
		case "content":
			targets[i] = &m.Content
		case "keywords":
			targets[i] = &m.Keywords
		case "createdAt":
			targets[i] = &m.CreatedAt
		case "updatedAt":
			targets[i] = &m.UpdatedAt
		}
	}
	return targets
}

var skillDefaultCols = []string{
	"id",
	"name",
	"description",
	"url",
	"content",
	"keywords",
	"createdAt",
	"updatedAt",
}

var skillPKCols = []string{
	"id",
}

var skillUniqueCols = []string{
	"id",
	"url",
}

func selectSkillCols(selects *SkillSelect, omits *SkillOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return skillDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Name || selects.Description || selects.Url || selects.Content || selects.Keywords || selects.CreatedAt || selects.UpdatedAt || selects.Dependencies != nil || selects.Dependents != nil || selects.Projects != nil || selects.References != nil || selects.Trashes != nil || selects.Agents != nil || selects.Apps != nil || selects.Users != nil || selects.Userblacklist != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
		{"description", selects != nil && selects.Description, omits != nil && omits.Description, false},
		{"url", selects != nil && selects.Url, omits != nil && omits.Url, false},
		{"content", selects != nil && selects.Content, omits != nil && omits.Content, false},
		{"keywords", selects != nil && selects.Keywords, omits != nil && omits.Keywords, false},
		{"createdAt", selects != nil && selects.CreatedAt, omits != nil && omits.CreatedAt, false},
		{"updatedAt", selects != nil && selects.UpdatedAt, omits != nil && omits.UpdatedAt, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *SkillSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Dependencies != nil || s.Dependents != nil || s.Projects != nil || s.References != nil || s.Trashes != nil || s.Agents != nil || s.Apps != nil || s.Users != nil || s.Userblacklist != nil
}

type SkillUpsertBuilder struct {
	*CreateBuilder[Skill, SkillSelect, SkillOmit]
}

func (b *SkillUpsertBuilder) Select(s SkillSelect) *SkillUpsertBuilder {
	b.selects = &s
	return b
}

func (b *SkillUpsertBuilder) Omit(o SkillOmit) *SkillUpsertBuilder {
	b.omits = &o
	return b
}

type SkillCreateBuilder struct {
	*CreateBuilder[Skill, SkillSelect, SkillOmit]
}

func (b *SkillCreateBuilder) Select(s SkillSelect) *SkillCreateBuilder {
	b.selects = &s
	return b
}

func (b *SkillCreateBuilder) Omit(o SkillOmit) *SkillCreateBuilder {
	b.omits = &o
	return b
}

func (b *SkillCreateBuilder) OnConflict(target UniqueConstraintTarget) *SkillConflictBuilder[SkillUpsertBuilder] {
	upsertBuilder := &SkillUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &SkillConflictBuilder[SkillUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *SkillCreateBuilder) SetId(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetName(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetDescription(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetUrl(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "url", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetContent(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetKeywords(v string) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "keywords", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetCreatedAt(v time.Time) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *SkillCreateBuilder) SetUpdatedAt(v time.Time) *SkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *SkillCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Skill]) *SkillCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillDelegate) Create() *SkillCreateBuilder {
	return &SkillCreateBuilder{
		CreateBuilder: &CreateBuilder[Skill, SkillSelect, SkillOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedSkillId          uint64 = 1 << 0
	providedSkillName        uint64 = 1 << 1
	providedSkillDescription uint64 = 1 << 2
	providedSkillUrl         uint64 = 1 << 3
	providedSkillContent     uint64 = 1 << 4
	providedSkillKeywords    uint64 = 1 << 5
	providedSkillCreatedAt   uint64 = 1 << 6
	providedSkillUpdatedAt   uint64 = 1 << 7
)

func assignmentsToSkillCreate(assignments []FieldAssignment) (SkillCreate, error) {
	var input SkillCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedSkillId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				errs.ValidateString("id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "name":
			provided |= providedSkillName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "description":
			provided |= providedSkillDescription
			if v, ok := a.Val.(string); ok {
				input.Description = v
				errs.ValidateString("description", v, true, 0, false, false)
			} else {
				errs.Add("description", a.Val, "type", "field description must be of type string")
			}
		case "url":
			provided |= providedSkillUrl
			if v, ok := a.Val.(string); ok {
				input.Url = v
				errs.ValidateString("url", v, true, 0, false, false)
			} else {
				errs.Add("url", a.Val, "type", "field url must be of type string")
			}
		case "content":
			provided |= providedSkillContent
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "keywords":
			provided |= providedSkillKeywords
			if v, ok := a.Val.(string); ok {
				input.Keywords = &v
				errs.ValidateString("keywords", v, false, 0, false, false)
			} else {
				errs.Add("keywords", a.Val, "type", "field keywords must be of type string")
			}
		case "createdAt":
			provided |= providedSkillCreatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		case "updatedAt":
			provided |= providedSkillUpdatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.UpdatedAt = v
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type time.Time")
			}
		}
	}
	if provided&providedSkillName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}
	if provided&providedSkillDescription == 0 {
		errs.Add("description", "", "required", "field Description is required")
	}
	if provided&providedSkillUrl == 0 {
		errs.Add("url", "", "required", "field Url is required")
	}
	if provided&providedSkillUpdatedAt == 0 {
		input.UpdatedAt = time.Now().Truncate(time.Microsecond)
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *SkillCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 8)
	vals = make([]any, 0, 8)
	cols = append(cols, "id")
	if s.Id != nil {
		vals = append(vals, *s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "name")
	vals = append(vals, s.Name)
	cols = append(cols, "description")
	vals = append(vals, s.Description)
	cols = append(cols, "url")
	vals = append(vals, s.Url)
	if s.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, *s.Content)
	}
	if s.Keywords != nil {
		cols = append(cols, "keywords")
		vals = append(vals, *s.Keywords)
	}
	cols = append(cols, "createdAt")
	if s.CreatedAt != nil {
		vals = append(vals, *s.CreatedAt)
	} else {
		vals = append(vals, time.Now())
	}
	cols = append(cols, "updatedAt")
	if !s.UpdatedAt.IsZero() {
		vals = append(vals, s.UpdatedAt)
	} else {
		vals = append(vals, time.Now().Truncate(time.Microsecond))
	}
	return
}

func partitionSkillInputs(dialect Dialect, inputs []*SkillCreate) [][]*SkillCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*SkillCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*SkillCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*SkillCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*SkillCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*SkillCreate{inputs}
}

func (d *SkillDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *SkillSelect, omits *SkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Skill, error) {
	input, err := assignmentsToSkillCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectSkillCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillCreateArgs) (*Skill, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectSkillCols(a.Select, omits)
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
			curr = func(c context.Context, a *SkillCreateArgs) (*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type SkillCreateManyBuilder struct {
	*CreateManyBuilder[Skill]
}

func (b *SkillCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *SkillConflictBuilder[SkillCreateManyBuilder] {
	return &SkillConflictBuilder[SkillCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type SkillCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Skill, SkillSelect, SkillOmit]
}

func (b *SkillCreateManyAndReturnBuilder) Select(s SkillSelect) *SkillCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *SkillCreateManyAndReturnBuilder) Omit(o SkillOmit) *SkillCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *SkillCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *SkillConflictBuilder[SkillCreateManyAndReturnBuilder] {
	return &SkillConflictBuilder[SkillCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToSkillRecordInputs(builders []*SkillCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *SkillDelegate) CreateMany(builders ...*SkillCreateBuilder) *SkillCreateManyBuilder {
	return &SkillCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Skill]{
			records:  createBuildersToSkillRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *SkillDelegate) CreateManyAndReturn(builders ...*SkillCreateBuilder) *SkillCreateManyAndReturnBuilder {
	return &SkillCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Skill, SkillSelect, SkillOmit]{
			records:  createBuildersToSkillRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToSkillCreateInputs(records []RecordInput) ([]*SkillCreate, error) {
	structs := make([]SkillCreate, len(records))
	inputs := make([]*SkillCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToSkillCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *SkillDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToSkillCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &SkillCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *SkillSelect, omits *SkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Skill, error) {
	inputs, err := recordsToSkillCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillCreateManyAndReturnArgs) ([]*Skill, error) {
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
			curr = func(c context.Context, a *SkillCreateManyAndReturnArgs) ([]*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *SkillSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Skill, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Skill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Skill.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Skill.loadRelations(ctx, []*Skill{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Skill", cols, returningCols, skillPKCols, conflictTarget, conflictAction, len(vals))
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

		var res Skill
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, skillPKCols)
}

func (d *SkillDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Skill, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Skill") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Skill")
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

	var res Skill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *SkillDelegate) buildBulkInsertSQL(q *Queries, batch []*SkillCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 8)
	for i, c := range skillDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Skill")
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
			case "id":
				if input.Id != nil && *input.Id != "" {
					vals = append(vals, *input.Id)
				} else {
					vals = append(vals, generateCUID())
				}
			case "name":
				vals = append(vals, input.Name)
			case "description":
				vals = append(vals, input.Description)
			case "url":
				vals = append(vals, input.Url)
			case "content":
				if input.Content != nil {
					vals = append(vals, *input.Content)
				} else {
					writeDefault = true
				}
			case "keywords":
				if input.Keywords != nil {
					vals = append(vals, *input.Keywords)
				} else {
					writeDefault = true
				}
			case "createdAt":
				if input.CreatedAt != nil {
					vals = append(vals, *input.CreatedAt)
				} else {
					vals = append(vals, time.Now())
				}
			case "updatedAt":
				if !input.UpdatedAt.IsZero() {
					vals = append(vals, input.UpdatedAt)
				} else {
					vals = append(vals, time.Now().Truncate(time.Microsecond))
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

func applySkillConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanSkillRows(rows *sql.Rows, returningCols []string) ([]*Skill, error) {
	var records []*Skill
	for rows.Next() {
		var res Skill
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

func (d *SkillDelegate) runCreateMany(ctx context.Context, inputs []*SkillCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionSkillInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillConflictClause(d.client.dialect, queryStr, vals, cols, skillPKCols, conflictTarget, conflictAction)

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

func (d *SkillDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*SkillCreate,
	selects *SkillSelect,
	omits *SkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Skill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionSkillInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Skill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Skill.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Skill.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectSkillCols(selects, omits, skillPKCols...)
	recordsOut := make([]*Skill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillConflictClause(d.client.dialect, queryStr, vals, cols, skillPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanSkillRows(rows, returningCols)
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

func (d *SkillDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*SkillCreate,
	selects *SkillSelect,
	omits *SkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Skill, error) {
	batches := partitionSkillInputs(d.client.dialect, inputs)
	returningCols := selectSkillCols(selects, omits, skillPKCols...)
	recordsOut := make([]*Skill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillConflictClause(d.client.dialect, queryStr, vals, cols, skillPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Skill") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Skill")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanSkillRows(rows, returningCols)
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

type SkillConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *SkillConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillConflictBuilder[B]) Update(fn func(u *SkillUpsert)) *B {
	cb.setAction(*SkillConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// SkillConflictUpdate creates a custom ConflictAction for Skill upsert conflicts.
func SkillConflictUpdate(fn func(u *SkillUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newSkillUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type SkillUpsert struct {
	Id          fieldUpsert[string]
	Name        fieldUpsert[string]
	Description fieldUpsert[string]
	Url         fieldUpsert[string]
	Content     fieldUpsert[string]
	Keywords    fieldUpsert[string]
	CreatedAt   fieldUpsert[time.Time]
	UpdatedAt   fieldUpsert[time.Time]
}

func newSkillUpsert(up *ConflictUpdate) *SkillUpsert {
	return &SkillUpsert{
		Id:          fieldUpsert[string]{column: "id", update: up},
		Name:        fieldUpsert[string]{column: "name", update: up},
		Description: fieldUpsert[string]{column: "description", update: up},
		Url:         fieldUpsert[string]{column: "url", update: up},
		Content:     fieldUpsert[string]{column: "content", update: up},
		Keywords:    fieldUpsert[string]{column: "keywords", update: up},
		CreatedAt:   fieldUpsert[time.Time]{column: "createdAt", update: up},
		UpdatedAt:   fieldUpsert[time.Time]{column: "updatedAt", update: up},
	}
}

type SkillUpdateBuilder struct {
	*UpdateBuilder[Skill, SkillSelect, SkillOmit]
}

type SkillUpdateManyBuilder struct {
	*UpdateManyBuilder[Skill]
}

type SkillUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Skill, SkillSelect, SkillOmit]
}

func (b *SkillUpdateBuilder) SetId(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetId(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetId(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetName(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetName(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetName(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetDescription(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetDescription(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetDescription(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetUrl(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "url", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetUrl(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "url", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetUrl(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "url", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetContent(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetContent(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetContent(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetKeywords(v string) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "keywords", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetKeywords(v string) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "keywords", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetKeywords(v string) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "keywords", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetCreatedAt(v time.Time) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetCreatedAt(v time.Time) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetCreatedAt(v time.Time) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *SkillUpdateBuilder) SetUpdatedAt(v time.Time) *SkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *SkillUpdateManyBuilder) SetUpdatedAt(v time.Time) *SkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) SetUpdatedAt(v time.Time) *SkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *SkillUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Skill]) *SkillUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Skill]) *SkillUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Skill]) *SkillUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillDelegate) Update(where UniquePredicate[Skill], additional ...PredicateOf[Skill]) *SkillUpdateBuilder {
	return &SkillUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Skill, SkillSelect, SkillOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *SkillDelegate) UpdateMany(preds ...PredicateOf[Skill]) *SkillUpdateManyBuilder {
	return &SkillUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Skill]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *SkillDelegate) UpdateManyAndReturn(preds ...PredicateOf[Skill]) *SkillUpdateManyAndReturnBuilder {
	return &SkillUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Skill, SkillSelect, SkillOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *SkillDelegate) buildUpdateSQL(preds []PredicateOf[Skill], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Skill")
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

func (d *SkillDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Skill], additional []PredicateOf[Skill], assignments []FieldAssignment, selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
	allWhere := make([]PredicateOf[Skill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillUpdateArgs) (*Skill, error) {
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
			curr = func(c context.Context, a *SkillUpdateArgs) (*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Skill], cols []string, vals []any, selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
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
		var res *Skill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Skill.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Skill.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillCols(selects, omits, skillPKCols...)
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
		return nil, &NotFoundError{Model: "Skill"}
	}

	var res Skill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Skill{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *SkillDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Skill], cols []string, vals []any) (int64, error) {
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

func (d *SkillDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Skill], cols []string, vals []any, selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Skill"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *SkillDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Skill], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToSkillUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &SkillUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *SkillUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *SkillDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Skill], assignments []FieldAssignment, selects *SkillSelect, omits *SkillOmit) ([]*Skill, error) {
	input, err := assignmentsToSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillUpdateManyAndReturnArgs) ([]*Skill, error) {
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
			curr = func(c context.Context, a *SkillUpdateManyAndReturnArgs) ([]*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Skill], cols []string, vals []any, selects *SkillSelect, omits *SkillOmit) ([]*Skill, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Skill]{Where: preds}, selects, omits)
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
		var res []*Skill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Skill.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Skill.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillCols(selects, omits, skillPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanSkillRows(rows, returningCols)
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

func (d *SkillDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Skill], cols []string, vals []any, selects *SkillSelect, omits *SkillOmit) ([]*Skill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Skill{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Skill]{Where: preds}, selects, omits)
}
func (d *SkillDelegate) FindUnique(where UniquePredicate[Skill], additional ...PredicateOf[Skill]) *FindUniqueBuilder[Skill, SkillSelect, SkillOmit] {
	return &FindUniqueBuilder[Skill, SkillSelect, SkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *SkillDelegate) FindFirst(preds ...PredicateOf[Skill]) *FindFirstBuilder[Skill, SkillSelect, SkillOmit] {
	return &FindFirstBuilder[Skill, SkillSelect, SkillOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *SkillDelegate) FindMany(preds ...PredicateOf[Skill]) *FindManyBuilder[Skill, SkillSelect, SkillOmit] {
	return &FindManyBuilder[Skill, SkillSelect, SkillOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *SkillDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Skill], additional []PredicateOf[Skill], selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
	allWhere := make([]PredicateOf[Skill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillFindUniqueArgs) (*Skill, error) {
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
			curr = func(c context.Context, a *SkillFindUniqueArgs) (*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Skill],
	selects *SkillSelect,
	omits *SkillOmit,
) (*Skill, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillFindFirstArgs) (*Skill, error) {
		return d.runFindFirst(c, QueryParams[Skill]{
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
			curr = func(c context.Context, a *SkillFindFirstArgs) (*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Skill],
	selects *SkillSelect,
	omits *SkillOmit,
) ([]*Skill, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillFindManyArgs) ([]*Skill, error) {
		return d.runFindMany(c, QueryParams[Skill]{
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
			curr = func(c context.Context, a *SkillFindManyArgs) ([]*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Skill], selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
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
	returningCols := selectSkillCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Skill{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *SkillDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Skill],
	selects *SkillSelect,
	omits *SkillOmit,
) (*Skill, error) {
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

func (d *SkillDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Skill],
	selects *SkillSelect,
	omits *SkillOmit,
) ([]*Skill, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, skillPKCols, skillUniqueCols, "Skill", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, skillPKCols, skillUniqueCols, isCursorQuery, params.Take)
	returningCols := selectSkillCols(selects, omits)

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

func (d *SkillDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Skill, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Skill", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Skill
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *SkillDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Skill, error) {
	query := buildSelectSQL(d.client, "Skill", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Skill, 0)
	for rows.Next() {
		var res Skill
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
func (d *SkillDelegate) DeleteMany(preds ...PredicateOf[Skill]) *DeleteManyBuilder[Skill] {
	return &DeleteManyBuilder[Skill]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *SkillDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Skill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &SkillDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *SkillDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Skill]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Skill")
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

func (d *SkillDelegate) Delete(where UniquePredicate[Skill], additional ...PredicateOf[Skill]) *DeleteBuilder[Skill, SkillSelect, SkillOmit] {
	return &DeleteBuilder[Skill, SkillSelect, SkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *SkillDelegate) executeDelete(ctx context.Context, where UniquePredicate[Skill], additional []PredicateOf[Skill], selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
	allWhere := make([]PredicateOf[Skill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillSelect()
	}

	args := &SkillDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillDeleteArgs) (*Skill, error) {
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
			curr = func(c context.Context, a *SkillDeleteArgs) (*Skill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runDelete(ctx context.Context, where []PredicateOf[Skill], selects *SkillSelect, omits *SkillOmit) (*Skill, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectSkillCols(selects, omits, skillPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Skill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Skill.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Skill"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Skill")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Skill]
			pkPreds = append(pkPreds, Predicate[Skill]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
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
	d.client.dialect.WriteQuotedIdent(&sb, "Skill")

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
		return nil, &NotFoundError{Model: "Skill"}
	}

	var row Skill
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *SkillDelegate) Count(preds ...PredicateOf[Skill]) *CountBuilder[Skill] {
	return &CountBuilder[Skill]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *SkillDelegate) executeCount(ctx context.Context, params QueryParams[Skill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &SkillCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *SkillCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Skill]{
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
			curr = func(c context.Context, a *SkillCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDelegate) runCount(ctx context.Context, params QueryParams[Skill]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Skill")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Skill")
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
func (d *SkillDelegate) loadRelations(ctx context.Context, records []*Skill, selects *SkillSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Dependencies != nil {
		relationSelects, relationOmits, relationParams := selects.Dependencies.GetRelationParams()
		returningCols := selectSkillDependencyCols(relationSelects, relationOmits, "dependencyId")
		// Inverse holds the FK: SkillDependency.dependencyId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"SkillDependency",
			"dependencyId",
			returningCols,
			scanInto(returningCols, (*SkillDependency).ScanFields),
			directKey(func(c *SkillDependency) string { return c.DependencyId }),
			appendMany(func(p *Skill) *[]*SkillDependency { return &p.Dependencies }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading dependencies: %w", err)
		}
		if err := d.client.SkillDependency.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Dependents != nil {
		relationSelects, relationOmits, relationParams := selects.Dependents.GetRelationParams()
		returningCols := selectSkillDependencyCols(relationSelects, relationOmits, "dependentId")
		// Inverse holds the FK: SkillDependency.dependentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"SkillDependency",
			"dependentId",
			returningCols,
			scanInto(returningCols, (*SkillDependency).ScanFields),
			directKey(func(c *SkillDependency) string { return c.DependentId }),
			appendMany(func(p *Skill) *[]*SkillDependency { return &p.Dependents }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading dependents: %w", err)
		}
		if err := d.client.SkillDependency.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Projects != nil {
		relationSelects, relationOmits, relationParams := selects.Projects.GetRelationParams()
		returningCols := selectProjectSkillCols(relationSelects, relationOmits, "skillId")
		// Inverse holds the FK: ProjectSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"ProjectSkill",
			"skillId",
			returningCols,
			scanInto(returningCols, (*ProjectSkill).ScanFields),
			directKey(func(c *ProjectSkill) string { return c.SkillId }),
			appendMany(func(p *Skill) *[]*ProjectSkill { return &p.Projects }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading projects: %w", err)
		}
		if err := d.client.ProjectSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.References != nil {
		relationSelects, relationOmits, relationParams := selects.References.GetRelationParams()
		returningCols := selectSkillReferenceCols(relationSelects, relationOmits, "skillId")
		// Inverse holds the FK: SkillReference.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"SkillReference",
			"skillId",
			returningCols,
			scanInto(returningCols, (*SkillReference).ScanFields),
			directKey(func(c *SkillReference) string { return c.SkillId }),
			appendMany(func(p *Skill) *[]*SkillReference { return &p.References }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading references: %w", err)
		}
		if err := d.client.SkillReference.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Trashes != nil {
		relationSelects, relationOmits, relationParams := selects.Trashes.GetRelationParams()
		returningCols := selectTrashSkillCols(relationSelects, relationOmits, "skillId")
		// Inverse holds the FK: TrashSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"TrashSkill",
			"skillId",
			returningCols,
			scanInto(returningCols, (*TrashSkill).ScanFields),
			directKey(func(c *TrashSkill) string { return c.SkillId }),
			appendMany(func(p *Skill) *[]*TrashSkill { return &p.Trashes }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading trashes: %w", err)
		}
		if err := d.client.TrashSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Agents != nil {
		relationSelects, relationOmits, relationParams := selects.Agents.GetRelationParams()
		returningCols := selectAgentSkillCols(relationSelects, relationOmits, "skillId")
		// Inverse holds the FK: AgentSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"AgentSkill",
			"skillId",
			returningCols,
			scanInto(returningCols, (*AgentSkill).ScanFields),
			directKey(func(c *AgentSkill) string { return c.SkillId }),
			appendMany(func(p *Skill) *[]*AgentSkill { return &p.Agents }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agents: %w", err)
		}
		if err := d.client.AgentSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Apps != nil {
		relationSelects, relationOmits, relationParams := selects.Apps.GetRelationParams()
		returningCols := selectAppSkillCols(relationSelects, relationOmits, "skillId")
		// Inverse holds the FK: AppSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"AppSkill",
			"skillId",
			returningCols,
			scanInto(returningCols, (*AppSkill).ScanFields),
			directKey(func(c *AppSkill) string { return c.SkillId }),
			appendMany(func(p *Skill) *[]*AppSkill { return &p.Apps }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading apps: %w", err)
		}
		if err := d.client.AppSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Users != nil {
		// Many-to-many through join table _UserSkills
		relationSelects, relationOmits, relationParams := selects.Users.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"User",
			"id",
			"_UserSkills",
			"skillId",
			"userId",
			returningCols,
			(*User).ScanFields,
			appendMany(func(p *Skill) *[]*User { return &p.Users }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading users: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Userblacklist != nil {
		// Many-to-many through join table _UserSkillBlacklist
		relationSelects, relationOmits, relationParams := selects.Userblacklist.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *Skill) string { return p.Id }),
			"User",
			"id",
			"_UserSkillBlacklist",
			"skillId",
			"userId",
			returningCols,
			(*User).ScanFields,
			appendMany(func(p *Skill) *[]*User { return &p.Userblacklist }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading userblacklist: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
