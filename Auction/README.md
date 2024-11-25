# Distributed auction system - How to run the system:

## Starting servers

Go to the root folder in 3 seperate shells and run these following commands in its seperate shell to start 3 servers:

go run server/server.go 0

go run server/server.go 1

go run server/server.go 2

## Starting the Client/Frontend pairs

For each client you must run the following two commands in 2 separate shells:

go run main.go 0

go run client/client.go 0

For each new client you want to create to enter the auction, you must repeat the above 2 commands just with a new id (i.e. instead of 0 write 1 and so on).