package llm

import "servify/apps/server/internal/platform/aiprovider"

func OpenAIDescriptor(enabled bool, model string) aiprovider.ProviderDescriptor {
	return aiprovider.ProviderDescriptor{
		ID:      "openai",
		Kind:    aiprovider.KindLLM,
		Driver:  "openai",
		Enabled: enabled,
		Capabilities: []aiprovider.CapabilityDeclaration{
			{Name: aiprovider.CapabilityChat, Enabled: true},
			{Name: aiprovider.CapabilityChatStream, Enabled: false},
			{Name: aiprovider.CapabilityEmbeddings, Enabled: false},
			{Name: aiprovider.CapabilityToolCalling, Enabled: false},
			{Name: aiprovider.CapabilityHealthCheck, Enabled: true},
		},
		Fallback: aiprovider.FallbackPolicy{
			Priority: 1,
		},
	}
}
