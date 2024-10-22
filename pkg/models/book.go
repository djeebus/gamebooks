package models

type Book struct {
	ID   string
	Name string
	Path string

	LuaPath   string
	StartPage string
}

type BookResult struct {
	Functions map[string]interface{}
}
