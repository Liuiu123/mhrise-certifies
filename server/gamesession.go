package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	matchmaking "npln.nintendo.net/npln-practice/proto/matchmaking/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameSessionServer struct {
	matchmaking.UnimplementedGameSessionServiceServer

	mu      sync.Mutex
	tickets map[string]*matchmaking.GameSessionCreationTicket
}

func NewGameSessionServer() *GameSessionServer {
	return &GameSessionServer{
		tickets: make(map[string]*matchmaking.GameSessionCreationTicket),
	}
}

func (s *GameSessionServer) CreateGameSessionCreationTicket(
	ctx context.Context,
	req *matchmaking.CreateGameSessionCreationTicketRequest,
) (*matchmaking.GameSessionCreationTicket, error) {

	if req == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	if req.GetGameSessionCreationTicket() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"game_session_creation_ticket is required",
		)
	}

	src := req.GetGameSessionCreationTicket()

	if src.GetMatchmakingConfig() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"matchmaking_config is required",
		)
	}

	ticket := &matchmaking.GameSessionCreationTicket{
		Name:              fmt.Sprintf("game-session-tickets/mhrise-%d", time.Now().UnixNano()),
		MatchmakingConfig: src.GetMatchmakingConfig(),
		UserDefinitions:   src.GetUserDefinitions(),

		// No 0.21.x, o estado inicial válido é PENDING.
		State: matchmaking.GameSessionCreationTicket_PENDING,
	}

	s.mu.Lock()
	s.tickets[ticket.GetName()] = ticket
	s.mu.Unlock()

	log.Printf(
		"[GameSession] CreateGameSessionCreationTicket name=%s config=%q users=%d state=PENDING",
		ticket.GetName(),
		ticket.GetMatchmakingConfig(),
		len(ticket.GetUserDefinitions()),
	)

	return ticket, nil
}

func (s *GameSessionServer) TrackGameSessionCreationTicket(
	req *matchmaking.TrackGameSessionCreationTicketRequest,
	stream matchmaking.GameSessionService_TrackGameSessionCreationTicketServer,
) error {

	if req == nil || req.GetName() == "" {
		return status.Error(
			codes.InvalidArgument,
			"name is required",
		)
	}

	name := req.GetName()

	log.Printf(
		"[GameSession] TrackGameSessionCreationTicket name=%s",
		name,
	)

	// Primeiro envio: PENDING.
	if err := s.sendTicket(stream, name); err != nil {
		return err
	}

	// Espera um pouco para simular processamento.
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()

	case <-time.After(1 * time.Second):
	}

	// Transição para SUCCEEDED.
	s.mu.Lock()

	ticket, exists := s.tickets[name]

	if !exists {
		s.mu.Unlock()

		return status.Error(
			codes.NotFound,
			"game session creation ticket not found",
		)
	}

	ticket.State = matchmaking.GameSessionCreationTicket_SUCCEEDED

	sessionName := fmt.Sprintf(
		"game-sessions/mhrise-%d",
		time.Now().UnixNano(),
	)

	ticket.GameSession = &matchmaking.GameSession{
		Name:                   sessionName,
		MaxParticipantCount:    4,
		CurrentParticipantCount: 0,
		CanParticipate:         true,
		IsPublic:               true,
		State:                  matchmaking.GameSession_ACTIVE,
	}

	s.mu.Unlock()

	log.Printf(
		"[GameSession] Ticket %s -> SUCCEEDED",
		name,
	)

	log.Printf(
		"[GameSession] GameSession criada: %s",
		sessionName,
	)

	// Segundo envio: SUCCEEDED + GameSession.
	return s.sendTicket(stream, name)
}

func (s *GameSessionServer) sendTicket(
	stream matchmaking.GameSessionService_TrackGameSessionCreationTicketServer,
	name string,
) error {

	s.mu.Lock()

	ticket, exists := s.tickets[name]

	if !exists {
		s.mu.Unlock()

		return status.Error(
			codes.NotFound,
			"game session creation ticket not found",
		)
	}

	// Cópia para não manter o mutex durante o envio.
	update := *ticket

	s.mu.Unlock()

	log.Printf(
		"[GameSession] Ticket update name=%s state=%s",
		update.GetName(),
		update.GetState().String(),
	)

	if err := stream.Send(&update); err != nil {
		return err
	}

	return nil
}

func (s *GameSessionServer) GetGameSession(
	ctx context.Context,
	req *matchmaking.GetGameSessionRequest,
) (*matchmaking.GameSession, error) {

	if req == nil || req.GetName() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"name is required",
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ticket := range s.tickets {

		session := ticket.GetGameSession()

		if session != nil &&
			session.GetName() == req.GetName() {

			log.Printf(
				"[GameSession] GetGameSession name=%s",
				req.GetName(),
			)

			return session, nil
		}
	}

	return nil, status.Error(
		codes.NotFound,
		"game session not found",
	)
}