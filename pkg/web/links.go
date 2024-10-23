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
	if l, ok := result.(*ast.Link); ok {
		stored, ok := pc.Get(linksKey).([]string)
		if !ok {
			stored = []string{}
		}
		stored = append(stored, string(l.Destination))
		pc.Set(linksKey, stored)
		l.Destination = append([]byte("?goto="), l.Destination...)
	}

	return result
}

func getLinksFromContext(context parser.Context) []string {
	links, ok := context.Get(linksKey).([]string)
	if !ok {
		return nil
	}
	return links
}
