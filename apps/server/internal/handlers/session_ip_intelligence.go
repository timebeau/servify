package handlers

import "strings"

type sessionIPIntelligence interface {
	DescribeIP(ip string) sessionIPDescription
}

type sessionIPDescription struct {
	NetworkLabel  string
	LocationLabel string
}

type heuristicSessionIPIntelligence struct{}

func (heuristicSessionIPIntelligence) DescribeIP(ip string) sessionIPDescription {
	return sessionIPDescription{
		NetworkLabel:  classifyNetworkLabel(ip),
		LocationLabel: classifyLocationLabel(ip),
	}
}

func describeSessionIP(provider sessionIPIntelligence, ip string) sessionIPDescription {
	if provider == nil {
		provider = heuristicSessionIPIntelligence{}
	}
	desc := provider.DescribeIP(ip)
	desc.NetworkLabel = strings.TrimSpace(desc.NetworkLabel)
	desc.LocationLabel = strings.TrimSpace(desc.LocationLabel)
	if desc.NetworkLabel == "" {
		desc.NetworkLabel = classifyNetworkLabel(ip)
	}
	if desc.LocationLabel == "" {
		desc.LocationLabel = classifyLocationLabel(ip)
	}
	return desc
}
