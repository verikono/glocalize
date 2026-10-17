package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Reference represents the database model
type Reference struct {
	Name         string            `db:"name" json:"name"`
	Description  string            `db:"description" json:"description"`
	AccessType   URIType           `db:"accessType" json:"accessType"`
	ResourceType ReferenceTypeType `db:"resourceType" json:"resourceType"`
	Content      *string           `db:"content" json:"content,omitempty"`
	Skills       []*SkillReference `json:"skills,omitempty"`
}

// ReferenceCreate contains model input fields for Reference creation operations.
//
// Fields for Reference:
//
//	name         string        required
//	description  string        required
//	accessType   URI           required
//	resourceType ReferenceType required
//	content      string        optional
type ReferenceCreate struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	AccessType   URIType           `json:"accessType"`
	ResourceType ReferenceTypeType `json:"resourceType"`
	Content      *string           `json:"content"`
}

// colMask returns a bit mask of columns that are set
func (s *ReferenceCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	mask |= 1 << 3
	if s.Content != nil {
		mask |= 1 << 4
	}
	return mask
}

// ReferenceUpdate contains model input fields for Reference update operations.
type ReferenceUpdate struct {
	Name         *string            `json:"name"`
	Description  *string            `json:"description"`
	AccessType   *URIType           `json:"accessType"`
	ResourceType *ReferenceTypeType `json:"resourceType"`
	Content      *string            `json:"content"`
}

func (u *ReferenceUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	if u.Description != nil {
		cols = append(cols, "description")
		vals = append(vals, u.Description)
	}
	if u.AccessType != nil {
		cols = append(cols, "accessType")
		vals = append(vals, u.AccessType)
	}
	if u.ResourceType != nil {
		cols = append(cols, "resourceType")
		vals = append(vals, u.ResourceType)
	}
	if u.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, u.Content)
	}
	return cols, vals
}

func assignmentsToReferenceUpdate(assignments []FieldAssignment) (ReferenceUpdate, error) {
	var input ReferenceUpdate
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
		case "description":
			if v, ok := a.Val.(string); ok {
				input.Description = &v
				errs.ValidateString("description", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Description = v
			} else {
				errs.Add("description", a.Val, "type", "field description must be of type string")
			}
		case "accessType":
			if v, ok := a.Val.(URIType); ok {
				input.AccessType = &v
				if !v.IsValid() {
					errs.Add("accessType", v, "enum", fmt.Sprintf("invalid enum value %q for field accessType", v))
				}
			} else if v, ok := a.Val.(*URIType); ok {
				input.AccessType = v
				if v != nil && !v.IsValid() {
					errs.Add("accessType", *v, "enum", fmt.Sprintf("invalid enum value %q for field accessType", *v))
				}
			} else {
				errs.Add("accessType", a.Val, "type", "field accessType must be of type URIType")
			}
		case "resourceType":
			if v, ok := a.Val.(ReferenceTypeType); ok {
				input.ResourceType = &v
				if !v.IsValid() {
					errs.Add("resourceType", v, "enum", fmt.Sprintf("invalid enum value %q for field resourceType", v))
				}
			} else if v, ok := a.Val.(*ReferenceTypeType); ok {
				input.ResourceType = v
				if v != nil && !v.IsValid() {
					errs.Add("resourceType", *v, "enum", fmt.Sprintf("invalid enum value %q for field resourceType", *v))
				}
			} else {
				errs.Add("resourceType", a.Val, "type", "field resourceType must be of type ReferenceTypeType")
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
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// ReferenceSelect specifies which scalar and relation fields to select for Reference.
//
// Selectable fields:
//
//	-- Scalars --
//	name         (bool)
//	description  (bool)
//	accessType   (bool)
//	resourceType (bool)
//	content      (bool)
//	-- Relations --
//	skills       ([]SkillReference)
type ReferenceSelect struct {
	Name         bool                      `json:"name"`
	Description  bool                      `json:"description"`
	AccessType   bool                      `json:"accessType"`
	ResourceType bool                      `json:"resourceType"`
	Content      bool                      `json:"content"`
	Skills       SkillReferenceSelectQuery `json:"skills,omitempty"`
}

var fullReferenceSelectVal = &ReferenceSelect{
	Name:         true,
	Description:  true,
	AccessType:   true,
	ResourceType: true,
	Content:      true,
}

func fullReferenceSelect() *ReferenceSelect {
	return fullReferenceSelectVal
}

func (s *ReferenceSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Name || s.Description || s.AccessType || s.ResourceType || s.Content
}

func (s *ReferenceSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type ReferenceOmit struct {
	Name         bool `json:"name"`
	Description  bool `json:"description"`
	AccessType   bool `json:"accessType"`
	ResourceType bool `json:"resourceType"`
	Content      bool `json:"content"`
}

type ReferenceSelectQuery interface {
	GetRelationParams() (*ReferenceSelect, *ReferenceOmit, QueryParams[Reference])
}

func (s *ReferenceSelect) GetRelationParams() (*ReferenceSelect, *ReferenceOmit, QueryParams[Reference]) {
	return s, nil, QueryParams[Reference]{}
}

type ReferenceQueryBuilder struct {
	selects *ReferenceSelect
	omits   *ReferenceOmit
	where   []PredicateOf[Reference]
	take    *int
	skip    *int
	orderBy []OrderBy[Reference]
	cursor  UniquePredicate[Reference]
}

func (b *ReferenceQueryBuilder) Where(preds ...PredicateOf[Reference]) *ReferenceQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *ReferenceQueryBuilder) Take(limit int) *ReferenceQueryBuilder {
	b.take = &limit
	return b
}

func (b *ReferenceQueryBuilder) Skip(offset int) *ReferenceQueryBuilder {
	b.skip = &offset
	return b
}

func (b *ReferenceQueryBuilder) OrderBy(orders ...OrderBy[Reference]) *ReferenceQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *ReferenceQueryBuilder) Cursor(where UniquePredicate[Reference]) *ReferenceQueryBuilder {
	b.cursor = where
	return b
}

func (b *ReferenceQueryBuilder) Select(s ReferenceSelect) *ReferenceQueryBuilder {
	b.selects = &s
	return b
}

func (b *ReferenceQueryBuilder) Omit(o ReferenceOmit) *ReferenceQueryBuilder {
	b.omits = &o
	return b
}

func (b *ReferenceQueryBuilder) GetRelationParams() (*ReferenceSelect, *ReferenceOmit, QueryParams[Reference]) {
	if b == nil {
		return nil, nil, QueryParams[Reference]{}
	}
	return b.selects, b.omits, QueryParams[Reference]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// ReferenceCreateArgs is the input argument passed to Reference Create extension hooks.
//
// Fields for Reference:
//
//	name         string        required
//	description  string        required
//	accessType   URI           required
//	resourceType ReferenceType required
//	content      string        optional
//
// Relations for Reference:
//
//	skills       ([]SkillReference)
type ReferenceCreateArgs struct {
	// Data contains the model fields to insert.
	Data *ReferenceCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *ReferenceSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// ReferenceCreateManyArgs is the input argument passed to Reference CreateMany extension hooks.
//
// Fields for Reference:
//
//	name         string        required
//	description  string        required
//	accessType   URI           required
//	resourceType ReferenceType required
//	content      string        optional
type ReferenceCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ReferenceCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ReferenceCreateManyArgs) AppendData(builders ...*ReferenceCreateBuilder) *ReferenceCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToReferenceCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ReferenceCreateManyAndReturnArgs is the input argument passed to Reference CreateManyAndReturn extension hooks.
//
// Fields for Reference:
//
//	name         string        required
//	description  string        required
//	accessType   URI           required
//	resourceType ReferenceType required
//	content      string        optional
//
// Relations for Reference:
//
//	skills       ([]SkillReference)
type ReferenceCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*ReferenceCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *ReferenceSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *ReferenceCreateManyAndReturnArgs) AppendData(builders ...*ReferenceCreateBuilder) *ReferenceCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToReferenceCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// ReferenceFindUniqueArgs is the input argument passed to Reference FindUnique extension hooks.
type ReferenceFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Reference]
	// Select specifies which scalar and relation fields to select and return.
	Select *ReferenceSelect
}

func (a *ReferenceFindUniqueArgs) SetWhere(unique UniquePredicate[Reference], additional ...PredicateOf[Reference]) *ReferenceFindUniqueArgs {
	a.Where = make([]PredicateOf[Reference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ReferenceFindFirstArgs is the input argument passed to Reference FindFirst extension hooks.
type ReferenceFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Reference]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Reference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ReferenceSelect
}

func (a *ReferenceFindFirstArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceFindFirstArgs {
	a.Where = preds
	return a
}

func (a *ReferenceFindFirstArgs) SetOrderBy(orders ...OrderBy[Reference]) *ReferenceFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *ReferenceFindFirstArgs) SetCursor(cursor UniquePredicate[Reference]) *ReferenceFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *ReferenceFindFirstArgs) SetSkip(n int) *ReferenceFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *ReferenceFindFirstArgs) SetTake(n int) *ReferenceFindFirstArgs {
	a.Take = &n
	return a
}

// ReferenceFindManyArgs is the input argument passed to Reference FindMany extension hooks.
type ReferenceFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Reference]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Reference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *ReferenceSelect
}

func (a *ReferenceFindManyArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceFindManyArgs {
	a.Where = preds
	return a
}

func (a *ReferenceFindManyArgs) SetOrderBy(orders ...OrderBy[Reference]) *ReferenceFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *ReferenceFindManyArgs) SetCursor(cursor UniquePredicate[Reference]) *ReferenceFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *ReferenceFindManyArgs) SetSkip(n int) *ReferenceFindManyArgs {
	a.Skip = &n
	return a
}

func (a *ReferenceFindManyArgs) SetTake(n int) *ReferenceFindManyArgs {
	a.Take = &n
	return a
}

// ReferenceCountArgs is the input argument passed to Reference Count extension hooks.
type ReferenceCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *ReferenceCountArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceCountArgs {
	a.Where = preds
	return a
}

func (a *ReferenceCountArgs) SetSkip(n int) *ReferenceCountArgs {
	a.Skip = &n
	return a
}

func (a *ReferenceCountArgs) SetTake(n int) *ReferenceCountArgs {
	a.Take = &n
	return a
}

// ReferenceDeleteArgs is the input argument passed to Reference Delete extension hooks.
type ReferenceDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Reference]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *ReferenceSelect
}

func (a *ReferenceDeleteArgs) SetWhere(unique UniquePredicate[Reference], additional ...PredicateOf[Reference]) *ReferenceDeleteArgs {
	a.Where = make([]PredicateOf[Reference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ReferenceDeleteManyArgs is the input argument passed to Reference DeleteMany extension hooks.
type ReferenceDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
}

func (a *ReferenceDeleteManyArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceDeleteManyArgs {
	a.Where = preds
	return a
}

// ReferenceUpdateArgs is the input argument passed to Reference Update extension hooks.
type ReferenceUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Reference]
	// Data contains the model fields to update.
	Data *ReferenceUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ReferenceSelect
}

func (a *ReferenceUpdateArgs) SetWhere(unique UniquePredicate[Reference], additional ...PredicateOf[Reference]) *ReferenceUpdateArgs {
	a.Where = make([]PredicateOf[Reference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// ReferenceUpdateManyArgs is the input argument passed to Reference UpdateMany extension hooks.
type ReferenceUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
	// Data contains the model fields to update.
	Data *ReferenceUpdate
}

func (a *ReferenceUpdateManyArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceUpdateManyArgs {
	a.Where = preds
	return a
}

// ReferenceUpdateManyAndReturnArgs is the input argument passed to Reference UpdateManyAndReturn extension hooks.
type ReferenceUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Reference]
	// Data contains the model fields to update.
	Data *ReferenceUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *ReferenceSelect
}

func (a *ReferenceUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Reference]) *ReferenceUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type ReferenceCreateQuery = func(ctx context.Context, args *ReferenceCreateArgs) (*Reference, error)
type ReferenceCreateManyQuery = func(ctx context.Context, args *ReferenceCreateManyArgs) (int64, error)
type ReferenceCreateManyAndReturnQuery = func(ctx context.Context, args *ReferenceCreateManyAndReturnArgs) ([]*Reference, error)
type ReferenceFindUniqueQuery = func(ctx context.Context, args *ReferenceFindUniqueArgs) (*Reference, error)
type ReferenceFindFirstQuery = func(ctx context.Context, args *ReferenceFindFirstArgs) (*Reference, error)
type ReferenceFindManyQuery = func(ctx context.Context, args *ReferenceFindManyArgs) ([]*Reference, error)
type ReferenceDeleteQuery = func(ctx context.Context, args *ReferenceDeleteArgs) (*Reference, error)
type ReferenceDeleteManyQuery = func(ctx context.Context, args *ReferenceDeleteManyArgs) (int64, error)
type ReferenceCountQuery = func(ctx context.Context, args *ReferenceCountArgs) (int64, error)
type ReferenceUpdateQuery = func(ctx context.Context, args *ReferenceUpdateArgs) (*Reference, error)
type ReferenceUpdateManyQuery = func(ctx context.Context, args *ReferenceUpdateManyArgs) (int64, error)
type ReferenceUpdateManyAndReturnQuery = func(ctx context.Context, args *ReferenceUpdateManyAndReturnArgs) ([]*Reference, error)

type ReferenceExtension struct {
	Create              func(ctx context.Context, args *ReferenceCreateArgs, next ReferenceCreateQuery) (*Reference, error)
	CreateMany          func(ctx context.Context, args *ReferenceCreateManyArgs, next ReferenceCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *ReferenceCreateManyAndReturnArgs, next ReferenceCreateManyAndReturnQuery) ([]*Reference, error)
	FindUnique          func(ctx context.Context, args *ReferenceFindUniqueArgs, next ReferenceFindUniqueQuery) (*Reference, error)
	FindFirst           func(ctx context.Context, args *ReferenceFindFirstArgs, next ReferenceFindFirstQuery) (*Reference, error)
	FindMany            func(ctx context.Context, args *ReferenceFindManyArgs, next ReferenceFindManyQuery) ([]*Reference, error)
	Delete              func(ctx context.Context, args *ReferenceDeleteArgs, next ReferenceDeleteQuery) (*Reference, error)
	DeleteMany          func(ctx context.Context, args *ReferenceDeleteManyArgs, next ReferenceDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *ReferenceCountArgs, next ReferenceCountQuery) (int64, error)
	Update              func(ctx context.Context, args *ReferenceUpdateArgs, next ReferenceUpdateQuery) (*Reference, error)
	UpdateMany          func(ctx context.Context, args *ReferenceUpdateManyArgs, next ReferenceUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *ReferenceUpdateManyAndReturnArgs, next ReferenceUpdateManyAndReturnQuery) ([]*Reference, error)
}

type ReferenceDelegate struct {
	client     *Queries
	extensions []ReferenceExtension
}

func (d *ReferenceDelegate) Use(exts ...ReferenceExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Reference) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "name":
			targets[i] = &m.Name
		case "description":
			targets[i] = &m.Description
		case "accessType":
			targets[i] = &m.AccessType
		case "resourceType":
			targets[i] = &m.ResourceType
		case "content":
			targets[i] = &m.Content
		}
	}
	return targets
}

var referenceDefaultCols = []string{
	"name",
	"description",
	"accessType",
	"resourceType",
	"content",
}

var referencePKCols = []string{}

var referenceUniqueCols = []string{
	"name",
}

func selectReferenceCols(selects *ReferenceSelect, omits *ReferenceOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return referenceDefaultCols
	}

	anySelected := selects != nil && (selects.Name || selects.Description || selects.AccessType || selects.ResourceType || selects.Content || selects.Skills != nil)

	specs := []colSpec{
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
		{"description", selects != nil && selects.Description, omits != nil && omits.Description, false},
		{"accessType", selects != nil && selects.AccessType, omits != nil && omits.AccessType, false},
		{"resourceType", selects != nil && selects.ResourceType, omits != nil && omits.ResourceType, false},
		{"content", selects != nil && selects.Content, omits != nil && omits.Content, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *ReferenceSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Skills != nil
}

type ReferenceUpsertBuilder struct {
	*CreateBuilder[Reference, ReferenceSelect, ReferenceOmit]
}

func (b *ReferenceUpsertBuilder) Select(s ReferenceSelect) *ReferenceUpsertBuilder {
	b.selects = &s
	return b
}

func (b *ReferenceUpsertBuilder) Omit(o ReferenceOmit) *ReferenceUpsertBuilder {
	b.omits = &o
	return b
}

type ReferenceCreateBuilder struct {
	*CreateBuilder[Reference, ReferenceSelect, ReferenceOmit]
}

func (b *ReferenceCreateBuilder) Select(s ReferenceSelect) *ReferenceCreateBuilder {
	b.selects = &s
	return b
}

func (b *ReferenceCreateBuilder) Omit(o ReferenceOmit) *ReferenceCreateBuilder {
	b.omits = &o
	return b
}

func (b *ReferenceCreateBuilder) OnConflict(target UniqueConstraintTarget) *ReferenceConflictBuilder[ReferenceUpsertBuilder] {
	upsertBuilder := &ReferenceUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &ReferenceConflictBuilder[ReferenceUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *ReferenceCreateBuilder) SetName(v string) *ReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *ReferenceCreateBuilder) SetDescription(v string) *ReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *ReferenceCreateBuilder) SetAccessType(v URIType) *ReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "accessType", Val: v})
	return b
}
func (b *ReferenceCreateBuilder) SetResourceType(v ReferenceTypeType) *ReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "resourceType", Val: v})
	return b
}
func (b *ReferenceCreateBuilder) SetContent(v string) *ReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *ReferenceCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Reference]) *ReferenceCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ReferenceDelegate) Create() *ReferenceCreateBuilder {
	return &ReferenceCreateBuilder{
		CreateBuilder: &CreateBuilder[Reference, ReferenceSelect, ReferenceOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedReferenceName         uint64 = 1 << 0
	providedReferenceDescription  uint64 = 1 << 1
	providedReferenceAccessType   uint64 = 1 << 2
	providedReferenceResourceType uint64 = 1 << 3
	providedReferenceContent      uint64 = 1 << 4
)

func assignmentsToReferenceCreate(assignments []FieldAssignment) (ReferenceCreate, error) {
	var input ReferenceCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "name":
			provided |= providedReferenceName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		case "description":
			provided |= providedReferenceDescription
			if v, ok := a.Val.(string); ok {
				input.Description = v
				errs.ValidateString("description", v, true, 0, false, false)
			} else {
				errs.Add("description", a.Val, "type", "field description must be of type string")
			}
		case "accessType":
			provided |= providedReferenceAccessType
			if v, ok := a.Val.(URIType); ok {
				input.AccessType = v
				if !v.IsValid() {
					errs.Add("accessType", v, "enum", fmt.Sprintf("invalid enum value %q for field accessType", v))
				}
			} else {
				errs.Add("accessType", a.Val, "type", "field accessType must be of type URIType")
			}
		case "resourceType":
			provided |= providedReferenceResourceType
			if v, ok := a.Val.(ReferenceTypeType); ok {
				input.ResourceType = v
				if !v.IsValid() {
					errs.Add("resourceType", v, "enum", fmt.Sprintf("invalid enum value %q for field resourceType", v))
				}
			} else {
				errs.Add("resourceType", a.Val, "type", "field resourceType must be of type ReferenceTypeType")
			}
		case "content":
			provided |= providedReferenceContent
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		}
	}
	if provided&providedReferenceName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}
	if provided&providedReferenceDescription == 0 {
		errs.Add("description", "", "required", "field Description is required")
	}
	if provided&providedReferenceAccessType == 0 {
		errs.Add("accessType", nil, "required", "field AccessType is required")
	}
	if provided&providedReferenceResourceType == 0 {
		errs.Add("resourceType", nil, "required", "field ResourceType is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *ReferenceCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 5)
	vals = make([]any, 0, 5)
	cols = append(cols, "name")
	vals = append(vals, s.Name)
	cols = append(cols, "description")
	vals = append(vals, s.Description)
	cols = append(cols, "accessType")
	vals = append(vals, s.AccessType)
	cols = append(cols, "resourceType")
	vals = append(vals, s.ResourceType)
	if s.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, *s.Content)
	}
	return
}

func partitionReferenceInputs(dialect Dialect, inputs []*ReferenceCreate) [][]*ReferenceCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*ReferenceCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*ReferenceCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*ReferenceCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*ReferenceCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*ReferenceCreate{inputs}
}

func (d *ReferenceDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *ReferenceSelect, omits *ReferenceOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Reference, error) {
	input, err := assignmentsToReferenceCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectReferenceCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ReferenceCreateArgs) (*Reference, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectReferenceCols(a.Select, omits)
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
			curr = func(c context.Context, a *ReferenceCreateArgs) (*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type ReferenceCreateManyBuilder struct {
	*CreateManyBuilder[Reference]
}

func (b *ReferenceCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *ReferenceConflictBuilder[ReferenceCreateManyBuilder] {
	return &ReferenceConflictBuilder[ReferenceCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type ReferenceCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Reference, ReferenceSelect, ReferenceOmit]
}

func (b *ReferenceCreateManyAndReturnBuilder) Select(s ReferenceSelect) *ReferenceCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *ReferenceCreateManyAndReturnBuilder) Omit(o ReferenceOmit) *ReferenceCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *ReferenceCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *ReferenceConflictBuilder[ReferenceCreateManyAndReturnBuilder] {
	return &ReferenceConflictBuilder[ReferenceCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToReferenceRecordInputs(builders []*ReferenceCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *ReferenceDelegate) CreateMany(builders ...*ReferenceCreateBuilder) *ReferenceCreateManyBuilder {
	return &ReferenceCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Reference]{
			records:  createBuildersToReferenceRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *ReferenceDelegate) CreateManyAndReturn(builders ...*ReferenceCreateBuilder) *ReferenceCreateManyAndReturnBuilder {
	return &ReferenceCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Reference, ReferenceSelect, ReferenceOmit]{
			records:  createBuildersToReferenceRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToReferenceCreateInputs(records []RecordInput) ([]*ReferenceCreate, error) {
	structs := make([]ReferenceCreate, len(records))
	inputs := make([]*ReferenceCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToReferenceCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *ReferenceDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToReferenceCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &ReferenceCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ReferenceCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ReferenceCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ReferenceSelect, omits *ReferenceOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Reference, error) {
	inputs, err := recordsToReferenceCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *ReferenceCreateManyAndReturnArgs) ([]*Reference, error) {
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
			curr = func(c context.Context, a *ReferenceCreateManyAndReturnArgs) ([]*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *ReferenceSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Reference, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Reference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Reference.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Reference.loadRelations(ctx, []*Reference{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Reference", cols, returningCols, referencePKCols, conflictTarget, conflictAction, len(vals))
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

		var res Reference
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, referencePKCols)
}

func (d *ReferenceDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Reference, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Reference") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Reference")
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

	var res Reference
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *ReferenceDelegate) buildBulkInsertSQL(q *Queries, batch []*ReferenceCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 5)
	for i, c := range referenceDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Reference")
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
				vals = append(vals, input.Name)
			case "description":
				vals = append(vals, input.Description)
			case "accessType":
				vals = append(vals, input.AccessType)
			case "resourceType":
				vals = append(vals, input.ResourceType)
			case "content":
				if input.Content != nil {
					vals = append(vals, *input.Content)
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

func applyReferenceConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanReferenceRows(rows *sql.Rows, returningCols []string) ([]*Reference, error) {
	var records []*Reference
	for rows.Next() {
		var res Reference
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

func (d *ReferenceDelegate) runCreateMany(ctx context.Context, inputs []*ReferenceCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionReferenceInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyReferenceConflictClause(d.client.dialect, queryStr, vals, cols, referencePKCols, conflictTarget, conflictAction)

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

func (d *ReferenceDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*ReferenceCreate,
	selects *ReferenceSelect,
	omits *ReferenceOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Reference, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionReferenceInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Reference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Reference.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Reference.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectReferenceCols(selects, omits, referencePKCols...)
	recordsOut := make([]*Reference, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyReferenceConflictClause(d.client.dialect, queryStr, vals, cols, referencePKCols, conflictTarget, conflictAction)

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

		scanned, err := scanReferenceRows(rows, returningCols)
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

func (d *ReferenceDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*ReferenceCreate,
	selects *ReferenceSelect,
	omits *ReferenceOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Reference, error) {
	batches := partitionReferenceInputs(d.client.dialect, inputs)
	returningCols := selectReferenceCols(selects, omits, referencePKCols...)
	recordsOut := make([]*Reference, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyReferenceConflictClause(d.client.dialect, queryStr, vals, cols, referencePKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Reference") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Reference")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, referencePKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, referencePKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanReferenceRows(rows, returningCols)
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

type ReferenceConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *ReferenceConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *ReferenceConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *ReferenceConflictBuilder[B]) Update(fn func(u *ReferenceUpsert)) *B {
	cb.setAction(*ReferenceConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// ReferenceConflictUpdate creates a custom ConflictAction for Reference upsert conflicts.
func ReferenceConflictUpdate(fn func(u *ReferenceUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newReferenceUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type ReferenceUpsert struct {
	Name         fieldUpsert[string]
	Description  fieldUpsert[string]
	AccessType   fieldUpsert[URIType]
	ResourceType fieldUpsert[ReferenceTypeType]
	Content      fieldUpsert[*string]
}

func newReferenceUpsert(up *ConflictUpdate) *ReferenceUpsert {
	return &ReferenceUpsert{
		Name:         fieldUpsert[string]{column: "name", update: up},
		Description:  fieldUpsert[string]{column: "description", update: up},
		AccessType:   fieldUpsert[URIType]{column: "accessType", update: up},
		ResourceType: fieldUpsert[ReferenceTypeType]{column: "resourceType", update: up},
		Content:      fieldUpsert[*string]{column: "content", update: up},
	}
}

type ReferenceUpdateBuilder struct {
	*UpdateBuilder[Reference, ReferenceSelect, ReferenceOmit]
}

type ReferenceUpdateManyBuilder struct {
	*UpdateManyBuilder[Reference]
}

type ReferenceUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Reference, ReferenceSelect, ReferenceOmit]
}

func (b *ReferenceUpdateBuilder) SetName(v string) *ReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *ReferenceUpdateManyBuilder) SetName(v string) *ReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) SetName(v string) *ReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}
func (b *ReferenceUpdateBuilder) SetDescription(v string) *ReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *ReferenceUpdateManyBuilder) SetDescription(v string) *ReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) SetDescription(v string) *ReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "description", Val: v})
	return b
}
func (b *ReferenceUpdateBuilder) SetAccessType(v URIType) *ReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "accessType", Val: v})
	return b
}

func (b *ReferenceUpdateManyBuilder) SetAccessType(v URIType) *ReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "accessType", Val: v})
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) SetAccessType(v URIType) *ReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "accessType", Val: v})
	return b
}
func (b *ReferenceUpdateBuilder) SetResourceType(v ReferenceTypeType) *ReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "resourceType", Val: v})
	return b
}

func (b *ReferenceUpdateManyBuilder) SetResourceType(v ReferenceTypeType) *ReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "resourceType", Val: v})
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) SetResourceType(v ReferenceTypeType) *ReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "resourceType", Val: v})
	return b
}
func (b *ReferenceUpdateBuilder) SetContent(v string) *ReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *ReferenceUpdateManyBuilder) SetContent(v string) *ReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) SetContent(v string) *ReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *ReferenceUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Reference]) *ReferenceUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ReferenceUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Reference]) *ReferenceUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *ReferenceUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Reference]) *ReferenceUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *ReferenceDelegate) Update(where UniquePredicate[Reference], additional ...PredicateOf[Reference]) *ReferenceUpdateBuilder {
	return &ReferenceUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Reference, ReferenceSelect, ReferenceOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *ReferenceDelegate) UpdateMany(preds ...PredicateOf[Reference]) *ReferenceUpdateManyBuilder {
	return &ReferenceUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Reference]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *ReferenceDelegate) UpdateManyAndReturn(preds ...PredicateOf[Reference]) *ReferenceUpdateManyAndReturnBuilder {
	return &ReferenceUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Reference, ReferenceSelect, ReferenceOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *ReferenceDelegate) buildUpdateSQL(preds []PredicateOf[Reference], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Reference")
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

func (d *ReferenceDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Reference], additional []PredicateOf[Reference], assignments []FieldAssignment, selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
	allWhere := make([]PredicateOf[Reference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToReferenceUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ReferenceUpdateArgs) (*Reference, error) {
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
			curr = func(c context.Context, a *ReferenceUpdateArgs) (*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Reference], cols []string, vals []any, selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
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
		var res *Reference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Reference.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Reference.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectReferenceCols(selects, omits, referencePKCols...)
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
		return nil, &NotFoundError{Model: "Reference"}
	}

	var res Reference
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Reference{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *ReferenceDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Reference], cols []string, vals []any) (int64, error) {
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

func (d *ReferenceDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Reference], cols []string, vals []any, selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Reference"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *ReferenceDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Reference], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToReferenceUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &ReferenceUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *ReferenceUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ReferenceUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *ReferenceDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Reference], assignments []FieldAssignment, selects *ReferenceSelect, omits *ReferenceOmit) ([]*Reference, error) {
	input, err := assignmentsToReferenceUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *ReferenceUpdateManyAndReturnArgs) ([]*Reference, error) {
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
			curr = func(c context.Context, a *ReferenceUpdateManyAndReturnArgs) ([]*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Reference], cols []string, vals []any, selects *ReferenceSelect, omits *ReferenceOmit) ([]*Reference, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Reference]{Where: preds}, selects, omits)
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
		var res []*Reference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Reference.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Reference.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectReferenceCols(selects, omits, referencePKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanReferenceRows(rows, returningCols)
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

func (d *ReferenceDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Reference], cols []string, vals []any, selects *ReferenceSelect, omits *ReferenceOmit) ([]*Reference, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Reference{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Reference]{Where: preds}, selects, omits)
}
func (d *ReferenceDelegate) FindUnique(where UniquePredicate[Reference], additional ...PredicateOf[Reference]) *FindUniqueBuilder[Reference, ReferenceSelect, ReferenceOmit] {
	return &FindUniqueBuilder[Reference, ReferenceSelect, ReferenceOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *ReferenceDelegate) FindFirst(preds ...PredicateOf[Reference]) *FindFirstBuilder[Reference, ReferenceSelect, ReferenceOmit] {
	return &FindFirstBuilder[Reference, ReferenceSelect, ReferenceOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *ReferenceDelegate) FindMany(preds ...PredicateOf[Reference]) *FindManyBuilder[Reference, ReferenceSelect, ReferenceOmit] {
	return &FindManyBuilder[Reference, ReferenceSelect, ReferenceOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *ReferenceDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Reference], additional []PredicateOf[Reference], selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
	allWhere := make([]PredicateOf[Reference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ReferenceFindUniqueArgs) (*Reference, error) {
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
			curr = func(c context.Context, a *ReferenceFindUniqueArgs) (*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Reference],
	selects *ReferenceSelect,
	omits *ReferenceOmit,
) (*Reference, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ReferenceFindFirstArgs) (*Reference, error) {
		return d.runFindFirst(c, QueryParams[Reference]{
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
			curr = func(c context.Context, a *ReferenceFindFirstArgs) (*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Reference],
	selects *ReferenceSelect,
	omits *ReferenceOmit,
) ([]*Reference, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *ReferenceFindManyArgs) ([]*Reference, error) {
		return d.runFindMany(c, QueryParams[Reference]{
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
			curr = func(c context.Context, a *ReferenceFindManyArgs) ([]*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Reference], selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
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
	returningCols := selectReferenceCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Reference{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *ReferenceDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Reference],
	selects *ReferenceSelect,
	omits *ReferenceOmit,
) (*Reference, error) {
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

func (d *ReferenceDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Reference],
	selects *ReferenceSelect,
	omits *ReferenceOmit,
) ([]*Reference, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, referencePKCols, referenceUniqueCols, "Reference", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, referencePKCols, referenceUniqueCols, isCursorQuery, params.Take)
	returningCols := selectReferenceCols(selects, omits)

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

func (d *ReferenceDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Reference, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Reference", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Reference
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *ReferenceDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Reference, error) {
	query := buildSelectSQL(d.client, "Reference", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Reference, 0)
	for rows.Next() {
		var res Reference
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
func (d *ReferenceDelegate) DeleteMany(preds ...PredicateOf[Reference]) *DeleteManyBuilder[Reference] {
	return &DeleteManyBuilder[Reference]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *ReferenceDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Reference]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &ReferenceDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *ReferenceDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *ReferenceDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Reference]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Reference")
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

func (d *ReferenceDelegate) Delete(where UniquePredicate[Reference], additional ...PredicateOf[Reference]) *DeleteBuilder[Reference, ReferenceSelect, ReferenceOmit] {
	return &DeleteBuilder[Reference, ReferenceSelect, ReferenceOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *ReferenceDelegate) executeDelete(ctx context.Context, where UniquePredicate[Reference], additional []PredicateOf[Reference], selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
	allWhere := make([]PredicateOf[Reference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullReferenceSelect()
	}

	args := &ReferenceDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *ReferenceDeleteArgs) (*Reference, error) {
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
			curr = func(c context.Context, a *ReferenceDeleteArgs) (*Reference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runDelete(ctx context.Context, where []PredicateOf[Reference], selects *ReferenceSelect, omits *ReferenceOmit) (*Reference, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectReferenceCols(selects, omits, referencePKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Reference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Reference.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Reference"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Reference")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Reference]

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
	d.client.dialect.WriteQuotedIdent(&sb, "Reference")

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
		return nil, &NotFoundError{Model: "Reference"}
	}

	var row Reference
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *ReferenceDelegate) Count(preds ...PredicateOf[Reference]) *CountBuilder[Reference] {
	return &CountBuilder[Reference]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *ReferenceDelegate) executeCount(ctx context.Context, params QueryParams[Reference]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &ReferenceCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *ReferenceCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Reference]{
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
			curr = func(c context.Context, a *ReferenceCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *ReferenceDelegate) runCount(ctx context.Context, params QueryParams[Reference]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Reference")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Reference")
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
func (d *ReferenceDelegate) loadRelations(ctx context.Context, records []*Reference, selects *ReferenceSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Skills != nil {
		relationSelects, relationOmits, relationParams := selects.Skills.GetRelationParams()
		returningCols := selectSkillReferenceCols(relationSelects, relationOmits, "referenceName")
		// Inverse holds the FK: SkillReference.referenceName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Reference) string { return p.Name }),
			"SkillReference",
			"referenceName",
			returningCols,
			scanInto(returningCols, (*SkillReference).ScanFields),
			directKey(func(c *SkillReference) string { return c.ReferenceName }),
			appendMany(func(p *Reference) *[]*SkillReference { return &p.Skills }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skills: %w", err)
		}
		if err := d.client.SkillReference.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
