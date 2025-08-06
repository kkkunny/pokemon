package i18n

import stlslices "github.com/kkkunny/stl/container/slices"

type Localisation struct {
	cache map[string]string
}

func NewLocalisation() *Localisation {
	return &Localisation{cache: make(map[string]string)}
}

func (loc *Localisation) Add(k string, v string) {
	loc.cache[k] = v
}

func (loc *Localisation) MultiAdd(kvs map[string]string) {
	for k, v := range kvs {
		loc.cache[k] = v
	}
}

func (loc *Localisation) Get(key string, defaultValue ...string) string {
	s, ok := loc.cache[key]
	if ok {
		return s
	}
	return stlslices.Last(defaultValue)
}
