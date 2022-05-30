package gomodel

import (
	"github.com/lib/pq"
	"time"
)

var (
	ACMId        = 10
	ACMName      = "Foo"
	ACMPages     = []string{"one", "two", "three"}
	ACMSomeInt   = 100500
	ACMUpdatedAt = time.Now()
	ACMDeletedAt = time.Now()
)

type InsertModel1 struct {
	Id        *int       `json:"id" db:"col~id;prk;req;seq;"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *InsertModel1) Table() string { return "test_model" }

// Model columns
func (m *InsertModel1) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *InsertModel1) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type InsertModel2 struct {
	Id        *int       `json:"id" db:"col~id;unq;req;seq;"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *InsertModel2) Table() string { return "test_model" }

// Model columns
func (m *InsertModel2) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *InsertModel2) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpsertModel1 struct {
	Id        *int       `json:"id" db:"col~id;req;seq;"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpsertModel1) Table() string { return "test_model" }

// Model columns
func (m *UpsertModel1) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpsertModel1) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpsertModel2 struct {
	Id        *int       `json:"id" db:"col~id;prk;req"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpsertModel2) Table() string { return "test_model" }

// Model columns
func (m *UpsertModel2) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpsertModel2) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpsertModel3 struct {
	Id        *int       `json:"id" db:"col~id;prk;req"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpsertModel3) Table() string { return "test_model" }

// Model columns
func (m *UpsertModel3) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpsertModel3) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpsertModel4 struct {
	Id        *int       `json:"id" db:"col~id;req"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpsertModel4) Table() string { return "test_model" }

// Model columns
func (m *UpsertModel4) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpsertModel4) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpsertModel5 struct {
	Id        *int       `json:"id" db:"col~id;prk;req"`
	ComplexId *int       `json:"complexId" db:"col~complex_id;prk;req"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpsertModel5) Table() string { return "test_model" }

// Model columns
func (m *UpsertModel5) Columns() []string {
	return []string{"id", "complex_id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpsertModel5) Values() []any {
	return []any{&m.Id, &m.ComplexId, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpdateModel1 struct {
	Id        *int       `json:"id" db:"col~id;prk;req;seq"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpdateModel1) Table() string { return "test_model" }

// Model columns
func (m *UpdateModel1) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpdateModel1) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type UpdateModel2 struct {
	Id        *int       `json:"id" db:"col~id;req;seq"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;seq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
	DeletedAt *time.Time `json:"deletedAt" db:"col~deleted_at;dat;"`
}

// Model table name
func (m *UpdateModel2) Table() string { return "test_model" }

// Model columns
func (m *UpdateModel2) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *UpdateModel2) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

type DeleteModel1 struct {
	Id        *int       `json:"id" db:"col~id;prk;req;seq"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;seq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
}

// Model table name
func (m *DeleteModel1) Table() string { return "test_model" }

// Model columns
func (m *DeleteModel1) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at"}
}

// Model values
func (m *DeleteModel1) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt}
}

type DeleteModel2 struct {
	Id        *int       `json:"id" db:"col~id;req;seq"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;unq;seq;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	UpdatedAt *time.Time `json:"updatedAt" db:"col~updated_at;uat;"`
}

// Model table name
func (m *DeleteModel2) Table() string { return "test_model" }

// Model columns
func (m *DeleteModel2) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "updated_at"}
}

// Model values
func (m *DeleteModel2) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.UpdatedAt}
}
