package main

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "train-ticket-app/pb/proto"

	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Example user
	user := &pb.User{FirstName: "Andrus", LastName: "Taylor", Email: "andrustaylor90@gmail.com"}

	// 1. Purchase Ticket
	receipt, err := client.PurchaseTicket(ctx, &pb.PurchaseRequest{User: user, From: "London", To: "France"})
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	fmt.Printf("Purchased Ticket: %+v\n", receipt)

	// 2. Get Receipt
	gotReceipt, err := client.GetReceipt(ctx, &pb.UserRequest{Email: "andrustaylor90@gmail.com"})
	if err != nil {
		log.Fatalf("could not get receipt: %v", err)
	}
	fmt.Printf("Got Receipt: %+v\n", gotReceipt)

	// 3. View Seats
	seatsA, err := client.ViewSeats(ctx, &pb.SectionRequest{Section: "A"})
	if err != nil {
		log.Fatalf("could not view seats: %v", err)
	}
	fmt.Printf("Seats in Section A: %+v\n", seatsA)

	// 4. Modify Seat (this will fail as the user is already removed, but included for completeness)
	modifyResponse, err := client.ModifySeat(ctx, &pb.ModifySeatRequest{Email: "andrustaylor90@gmail.com", NewSeat: "B1"})
	if err != nil {
		fmt.Printf("could not modify seat (as expected since user is removed): %v\n", err)
	} else {
		fmt.Printf("Modify Seat Response: %+v\n", modifyResponse)
	}

	// 5. Get Receipt
	gotReceiptAfterUpdating, err := client.GetReceipt(ctx, &pb.UserRequest{Email: "andrustaylor90@gmail.com"})
	if err != nil {
		log.Fatalf("could not get receipt: %v", err)
	}
	fmt.Printf("Got Receipt: %+v\n", gotReceiptAfterUpdating)

	// 6. Remove User
	removeResponse, err := client.RemoveUser(ctx, &pb.UserRequest{Email: "andrustaylor90@gmail.com"})
	if err != nil {
		log.Fatalf("could not remove user: %v", err)
	}
	fmt.Printf("Remove User Response: %+v\n", removeResponse)

	// 7. Try to Modify Seat (this will fail as the user is already removed, but included for completeness)
	modifyResponseAgain, err := client.ModifySeat(ctx, &pb.ModifySeatRequest{Email: "andrustaylor90@gmail.com", NewSeat: "B1"})
	if err != nil {
		fmt.Printf("could not modify seat (as expected since user is removed): %v\n", err)
	} else {
		fmt.Printf("Modify Seat Response: %+v\n", modifyResponseAgain)
	}
}
