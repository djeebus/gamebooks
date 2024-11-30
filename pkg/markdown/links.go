package markdown

import (
	"gamebooks/pkg/models"
	bookRepo "gamebooks/pkg/repo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"path/filepath"
)

type LinkTracker struct {
	repo bookRepo.Repo
}

func NewLinkTracker(r bookRepo.Repo) *LinkTracker {
	return &LinkTracker{repo: r}
}

func (lt *LinkTracker) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(newLinkTrackerParser(lt.repo), 1),
		),
	)
}

type linkTrackingParser struct {
	repo    bookRepo.Repo
	wrapped parser.InlineParser
}

var _ parser.InlineParser = new(linkTrackingParser)

func newLinkTrackerParser(repo bookRepo.Repo) *linkTrackingParser {
	wrapped := parser.NewLinkParser()
	return &linkTrackingParser{repo: repo, wrapped: wrapped}
}

func (l linkTrackingParser) Trigger() []byte {
	return l.wrapped.Trigger()
}

var linkPageIDs = parser.NewContextKey()
var linkCommands = parser.NewContextKey()

func appendToKey[T any](pc parser.Context, key parser.ContextKey, value T) {
	stored, ok := pc.Get(key).([]T)
	if !ok {
		stored = []T{}
	}

	stored = append(stored, value)
	pc.Set(key, stored)
}

func (l linkTrackingParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	result := l.wrapped.Parse(parent, block, pc)
	link, ok := result.(*ast.Link)
	if !ok {
		return result
	}

	destination := string(link.Destination)
	if len(destination) == 0 {
		return result
	}

	if destination[0] == '!' {
		appendToKey(pc, linkCommands, destination)
		link.Destination = append([]byte("?cmd="), link.Destination[1:]...)
		return link
	}

	book := GetCurrentBook(pc)
	currentPageID := GetCurrentPageID(pc)

	pageID := string(link.Destination)

	currentPageID = filepath.Join(book.Path, currentPageID)
	currentPageDir := filepath.Dir(currentPageID)
	newPageID := filepath.Join(currentPageDir, pageID)
	newPageID, err := filepath.Rel(book.Path, newPageID)
	if err != nil {
		panic("failed to generate link")
	}

	appendToKey(pc, linkPageIDs, newPageID)

	link.Destination = append([]byte("?goto="), []byte(newPageID)...)

	return result
}

func GetCommandsFromContext(context parser.Context) []string {
	commands, ok := context.Get(linkCommands).([]string)
	if !ok {
		return nil
	}
	return commands
}

func GetLinksFromContext(context parser.Context) []string {
	links, ok := context.Get(linkPageIDs).([]string)
	if !ok {
		return nil
	}
	return links
}

var currentPageKey = parser.NewContextKey()

func GetCurrentPageID(pc parser.Context) string {
	pagePath, ok := pc.Get(currentPageKey).(string)
	if !ok {
		panic("failed to find current page in context")
	}

	return pagePath
}

func SetCurrentPageID(pc parser.Context, path string) {
	pc.Set(currentPageKey, path)
}

var currentBookKey = parser.NewContextKey()

func SetCurrentBook(pc parser.Context, book *models.Book) {
	pc.Set(currentBookKey, book)
}

func GetCurrentBook(pc parser.Context) *models.Book {
	book, ok := pc.Get(currentBookKey).(*models.Book)
	if !ok {
		panic("failed to find current book in context")
	}
	return book
}
