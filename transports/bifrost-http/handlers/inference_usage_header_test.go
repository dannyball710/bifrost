package handlers

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/valyala/fasthttp"
)

func TestApplyUsageHeaderForChatResponse(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	resp := &schemas.BifrostChatResponse{
		Usage: &schemas.BifrostLLMUsage{
			PromptTokens:     12,
			CompletionTokens: 8,
			TotalTokens:      20,
			Cost: &schemas.BifrostCost{
				TotalCost: 0.42,
			},
		},
	}

	applyUsageHeader(ctx, resp, nil)

	assertUsageHeader(t, ctx.Response.Header.Peek("X-Usage"), usageHeaderPayload{
		InputTokens:  12,
		OutputTokens: 8,
		TotalTokens:  20,
		Cost:         0.42,
	})
}

func TestApplyUsageHeaderForResponsesResponse(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	resp := &schemas.BifrostResponsesResponse{
		Usage: &schemas.ResponsesResponseUsage{
			InputTokens:  30,
			OutputTokens: 15,
			TotalTokens:  45,
			Cost: &schemas.BifrostCost{
				TotalCost: 1.25,
			},
		},
	}

	applyUsageHeader(ctx, resp, nil)

	assertUsageHeader(t, ctx.Response.Header.Peek("X-Usage"), usageHeaderPayload{
		InputTokens:  30,
		OutputTokens: 15,
		TotalTokens:  45,
		Cost:         1.25,
	})
}

func TestApplyUsageHeaderForSpeechResponse(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	resp := &schemas.BifrostSpeechResponse{
		Usage: &schemas.SpeechUsage{
			InputTokens:  21,
			OutputTokens: 0,
			TotalTokens:  21,
		},
	}

	applyUsageHeader(ctx, resp, nil)

	assertUsageHeader(t, ctx.Response.Header.Peek("X-Usage"), usageHeaderPayload{
		InputTokens:  21,
		OutputTokens: 0,
		TotalTokens:  21,
		Cost:         0,
	})
}

func TestApplyUsageHeaderSkipsUnsupportedPayload(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}

	applyUsageHeader(ctx, map[string]any{"ok": true}, nil)

	if got := ctx.Response.Header.Peek("X-Usage"); len(got) != 0 {
		t.Fatalf("expected X-Usage to be absent, got %q", string(got))
	}
}

func assertUsageHeader(t *testing.T, raw []byte, want usageHeaderPayload) {
	t.Helper()
	if len(raw) == 0 {
		t.Fatal("expected X-Usage header to be set")
	}

	var got usageHeaderPayload
	if err := sonic.Unmarshal(raw, &got); err != nil {
		t.Fatalf("failed to unmarshal X-Usage header: %v", err)
	}

	if got != want {
		t.Fatalf("unexpected X-Usage header: got %+v want %+v", got, want)
	}
}
