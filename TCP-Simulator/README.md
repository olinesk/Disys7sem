#### What are packages in your implementation? What data structure do you use to transmit data and meta-data?
To verify our connections we send bools via channels. We generate random charsets of data and store this in strings and send it via string channels. 
We have chosen only to respresent the seq-ack-seq number as metadata which is stored as integers.

#### Does your implementation use threads or processes? Why is it not realistic to use threads?
We use threads because we are imitating a client and a server and the processes in between. It is not realistic to use threads since the protocol should run on a network, which we do not have here as we are running the 'processes' locally.

#### In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
We are not handling message re-ordering, our program is designed to run first-in-first-out principle as of now. To handle re-ordering we would modify the program to send a sequence number along with each SYN or data packet, which could help identify the correct order.

#### In case messages can be delayed or lost, how does your implementation handle message loss?
We have simulated a packet loss of 15%, if this is the case the error-message is printed and and the 'processes' stop. 

#### Why is the 3-way handshake important?
The 3-way handshake is important because is verifies that the client and server is connected and that data is ready to be send and received accordingly.