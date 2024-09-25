package types

import (
	"fmt"

	"github.com/samber/lo"
)

type Datum struct {
	val any
}

func (d *Datum) String() string {
	if str, ok := d.val.(string); ok {
		return str
	}
	return fmt.Sprint(d.val)
}

func MakeDatums(vals ...any) []*Datum {
	return lo.Map(vals, func(v any, _ int) *Datum {
		return &Datum{val: v}
	})
}
