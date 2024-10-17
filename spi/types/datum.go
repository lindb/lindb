package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/samber/lo"
)

type Datum struct {
	val any
}

func MakeDatums(vals ...any) []*Datum {
	return lo.Map(vals, func(v any, _ int) *Datum {
		return &Datum{val: v}
	})
}

func (d *Datum) String() string {
	if str, ok := d.val.(string); ok {
		return str
	}
	return fmt.Sprint(d.val)
}

func (d *Datum) Float() float64 {
	v := reflect.ValueOf(d.val)
	fmt.Printf("datum v=%v,type=%v\n", d.val, v.Kind())
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	default:
		return 0 // TODO:set 0?
	}
}

func (d *Datum) Int() int64 {
	v := reflect.ValueOf(d.val)
	fmt.Printf("datum v=%v,type=%v\n", d.val, v.Kind())
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Int()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	default:
		return 0 // TODO:set 0?
	}
}

func (d *Datum) Duration() time.Duration {
	return d.val.(time.Duration)
}
