package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

// forked from upstream
func luaOpenLibraries(l *lua.State, book *models.Book, preloaded ...lua.RegistryFunction) {
	libs := []lua.RegistryFunction{
		{"_G", lua.BaseOpen},
		{"package", luaPackageOpen(book.Path)}, // changed!
		// {"coroutine", CoroutineOpen},
		{"table", lua.TableOpen},
		{"io", lua.IOOpen},
		{"os", lua.OSOpen},
		{"string", lua.StringOpen},
		{"bit32", lua.Bit32Open},
		{"math", lua.MathOpen},
		{"debug", lua.DebugOpen},
	}
	for _, lib := range libs {
		lua.Require(l, lib.Name, lib.Function, true)
		l.Pop(1)
	}
	lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")
	for _, lib := range preloaded {
		l.PushGoFunction(lib.Function)
		l.SetField(-2, lib.Name)
	}
	l.Pop(1)
}

var packageLibrary = []lua.RegistryFunction{
	{"loadlib", func(l *lua.State) int {
		_ = lua.CheckString(l, 1) // path
		_ = lua.CheckString(l, 2) // init
		l.PushNil()
		l.PushString("dynamic libraries not enabled; check your Lua installation")
		l.PushString("absent")
		return 3 // Return nil, error message, and where.
	}},
	{"searchpath", func(l *lua.State) int {
		name := lua.CheckString(l, 1)
		path := lua.CheckString(l, 2)
		sep := lua.OptString(l, 3, ".")
		dirSep := lua.OptString(l, 4, string(filepath.Separator))
		f, err := searchPath(l, name, path, sep, dirSep)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		l.PushString(f)
		return 1
	}},
}

// PackageOpen opens the package library. Usually passed to Require.
func luaPackageOpen(luaRoot string) lua.Function {
	return func(l *lua.State) int {
		lua.NewLibrary(l, packageLibrary)
		createSearchersTable(l)
		l.SetField(-2, "searchers")
		setPath(l, "path", "LUA_PATH", luaRoot+"/?.lua")
		l.PushString(fmt.Sprintf("%c\n%c\n?\n!\n-\n", filepath.Separator, pathListSeparator))
		l.SetField(-2, "config")
		lua.SubTable(l, lua.RegistryIndex, "_LOADED")
		l.SetField(-2, "loaded")
		lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")
		l.SetField(-2, "preload")
		l.PushGlobalTable()
		l.PushValue(-2)
		lua.SetFunctions(l, []lua.RegistryFunction{{"require", func(l *lua.State) int {
			name := lua.CheckString(l, 1)
			l.SetTop(1)
			l.Field(lua.RegistryIndex, "_LOADED")
			l.Field(2, name)
			if l.ToBoolean(-1) {
				return 1
			}
			l.Pop(1)
			findLoader(l, name)
			l.PushString(name)
			l.Insert(-2)
			l.Call(2, 1)
			if !l.IsNil(-1) {
				l.SetField(2, name)
			}
			l.Field(2, name)
			if l.IsNil(-1) {
				l.PushBoolean(true)
				l.PushValue(-1)
				l.SetField(2, name)
			}
			return 1
		}}}, 1)
		l.Pop(1)
		return 1
	}
}

var defaultPath = "./?.lua" // TODO "${LUA_LDIR}?.lua;${LUA_LDIR}?/init.lua;./?.lua"

const pathListSeparator = ';'

func noEnv(l *lua.State) bool {
	l.Field(lua.RegistryIndex, "LUA_NOENV")
	b := l.ToBoolean(-1)
	l.Pop(1)
	return b
}

func setPath(l *lua.State, field, env, def string) {
	if path := os.Getenv(env); path == "" || noEnv(l) {
		l.PushString(def)
	} else {
		o := fmt.Sprintf("%c%c", pathListSeparator, pathListSeparator)
		n := fmt.Sprintf("%c%s%c", pathListSeparator, def, pathListSeparator)
		path = strings.Replace(path, o, n, -1)
		l.PushString(path)
	}
	l.SetField(-2, field)
}

func findLoader(l *lua.State, name string) {
	var msg string
	if l.Field(lua.UpValueIndex(1), "searchers"); !l.IsTable(3) {
		lua.Errorf(l, "'package.searchers' must be a table")
	}
	for i := 1; ; i++ {
		if l.RawGetInt(3, i); l.IsNil(-1) {
			l.Pop(1)
			l.PushString(msg)
			lua.Errorf(l, "module '%s' not found: %s", name, msg)
		}
		l.PushString(name)
		if l.Call(1, 2); l.IsFunction(-2) {
			return
		} else if l.IsString(-2) {
			msg += lua.CheckString(l, -2)
		}
		l.Pop(2)
	}
}

func createSearchersTable(l *lua.State) {
	searchers := []lua.Function{searcherPreload, searcherLua}
	l.CreateTable(len(searchers), 0)
	for i, s := range searchers {
		l.PushValue(-2)
		l.PushGoClosure(s, 1)
		l.RawSetInt(-2, i+1)
	}
}

func searcherPreload(l *lua.State) int {
	name := lua.CheckString(l, 1)
	l.Field(lua.RegistryIndex, "_PRELOAD")
	l.Field(-1, name)
	if l.IsNil(-1) {
		l.PushString(fmt.Sprintf("\n\tno field package.preload['%s']", name))
	}
	return 1
}

func searcherLua(l *lua.State) int {
	name := lua.CheckString(l, 1)
	filename, err := findFile(l, name, "path", string(filepath.Separator))
	if err != nil {
		return 1 // Module not found in this path.
	}
	return checkLoad(l, lua.LoadFile(l, filename, "") == nil, filename)
}

func findFile(l *lua.State, name, field, dirSep string) (string, error) {
	l.Field(lua.UpValueIndex(1), field)
	path, ok := l.ToString(-1)
	if !ok {
		lua.Errorf(l, "'package.%s' must be a string", field)
	}
	return searchPath(l, name, path, ".", dirSep)
}

func checkLoad(l *lua.State, loaded bool, fileName string) int {
	if loaded { // Module loaded successfully?
		l.PushString(fileName) // Second argument to module.
		return 2               // Return open function & file name.
	}
	m := lua.CheckString(l, 1)
	e := lua.CheckString(l, -1)
	lua.Errorf(l, "error loading module '%s' from file '%s':\n\t%s", m, fileName, e)
	panic("unreachable")
}

func searchPath(l *lua.State, name, path, sep, dirSep string) (string, error) {
	var msg string
	if sep != "" {
		name = strings.Replace(name, sep, dirSep, -1) // Replace sep by dirSep.
	}
	path = strings.Replace(path, string(pathListSeparator), string(filepath.ListSeparator), -1)
	for _, template := range filepath.SplitList(path) {
		if template != "" {
			filename := strings.Replace(template, "?", name, -1)
			if readable(filename) {
				return filename, nil
			}
			msg = fmt.Sprintf("%s\n\tno file '%s'", msg, filename)
		}
	}
	return "", errors.New(msg)
}

func readable(filename string) bool {
	f, err := os.Open(filename)
	if f != nil {
		f.Close()
	}
	return err == nil
}
