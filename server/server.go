package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	pb "train-ticket-app/pb/proto"

	"google.golang.org/grpc"
)

// Example seat numbers
var seats_SectionA = []string{"A1", "A2", "A3", "A4"}
var seats_SectionB = []string{"B1", "B2", "B3", "B4"}

type trainServer struct {
	pb.UnimplementedTrainServiceServer
	receipts map[string]*pb.Receipt
	seats    map[string]string // Maps email to seat
	sectionA []string          // Seats in Section A
	sectionB []string          // Seats in Section B
	mu       sync.Mutex
}

func newServer() *trainServer {
	return &trainServer{
		receipts: make(map[string]*pb.Receipt),
		seats:    make(map[string]string),
		sectionA: seats_SectionA,
		sectionB: seats_SectionB,
	}
}

// Implement the PurchaseTicket RPC
func (s *trainServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already has a ticket
	if _, exists := s.receipts[req.User.Email]; exists {
		return nil, fmt.Errorf("user already has a ticket")
	}

	// Assign seat (simple round-robin allocation for demonstration)
	var seat string
	if len(s.sectionA) > 0 {
		seat, s.sectionA = s.sectionA[0], s.sectionA[1:]
	} else if len(s.sectionB) > 0 {
		seat, s.sectionB = s.sectionB[0], s.sectionB[1:]
	} else {
		return nil, fmt.Errorf("no seats available")
	}

	receipt := &pb.Receipt{
		User:  req.User,
		From:  req.From,
		To:    req.To,
		Price: 20.00,
		Seat:  seat,
	}

	s.receipts[req.User.Email] = receipt
	s.seats[req.User.Email] = seat

	return receipt, nil
}

// Implement the GetReceipt RPC
func (s *trainServer) GetReceipt(ctx context.Context, req *pb.UserRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.receipts[req.Email]
	if !exists {
		return nil, fmt.Errorf("no receipt found for user")
	}

	return receipt, nil
}

// Implement the ViewSeats RPC
func (s *trainServer) ViewSeats(ctx context.Context, req *pb.SectionRequest) (*pb.SeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Section != "A" && req.Section != "B" {
		return nil, fmt.Errorf("no section for %v", req.Section)
	}

	var users []*pb.User
	for email, seat := range s.seats {
		if (req.Section == "A" && seat[0] == 'A') || (req.Section == "B" && seat[0] == 'B') {
			receipt := s.receipts[email]
			users = append(users, receipt.User)
		}
	}

	return &pb.SeatResponse{Users: users}, nil
}

// Implement the RemoveUser RPC
func (s *trainServer) RemoveUser(ctx context.Context, req *pb.UserRequest) (*pb.GenericResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.receipts[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Return the seat to the pool
	seat := receipt.Seat
	if seat[0] == 'A' {
		s.sectionA = append(s.sectionA, seat)
	} else {
		s.sectionB = append(s.sectionB, seat)
	}

	delete(s.receipts, req.Email)
	delete(s.seats, req.Email)

	return &pb.GenericResponse{Message: "User removed successfully"}, nil
}

// Implement the ModifySeat RPC
func (s *trainServer) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.GenericResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.receipts[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Check if the new seat is not taken by others
	newSeat := req.NewSeat
	for _, seat := range s.seats {
		if seat == newSeat {
			return nil, fmt.Errorf("seat %s is already taken", newSeat)
		}
	}

	// Check if the new seat is exists in seat array
	isExist := false
	if req.NewSeat[0] == 'A' {
		for _, seat := range seats_SectionA {
			if seat == newSeat {
				isExist = true
			}
		}
	} else if req.NewSeat[0] == 'B' {
		for _, seat := range seats_SectionB {
			if seat == newSeat {
				isExist = true
			}
		}
	}

	if !isExist {
		return nil, fmt.Errorf("seat %v is not existing.", newSeat)
	}

	// Update the seat
	oldSeat := receipt.Seat
	receipt.Seat = newSeat
	s.seats[req.Email] = newSeat

	// Return the old seat to the pool
	if oldSeat[0] == 'A' {
		s.sectionA = append(s.sectionA, oldSeat)
	} else {
		s.sectionB = append(s.sectionB, oldSeat)
	}

	return &pb.GenericResponse{Message: "Seat modified successfully"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTrainServiceServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
