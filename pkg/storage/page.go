package storage

func PageStorage(s Storage, pageID string) Storage {
	return pageStorage{s: s, pageID: pageID}
}

type pageStorage struct {
	s      Storage
	pageID string
}

func (p pageStorage) makePageKey(key string) string {
	return p.pageID + "||" + key
}

func (p pageStorage) Get(key string) interface{} {
	key = p.makePageKey(key)
	return p.s.Get(key)
}

func (p pageStorage) Set(key string, value interface{}) {
	key = p.makePageKey(key)
	p.s.Set(key, value)
}
