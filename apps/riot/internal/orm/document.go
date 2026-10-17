package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Document represents the database model
type Document struct {
	Id        string `db:"id" json:"id"`
	Kind      string `db:"kind" json:"kind"`
	Body      []byte `db:"body" json:"body"`
	UpdatedAt int32  `db:"updated_at" json:"updatedAt"`
}

// DocumentCreate contains model input fields for Document creation operations.
//
// Fields for Document:
//
//	id        string required
//	kind      string default:
//	body      []byte required
//	updatedAt int32  required
type DocumentCreate struct {
	Id        string  `json:"id"`
	Kind      *string `json:"kind"`
	Body      []byte  `json:"body"`
	UpdatedAt int32   `json:"updatedAt"`
}

// colMask returns a bit mask of columns that are set
func (s *DocumentCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	if s.Kind != nil {
		mask |= 1 << 1
	}
	mask |= 1 << 2
	mask |= 1 << 3
	return mask
}

// DocumentUpdate contains model input fields for Document update operations.
type DocumentUpdate struct {
	Id        *string `json:"id"`
	Kind      *string `json:"kind"`
	Body      *[]byte `json:"body"`
	UpdatedAt *int32  `json:"updatedAt"`
}

func (u *DocumentUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, u.Id)
	}
	if u.Kind != nil {
		cols = append(cols, "kind")
		vals = append(vals, u.Kind)
	}
	if u.Body != nil {
		cols = append(cols, "body")
		vals = append(vals, u.Body)
	}
	if u.UpdatedAt != nil {
		cols = append(cols, "updated_at")
		vals = append(vals, u.UpdatedAt)
	}
	return cols, vals
}

func assignmentsToDocumentUpdate(assignments []FieldAssignment) (DocumentUpdate, error) {
	var input DocumentUpdate
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
		case "kind":
			if v, ok := a.Val.(string); ok {
				input.Kind = &v
				errs.ValidateString("kind", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Kind = v
			} else {
				errs.Add("kind", a.Val, "type", "field kind must be of type string")
			}
		case "body":
			if v, ok := a.Val.([]byte); ok {
				input.Body = &v
			} else if v, ok := a.Val.(*[]byte); ok {
				input.Body = v
			} else {
				errs.Add("body", a.Val, "type", "field body must be of type []byte")
			}
		case "updated_at":
			if v, ok := a.Val.(int32); ok {
				input.UpdatedAt = &v
			} else if v, ok := a.Val.(*int32); ok {
				input.UpdatedAt = v
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type int32")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// DocumentSelect specifies which scalar and relation fields to select for Document.
//
// Selectable fields:
//
//	-- Scalars --
//	id        (bool)
//	kind      (bool)
//	body      (bool)
//	updatedAt (bool)
type DocumentSelect struct {
	Id        bool `json:"id"`
	Kind      bool `json:"kind"`
	Body      bool `json:"body"`
	UpdatedAt bool `json:"updatedAt"`
}

var fullDocumentSelectVal = &DocumentSelect{
	Id:        true,
	Kind:      true,
	Body:      true,
	UpdatedAt: true,
}

func fullDocumentSelect() *DocumentSelect {
	return fullDocumentSelectVal
}

func (s *DocumentSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.Kind || s.Body || s.UpdatedAt
}

func (s *DocumentSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type DocumentOmit struct {
	Id        bool `json:"id"`
	Kind      bool `json:"kind"`
	Body      bool `json:"body"`
	UpdatedAt bool `json:"updatedAt"`
}

type DocumentSelectQuery interface {
	GetRelationParams() (*DocumentSelect, *DocumentOmit, QueryParams[Document])
}

func (s *DocumentSelect) GetRelationParams() (*DocumentSelect, *DocumentOmit, QueryParams[Document]) {
	return s, nil, QueryParams[Document]{}
}

type DocumentQueryBuilder struct {
	selects *DocumentSelect
	omits   *DocumentOmit
	where   []PredicateOf[Document]
	take    *int
	skip    *int
	orderBy []OrderBy[Document]
	cursor  UniquePredicate[Document]
}

func (b *DocumentQueryBuilder) Where(preds ...PredicateOf[Document]) *DocumentQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *DocumentQueryBuilder) Take(limit int) *DocumentQueryBuilder {
	b.take = &limit
	return b
}

func (b *DocumentQueryBuilder) Skip(offset int) *DocumentQueryBuilder {
	b.skip = &offset
	return b
}

func (b *DocumentQueryBuilder) OrderBy(orders ...OrderBy[Document]) *DocumentQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *DocumentQueryBuilder) Cursor(where UniquePredicate[Document]) *DocumentQueryBuilder {
	b.cursor = where
	return b
}

func (b *DocumentQueryBuilder) Select(s DocumentSelect) *DocumentQueryBuilder {
	b.selects = &s
	return b
}

func (b *DocumentQueryBuilder) Omit(o DocumentOmit) *DocumentQueryBuilder {
	b.omits = &o
	return b
}

func (b *DocumentQueryBuilder) GetRelationParams() (*DocumentSelect, *DocumentOmit, QueryParams[Document]) {
	if b == nil {
		return nil, nil, QueryParams[Document]{}
	}
	return b.selects, b.omits, QueryParams[Document]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// DocumentCreateArgs is the input argument passed to Document Create extension hooks.
//
// Fields for Document:
//
//	id        string required
//	kind      string default:
//	body      []byte required
//	updatedAt int32  required
type DocumentCreateArgs struct {
	// Data contains the model fields to insert.
	Data *DocumentCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *DocumentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// DocumentCreateManyArgs is the input argument passed to Document CreateMany extension hooks.
//
// Fields for Document:
//
//	id        string required
//	kind      string default:
//	body      []byte required
//	updatedAt int32  required
type DocumentCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*DocumentCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *DocumentCreateManyArgs) AppendData(builders ...*DocumentCreateBuilder) *DocumentCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToDocumentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// DocumentCreateManyAndReturnArgs is the input argument passed to Document CreateManyAndReturn extension hooks.
//
// Fields for Document:
//
//	id        string required
//	kind      string default:
//	body      []byte required
//	updatedAt int32  required
type DocumentCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*DocumentCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *DocumentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *DocumentCreateManyAndReturnArgs) AppendData(builders ...*DocumentCreateBuilder) *DocumentCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToDocumentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// DocumentFindUniqueArgs is the input argument passed to Document FindUnique extension hooks.
type DocumentFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Document]
	// Select specifies which scalar and relation fields to select and return.
	Select *DocumentSelect
}

func (a *DocumentFindUniqueArgs) SetWhere(unique UniquePredicate[Document], additional ...PredicateOf[Document]) *DocumentFindUniqueArgs {
	a.Where = make([]PredicateOf[Document], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// DocumentFindFirstArgs is the input argument passed to Document FindFirst extension hooks.
type DocumentFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Document]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Document]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *DocumentSelect
}

func (a *DocumentFindFirstArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentFindFirstArgs {
	a.Where = preds
	return a
}

func (a *DocumentFindFirstArgs) SetOrderBy(orders ...OrderBy[Document]) *DocumentFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *DocumentFindFirstArgs) SetCursor(cursor UniquePredicate[Document]) *DocumentFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *DocumentFindFirstArgs) SetSkip(n int) *DocumentFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *DocumentFindFirstArgs) SetTake(n int) *DocumentFindFirstArgs {
	a.Take = &n
	return a
}

// DocumentFindManyArgs is the input argument passed to Document FindMany extension hooks.
type DocumentFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Document]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Document]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *DocumentSelect
}

func (a *DocumentFindManyArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentFindManyArgs {
	a.Where = preds
	return a
}

func (a *DocumentFindManyArgs) SetOrderBy(orders ...OrderBy[Document]) *DocumentFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *DocumentFindManyArgs) SetCursor(cursor UniquePredicate[Document]) *DocumentFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *DocumentFindManyArgs) SetSkip(n int) *DocumentFindManyArgs {
	a.Skip = &n
	return a
}

func (a *DocumentFindManyArgs) SetTake(n int) *DocumentFindManyArgs {
	a.Take = &n
	return a
}

// DocumentCountArgs is the input argument passed to Document Count extension hooks.
type DocumentCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *DocumentCountArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentCountArgs {
	a.Where = preds
	return a
}

func (a *DocumentCountArgs) SetSkip(n int) *DocumentCountArgs {
	a.Skip = &n
	return a
}

func (a *DocumentCountArgs) SetTake(n int) *DocumentCountArgs {
	a.Take = &n
	return a
}

// DocumentDeleteArgs is the input argument passed to Document Delete extension hooks.
type DocumentDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Document]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *DocumentSelect
}

func (a *DocumentDeleteArgs) SetWhere(unique UniquePredicate[Document], additional ...PredicateOf[Document]) *DocumentDeleteArgs {
	a.Where = make([]PredicateOf[Document], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// DocumentDeleteManyArgs is the input argument passed to Document DeleteMany extension hooks.
type DocumentDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
}

func (a *DocumentDeleteManyArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentDeleteManyArgs {
	a.Where = preds
	return a
}

// DocumentUpdateArgs is the input argument passed to Document Update extension hooks.
type DocumentUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Document]
	// Data contains the model fields to update.
	Data *DocumentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *DocumentSelect
}

func (a *DocumentUpdateArgs) SetWhere(unique UniquePredicate[Document], additional ...PredicateOf[Document]) *DocumentUpdateArgs {
	a.Where = make([]PredicateOf[Document], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// DocumentUpdateManyArgs is the input argument passed to Document UpdateMany extension hooks.
type DocumentUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
	// Data contains the model fields to update.
	Data *DocumentUpdate
}

func (a *DocumentUpdateManyArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentUpdateManyArgs {
	a.Where = preds
	return a
}

// DocumentUpdateManyAndReturnArgs is the input argument passed to Document UpdateManyAndReturn extension hooks.
type DocumentUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Document]
	// Data contains the model fields to update.
	Data *DocumentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *DocumentSelect
}

func (a *DocumentUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Document]) *DocumentUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type DocumentCreateQuery = func(ctx context.Context, args *DocumentCreateArgs) (*Document, error)
type DocumentCreateManyQuery = func(ctx context.Context, args *DocumentCreateManyArgs) (int64, error)
type DocumentCreateManyAndReturnQuery = func(ctx context.Context, args *DocumentCreateManyAndReturnArgs) ([]*Document, error)
type DocumentFindUniqueQuery = func(ctx context.Context, args *DocumentFindUniqueArgs) (*Document, error)
type DocumentFindFirstQuery = func(ctx context.Context, args *DocumentFindFirstArgs) (*Document, error)
type DocumentFindManyQuery = func(ctx context.Context, args *DocumentFindManyArgs) ([]*Document, error)
type DocumentDeleteQuery = func(ctx context.Context, args *DocumentDeleteArgs) (*Document, error)
type DocumentDeleteManyQuery = func(ctx context.Context, args *DocumentDeleteManyArgs) (int64, error)
type DocumentCountQuery = func(ctx context.Context, args *DocumentCountArgs) (int64, error)
type DocumentUpdateQuery = func(ctx context.Context, args *DocumentUpdateArgs) (*Document, error)
type DocumentUpdateManyQuery = func(ctx context.Context, args *DocumentUpdateManyArgs) (int64, error)
type DocumentUpdateManyAndReturnQuery = func(ctx context.Context, args *DocumentUpdateManyAndReturnArgs) ([]*Document, error)

type DocumentExtension struct {
	Create              func(ctx context.Context, args *DocumentCreateArgs, next DocumentCreateQuery) (*Document, error)
	CreateMany          func(ctx context.Context, args *DocumentCreateManyArgs, next DocumentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *DocumentCreateManyAndReturnArgs, next DocumentCreateManyAndReturnQuery) ([]*Document, error)
	FindUnique          func(ctx context.Context, args *DocumentFindUniqueArgs, next DocumentFindUniqueQuery) (*Document, error)
	FindFirst           func(ctx context.Context, args *DocumentFindFirstArgs, next DocumentFindFirstQuery) (*Document, error)
	FindMany            func(ctx context.Context, args *DocumentFindManyArgs, next DocumentFindManyQuery) ([]*Document, error)
	Delete              func(ctx context.Context, args *DocumentDeleteArgs, next DocumentDeleteQuery) (*Document, error)
	DeleteMany          func(ctx context.Context, args *DocumentDeleteManyArgs, next DocumentDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *DocumentCountArgs, next DocumentCountQuery) (int64, error)
	Update              func(ctx context.Context, args *DocumentUpdateArgs, next DocumentUpdateQuery) (*Document, error)
	UpdateMany          func(ctx context.Context, args *DocumentUpdateManyArgs, next DocumentUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *DocumentUpdateManyAndReturnArgs, next DocumentUpdateManyAndReturnQuery) ([]*Document, error)
}

type DocumentDelegate struct {
	client     *Queries
	extensions []DocumentExtension
}

func (d *DocumentDelegate) Use(exts ...DocumentExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Document) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "kind":
			targets[i] = &m.Kind
		case "body":
			targets[i] = &m.Body
		case "updated_at":
			targets[i] = &m.UpdatedAt
		}
	}
	return targets
}

var documentDefaultCols = []string{
	"id",
	"kind",
	"body",
	"updated_at",
}

var documentPKCols = []string{
	"id",
	"kind",
}

var documentUniqueCols = []string{}

func selectDocumentCols(selects *DocumentSelect, omits *DocumentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return documentDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Kind || selects.Body || selects.UpdatedAt)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"kind", selects != nil && selects.Kind, omits != nil && omits.Kind, false},
		{"body", selects != nil && selects.Body, omits != nil && omits.Body, false},
		{"updated_at", selects != nil && selects.UpdatedAt, omits != nil && omits.UpdatedAt, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *DocumentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return false
}

type DocumentUpsertBuilder struct {
	*CreateBuilder[Document, DocumentSelect, DocumentOmit]
}

func (b *DocumentUpsertBuilder) Select(s DocumentSelect) *DocumentUpsertBuilder {
	b.selects = &s
	return b
}

func (b *DocumentUpsertBuilder) Omit(o DocumentOmit) *DocumentUpsertBuilder {
	b.omits = &o
	return b
}

type DocumentCreateBuilder struct {
	*CreateBuilder[Document, DocumentSelect, DocumentOmit]
}

func (b *DocumentCreateBuilder) Select(s DocumentSelect) *DocumentCreateBuilder {
	b.selects = &s
	return b
}

func (b *DocumentCreateBuilder) Omit(o DocumentOmit) *DocumentCreateBuilder {
	b.omits = &o
	return b
}

func (b *DocumentCreateBuilder) OnConflict(target UniqueConstraintTarget) *DocumentConflictBuilder[DocumentUpsertBuilder] {
	upsertBuilder := &DocumentUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &DocumentConflictBuilder[DocumentUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *DocumentCreateBuilder) SetId(v string) *DocumentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *DocumentCreateBuilder) SetKind(v string) *DocumentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "kind", Val: v})
	return b
}
func (b *DocumentCreateBuilder) SetBody(v []byte) *DocumentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "body", Val: v})
	return b
}
func (b *DocumentCreateBuilder) SetUpdatedAt(v int32) *DocumentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updated_at", Val: v})
	return b
}

func (b *DocumentCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Document]) *DocumentCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *DocumentDelegate) Create() *DocumentCreateBuilder {
	return &DocumentCreateBuilder{
		CreateBuilder: &CreateBuilder[Document, DocumentSelect, DocumentOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedDocumentId        uint64 = 1 << 0
	providedDocumentKind      uint64 = 1 << 1
	providedDocumentBody      uint64 = 1 << 2
	providedDocumentUpdatedAt uint64 = 1 << 3
)

func assignmentsToDocumentCreate(assignments []FieldAssignment) (DocumentCreate, error) {
	var input DocumentCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedDocumentId
			if v, ok := a.Val.(string); ok {
				input.Id = v
				errs.ValidateString("id", v, true, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "kind":
			provided |= providedDocumentKind
			if v, ok := a.Val.(string); ok {
				input.Kind = &v
				errs.ValidateString("kind", v, false, 0, false, false)
			} else {
				errs.Add("kind", a.Val, "type", "field kind must be of type string")
			}
		case "body":
			provided |= providedDocumentBody
			if v, ok := a.Val.([]byte); ok {
				input.Body = v
			} else {
				errs.Add("body", a.Val, "type", "field body must be of type []byte")
			}
		case "updated_at":
			provided |= providedDocumentUpdatedAt
			if v, ok := a.Val.(int32); ok {
				input.UpdatedAt = v
				errs.ValidateInt32("updatedAt", v, "")
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type int32")
			}
		}
	}
	if provided&providedDocumentId == 0 {
		errs.Add("id", "", "required", "field Id is required")
	}
	if provided&providedDocumentBody == 0 {
		errs.Add("body", nil, "required", "field Body is required")
	}
	if provided&providedDocumentUpdatedAt == 0 {
		errs.Add("updated_at", nil, "required", "field UpdatedAt is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *DocumentCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 4)
	vals = make([]any, 0, 4)
	cols = append(cols, "id")
	vals = append(vals, s.Id)
	if s.Kind != nil {
		cols = append(cols, "kind")
		vals = append(vals, *s.Kind)
	}
	cols = append(cols, "body")
	vals = append(vals, s.Body)
	cols = append(cols, "updated_at")
	vals = append(vals, s.UpdatedAt)
	return
}

func partitionDocumentInputs(dialect Dialect, inputs []*DocumentCreate) [][]*DocumentCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*DocumentCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*DocumentCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*DocumentCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*DocumentCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*DocumentCreate{inputs}
}

func (d *DocumentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *DocumentSelect, omits *DocumentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Document, error) {
	input, err := assignmentsToDocumentCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectDocumentCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *DocumentCreateArgs) (*Document, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectDocumentCols(a.Select, omits)
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
			curr = func(c context.Context, a *DocumentCreateArgs) (*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type DocumentCreateManyBuilder struct {
	*CreateManyBuilder[Document]
}

func (b *DocumentCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *DocumentConflictBuilder[DocumentCreateManyBuilder] {
	return &DocumentConflictBuilder[DocumentCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type DocumentCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Document, DocumentSelect, DocumentOmit]
}

func (b *DocumentCreateManyAndReturnBuilder) Select(s DocumentSelect) *DocumentCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *DocumentCreateManyAndReturnBuilder) Omit(o DocumentOmit) *DocumentCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *DocumentCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *DocumentConflictBuilder[DocumentCreateManyAndReturnBuilder] {
	return &DocumentConflictBuilder[DocumentCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToDocumentRecordInputs(builders []*DocumentCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *DocumentDelegate) CreateMany(builders ...*DocumentCreateBuilder) *DocumentCreateManyBuilder {
	return &DocumentCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Document]{
			records:  createBuildersToDocumentRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *DocumentDelegate) CreateManyAndReturn(builders ...*DocumentCreateBuilder) *DocumentCreateManyAndReturnBuilder {
	return &DocumentCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Document, DocumentSelect, DocumentOmit]{
			records:  createBuildersToDocumentRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToDocumentCreateInputs(records []RecordInput) ([]*DocumentCreate, error) {
	structs := make([]DocumentCreate, len(records))
	inputs := make([]*DocumentCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToDocumentCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *DocumentDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToDocumentCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &DocumentCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *DocumentCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *DocumentCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *DocumentSelect, omits *DocumentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Document, error) {
	inputs, err := recordsToDocumentCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *DocumentCreateManyAndReturnArgs) ([]*Document, error) {
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
			curr = func(c context.Context, a *DocumentCreateManyAndReturnArgs) ([]*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *DocumentSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Document, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Document
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Document.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Document.loadRelations(ctx, []*Document{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Document", cols, returningCols, documentPKCols, conflictTarget, conflictAction, len(vals))
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

		var res Document
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, documentPKCols)
}

func (d *DocumentDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Document, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("documents") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "documents")
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

	var res Document
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *DocumentDelegate) buildBulkInsertSQL(q *Queries, batch []*DocumentCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 4)
	for i, c := range documentDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "documents")
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
				vals = append(vals, input.Id)
			case "kind":
				if input.Kind != nil {
					vals = append(vals, *input.Kind)
				} else {
					writeDefault = true
				}
			case "body":
				vals = append(vals, input.Body)
			case "updated_at":
				vals = append(vals, input.UpdatedAt)
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

func applyDocumentConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanDocumentRows(rows *sql.Rows, returningCols []string) ([]*Document, error) {
	var records []*Document
	for rows.Next() {
		var res Document
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

func (d *DocumentDelegate) runCreateMany(ctx context.Context, inputs []*DocumentCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionDocumentInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyDocumentConflictClause(d.client.dialect, queryStr, vals, cols, documentPKCols, conflictTarget, conflictAction)

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

func (d *DocumentDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*DocumentCreate,
	selects *DocumentSelect,
	omits *DocumentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Document, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionDocumentInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Document
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Document.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Document.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectDocumentCols(selects, omits, documentPKCols...)
	recordsOut := make([]*Document, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyDocumentConflictClause(d.client.dialect, queryStr, vals, cols, documentPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanDocumentRows(rows, returningCols)
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

func (d *DocumentDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*DocumentCreate,
	selects *DocumentSelect,
	omits *DocumentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Document, error) {
	batches := partitionDocumentInputs(d.client.dialect, inputs)
	returningCols := selectDocumentCols(selects, omits, documentPKCols...)
	recordsOut := make([]*Document, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyDocumentConflictClause(d.client.dialect, queryStr, vals, cols, documentPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("documents") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "documents")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, documentPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, documentPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanDocumentRows(rows, returningCols)
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

type DocumentConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *DocumentConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *DocumentConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *DocumentConflictBuilder[B]) Update(fn func(u *DocumentUpsert)) *B {
	cb.setAction(*DocumentConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// DocumentConflictUpdate creates a custom ConflictAction for Document upsert conflicts.
func DocumentConflictUpdate(fn func(u *DocumentUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newDocumentUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type DocumentUpsert struct {
	Id        fieldUpsert[string]
	Kind      fieldUpsert[string]
	Body      fieldUpsert[[]byte]
	UpdatedAt numericFieldUpsert[int32]
}

func newDocumentUpsert(up *ConflictUpdate) *DocumentUpsert {
	return &DocumentUpsert{
		Id:   fieldUpsert[string]{column: "id", update: up},
		Kind: fieldUpsert[string]{column: "kind", update: up},
		Body: fieldUpsert[[]byte]{column: "body", update: up},
		UpdatedAt: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "updated_at", update: up},
			tableName:   "documents",
		},
	}
}

type DocumentUpdateBuilder struct {
	*UpdateBuilder[Document, DocumentSelect, DocumentOmit]
}

type DocumentUpdateManyBuilder struct {
	*UpdateManyBuilder[Document]
}

type DocumentUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Document, DocumentSelect, DocumentOmit]
}

func (b *DocumentUpdateBuilder) SetId(v string) *DocumentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *DocumentUpdateManyBuilder) SetId(v string) *DocumentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *DocumentUpdateManyAndReturnBuilder) SetId(v string) *DocumentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *DocumentUpdateBuilder) SetKind(v string) *DocumentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "kind", Val: v})
	return b
}

func (b *DocumentUpdateManyBuilder) SetKind(v string) *DocumentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "kind", Val: v})
	return b
}

func (b *DocumentUpdateManyAndReturnBuilder) SetKind(v string) *DocumentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "kind", Val: v})
	return b
}
func (b *DocumentUpdateBuilder) SetBody(v []byte) *DocumentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "body", Val: v})
	return b
}

func (b *DocumentUpdateManyBuilder) SetBody(v []byte) *DocumentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "body", Val: v})
	return b
}

func (b *DocumentUpdateManyAndReturnBuilder) SetBody(v []byte) *DocumentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "body", Val: v})
	return b
}
func (b *DocumentUpdateBuilder) SetUpdatedAt(v int32) *DocumentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updated_at", Val: v})
	return b
}

func (b *DocumentUpdateManyBuilder) SetUpdatedAt(v int32) *DocumentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updated_at", Val: v})
	return b
}

func (b *DocumentUpdateManyAndReturnBuilder) SetUpdatedAt(v int32) *DocumentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updated_at", Val: v})
	return b
}

func (b *DocumentUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Document]) *DocumentUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *DocumentUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Document]) *DocumentUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *DocumentUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Document]) *DocumentUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *DocumentDelegate) Update(where UniquePredicate[Document], additional ...PredicateOf[Document]) *DocumentUpdateBuilder {
	return &DocumentUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Document, DocumentSelect, DocumentOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *DocumentDelegate) UpdateMany(preds ...PredicateOf[Document]) *DocumentUpdateManyBuilder {
	return &DocumentUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Document]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *DocumentDelegate) UpdateManyAndReturn(preds ...PredicateOf[Document]) *DocumentUpdateManyAndReturnBuilder {
	return &DocumentUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Document, DocumentSelect, DocumentOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *DocumentDelegate) buildUpdateSQL(preds []PredicateOf[Document], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "documents")
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

func (d *DocumentDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Document], additional []PredicateOf[Document], assignments []FieldAssignment, selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
	allWhere := make([]PredicateOf[Document], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToDocumentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *DocumentUpdateArgs) (*Document, error) {
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
			curr = func(c context.Context, a *DocumentUpdateArgs) (*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Document], cols []string, vals []any, selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
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
		var res *Document
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Document.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Document.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectDocumentCols(selects, omits, documentPKCols...)
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
		return nil, &NotFoundError{Model: "Document"}
	}

	var res Document
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Document{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *DocumentDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Document], cols []string, vals []any) (int64, error) {
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

func (d *DocumentDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Document], cols []string, vals []any, selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Document"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *DocumentDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Document], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToDocumentUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &DocumentUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *DocumentUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *DocumentUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *DocumentDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Document], assignments []FieldAssignment, selects *DocumentSelect, omits *DocumentOmit) ([]*Document, error) {
	input, err := assignmentsToDocumentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *DocumentUpdateManyAndReturnArgs) ([]*Document, error) {
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
			curr = func(c context.Context, a *DocumentUpdateManyAndReturnArgs) ([]*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Document], cols []string, vals []any, selects *DocumentSelect, omits *DocumentOmit) ([]*Document, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Document]{Where: preds}, selects, omits)
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
		var res []*Document
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Document.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Document.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectDocumentCols(selects, omits, documentPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanDocumentRows(rows, returningCols)
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

func (d *DocumentDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Document], cols []string, vals []any, selects *DocumentSelect, omits *DocumentOmit) ([]*Document, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Document{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Document]{Where: preds}, selects, omits)
}
func (d *DocumentDelegate) FindUnique(where UniquePredicate[Document], additional ...PredicateOf[Document]) *FindUniqueBuilder[Document, DocumentSelect, DocumentOmit] {
	return &FindUniqueBuilder[Document, DocumentSelect, DocumentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *DocumentDelegate) FindFirst(preds ...PredicateOf[Document]) *FindFirstBuilder[Document, DocumentSelect, DocumentOmit] {
	return &FindFirstBuilder[Document, DocumentSelect, DocumentOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *DocumentDelegate) FindMany(preds ...PredicateOf[Document]) *FindManyBuilder[Document, DocumentSelect, DocumentOmit] {
	return &FindManyBuilder[Document, DocumentSelect, DocumentOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *DocumentDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Document], additional []PredicateOf[Document], selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
	allWhere := make([]PredicateOf[Document], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *DocumentFindUniqueArgs) (*Document, error) {
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
			curr = func(c context.Context, a *DocumentFindUniqueArgs) (*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Document],
	selects *DocumentSelect,
	omits *DocumentOmit,
) (*Document, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *DocumentFindFirstArgs) (*Document, error) {
		return d.runFindFirst(c, QueryParams[Document]{
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
			curr = func(c context.Context, a *DocumentFindFirstArgs) (*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Document],
	selects *DocumentSelect,
	omits *DocumentOmit,
) ([]*Document, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *DocumentFindManyArgs) ([]*Document, error) {
		return d.runFindMany(c, QueryParams[Document]{
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
			curr = func(c context.Context, a *DocumentFindManyArgs) ([]*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Document], selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
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
	returningCols := selectDocumentCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Document{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *DocumentDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Document],
	selects *DocumentSelect,
	omits *DocumentOmit,
) (*Document, error) {
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

func (d *DocumentDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Document],
	selects *DocumentSelect,
	omits *DocumentOmit,
) ([]*Document, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, documentPKCols, documentUniqueCols, "documents", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, documentPKCols, documentUniqueCols, isCursorQuery, params.Take)
	returningCols := selectDocumentCols(selects, omits)

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

func (d *DocumentDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Document, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "documents", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Document
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *DocumentDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Document, error) {
	query := buildSelectSQL(d.client, "documents", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Document, 0)
	for rows.Next() {
		var res Document
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
func (d *DocumentDelegate) DeleteMany(preds ...PredicateOf[Document]) *DeleteManyBuilder[Document] {
	return &DeleteManyBuilder[Document]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *DocumentDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Document]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &DocumentDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *DocumentDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *DocumentDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Document]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "documents")
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

func (d *DocumentDelegate) Delete(where UniquePredicate[Document], additional ...PredicateOf[Document]) *DeleteBuilder[Document, DocumentSelect, DocumentOmit] {
	return &DeleteBuilder[Document, DocumentSelect, DocumentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *DocumentDelegate) executeDelete(ctx context.Context, where UniquePredicate[Document], additional []PredicateOf[Document], selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
	allWhere := make([]PredicateOf[Document], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullDocumentSelect()
	}

	args := &DocumentDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *DocumentDeleteArgs) (*Document, error) {
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
			curr = func(c context.Context, a *DocumentDeleteArgs) (*Document, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runDelete(ctx context.Context, where []PredicateOf[Document], selects *DocumentSelect, omits *DocumentOmit) (*Document, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectDocumentCols(selects, omits, documentPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Document
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Document.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Document"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "documents")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Document]
			pkPreds = append(pkPreds, Predicate[Document]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
				},
			})
			pkPreds = append(pkPreds, Predicate[Document]{
				Data: PredicateData{
					Column:   "kind",
					Operator: "=",
					Value:    res.Kind,
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
	d.client.dialect.WriteQuotedIdent(&sb, "documents")

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
		return nil, &NotFoundError{Model: "Document"}
	}

	var row Document
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *DocumentDelegate) Count(preds ...PredicateOf[Document]) *CountBuilder[Document] {
	return &CountBuilder[Document]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *DocumentDelegate) executeCount(ctx context.Context, params QueryParams[Document]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &DocumentCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *DocumentCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Document]{
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
			curr = func(c context.Context, a *DocumentCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *DocumentDelegate) runCount(ctx context.Context, params QueryParams[Document]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "documents")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "documents")
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
func (d *DocumentDelegate) loadRelations(ctx context.Context, records []*Document, selects *DocumentSelect) error {
	_ = ctx
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
