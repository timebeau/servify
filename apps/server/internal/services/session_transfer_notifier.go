package services

import "time"

type sessionTransferRealtimeSink interface {
	SendToSession(sessionID string, message WebSocketMessage)
}

type sessionTransferNotifier interface {
	NotifyTransfer(sessionID string, agentID uint, message string, at time.Time)
	NotifyWaiting(sessionID string, message string, at time.Time)
}

type websocketSessionTransferNotifier struct {
	sink sessionTransferRealtimeSink
}

func NewSessionTransferNotifier(sink sessionTransferRealtimeSink) sessionTransferNotifier {
	if sink == nil {
		return nil
	}
	return &websocketSessionTransferNotifier{sink: sink}
}

func (n *websocketSessionTransferNotifier) NotifyTransfer(sessionID string, agentID uint, message string, at time.Time) {
	n.sink.SendToSession(sessionID, WebSocketMessage{
		Type: "transfer_notification",
		Data: map[string]interface{}{
			"message":   message,
			"agent_id":  agentID,
			"timestamp": at,
		},
	})
}

func (n *websocketSessionTransferNotifier) NotifyWaiting(sessionID string, message string, at time.Time) {
	n.sink.SendToSession(sessionID, WebSocketMessage{
		Type: "waiting_notification",
		Data: map[string]interface{}{
			"message":   message,
			"timestamp": at,
		},
	})
}
