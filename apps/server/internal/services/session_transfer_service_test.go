//go:build integration
// +build integration

package services

import (
	"context"
	"strings"
	"testing"
	"time"

	agentdelivery "servify/apps/server/internal/modules/agent/delivery"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	routingapp "servify/apps/server/internal/modules/routing/application"
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"
	routinginfra "servify/apps/server/internal/modules/routing/infra"
	ticketdelivery "servify/apps/server/internal/modules/ticket/delivery"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"servify/apps/server/internal/models"
)

type stubAIForTransfer struct{}

func (s stubAIForTransfer) ProcessQuery(ctx context.Context, query string, sessionID string) (*AIResponse, error) {
	return &AIResponse{Content: "ok", Confidence: 1, Source: "ai"}, nil
}
func (s stubAIForTransfer) ShouldTransferToHuman(query string, sessionHistory []models.Message) bool {
	return false
}
func (s stubAIForTransfer) GetSessionSummary(messages []models.Message) (string, error) {
	return "summary", nil
}
func (s stubAIForTransfer) InitializeKnowledgeBase() {}
func (s stubAIForTransfer) GetStatus(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{"type": "stub"}
}

func newTestDBForSessionTransfer(t *testing.T) *gorm.DB {
	t.Helper()

	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	dsn := "file:session_transfer_" + name + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.Session{},
		&models.Message{},
		&models.TransferRecord{},
		&models.WaitingRecord{},
		&models.Ticket{},
		&models.TicketStatus{},
	); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func TestSessionTransferService_ToHuman_NoAgents_GoesWaiting(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.Session{ID: "s1", UserID: 1, Status: "active", Platform: "web", StartedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}

	agentSvc := NewAgentService(db, logger)
	transferSvc := NewSessionTransferService(db, logger, stubAIForTransfer{}, agentSvc, nil)

	res, err := transferSvc.TransferToHuman(context.Background(), &TransferRequest{SessionID: "s1", Reason: "need help"})
	if err != nil {
		t.Fatalf("TransferToHuman: %v", err)
	}
	if !res.IsWaiting {
		t.Fatalf("expected waiting result: %+v", res)
	}

	var sess models.Session
	if err := db.First(&sess, "id = ?", "s1").Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if sess.Status != "active" || sess.AgentID != nil {
		t.Fatalf("expected session active & unassigned; got status=%s agent=%v", sess.Status, sess.AgentID)
	}

	var waiting models.WaitingRecord
	if err := db.Where("session_id = ? AND status = ?", "s1", "waiting").First(&waiting).Error; err != nil {
		t.Fatalf("expected waiting record: %v", err)
	}
}

func TestSessionTransferService_ToHuman_NoAgents_GoesWaitingViaRoutingAdapter(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.Session{ID: "s-routing-wait", UserID: 1, Status: "active", Platform: "web", StartedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}

	agentSvc := NewAgentService(db, logger)
	bus := eventbus.NewInMemoryBus()
	routingSvc := routingapp.NewService(routinginfra.NewGormRepository(db), bus)
	transferSvc := NewSessionTransferServiceWithAdapters(db, logger, stubAIForTransfer{}, agentSvc, nil, SessionTransferAdapters{
		Routing:      routingdelivery.NewSessionTransferAdapter(routingSvc, bus),
		Tickets:      ticketdelivery.NewRuntimeAdapter(bus),
		Conversation: conversationdelivery.NewRuntimeAdapter(db, bus),
		Agents:       agentdelivery.NewTransferRuntimeAdapter(),
	})

	res, err := transferSvc.TransferToHuman(context.Background(), &TransferRequest{
		SessionID:    "s-routing-wait",
		Reason:       "need help",
		TargetSkills: []string{"billing"},
		Priority:     "high",
		Notes:        "vip",
	})
	if err != nil {
		t.Fatalf("TransferToHuman: %v", err)
	}
	if !res.IsWaiting {
		t.Fatalf("expected waiting result: %+v", res)
	}

	var waiting models.WaitingRecord
	if err := db.Where("session_id = ? AND status = ?", "s-routing-wait", "waiting").First(&waiting).Error; err != nil {
		t.Fatalf("expected waiting record: %v", err)
	}
	if waiting.Reason != "need help" || waiting.Priority != "high" || waiting.Notes != "vip" {
		t.Fatalf("unexpected waiting record: %+v", waiting)
	}

	var message models.Message
	if err := db.Where("session_id = ? AND type = ?", "s-routing-wait", "system").First(&message).Error; err != nil {
		t.Fatalf("expected waiting message: %v", err)
	}
}

func TestSessionTransferService_ToHuman_AssignsAgent_RecordsTransfer(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.User{ID: 2, Username: "a1", Name: "a1", Email: "a1@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("seed agent user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 2, Status: "offline", MaxConcurrent: 5, CurrentLoad: 0}).Error; err != nil {
		t.Fatalf("seed agent: %v", err)
	}

	ticketID := uint(42)
	if err := db.Create(&models.Ticket{ID: ticketID, Title: "t1", Description: "d", CustomerID: 1, Status: "open", Priority: "normal", Source: "web", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed ticket: %v", err)
	}

	if err := db.Create(&models.Session{ID: "s2", UserID: 1, TicketID: &ticketID, Status: "active", Platform: "web", StartedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}

	agentSvc := NewAgentService(db, logger)
	if err := agentSvc.AgentGoOnline(context.Background(), 2); err != nil {
		t.Fatalf("AgentGoOnline: %v", err)
	}
	transferSvc := NewSessionTransferService(db, logger, stubAIForTransfer{}, agentSvc, nil)

	res, err := transferSvc.TransferToHuman(context.Background(), &TransferRequest{SessionID: "s2", Reason: "need help"})
	if err != nil {
		t.Fatalf("TransferToHuman: %v", err)
	}
	if res.IsWaiting || res.NewAgentID != 2 {
		t.Fatalf("expected assigned to agent 2: %+v", res)
	}

	var sess models.Session
	if err := db.First(&sess, "id = ?", "s2").Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if sess.AgentID == nil || *sess.AgentID != 2 || sess.Status != "active" {
		t.Fatalf("expected session assigned active; got status=%s agent=%v", sess.Status, sess.AgentID)
	}

	var tr models.TransferRecord
	if err := db.Where("session_id = ?", "s2").First(&tr).Error; err != nil {
		t.Fatalf("expected transfer record: %v", err)
	}
	if tr.ToAgentID == nil || *tr.ToAgentID != 2 {
		t.Fatalf("expected to_agent_id=2: %+v", tr)
	}

	var agent models.Agent
	if err := db.First(&agent, "user_id = ?", 2).Error; err != nil {
		t.Fatalf("load agent: %v", err)
	}
	if agent.CurrentLoad != 1 {
		t.Fatalf("expected agent load 1, got %d", agent.CurrentLoad)
	}

	var ticket models.Ticket
	if err := db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		t.Fatalf("load ticket: %v", err)
	}
	if ticket.AgentID == nil || *ticket.AgentID != 2 {
		t.Fatalf("expected ticket assigned to agent 2; got %v", ticket.AgentID)
	}
	if ticket.Status != "assigned" {
		t.Fatalf("expected ticket status assigned; got %s", ticket.Status)
	}
}

func TestSessionTransferService_ToHuman_AssignsAgent_RecordsTransferViaRoutingAdapter(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.User{ID: 2, Username: "a1", Name: "a1", Email: "a1@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("seed agent user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 2, Status: "offline", MaxConcurrent: 5, CurrentLoad: 0}).Error; err != nil {
		t.Fatalf("seed agent: %v", err)
	}
	ticketID := uint(43)
	if err := db.Create(&models.Ticket{ID: ticketID, Title: "t-routing", Description: "d", CustomerID: 1, Status: "open", Priority: "normal", Source: "web", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed ticket: %v", err)
	}
	if err := db.Create(&models.Session{ID: "s3", UserID: 1, TicketID: &ticketID, Status: "active", Platform: "web", StartedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}

	agentSvc := NewAgentService(db, logger)
	if err := agentSvc.AgentGoOnline(context.Background(), 2); err != nil {
		t.Fatalf("AgentGoOnline: %v", err)
	}

	bus := eventbus.NewInMemoryBus()
	routingSvc := routingapp.NewService(routinginfra.NewGormRepository(db), bus)
	transferSvc := NewSessionTransferServiceWithAdapters(db, logger, stubAIForTransfer{}, agentSvc, nil, SessionTransferAdapters{
		Routing:      routingdelivery.NewSessionTransferAdapter(routingSvc, bus),
		Tickets:      ticketdelivery.NewRuntimeAdapter(bus),
		Conversation: conversationdelivery.NewRuntimeAdapter(db, bus),
		Agents:       agentdelivery.NewTransferRuntimeAdapter(),
	})

	res, err := transferSvc.TransferToHuman(context.Background(), &TransferRequest{
		SessionID: "s3",
		Reason:    "need help",
		Notes:     "priority customer",
	})
	if err != nil {
		t.Fatalf("TransferToHuman: %v", err)
	}
	if res.IsWaiting || res.NewAgentID != 2 || res.Summary != "summary" {
		t.Fatalf("expected assigned to agent 2: %+v", res)
	}

	var tr models.TransferRecord
	if err := db.Where("session_id = ?", "s3").First(&tr).Error; err != nil {
		t.Fatalf("expected transfer record: %v", err)
	}
	if tr.ToAgentID == nil || *tr.ToAgentID != 2 || tr.SessionSummary != "summary" || tr.Notes != "priority customer" {
		t.Fatalf("unexpected transfer record: %+v", tr)
	}

	var sess models.Session
	if err := db.First(&sess, "id = ?", "s3").Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if sess.AgentID == nil || *sess.AgentID != 2 || sess.Status != "active" || sess.EndedAt != nil {
		t.Fatalf("unexpected session after transfer: %+v", sess)
	}

	var ticket models.Ticket
	if err := db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		t.Fatalf("load ticket: %v", err)
	}
	if ticket.AgentID == nil || *ticket.AgentID != 2 || ticket.Status != "assigned" {
		t.Fatalf("unexpected ticket after transfer: %+v", ticket)
	}

	var agent models.Agent
	if err := db.First(&agent, "user_id = ?", 2).Error; err != nil {
		t.Fatalf("load agent: %v", err)
	}
	if agent.CurrentLoad != 1 {
		t.Fatalf("unexpected agent load after transfer: %+v", agent)
	}
}

func TestSessionTransferService_CancelWaitingRecord_ViaRoutingAdapter(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.WaitingRecord{
		SessionID:    "s-cancel",
		Reason:       "need help",
		Status:       "waiting",
		QueuedAt:     time.Now(),
		TargetSkills: "billing",
		CreatedAt:    time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed waiting record: %v", err)
	}

	bus := eventbus.NewInMemoryBus()
	routingSvc := routingapp.NewService(routinginfra.NewGormRepository(db), bus)
	transferSvc := NewSessionTransferServiceWithAdapters(
		db,
		logger,
		stubAIForTransfer{},
		NewAgentService(db, logger),
		nil,
		SessionTransferAdapters{
			Routing:      routingdelivery.NewSessionTransferAdapter(routingSvc, bus),
			Conversation: conversationdelivery.NewRuntimeAdapter(db, bus),
		},
	)

	if err := transferSvc.CancelWaitingRecord(context.Background(), "s-cancel", 1, "user_left"); err != nil {
		t.Fatalf("CancelWaitingRecord: %v", err)
	}

	var waiting models.WaitingRecord
	if err := db.Where("session_id = ?", "s-cancel").First(&waiting).Error; err != nil {
		t.Fatalf("load waiting record: %v", err)
	}
	if waiting.Status != "cancelled" || waiting.Notes != "user_left" {
		t.Fatalf("unexpected waiting record after cancel: %+v", waiting)
	}

	var message models.Message
	if err := db.Where("session_id = ? AND type = ?", "s-cancel", "system").First(&message).Error; err != nil {
		t.Fatalf("load cancel message: %v", err)
	}
	if !strings.Contains(message.Content, "user_left") {
		t.Fatalf("unexpected cancel message: %+v", message)
	}
}

func TestSessionTransferService_ProcessWaitingQueue_TransfersWaitingViaRoutingAdapter(t *testing.T) {
	db := newTestDBForSessionTransfer(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{ID: 1, Username: "u1", Name: "u1", Email: "u1@example.com", Role: "customer"}).Error; err != nil {
		t.Fatalf("seed customer: %v", err)
	}
	if err := db.Create(&models.User{ID: 2, Username: "a1", Name: "a1", Email: "a1@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("seed agent user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 2, Status: "offline", MaxConcurrent: 5, CurrentLoad: 0}).Error; err != nil {
		t.Fatalf("seed agent: %v", err)
	}
	if err := db.Create(&models.Session{
		ID:        "s-process-wait",
		UserID:    1,
		Status:    "active",
		Platform:  "web",
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}
	if err := db.Create(&models.WaitingRecord{
		SessionID:    "s-process-wait",
		Reason:       "need help",
		Notes:        "from queue",
		Priority:     "high",
		Status:       "waiting",
		QueuedAt:     time.Now(),
		TargetSkills: "billing",
		CreatedAt:    time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed waiting record: %v", err)
	}

	agentSvc := NewAgentService(db, logger)
	if err := agentSvc.AgentGoOnline(context.Background(), 2); err != nil {
		t.Fatalf("AgentGoOnline: %v", err)
	}

	bus := eventbus.NewInMemoryBus()
	routingSvc := routingapp.NewService(routinginfra.NewGormRepository(db), bus)
	transferSvc := NewSessionTransferServiceWithAdapters(
		db,
		logger,
		stubAIForTransfer{},
		agentSvc,
		nil,
		SessionTransferAdapters{
			Routing: routingdelivery.NewSessionTransferAdapter(routingSvc, bus),
		},
	)

	if err := transferSvc.ProcessWaitingQueue(context.Background()); err != nil {
		t.Fatalf("ProcessWaitingQueue: %v", err)
	}

	var waiting models.WaitingRecord
	if err := db.Where("session_id = ?", "s-process-wait").First(&waiting).Error; err != nil {
		t.Fatalf("load waiting record: %v", err)
	}
	if waiting.Status != "transferred" || waiting.AssignedTo == nil || *waiting.AssignedTo != 2 || waiting.AssignedAt == nil {
		t.Fatalf("unexpected waiting record after transfer: %+v", waiting)
	}

	var tr models.TransferRecord
	if err := db.Where("session_id = ?", "s-process-wait").First(&tr).Error; err != nil {
		t.Fatalf("load transfer record: %v", err)
	}
	if tr.ToAgentID == nil || *tr.ToAgentID != 2 {
		t.Fatalf("unexpected transfer record: %+v", tr)
	}
}
