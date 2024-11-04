package markdown

import (
	"gamebooks/pkg/models"
	bookRepo "gamebooks/pkg/repo"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

var linksKey = parser.NewContextKey()

func (l linkTrackingParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	result := l.wrapped.Parse(parent, block, pc)
	link, ok := result.(*ast.Link)
	if !ok {
		return result
	}

	stored, ok := pc.Get(linksKey).([]string)
	if !ok {
		stored = []string{}
	}

	destination := string(link.Destination)
	if len(destination) == 0 {
		return result
	}

	if destination[0] == '!' {
		link.Destination = append([]byte("?cmd="), link.Destination[1:]...)
		return link
	}

	book := GetCurrentBook(pc)
	pagePath := GetCurrentPagePath(pc)

	pageID := string(link.Destination)
	path, err := l.repo.FindPagePath(book, pagePath, pageID)
	if err != nil {
		log.Warn().
			Err(err).
			Str("page_id", pageID).
			Msg("page cannot be found")
	}

	stored = append(stored, path)
	pc.Set(linksKey, stored)

	link.Destination = append([]byte("?goto="), []byte(path)...)

	return result
}

func GetLinksFromContext(context parser.Context) []string {
	links, ok := context.Get(linksKey).([]string)
	if !ok {
		return nil
	}
	return links
}

var currentPageKey = parser.NewContextKey()

func GetCurrentPagePath(pc parser.Context) string {
	pagePath, ok := pc.Get(currentPageKey).(string)
	if !ok {
		panic("failed to find current page in context")
	}

	return pagePath
}

func SetCurrentPagePath(pc parser.Context, path string) {
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
