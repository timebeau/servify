package server

import (
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"
	realtimeplatform "servify/apps/server/internal/platform/realtime"
)

type routingTransferNotifier struct {
	gateway realtimeplatform.RealtimeGateway
}

func newRoutingTransferNotifier(gateway realtimeplatform.RealtimeGateway) *routingTransferNotifier {
	if gateway == nil {
		return nil
	}
	return &routingTransferNotifier{gateway: gateway}
}

func (n *routingTransferNotifier) SendToSession(sessionID string, message routingdelivery.Notification) {
	n.gateway.SendToSession(sessionID, realtimeplatform.Message{
		Type:      message.Type,
		Data:      message.Data,
		SessionID: message.SessionID,
		Timestamp: message.Timestamp,
	})
}
