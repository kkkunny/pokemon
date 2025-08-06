package context

import (
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/util/i18n"
)

type Context interface {
	Config() *config.Config
	Localisation() *i18n.Localisation
}

type _Context struct {
	cfg *config.Config
	loc *i18n.Localisation
}

func NewContext(cfg *config.Config, loc *i18n.Localisation) Context {
	return &_Context{
		cfg: cfg,
		loc: loc,
	}
}

func (ctx *_Context) Config() *config.Config {
	return ctx.cfg
}

func (ctx *_Context) Localisation() *i18n.Localisation {
	return ctx.loc
}
