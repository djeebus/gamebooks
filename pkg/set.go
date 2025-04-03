package pkg

type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Remove(key string) {
	delete(s, key)
}

func (s Set) Without(other Set) Set {
	out := NewSet()
	for k := range s {
		out.Add(k)
	}
	for k := range other {
		out.Remove(k)
	}

	return out
}
