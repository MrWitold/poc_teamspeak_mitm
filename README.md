# poc_teamspeak_mitm
### Possible way to execute man in the middle attack on teamspeak app
#### Only for educational purposes and for showcase how dangerous is unencrypted traffic

# Warning
It's only proof of concept - may contains a lot of hardcoded variables, bad and inefficient practices 

# Description
The above application allows to listen and modify **unencrypted voice packets**, based on proxy server but also is able to read voice packets from packet dump.

The possible ways usage of application - (Not all included but quite easy to implement):
* Listen voice of users located at other channel, with out of any signs. (Capture outcoming and incoming voice packets for victim, clones and send them to attacker)
* Interupt voice packets (victim will be no able to listen voice of certain users - block incomming voice packet on proxy)
* Inject voice packets (attacker is able to talk as a victm)
* Replay voice packets from packet dump

# Requiremets
* To make any modification of traffic - **attacker and victim** must be connected to server be the proxy
* To listen - App runs in the location where we capture the traffic and **attacker** must be connected to server be the proxy

# Details
* App listen for UDP traffic and then create UDP connection to Teamspeak server
* Each connection has own connotation for Teamspeak server and requests are transferred using channels
* App recognize pattern for voice communication and save need information for later packet injection 
* To listen voice packets - app check for defined udp traffic from certain ip address and redirects them into attacker connection
* Befor sending traffic to the attacker, packet are prepared. (Because be receiving others users packet, we will not hear anything cause packet is not defined for us - that's why we chang "receiver_id","packet_id", "count_id"). Also we modify packets depends on are they incomming or outcomming - so if we want to listen we need to change outcomming traffic into incomming for us
* We do the same as in the step above in case when we want to inject the traffic into other user comunication but the opposite way. (victim doesn't see anything because form servers point of view its looks like he is talking, but his client don't have this information - that's why others see it like he is talking but he doesn't )
* To listen voice packets from dump file, we load the file and inject packets into attacker communication (In this case difficulty is correct speed rate of sending packets).

# Diagram packet flow
[!Diagram](poc_teamspeak_mitm.png)

# How to use
1. Set teamspeak server address ip at main:79
1. `go run main.go`
1. Connect attacker and victim to server by `localhost:324` (if run locally)
1. `output@client_id`     (client_id=32)
1. `attacker@attacker_ip`  (attacker_ip=10.0.0.1:53203)
1. `victim@victim_ip`  (attacker_ip=10.0.0.2:54249)
1. `5` (To start or stop listening)

# How voice teamspeak packets are built, and how to live transcode them for other user.
Soon
