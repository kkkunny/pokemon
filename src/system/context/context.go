package context

import (
	"github.com/kkkunny/pokemon/src/util/i18n"
)

type Context interface {
	Localisation() *i18n.Localisation
}

type _Context struct {
	loc *i18n.Localisation
}

func NewContext(loc *i18n.Localisation) Context {
	return &_Context{
		loc: loc,
	}
}

func (ctx *_Context) Localisation() *i18n.Localisation {
	return ctx.loc
}
