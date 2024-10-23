package executor

import (
	"gamebooks/pkg/models"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
)

type bookResult struct {
	state *lua.State

	name, startPage string
}

func newBookResult(l *lua.State) (bookResult, error) {
	var (
		ok bool
		r  = bookResult{state: l}
	)

	l.Field(-1, "name")
	r.name, ok = l.ToString(-1)
	if !ok {
		return r, errors.Wrap(ErrNoField, "failed to load name")
	}
	l.Pop(1)

	l.Field(-1, "start_page")
	r.startPage, ok = l.ToString(-1)
	if !ok {
		return r, errors.Wrap(ErrNoField, "failed to load start_page")
	}
	l.Pop(1)

	return r, nil
}

func (b bookResult) GetName() string {
	return b.name
}

func (b bookResult) GetStartPage() string {
	return b.startPage
}

func (b bookResult) OnStart() error {
	b.state.Field(-1, "on_start")

	if !b.state.IsFunction(-1) {
		return errors.New("on_start is not a function")
	}

	b.state.Call(0, 0)
	return nil
}

func (b bookResult) OnPage(page *models.PageResult) error {
	b.state.Field(-1, "on_page")

	if !b.state.IsFunction(-1) {
		return errors.New("on_page is not a function")
	}

	b.state.NewTable()
	b.state.PushString(page.Title)
	b.state.SetField(-2, "title")
	b.state.PushString(page.Markdown)
	b.state.SetField(-2, "markdown")

	b.state.Call(1, 1)

	var ok bool

	b.state.Field(-1, "title")
	page.Title, ok = b.state.ToString(-1)
	if !ok {
		return errors.Wrap(ErrNoField, "failed to find title")
	}
	b.state.Pop(1)

	b.state.Field(-1, "markdown")
	page.Markdown, ok = b.state.ToString(-1)
	if !ok {
		return errors.Wrap(ErrNoField, "failed to find markdown")
	}
	b.state.Pop(1)

	return nil
}

var _ models.BookResult = new(bookResult)
