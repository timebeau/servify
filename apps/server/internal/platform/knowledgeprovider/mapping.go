package knowledgeprovider

import "strings"

type Namespace struct {
	TenantID    string
	KnowledgeID string
}

func ResolveNamespace(defaultTenantID, defaultKnowledgeID string, tenantID string, knowledgeID string) Namespace {
	resolvedTenant := strings.TrimSpace(tenantID)
	if resolvedTenant == "" {
		resolvedTenant = strings.TrimSpace(defaultTenantID)
	}

	resolvedKnowledge := strings.TrimSpace(knowledgeID)
	if resolvedKnowledge == "" {
		resolvedKnowledge = strings.TrimSpace(defaultKnowledgeID)
	}

	if resolvedKnowledge == "" && resolvedTenant != "" {
		resolvedKnowledge = resolvedTenant
	}

	return Namespace{
		TenantID:    resolvedTenant,
		KnowledgeID: resolvedKnowledge,
	}
}
