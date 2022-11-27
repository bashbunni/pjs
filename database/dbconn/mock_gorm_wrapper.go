package dbconn

import (
	"errors"
	"reflect"
)

type MockGormWrapper interface {
	GormWrapper
	Created() []interface{}
	Chain() *queryChain
	SetAutoMigrateError(error) MockGormWrapper
	SetError(error) MockGormWrapper
	SetResult(interface{}) MockGormWrapper
}

type mockGormWrapper struct {
	automigrateError         error
	error                    error
	triedToFind, triedToSave interface{}
	created                  []interface{}
	deleted                  []interface{}
	chain                    *queryChain
	result                   interface{}
}

type queryChain struct {
	Where whereQuery
}

type whereQuery struct {
	Query interface{}
	Args  []interface{}
	First firstSelect
}

type firstSelect struct {
	Conds []interface{}
}

func Mock() MockGormWrapper {
	return &mockGormWrapper{}
}

func (w *mockGormWrapper) Created() []interface{} {
	return w.created
}

func (w *mockGormWrapper) Chain() *queryChain {
	return w.chain
}

func (w *mockGormWrapper) SetAutoMigrateError(e error) MockGormWrapper {
	w.automigrateError = e
	return w
}

func (w *mockGormWrapper) SetError(e error) MockGormWrapper {
	w.error = e
	return w
}

func (w *mockGormWrapper) SetResult(r interface{}) MockGormWrapper {
	w.result = r
	return w
}

func (w *mockGormWrapper) Error() error {
	return w.error
}

func (w *mockGormWrapper) AutoMigrate(...interface{}) error {
	return w.automigrateError
}

func (w *mockGormWrapper) Create(value interface{}) GormWrapper {
	if w.error == nil {
		w.created = append(w.created, value)
	}
	return w
}

func (w *mockGormWrapper) Delete(value interface{}, args ...interface{}) GormWrapper {
	if w.error == nil {
		w.deleted = append(w.deleted, value)
	}
	return w
}

func (w *mockGormWrapper) Where(query interface{}, args ...interface{}) GormWrapper {
	w.chain = &queryChain{
		Where: whereQuery{
			Query: query,
			Args:  args,
		},
	}
	return w
}

type MissingQuery error

type err struct {
	i error
}

func (e err) Error() string { return e.i.Error() }

var ErrMissingQuery = MissingQuery(errors.New("need to call query first"))

func (w *mockGormWrapper) First(dest interface{}, conds ...interface{}) GormWrapper {
	if w.chain == nil {
		w.error = ErrMissingQuery
		return w
	}

	w.chain.Where.First = firstSelect{conds}
	err := Replace(dest, w.result)
	if w.error == nil {
		w.error = err
	}

	return w
}

func (w *mockGormWrapper) Find(dest interface{}, conds ...interface{}) GormWrapper {
	w.triedToFind = dest
	return w
}

func (w *mockGormWrapper) Save(value interface{}) GormWrapper {
	w.triedToSave = value
	return w
}

func Replace(i, v interface{}) error {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}

	val = val.Elem()

	newVal := reflect.Indirect(reflect.ValueOf(v))

	if !val.Type().AssignableTo(newVal.Type()) {
		return errors.New("mismatched types")
	}

	val.Set(newVal)
	return nil
}
