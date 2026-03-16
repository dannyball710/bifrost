package integrations

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestApplyUsageHeaderForGenAIResponses(t *testing.T) {
	router := &GenericRouter{}
	ctx := &fasthttp.RequestCtx{}
	resp := &schemas.BifrostResponsesResponse{
		Usage: &schemas.ResponsesResponseUsage{
			InputTokens:  12,
			OutputTokens: 7,
			TotalTokens:  19,
			Cost: &schemas.BifrostCost{
				TotalCost: 1.25,
			},
		},
	}

	router.applyUsageHeader(ctx, RouteConfigTypeGenAI, resp)

	assertUsageHeader(t, ctx.Response.Header.Peek("X-Usage"), usageHeaderPayload{
		InputTokens:  12,
		OutputTokens: 7,
		TotalTokens:  19,
		Cost:         1.25,
	})
}

func TestApplyUsageHeaderForGenAICountTokens(t *testing.T) {
	router := &GenericRouter{}
	ctx := &fasthttp.RequestCtx{}
	outputTokens := 4
	totalTokens := 21
	resp := &schemas.BifrostCountTokensResponse{
		InputTokens:  17,
		OutputTokens: &outputTokens,
		TotalTokens:  &totalTokens,
	}

	router.applyUsageHeader(ctx, RouteConfigTypeGenAI, resp)

	assertUsageHeader(t, ctx.Response.Header.Peek("X-Usage"), usageHeaderPayload{
		InputTokens:  17,
		OutputTokens: 4,
		TotalTokens:  21,
		Cost:         0,
	})
}

func TestApplyUsageHeaderSkipsNonGenAI(t *testing.T) {
	router := &GenericRouter{}
	ctx := &fasthttp.RequestCtx{}
	resp := &schemas.BifrostResponsesResponse{
		Usage: &schemas.ResponsesResponseUsage{
			InputTokens:  12,
			OutputTokens: 7,
			TotalTokens:  19,
		},
	}

	router.applyUsageHeader(ctx, RouteConfigTypeOpenAI, resp)

	require.Empty(t, ctx.Response.Header.Peek("X-Usage"))
}

func assertUsageHeader(t *testing.T, headerValue []byte, want usageHeaderPayload) {
	t.Helper()
	require.NotEmpty(t, headerValue)

	var got usageHeaderPayload
	require.NoError(t, sonic.Unmarshal(headerValue, &got))
	require.Equal(t, want, got)
}
