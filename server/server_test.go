package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "train-ticket-app/pb/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Buffer size for the listener
const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTrainServiceServer(s, newServer())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestPurchaseTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	tests := []struct {
		name    string
		user    *pb.User
		from    string
		to      string
		wantErr bool
	}{
		{"ValidPurchase", &pb.User{FirstName: "Andrus", LastName: "Taylor", Email: "andrustaylor90@gmail.com"}, "London", "France", false},
		{"DuplicatePurchase", &pb.User{FirstName: "Andrus", LastName: "Taylor", Email: "andrustaylor90@gmail.com"}, "London", "France", true},
		// More test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: tc.user, From: tc.from, To: tc.to})
			if (err != nil) != tc.wantErr {
				t.Errorf("PurchaseTicket(%v) got err: %v, wantErr %v", tc.user, err, tc.wantErr)
			}
		})
	}
}

func TestGetReceipt(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: &pb.User{FirstName: "Alice", LastName: "Smith", Email: "alice@gmail.com"}, From: "London", To: "France"})

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"ValidRequest", "alice@gmail.com", false},
		{"InvalidRequest", "bob@gmail.com", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.GetReceipt(ctx, &pb.UserRequest{Email: tc.email})

			if (err != nil) != tc.wantErr {
				t.Errorf("GetReceipt(%v) got err: %v, wantErr %v", tc.email, err, tc.wantErr)
			}
		})
	}
}

func TestViewSeats(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	// Pre-populate with some data
	client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: &pb.User{FirstName: "Alice", LastName: "Smith", Email: "alice@gmail.com"}, From: "London", To: "France"})

	tests := []struct {
		name    string
		section string
		wantErr bool
	}{
		{"ViewSectionA", "A", false},
		{"ViewSectionB", "B", false},
		{"InvalidSection", "C", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.ViewSeats(ctx, &pb.SectionRequest{Section: tc.section})
			if (err != nil) != tc.wantErr {
				t.Errorf("ViewSeats(%v) got err: %v, wantErr %v", tc.section, err, tc.wantErr)
			}
		})
	}
}

func TestRemoveUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: &pb.User{FirstName: "Bob", LastName: "Johnson", Email: "bob@gmail.com"}, From: "London", To: "France"})

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"ValidRemoval", "bob@gmail.com", false},
		{"InvalidRemoval", "charlie@gmail.com", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.RemoveUser(ctx, &pb.UserRequest{Email: tc.email})
			if (err != nil) != tc.wantErr {
				t.Errorf("RemoveUser(%v) got err: %v, wantErr %v", tc.email, err, tc.wantErr)
			}
		})
	}
}

func TestModifySeat(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: &pb.User{FirstName: "Eve", LastName: "Williams", Email: "eve@gmail.com"}, From: "London", To: "France"})

	tests := []struct {
		name    string
		email   string
		newSeat string
		wantErr bool
	}{
		{"ValidSeatChange", "eve@gmail.com", "B2", false},
		{"InvalidSeatChange", "eve@gmail.com", "B5", true},
		{"NonExistentUser", "unknown@gmail.com", "A1", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.ModifySeat(ctx, &pb.ModifySeatRequest{Email: tc.email, NewSeat: tc.newSeat})
			if (err != nil) != tc.wantErr {
				t.Errorf("ModifySeat(%v, %v) got err: %v, wantErr %v", tc.email, tc.newSeat, err, tc.wantErr)
			}
		})
	}
}
