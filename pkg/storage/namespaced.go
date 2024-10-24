package storage

func NamespacedStorage(s Storage, namespace string) Storage {
	return namespacedStorage{s: s, namespace: namespace}
}

type namespacedStorage struct {
	s         Storage
	namespace string
}

func (p namespacedStorage) Remove(key string) {
	key = p.makePageKey(key)
	p.s.Remove(key)
}

func (p namespacedStorage) makePageKey(key string) string {
	return p.namespace + "||" + key
}

func (p namespacedStorage) Get(key string) interface{} {
	key = p.makePageKey(key)
	return p.s.Get(key)
}

func (p namespacedStorage) Set(key string, value interface{}) {
	key = p.makePageKey(key)
	p.s.Set(key, value)
}
