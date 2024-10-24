package web

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type LinkTracker struct{}

func NewLinkTracker() *LinkTracker {
	return new(LinkTracker)
}

func (lt *LinkTracker) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(newLinkTrackerParser(), 1),
		),
	)
}

type linkTrackingParser struct {
	wrapped parser.InlineParser
}

var _ parser.InlineParser = new(linkTrackingParser)

func newLinkTrackerParser() *linkTrackingParser {
	wrapped := parser.NewLinkParser()
	return &linkTrackingParser{wrapped: wrapped}
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

	stored = append(stored, destination)
	pc.Set(linksKey, stored)
	link.Destination = append([]byte("?goto="), link.Destination...)

	return result
}

func getLinksFromContext(context parser.Context) []string {
	links, ok := context.Get(linksKey).([]string)
	if !ok {
		return nil
	}
	return links
}
