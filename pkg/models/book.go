package models

type Book struct {
	ID   string
	Path string

	LuaPath string
}

type BookResult struct {
	Name      string
	StartPage string
}
