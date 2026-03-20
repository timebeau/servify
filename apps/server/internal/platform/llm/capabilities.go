package llm

import "servify/apps/server/internal/platform/aiprovider"

func OpenAIDescriptor(enabled bool, model string) aiprovider.ProviderDescriptor {
	return aiprovider.ProviderDescriptor{
		ID:      "openai",
		Kind:    aiprovider.KindLLM,
		Driver:  "openai",
		Enabled: enabled,
		Capabilities: []aiprovider.CapabilityDeclaration{
			{Name: aiprovider.CapabilityChat, Enabled: true, Metadata: map[string]any{"default_model": model}},
			{Name: aiprovider.CapabilityChatStream, Enabled: false},
			{Name: aiprovider.CapabilityEmbeddings, Enabled: false},
			{Name: aiprovider.CapabilityToolCalling, Enabled: true},
			{Name: aiprovider.CapabilityHealthCheck, Enabled: true},
		},
		Fallback: aiprovider.FallbackPolicy{
			Priority: 1,
		},
	}
}

func AnthropicDescriptor(enabled bool, model string) aiprovider.ProviderDescriptor {
	return aiprovider.ProviderDescriptor{
		ID:      "anthropic",
		Kind:    aiprovider.KindLLM,
		Driver:  "anthropic",
		Enabled: enabled,
		Capabilities: []aiprovider.CapabilityDeclaration{
			{Name: aiprovider.CapabilityChat, Enabled: true, Metadata: map[string]any{"default_model": model}},
			{Name: aiprovider.CapabilityChatStream, Enabled: false},
			{Name: aiprovider.CapabilityEmbeddings, Enabled: false},
			{Name: aiprovider.CapabilityToolCalling, Enabled: true},
			{Name: aiprovider.CapabilityHealthCheck, Enabled: true},
		},
		Fallback: aiprovider.FallbackPolicy{
			Priority:   2,
			FallbackTo: []string{"openai"},
		},
	}
}
