# Sr. Software Engineer – Backend - Coding assessment

## Requirements: 
1.	Code must be published in Github with a link we can access (use public repo).
2.	Code must compile with some effort on unit tests, doesn’t have to be 100%, but it shouldn’t be 0%.
3.	Here is the application we want to build for the interview
4.	Please code this with Golang and gRPC
5.	Adding a persistence layer will be cumbersome, so just store the data in your current session/in memory.
6.	The results can be in the console output

## App to be coded
Background: All API referenced are gRPC APIs, not REST ones.
I want to board a train from London to France. The train ticket will cost $20.  
1.	Create API where you can submit a purchase for a ticket.  Details included in the receipt are : From, To, User , price paid. ( User should include first and last name, email address )
2.	The user is allocated a seat in the train.  Assume the train has only 2 sections, section A and section B.
3.	An API that shows the details of the receipt for the user
4.	An API that lets you view the users and seat they are allocated by the requested section
5.	An API to remove a user from the train
6.	An API to modify a user’s seat

## Test
```
go test .\server\
```

## Run Server
```
go run .\server\
```

## Run Client
```
go run .\client\
```