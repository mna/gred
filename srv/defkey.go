package srv

import (
	"time"

	"github.com/PuerkitoBio/gred/vals"
)

var dv defVal

var (
	_ Key         = (*defKey)(nil)
	_ vals.String = (*defVal)(nil)
)

type defKey string

func (d defKey) Lock()                             {}
func (d defKey) Unlock()                           {}
func (d defKey) RLock()                            {}
func (d defKey) RUnlock()                          {}
func (d defKey) Expire(_ time.Time, _ func()) bool { return true }
func (d defKey) TTL() time.Duration                { return 0 }
func (d defKey) Abort() bool                       { return true }
func (d defKey) Val() vals.Value                   { return dv }

func (d defKey) Name() string { return string(d) }

type defVal struct{}

func (d defVal) Type() string { panic("Type called on defKey value") }

// String implementation
func (d defVal) Append(_ string) int64      { return 0 }
func (d defVal) Get() string                { return "" }
func (d defVal) GetRange(_, _ int64) string { return "" }
func (d defVal) GetSet(_ string) string     { return "" }
func (d defVal) Set(_ string)               {}
func (d defVal) StrLen() int64              { return 0 }
