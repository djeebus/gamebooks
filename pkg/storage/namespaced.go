package storage

func NamespacedStorage(s Storage, namespace string) Storage {
	return pageStorage{s: s, namespace: namespace}
}

type pageStorage struct {
	s         Storage
	namespace string
}

func (p pageStorage) makePageKey(key string) string {
	return p.namespace + "||" + key
}

func (p pageStorage) Get(key string) interface{} {
	key = p.makePageKey(key)
	return p.s.Get(key)
}

func (p pageStorage) Set(key string, value interface{}) {
	key = p.makePageKey(key)
	p.s.Set(key, value)
}
