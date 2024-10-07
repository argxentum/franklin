package functions

import (
	"assistant/pkg/api/irc"
	"assistant/pkg/api/style"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func (f *summaryFunction) parseMastodon(_ *irc.Event, doc *goquery.Document) (*summary, error) {
	titleAttr, _ := doc.Find("meta[property='og:title']").First().Attr("content")
	title := strings.TrimSpace(titleAttr)
	descriptionAttr, _ := doc.Find("html meta[property='og:description']").First().Attr("content")
	description := strings.TrimSpace(descriptionAttr)

	if len(description) > maximumDescriptionLength {
		description = description[:maximumDescriptionLength] + "..."
	}

	return &summary{text: fmt.Sprintf("%s • %s", style.Bold(description), title)}, nil
}
