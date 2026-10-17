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

// Agent represents the database model
type Agent struct {
	Id           string        `db:"id" json:"id"`
	Name         string        `db:"name" json:"name"`
	Description  string        `db:"description" json:"description"`
	Content      string        `db:"content" json:"content"`
	AllowedTools string        `db:"allowed-tools" json:"allowedTools"`
	Metadata     string        `db:"metadata" json:"metadata"`
	CreatedAt    time.Time     `db:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `db:"updatedAt" json:"updatedAt"`
	Trashes      []*TrashAgent `json:"trashes,omitempty"`
	Skills       []*AgentSkill `json:"skills,omitempty"`
	Projects     []*Project    `json:"projects,omitempty"`
}

// AgentCreate contains model input fields for Agent creation operations.
//
// Fields for Agent:
//
//	id           string    default: cuid()
//	name         string    required
//	description  string    required
//	content      string    default:
//	allowedTools string    required
//	metadata     string    required
//	createdAt    time.Time default: now()
//	updatedAt    time.Time required
type AgentCreate struct {
	Id           *string    `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Content      *string    `json:"content"`
	AllowedTools string     `json:"allowedTools"`
	Metadata     string     `json:"metadata"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// colMask returns a bit mask of columns that are set
func (s *AgentCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	if s.Content != nil {
		mask |= 1 << 3
	}
	mask |= 1 << 4
	mask |= 1 << 5
	mask |= 1 << 6
	mask |= 1 << 7
	return mask
}

// AgentUpdate contains model input fields for Agent update operations.
type AgentUpdate struct {
	Id           *string    `json:"id"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	Content      *string    `json:"content"`
	AllowedTools *string    `json:"allowedTools"`
	Metadata     *string    `json:"metadata"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

func (u *AgentUpdate) ToColsVals() ([]string, []any) {
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
	if u.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, u.Content)
	}
	if u.AllowedTools != nil {
		cols = append(cols, "allowed-tools")
		vals = append(vals, u.AllowedTools)
	}
	if u.Metadata != nil {
		cols = append(cols, "metadata")
		vals = append(vals, u.Metadata)
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

func assignmentsToAgentUpdate(assignments []FieldAssignment) (AgentUpdate, error) {
	var input AgentUpdate
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
		case "content":
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Content = v
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "allowed-tools":
			if v, ok := a.Val.(string); ok {
				input.AllowedTools = &v
				errs.ValidateString("allowedTools", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.AllowedTools = v
			} else {
				errs.Add("allowedTools", a.Val, "type", "field allowedTools must be of type string")
			}
		case "metadata":
			if v, ok := a.Val.(string); ok {
				input.Metadata = &v
				errs.ValidateString("metadata", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Metadata = v
			} else {
				errs.Add("metadata", a.Val, "type", "field metadata must be of type string")
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

// AgentSelect specifies which scalar and relation fields to select for Agent.
//
// Selectable fields:
//
//	-- Scalars --
//	id           (bool)
//	name         (bool)
//	description  (bool)
//	content      (bool)
//	allowedTools (bool)
//	metadata     (bool)
//	createdAt    (bool)
//	updatedAt    (bool)
//	-- Relations --
//	trashes      ([]TrashAgent)
//	skills       ([]AgentSkill)
//	projects     ([]Project)
type AgentSelect struct {
	Id           bool                  `json:"id"`
	Name         bool                  `json:"name"`
	Description  bool                  `json:"description"`
	Content      bool                  `json:"content"`
	AllowedTools bool                  `json:"allowedTools"`
	Metadata     bool                  `json:"metadata"`
	CreatedAt    bool                  `json:"createdAt"`
	UpdatedAt    bool                  `json:"updatedAt"`
	Trashes      TrashAgentSelectQuery `json:"trashes,omitempty"`
	Skills       AgentSkillSelectQuery `json:"skills,omitempty"`
	Projects     ProjectSelectQuery    `json:"projects,omitempty"`
}

var fullAgentSelectVal = &AgentSelect{
	Id:           true,
	Name:         true,
	Description:  true,
	Content:      true,
	AllowedTools: true,
	Metadata:     true,
	CreatedAt:    true,
	UpdatedAt:    true,
}

func fullAgentSelect() *AgentSelect {
	return fullAgentSelectVal
}

func (s *AgentSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.Name || s.Description || s.Content || s.AllowedTools || s.Metadata || s.CreatedAt || s.UpdatedAt
}

func (s *AgentSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type AgentOmit struct {
	Id           bool `json:"id"`
	Name         bool `json:"name"`
	Description  bool `json:"description"`
	Content      bool `json:"content"`
	AllowedTools bool `json:"allowedTools"`
	Metadata     bool `json:"metadata"`
	CreatedAt    bool `json:"createdAt"`
	UpdatedAt    bool `json:"updatedAt"`
}

type AgentSelectQuery interface {
	GetRelationParams() (*AgentSelect, *AgentOmit, QueryParams[Agent])
}

func (s *AgentSelect) GetRelationParams() (*AgentSelect, *AgentOmit, QueryParams[Agent]) {
	return s, nil, QueryParams[Agent]{}
}

type AgentQueryBuilder struct {
	selects *AgentSelect
	omits   *AgentOmit
	where   []PredicateOf[Agent]
	take    *int
	skip    *int
	orderBy []OrderBy[Agent]
	cursor  UniquePredicate[Agent]
}

func (b *AgentQueryBuilder) Where(preds ...PredicateOf[Agent]) *AgentQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AgentQueryBuilder) Take(limit int) *AgentQueryBuilder {
	b.take = &limit
	return b
}

func (b *AgentQueryBuilder) Skip(offset int) *AgentQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AgentQueryBuilder) OrderBy(orders ...OrderBy[Agent]) *AgentQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AgentQueryBuilder) Cursor(where UniquePredicate[Agent]) *AgentQueryBuilder {
	b.cursor = where
	return b
}

func (b *AgentQueryBuilder) Select(s AgentSelect) *AgentQueryBuilder {
	b.selects = &s
	return b
}

func (b *AgentQueryBuilder) Omit(o AgentOmit) *AgentQueryBuilder {
	b.omits = &o
	return b
}

func (b *AgentQueryBuilder) GetRelationParams() (*AgentSelect, *AgentOmit, QueryParams[Agent]) {
	if b == nil {
		return nil, nil, QueryParams[Agent]{}
	}
	return b.selects, b.omits, QueryParams[Agent]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// AgentCreateArgs is the input argument passed to Agent Create extension hooks.
//
// Fields for Agent:
//
//	id           string    default: cuid()
//	name         string    required
//	description  string    required
//	content      string    default:
//	allowedTools string    required
//	metadata     string    required
//	createdAt    time.Time default: now()
//	updatedAt    time.Time required
//
// Relations for Agent:
//
//	trashes      ([]TrashAgent)
//	skills       ([]AgentSkill)
//	projects     ([]Project)
type AgentCreateArgs struct {
	// Data contains the model fields to insert.
	Data *AgentCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *AgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// AgentCreateManyArgs is the input argument passed to Agent CreateMany extension hooks.
//
// Fields for Agent:
//
//	id           string    default: cuid()
//	name         string    required
//	description  string    required
//	content      string    default:
//	allowedTools string    required
//	metadata     string    required
//	createdAt    time.Time default: now()
//	updatedAt    time.Time required
type AgentCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AgentCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AgentCreateManyArgs) AppendData(builders ...*AgentCreateBuilder) *AgentCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AgentCreateManyAndReturnArgs is the input argument passed to Agent CreateManyAndReturn extension hooks.
//
// Fields for Agent:
//
//	id           string    default: cuid()
//	name         string    required
//	description  string    required
//	content      string    default:
//	allowedTools string    required
//	metadata     string    required
//	createdAt    time.Time default: now()
//	updatedAt    time.Time required
//
// Relations for Agent:
//
//	trashes      ([]TrashAgent)
//	skills       ([]AgentSkill)
//	projects     ([]Project)
type AgentCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AgentCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *AgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AgentCreateManyAndReturnArgs) AppendData(builders ...*AgentCreateBuilder) *AgentCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AgentFindUniqueArgs is the input argument passed to Agent FindUnique extension hooks.
type AgentFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Agent]
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSelect
}

func (a *AgentFindUniqueArgs) SetWhere(unique UniquePredicate[Agent], additional ...PredicateOf[Agent]) *AgentFindUniqueArgs {
	a.Where = make([]PredicateOf[Agent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentFindFirstArgs is the input argument passed to Agent FindFirst extension hooks.
type AgentFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Agent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Agent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSelect
}

func (a *AgentFindFirstArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentFindFirstArgs {
	a.Where = preds
	return a
}

func (a *AgentFindFirstArgs) SetOrderBy(orders ...OrderBy[Agent]) *AgentFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *AgentFindFirstArgs) SetCursor(cursor UniquePredicate[Agent]) *AgentFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *AgentFindFirstArgs) SetSkip(n int) *AgentFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *AgentFindFirstArgs) SetTake(n int) *AgentFindFirstArgs {
	a.Take = &n
	return a
}

// AgentFindManyArgs is the input argument passed to Agent FindMany extension hooks.
type AgentFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Agent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Agent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSelect
}

func (a *AgentFindManyArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentFindManyArgs {
	a.Where = preds
	return a
}

func (a *AgentFindManyArgs) SetOrderBy(orders ...OrderBy[Agent]) *AgentFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *AgentFindManyArgs) SetCursor(cursor UniquePredicate[Agent]) *AgentFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *AgentFindManyArgs) SetSkip(n int) *AgentFindManyArgs {
	a.Skip = &n
	return a
}

func (a *AgentFindManyArgs) SetTake(n int) *AgentFindManyArgs {
	a.Take = &n
	return a
}

// AgentCountArgs is the input argument passed to Agent Count extension hooks.
type AgentCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *AgentCountArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentCountArgs {
	a.Where = preds
	return a
}

func (a *AgentCountArgs) SetSkip(n int) *AgentCountArgs {
	a.Skip = &n
	return a
}

func (a *AgentCountArgs) SetTake(n int) *AgentCountArgs {
	a.Take = &n
	return a
}

// AgentDeleteArgs is the input argument passed to Agent Delete extension hooks.
type AgentDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Agent]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *AgentSelect
}

func (a *AgentDeleteArgs) SetWhere(unique UniquePredicate[Agent], additional ...PredicateOf[Agent]) *AgentDeleteArgs {
	a.Where = make([]PredicateOf[Agent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentDeleteManyArgs is the input argument passed to Agent DeleteMany extension hooks.
type AgentDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
}

func (a *AgentDeleteManyArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentDeleteManyArgs {
	a.Where = preds
	return a
}

// AgentUpdateArgs is the input argument passed to Agent Update extension hooks.
type AgentUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Agent]
	// Data contains the model fields to update.
	Data *AgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AgentSelect
}

func (a *AgentUpdateArgs) SetWhere(unique UniquePredicate[Agent], additional ...PredicateOf[Agent]) *AgentUpdateArgs {
	a.Where = make([]PredicateOf[Agent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentUpdateManyArgs is the input argument passed to Agent UpdateMany extension hooks.
type AgentUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
	// Data contains the model fields to update.
	Data *AgentUpdate
}

func (a *AgentUpdateManyArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentUpdateManyArgs {
	a.Where = preds
	return a
}

// AgentUpdateManyAndReturnArgs is the input argument passed to Agent UpdateManyAndReturn extension hooks.
type AgentUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Agent]
	// Data contains the model fields to update.
	Data *AgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AgentSelect
}

func (a *AgentUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Agent]) *AgentUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type AgentCreateQuery = func(ctx context.Context, args *AgentCreateArgs) (*Agent, error)
type AgentCreateManyQuery = func(ctx context.Context, args *AgentCreateManyArgs) (int64, error)
type AgentCreateManyAndReturnQuery = func(ctx context.Context, args *AgentCreateManyAndReturnArgs) ([]*Agent, error)
type AgentFindUniqueQuery = func(ctx context.Context, args *AgentFindUniqueArgs) (*Agent, error)
type AgentFindFirstQuery = func(ctx context.Context, args *AgentFindFirstArgs) (*Agent, error)
type AgentFindManyQuery = func(ctx context.Context, args *AgentFindManyArgs) ([]*Agent, error)
type AgentDeleteQuery = func(ctx context.Context, args *AgentDeleteArgs) (*Agent, error)
type AgentDeleteManyQuery = func(ctx context.Context, args *AgentDeleteManyArgs) (int64, error)
type AgentCountQuery = func(ctx context.Context, args *AgentCountArgs) (int64, error)
type AgentUpdateQuery = func(ctx context.Context, args *AgentUpdateArgs) (*Agent, error)
type AgentUpdateManyQuery = func(ctx context.Context, args *AgentUpdateManyArgs) (int64, error)
type AgentUpdateManyAndReturnQuery = func(ctx context.Context, args *AgentUpdateManyAndReturnArgs) ([]*Agent, error)

type AgentExtension struct {
	Create              func(ctx context.Context, args *AgentCreateArgs, next AgentCreateQuery) (*Agent, error)
	CreateMany          func(ctx context.Context, args *AgentCreateManyArgs, next AgentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *AgentCreateManyAndReturnArgs, next AgentCreateManyAndReturnQuery) ([]*Agent, error)
	FindUnique          func(ctx context.Context, args *AgentFindUniqueArgs, next AgentFindUniqueQuery) (*Agent, error)
	FindFirst           func(ctx context.Context, args *AgentFindFirstArgs, next AgentFindFirstQuery) (*Agent, error)
	FindMany            func(ctx context.Context, args *AgentFindManyArgs, next AgentFindManyQuery) ([]*Agent, error)
	Delete              func(ctx context.Context, args *AgentDeleteArgs, next AgentDeleteQuery) (*Agent, error)
	DeleteMany          func(ctx context.Context, args *AgentDeleteManyArgs, next AgentDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *AgentCountArgs, next AgentCountQuery) (int64, error)
	Update              func(ctx context.Context, args *AgentUpdateArgs, next AgentUpdateQuery) (*Agent, error)
	UpdateMany          func(ctx context.Context, args *AgentUpdateManyArgs, next AgentUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *AgentUpdateManyAndReturnArgs, next AgentUpdateManyAndReturnQuery) ([]*Agent, error)
}

type AgentDelegate struct {
	client     *Queries
	extensions []AgentExtension
}

func (d *AgentDelegate) Use(exts ...AgentExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Agent) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "name":
			targets[i] = &m.Name
		case "description":
			targets[i] = &m.Description
		case "content":
			targets[i] = &m.Content
		case "allowed-tools":
			targets[i] = &m.AllowedTools
		case "metadata":
			targets[i] = &m.Metadata
		case "createdAt":
			targets[i] = &m.CreatedAt
		case "updatedAt":
			targets[i] = &m.UpdatedAt
		}
	}
	return targets
}

var agentDefaultCols = []string{
	"id",
	"name",
	"description",
	"content",
	"allowed-tools",
	"metadata",
	"createdAt",
	"updatedAt",
}

var agentPKCols = []string{
	"id",
}

var agentUniqueCols = []string{
	"id",
	"name",
}

func selectAgentCols(selects *AgentSelect, omits *AgentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return agentDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Name || selects.Description || selects.Content || selects.AllowedTools || selects.Metadata || selects.CreatedAt || selects.UpdatedAt || selects.Trashes != nil || selects.Skills != nil || selects.Projects != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
		{"description", selects != nil && selects.Description, omits != nil && omits.Description, false},
		{"content", selects != nil && selects.Content, omits != nil && omits.Content, false},
		{"allowed-tools", selects != nil && selects.AllowedTools, omits != nil && omits.AllowedTools, false},
		{"metadata", selects != nil && selects.Metadata, omits != nil && omits.Metadata, false},
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

func (s *AgentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Trashes != nil || s.Skills != nil || s.Projects != nil
}

type AgentUpsertBuilder struct {
	*CreateBuilder[Agent, AgentSelect, AgentOmit]
}

func (b *AgentUpsertBuilder) Select(s AgentSelect) *AgentUpsertBuilder {
	b.selects = &s
	return b
}

func (b *AgentUpsertBuilder) Omit(o AgentOmit) *AgentUpsertBuilder {
	b.omits = &o
	return b
}

type AgentCreateBuilder struct {
	*CreateBuilder[Agent, AgentSelect, AgentOmit]
}

func (b *AgentCreateBuilder) Select(s AgentSelect) *AgentCreateBuilder {
	b.selects = &s
	return b
}

func (b *AgentCreateBuilder) Omit(o AgentOmit) *AgentCreateBuilder {
	b.omits = &o
	return b
}

func (b *AgentCreateBuilder) OnConflict(target UniqueConstraintTarget) *AgentConflictBuilder[AgentUpsertBuilder] {
	upsertBuilder := &AgentUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &AgentConflictBuilder[AgentUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AgentCreateBuilder) SetId(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetName(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetDescription(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetContent(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetAllowedTools(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "allowed-tools", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetMetadata(v string) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "metadata", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetCreatedAt(v time.Time) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *AgentCreateBuilder) SetUpdatedAt(v time.Time) *AgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *AgentCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Agent]) *AgentCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AgentDelegate) Create() *AgentCreateBuilder {
	return &AgentCreateBuilder{
		CreateBuilder: &CreateBuilder[Agent, AgentSelect, AgentOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedAgentId           uint64 = 1 << 0
	providedAgentName         uint64 = 1 << 1
	providedAgentDescription  uint64 = 1 << 2
	providedAgentContent      uint64 = 1 << 3
	providedAgentAllowedTools uint64 = 1 << 4
	providedAgentMetadata     uint64 = 1 << 5
	providedAgentCreatedAt    uint64 = 1 << 6
	providedAgentUpdatedAt    uint64 = 1 << 7
)

func assignmentsToAgentCreate(assignments []FieldAssignment) (AgentCreate, error) {
	var input AgentCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedAgentId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				errs.ValidateString("id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "name":
			provided |= providedAgentName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "description":
			provided |= providedAgentDescription
			if v, ok := a.Val.(string); ok {
				input.Description = v
				errs.ValidateString("description", v, true, 0, false, false)
			} else {
				errs.Add("description", a.Val, "type", "field description must be of type string")
			}
		case "content":
			provided |= providedAgentContent
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "allowed-tools":
			provided |= providedAgentAllowedTools
			if v, ok := a.Val.(string); ok {
				input.AllowedTools = v
				errs.ValidateString("allowedTools", v, true, 0, false, false)
			} else {
				errs.Add("allowedTools", a.Val, "type", "field allowedTools must be of type string")
			}
		case "metadata":
			provided |= providedAgentMetadata
			if v, ok := a.Val.(string); ok {
				input.Metadata = v
				errs.ValidateString("metadata", v, true, 0, false, false)
			} else {
				errs.Add("metadata", a.Val, "type", "field metadata must be of type string")
			}
		case "createdAt":
			provided |= providedAgentCreatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		case "updatedAt":
			provided |= providedAgentUpdatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.UpdatedAt = v
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type time.Time")
			}
		}
	}
	if provided&providedAgentName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}
	if provided&providedAgentDescription == 0 {
		errs.Add("description", "", "required", "field Description is required")
	}
	if provided&providedAgentAllowedTools == 0 {
		errs.Add("allowed-tools", "", "required", "field AllowedTools is required")
	}
	if provided&providedAgentMetadata == 0 {
		errs.Add("metadata", "", "required", "field Metadata is required")
	}
	if provided&providedAgentUpdatedAt == 0 {
		input.UpdatedAt = time.Now().Truncate(time.Microsecond)
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AgentCreate) ToColsVals() (cols []string, vals []any) {
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
	if s.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, *s.Content)
	}
	cols = append(cols, "allowed-tools")
	vals = append(vals, s.AllowedTools)
	cols = append(cols, "metadata")
	vals = append(vals, s.Metadata)
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

func partitionAgentInputs(dialect Dialect, inputs []*AgentCreate) [][]*AgentCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AgentCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AgentCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AgentCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AgentCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AgentCreate{inputs}
}

func (d *AgentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AgentSelect, omits *AgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Agent, error) {
	input, err := assignmentsToAgentCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAgentCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentCreateArgs) (*Agent, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectAgentCols(a.Select, omits)
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
			curr = func(c context.Context, a *AgentCreateArgs) (*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type AgentCreateManyBuilder struct {
	*CreateManyBuilder[Agent]
}

func (b *AgentCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AgentConflictBuilder[AgentCreateManyBuilder] {
	return &AgentConflictBuilder[AgentCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AgentCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Agent, AgentSelect, AgentOmit]
}

func (b *AgentCreateManyAndReturnBuilder) Select(s AgentSelect) *AgentCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *AgentCreateManyAndReturnBuilder) Omit(o AgentOmit) *AgentCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *AgentCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AgentConflictBuilder[AgentCreateManyAndReturnBuilder] {
	return &AgentConflictBuilder[AgentCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToAgentRecordInputs(builders []*AgentCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *AgentDelegate) CreateMany(builders ...*AgentCreateBuilder) *AgentCreateManyBuilder {
	return &AgentCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Agent]{
			records:  createBuildersToAgentRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AgentDelegate) CreateManyAndReturn(builders ...*AgentCreateBuilder) *AgentCreateManyAndReturnBuilder {
	return &AgentCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Agent, AgentSelect, AgentOmit]{
			records:  createBuildersToAgentRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToAgentCreateInputs(records []RecordInput) ([]*AgentCreate, error) {
	structs := make([]AgentCreate, len(records))
	inputs := make([]*AgentCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToAgentCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *AgentDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToAgentCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &AgentCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AgentSelect, omits *AgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Agent, error) {
	inputs, err := recordsToAgentCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentCreateManyAndReturnArgs) ([]*Agent, error) {
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
			curr = func(c context.Context, a *AgentCreateManyAndReturnArgs) ([]*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *AgentSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Agent, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Agent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Agent.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Agent.loadRelations(ctx, []*Agent{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Agent", cols, returningCols, agentPKCols, conflictTarget, conflictAction, len(vals))
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

		var res Agent
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, agentPKCols)
}

func (d *AgentDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Agent, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Agent") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Agent")
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

	var res Agent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *AgentDelegate) buildBulkInsertSQL(q *Queries, batch []*AgentCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 8)
	for i, c := range agentDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Agent")
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
			case "content":
				if input.Content != nil {
					vals = append(vals, *input.Content)
				} else {
					writeDefault = true
				}
			case "allowed-tools":
				vals = append(vals, input.AllowedTools)
			case "metadata":
				vals = append(vals, input.Metadata)
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

func applyAgentConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanAgentRows(rows *sql.Rows, returningCols []string) ([]*Agent, error) {
	var records []*Agent
	for rows.Next() {
		var res Agent
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

func (d *AgentDelegate) runCreateMany(ctx context.Context, inputs []*AgentCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAgentInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentConflictClause(d.client.dialect, queryStr, vals, cols, agentPKCols, conflictTarget, conflictAction)

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

func (d *AgentDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AgentCreate,
	selects *AgentSelect,
	omits *AgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Agent, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAgentInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Agent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Agent.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Agent.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectAgentCols(selects, omits, agentPKCols...)
	recordsOut := make([]*Agent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentConflictClause(d.client.dialect, queryStr, vals, cols, agentPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanAgentRows(rows, returningCols)
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

func (d *AgentDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*AgentCreate,
	selects *AgentSelect,
	omits *AgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Agent, error) {
	batches := partitionAgentInputs(d.client.dialect, inputs)
	returningCols := selectAgentCols(selects, omits, agentPKCols...)
	recordsOut := make([]*Agent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentConflictClause(d.client.dialect, queryStr, vals, cols, agentPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Agent") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Agent")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, agentPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, agentPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanAgentRows(rows, returningCols)
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

type AgentConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AgentConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AgentConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AgentConflictBuilder[B]) Update(fn func(u *AgentUpsert)) *B {
	cb.setAction(*AgentConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// AgentConflictUpdate creates a custom ConflictAction for Agent upsert conflicts.
func AgentConflictUpdate(fn func(u *AgentUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newAgentUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type AgentUpsert struct {
	Id           fieldUpsert[string]
	Name         fieldUpsert[string]
	Description  fieldUpsert[string]
	Content      fieldUpsert[string]
	AllowedTools fieldUpsert[string]
	Metadata     fieldUpsert[string]
	CreatedAt    fieldUpsert[time.Time]
	UpdatedAt    fieldUpsert[time.Time]
}

func newAgentUpsert(up *ConflictUpdate) *AgentUpsert {
	return &AgentUpsert{
		Id:           fieldUpsert[string]{column: "id", update: up},
		Name:         fieldUpsert[string]{column: "name", update: up},
		Description:  fieldUpsert[string]{column: "description", update: up},
		Content:      fieldUpsert[string]{column: "content", update: up},
		AllowedTools: fieldUpsert[string]{column: "allowed-tools", update: up},
		Metadata:     fieldUpsert[string]{column: "metadata", update: up},
		CreatedAt:    fieldUpsert[time.Time]{column: "createdAt", update: up},
		UpdatedAt:    fieldUpsert[time.Time]{column: "updatedAt", update: up},
	}
}

type AgentUpdateBuilder struct {
	*UpdateBuilder[Agent, AgentSelect, AgentOmit]
}

type AgentUpdateManyBuilder struct {
	*UpdateManyBuilder[Agent]
}

type AgentUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Agent, AgentSelect, AgentOmit]
}

func (b *AgentUpdateBuilder) SetId(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetId(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetId(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetName(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetName(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetName(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetDescription(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetDescription(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetDescription(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetContent(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetContent(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetContent(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetAllowedTools(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "allowed-tools", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetAllowedTools(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "allowed-tools", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetAllowedTools(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "allowed-tools", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetMetadata(v string) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "metadata", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetMetadata(v string) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "metadata", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetMetadata(v string) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "metadata", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetCreatedAt(v time.Time) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetCreatedAt(v time.Time) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetCreatedAt(v time.Time) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *AgentUpdateBuilder) SetUpdatedAt(v time.Time) *AgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *AgentUpdateManyBuilder) SetUpdatedAt(v time.Time) *AgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) SetUpdatedAt(v time.Time) *AgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}

func (b *AgentUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Agent]) *AgentUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AgentUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Agent]) *AgentUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AgentUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Agent]) *AgentUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AgentDelegate) Update(where UniquePredicate[Agent], additional ...PredicateOf[Agent]) *AgentUpdateBuilder {
	return &AgentUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Agent, AgentSelect, AgentOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *AgentDelegate) UpdateMany(preds ...PredicateOf[Agent]) *AgentUpdateManyBuilder {
	return &AgentUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Agent]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *AgentDelegate) UpdateManyAndReturn(preds ...PredicateOf[Agent]) *AgentUpdateManyAndReturnBuilder {
	return &AgentUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Agent, AgentSelect, AgentOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *AgentDelegate) buildUpdateSQL(preds []PredicateOf[Agent], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Agent")
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

func (d *AgentDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Agent], additional []PredicateOf[Agent], assignments []FieldAssignment, selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
	allWhere := make([]PredicateOf[Agent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentUpdateArgs) (*Agent, error) {
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
			curr = func(c context.Context, a *AgentUpdateArgs) (*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Agent], cols []string, vals []any, selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
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
		var res *Agent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Agent.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Agent.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAgentCols(selects, omits, agentPKCols...)
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
		return nil, &NotFoundError{Model: "Agent"}
	}

	var res Agent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Agent{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *AgentDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Agent], cols []string, vals []any) (int64, error) {
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

func (d *AgentDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Agent], cols []string, vals []any, selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Agent"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *AgentDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Agent], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToAgentUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &AgentUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *AgentUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *AgentDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Agent], assignments []FieldAssignment, selects *AgentSelect, omits *AgentOmit) ([]*Agent, error) {
	input, err := assignmentsToAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentUpdateManyAndReturnArgs) ([]*Agent, error) {
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
			curr = func(c context.Context, a *AgentUpdateManyAndReturnArgs) ([]*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Agent], cols []string, vals []any, selects *AgentSelect, omits *AgentOmit) ([]*Agent, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Agent]{Where: preds}, selects, omits)
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
		var res []*Agent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Agent.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Agent.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAgentCols(selects, omits, agentPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanAgentRows(rows, returningCols)
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

func (d *AgentDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Agent], cols []string, vals []any, selects *AgentSelect, omits *AgentOmit) ([]*Agent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Agent{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Agent]{Where: preds}, selects, omits)
}
func (d *AgentDelegate) FindUnique(where UniquePredicate[Agent], additional ...PredicateOf[Agent]) *FindUniqueBuilder[Agent, AgentSelect, AgentOmit] {
	return &FindUniqueBuilder[Agent, AgentSelect, AgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AgentDelegate) FindFirst(preds ...PredicateOf[Agent]) *FindFirstBuilder[Agent, AgentSelect, AgentOmit] {
	return &FindFirstBuilder[Agent, AgentSelect, AgentOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AgentDelegate) FindMany(preds ...PredicateOf[Agent]) *FindManyBuilder[Agent, AgentSelect, AgentOmit] {
	return &FindManyBuilder[Agent, AgentSelect, AgentOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AgentDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Agent], additional []PredicateOf[Agent], selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
	allWhere := make([]PredicateOf[Agent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentFindUniqueArgs) (*Agent, error) {
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
			curr = func(c context.Context, a *AgentFindUniqueArgs) (*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Agent],
	selects *AgentSelect,
	omits *AgentOmit,
) (*Agent, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AgentFindFirstArgs) (*Agent, error) {
		return d.runFindFirst(c, QueryParams[Agent]{
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
			curr = func(c context.Context, a *AgentFindFirstArgs) (*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Agent],
	selects *AgentSelect,
	omits *AgentOmit,
) ([]*Agent, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AgentFindManyArgs) ([]*Agent, error) {
		return d.runFindMany(c, QueryParams[Agent]{
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
			curr = func(c context.Context, a *AgentFindManyArgs) ([]*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Agent], selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
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
	returningCols := selectAgentCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Agent{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *AgentDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Agent],
	selects *AgentSelect,
	omits *AgentOmit,
) (*Agent, error) {
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

func (d *AgentDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Agent],
	selects *AgentSelect,
	omits *AgentOmit,
) ([]*Agent, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, agentPKCols, agentUniqueCols, "Agent", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, agentPKCols, agentUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAgentCols(selects, omits)

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

func (d *AgentDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Agent, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Agent", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Agent
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *AgentDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Agent, error) {
	query := buildSelectSQL(d.client, "Agent", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Agent, 0)
	for rows.Next() {
		var res Agent
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
func (d *AgentDelegate) DeleteMany(preds ...PredicateOf[Agent]) *DeleteManyBuilder[Agent] {
	return &DeleteManyBuilder[Agent]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AgentDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Agent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &AgentDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *AgentDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Agent]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Agent")
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

func (d *AgentDelegate) Delete(where UniquePredicate[Agent], additional ...PredicateOf[Agent]) *DeleteBuilder[Agent, AgentSelect, AgentOmit] {
	return &DeleteBuilder[Agent, AgentSelect, AgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *AgentDelegate) executeDelete(ctx context.Context, where UniquePredicate[Agent], additional []PredicateOf[Agent], selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
	allWhere := make([]PredicateOf[Agent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSelect()
	}

	args := &AgentDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentDeleteArgs) (*Agent, error) {
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
			curr = func(c context.Context, a *AgentDeleteArgs) (*Agent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runDelete(ctx context.Context, where []PredicateOf[Agent], selects *AgentSelect, omits *AgentOmit) (*Agent, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectAgentCols(selects, omits, agentPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Agent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Agent.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Agent"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Agent")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Agent]
			pkPreds = append(pkPreds, Predicate[Agent]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "Agent")

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
		return nil, &NotFoundError{Model: "Agent"}
	}

	var row Agent
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *AgentDelegate) Count(preds ...PredicateOf[Agent]) *CountBuilder[Agent] {
	return &CountBuilder[Agent]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AgentDelegate) executeCount(ctx context.Context, params QueryParams[Agent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &AgentCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *AgentCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Agent]{
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
			curr = func(c context.Context, a *AgentCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentDelegate) runCount(ctx context.Context, params QueryParams[Agent]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Agent")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Agent")
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
func (d *AgentDelegate) loadRelations(ctx context.Context, records []*Agent, selects *AgentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Trashes != nil {
		relationSelects, relationOmits, relationParams := selects.Trashes.GetRelationParams()
		returningCols := selectTrashAgentCols(relationSelects, relationOmits, "agentId")
		// Inverse holds the FK: TrashAgent.agentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Agent) string { return p.Id }),
			"TrashAgent",
			"agentId",
			returningCols,
			scanInto(returningCols, (*TrashAgent).ScanFields),
			directKey(func(c *TrashAgent) string { return c.AgentId }),
			appendMany(func(p *Agent) *[]*TrashAgent { return &p.Trashes }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading trashes: %w", err)
		}
		if err := d.client.TrashAgent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skills != nil {
		relationSelects, relationOmits, relationParams := selects.Skills.GetRelationParams()
		returningCols := selectAgentSkillCols(relationSelects, relationOmits, "agentId")
		// Inverse holds the FK: AgentSkill.agentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Agent) string { return p.Id }),
			"AgentSkill",
			"agentId",
			returningCols,
			scanInto(returningCols, (*AgentSkill).ScanFields),
			directKey(func(c *AgentSkill) string { return c.AgentId }),
			appendMany(func(p *Agent) *[]*AgentSkill { return &p.Skills }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skills: %w", err)
		}
		if err := d.client.AgentSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Projects != nil {
		// Many-to-many through join table _ProjectAgents
		relationSelects, relationOmits, relationParams := selects.Projects.GetRelationParams()
		returningCols := selectProjectCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *Agent) string { return p.Id }),
			"Project",
			"name",
			"_ProjectAgents",
			"agentId",
			"projectName",
			returningCols,
			(*Project).ScanFields,
			appendMany(func(p *Agent) *[]*Project { return &p.Projects }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading projects: %w", err)
		}
		if err := d.client.Project.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
