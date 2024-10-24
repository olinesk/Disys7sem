# Group 7 - Mandatory Handin 3

---

## How to run Chitty Chat

---

### 1. Setup the server

Start by setting up a server in your terminal, where users can communicate on:

    go run server/server.go -port=8081

### 2. Setup the Participants (clients)

In this step you will need three additional terminals, and in each terminal run one of the following commands.

In your 2nd terminal:

    go run client/client.go -uPort=8082 -sPort=8081 -name User1

In your 3rd terminal:

    go run client/client.go -uPort=8082 -sPort=8081 -name User2

In your 4th terminal:

    go run client/client.go -uPort=8082 -sPort=8081 -name User3

Now you can send messages to each client!

### Leaving Chitty Chat

If you are done sending messages, you can simple just write `exit` and the connection with be cut.

---
