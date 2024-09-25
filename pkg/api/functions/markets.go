package functions

import (
	"assistant/pkg/api/context"
	"assistant/pkg/api/irc"
	"assistant/pkg/api/style"
	"assistant/pkg/config"
	"assistant/pkg/log"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"net/url"
	"slices"
	"strings"
)

const marketsFunctionName = "markets"

type marketsFunction struct {
	FunctionStub
}

func NewMarketsFunction(ctx context.Context, cfg *config.Config, irc irc.IRC) (Function, error) {
	stub, err := newFunctionStub(ctx, cfg, irc, marketsFunctionName)
	if err != nil {
		return nil, err
	}

	return &marketsFunction{
		FunctionStub: stub,
	}, nil
}

func (f *marketsFunction) MayExecute(e *irc.Event) bool {
	return f.isValid(e, 0)
}

func (f *marketsFunction) Execute(e *irc.Event) {
	tokens := Tokens(e.Message())
	region := "US"
	if len(tokens) > 1 {
		region = tokens[1]
	}

	logger := log.Logger()
	logger.Infof(e, "⚡ [%s/%s] markets %s", e.From, e.ReplyTarget(), region)

	query := url.QueryEscape(fmt.Sprintf("stock markets %s", region))
	doc, err := getDocument(fmt.Sprintf(bingSearchURL, query), true)
	if err != nil {
		logger.Warningf(e, "unable to retrieve %s stock market information: %s", region, err)
		f.Replyf(e, "Unable to retrieve %s stock market information.", style.Bold(region))
		return
	}

	t := createDefaultTable()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft},
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignLeft},
	})

	section := doc.Find("div.finmkt").First()
	title := strings.TrimSpace(section.Find("h2").First().Text())

	markets := section.Find("div.finind_ind").First()
	markets.Find("div.finind_item").Each(func(i int, s *goquery.Selection) {
		ticker := strings.TrimSpace(s.Find("div.finind_ticker").First().Text())
		val := s.Find("div.finind_val").First()
		value := strings.TrimSpace(val.Text())
		change := strings.TrimSpace(val.Next().Text())

		if len(ticker) == 0 || len(value) == 0 {
			logger.Warningf(e, "skipping invalid stock market information: %s %s %s", ticker, value, change)
			return
		}

		styledChange := change
		if strings.HasPrefix(change, "▼") {
			styledChange = style.ColorForeground(change, style.ColorRed)
		} else if strings.HasPrefix(change, "▲") {
			styledChange = style.ColorForeground(change, style.ColorGreen)
		}

		t.AppendRow([]any{style.Bold(ticker), value, styledChange})
	})

	messages := strings.Split(t.Render(), "\n")

	if len(messages) == 0 {
		logger.Warningf(e, "no stock market information found")
		f.Replyf(e, "Unable to retrieve stock market information.")
		return
	}

	messages = slices.Insert(messages, 0, style.Bold(title))

	f.SendMessages(e, e.ReplyTarget(), messages)
}
