package gomodel

import (
	"net/http"
	"reflect"
	"slices"
	"sync"

	"github.com/dimonrus/porterr"
)

// Wizard complete all model fields step by step
type Wizard[T any] struct {
	// internal mutex
	m sync.RWMutex
	// Filed order
	columns []string
	// Pointer for values
	values []any
	// model for wizard
	model *T
	// current field position
	// 16-9 bits already completed fields
	// 8-1 bits current position
	cursor uint16
}

// Next increment field cursor
func (w *Wizard[T]) Next() bool {
	w.m.Lock()
	hasNext := w.next()
	w.m.Unlock()
	return hasNext
}

// Next increment field cursor
func (w *Wizard[T]) next() bool {
	next := w.cursor&255 + 1
	left := w.cursor>>8 | next
	w.cursor = left<<8 | next
	return int(w.cursor>>8) < len(w.values)
}

// Set model value
func (w *Wizard[T]) Set(value any) (next bool, e porterr.IError) {
	w.m.Lock()
	defer w.m.Unlock()
	if int(w.cursor&255) >= len(w.columns) {
		e = porterr.NewF(porterr.PortErrorType, "can't set values more then columns initiated: %v", len(w.columns)).HTTP(http.StatusBadRequest)
		return
	}
	if reflect.TypeOf(w.values[w.cursor&255]) != reflect.TypeOf(value) {
		e = porterr.NewF(porterr.PortErrorType, "type of value is not: %T", w.values[w.cursor&255]).HTTP(http.StatusBadRequest)
		return
	}
	w.values[w.cursor&255] = value
	next = w.next()
	return
}

// Get model
func (w *Wizard[T]) Get() *T {
	w.m.RLock()
	defer w.m.RUnlock()
	var model = new(T)
	iModel := interface{}(model).(IModel)
	meta := PrepareMetaModel(iModel)
	if meta == nil {
		return nil
	}
	for i := range meta.Fields {
		for j, column := range w.columns {
			if meta.Fields[i].Column == column {
				//*(meta.Fields[i].Value).(*T) = *w.values[j].(*T)
				va := reflect.ValueOf(w.values[j])
				reflect.ValueOf(meta.Fields[i].Value).Elem().Set(va.Elem())
				break
			}
		}
	}
	return model
}

// SetByName set model value by name
func (w *Wizard[T]) SetByName(field string, value any) (e porterr.IError) {
	w.m.Lock()
	defer w.m.Unlock()
	index := slices.Index(w.columns, field)
	if index < 0 {
		e = porterr.New(porterr.PortErrorIteration, "no such field found: "+field).HTTP(http.StatusBadRequest)
		return
	}
	w.values[index] = value
	left := w.cursor>>8 | uint16(index)
	w.cursor = left<<8 | uint16(index)
	return
}

// IsComplete check if all values filled
func (w *Wizard[T]) IsComplete() bool {
	w.m.RLock()
	isComplete := int(w.cursor>>8) >= len(w.values)
	w.m.RUnlock()
	return isComplete
}

// ResetCursor reset cursor
func (w *Wizard[T]) ResetCursor() {
	w.m.Lock()
	w.cursor = w.cursor >> 8 << 8
	w.m.Unlock()
	return
}

// GetCursor get cursor
func (w *Wizard[T]) GetCursor() uint16 {
	w.m.RLock()
	cursor := w.cursor & 255
	w.m.RUnlock()
	return cursor
}

// NewWizard init wizard object
// order - provide names of field should be completed step by step
func NewWizard[T any](order ...string) *Wizard[T] {
	var model = new(T)
	iModel := interface{}(model).(IModel)
	columns := (iModel).(IModel).Columns()
	values := (iModel).(IModel).Values()
	if len(order) == 0 {
		order = columns
	}
	wizard := &Wizard[T]{columns: order, model: model}
	for _, s := range order {
		for i, column := range columns {
			if s == column {
				wizard.values = append(wizard.values, values[i])
				break
			}
		}
	}
	return wizard
}
