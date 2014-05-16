package srv

import "github.com/PuerkitoBio/gred/vals"

type Key interface {
	RWLocker
	Expirer

	Val() vals.Value
	Name() string
}
