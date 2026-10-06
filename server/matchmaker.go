package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	matchmaking "npln.nintendo.net/npln-practice/proto/matchmaking/v1"
)

type MatchmakerServer struct {
	matchmaking.UnimplementedMatchmakerServer
	state *LabState
}

func NewMatchmakerServer() *MatchmakerServer {
	return &MatchmakerServer{
		state: NewLabState(),
	}
}

func (s *MatchmakerServer) CreateMatchmakingTicket(
	ctx context.Context,
	req *matchmaking.CreateMatchmakingTicketRequest,
) (*matchmaking.MatchmakingTicket, error) {

	if req.GetMatchmakingTicket() == nil {
		return nil, fmt.Errorf("matchmaking_ticket is required")
	}

	src := req.GetMatchmakingTicket()

	ticket := &matchmaking.MatchmakingTicket{
		Name:              s.state.nextID("tickets/mhrise"),
		MatchmakingConfig: src.GetMatchmakingConfig(),
		UserDefinitions:   src.GetUserDefinitions(),
		State:             matchmaking.MatchmakingTicket_SEARCHING,
	}

	s.state.mu.Lock()

	s.state.tickets[ticket.GetName()] = &TicketState{
		Name:       ticket.GetName(),
		State:      matchmaking.MatchmakingTicket_SEARCHING.String(),
		TrackCount: 0,
		CreatedAt:  time.Now(),
	}

	s.state.mu.Unlock()

	log.Printf(
		"[Matchmaker] CreateMatchmakingTicket name=%s users=%d config=%q state=SEARCHING",
		ticket.GetName(),
		len(ticket.GetUserDefinitions()),
		ticket.GetMatchmakingConfig(),
	)

	return ticket, nil
}

func (s *MatchmakerServer) TrackMatchmakingTicket(
	req *matchmaking.TrackMatchmakingTicketRequest,
	stream matchmaking.Matchmaker_TrackMatchmakingTicketServer,
) error {

	if req == nil || req.GetName() == "" {
		return fmt.Errorf("name is required")
	}

	name := req.GetName()

	for {
		s.state.mu.Lock()

		ticketState, ok := s.state.tickets[name]

		if !ok {
			s.state.mu.Unlock()
			return fmt.Errorf("ticket %q not found", name)
		}

		ticketState.TrackCount++

		var state matchmaking.MatchmakingTicket_State

		switch ticketState.TrackCount {
		case 1:
			state = matchmaking.MatchmakingTicket_SEARCHING

		case 2:
			state = matchmaking.MatchmakingTicket_PLACING

		default:
			state = matchmaking.MatchmakingTicket_SUCCEEDED
		}

		ticketState.State = state.String()

		trackCount := ticketState.TrackCount

		s.state.mu.Unlock()

		ticket := &matchmaking.MatchmakingTicket{
			Name:  name,
			State: state,
		}

		log.Printf(
			"[Matchmaker] TrackMatchmakingTicket name=%s track=%d state=%s",
			name,
			trackCount,
			state.String(),
		)

		if err := stream.Send(ticket); err != nil {
			return err
		}

		if state == matchmaking.MatchmakingTicket_SUCCEEDED {
			log.Printf(
				"[Matchmaker] Ticket %s concluido.",
				name,
			)

			return nil
		}

		select {
		case <-stream.Context().Done():
			return stream.Context().Err()

		case <-time.After(1 * time.Second):
		}
	}
}

func (s *MatchmakerServer) CancelMatchmakingTicket(
	ctx context.Context,
	req *matchmaking.CancelMatchmakingTicketRequest,
) (*emptypb.Empty, error) {

	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	s.state.mu.Lock()
	defer s.state.mu.Unlock()

	if ticket, ok := s.state.tickets[req.GetName()]; ok {
		ticket.State = matchmaking.MatchmakingTicket_CANCELLED.String()

		log.Printf(
			"[Matchmaker] CancelMatchmakingTicket name=%s state=CANCELLED",
			req.GetName(),
		)
	}

	return &emptypb.Empty{}, nil
}

func (s *MatchmakerServer) CreateAcceptance(
	ctx context.Context,
	req *matchmaking.CreateAcceptanceRequest,
) (*matchmaking.Acceptance, error) {

	var a *matchmaking.Acceptance

	if req != nil {
		a = req.GetAcceptance()
	}

	if a == nil {
		a = &matchmaking.Acceptance{}
	}

	a.Name = s.state.nextID("acceptances/mhrise")

	log.Printf(
		"[Matchmaker] CreateAcceptance name=%s",
		a.GetName(),
	)

	return a, nil
}