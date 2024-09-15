package functions

import (
	"assistant/config"
	"assistant/pkg/api/context"
	"assistant/pkg/api/core"
	"strings"
)

const echoFunctionName = "echo"

type echoFunction struct {
	stub
}

func NewEchoFunction(ctx context.Context, cfg *config.Config, irc core.IRC) (Function, error) {
	stub, err := newFunctionStub(ctx, cfg, irc, echoFunctionName)
	if err != nil {
		return nil, err
	}

	return &echoFunction{
		stub: stub,
	}, nil
}

func (f *echoFunction) ShouldExecute(e *core.Event) bool {
	ok, _ := f.verifyInput(e, 1)
	return ok
}

func (f *echoFunction) Execute(e *core.Event) error {
	tokens := parseTokens(e.Message())
	f.irc.SendMessage(e.ReplyTarget(), strings.Join(tokens[1:], " "))
	return nil
}
