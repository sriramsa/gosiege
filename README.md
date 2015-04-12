*(Under development, not ready for public consumption yet. please move on..)*

GoSiege
=======
A cluster based distributed load generator based on siege stress testing tool. Creates a 
siege cluster that can distribute it's load and target a web server for stress. Can scale 
linearly based on the number of nodes in the cluster. Is fault-tolerant and can adjust 
load dynamically based on complete node failures or nodes with reduced capacity.
Provides UI to administer the load and dynamically scale up or down the stress limits
Can also run on a single node.


Design
=======
#### Terminology
**GoSiege Cluster** - A set of nodes configured to generate load against a target.  
**Siege Node**      - An individual machine(or container) with gosiege installed.  
**Session**         - A single load testing session with a set of nodes and targets.  
**GoSiege state**   - Current state of the GoSiege cluster arranged as topics. Contains 
                      current node and session information.  
**Session State**   - Set of information pertaining to a session.

#### Components
1. **GoSiege service daemon**  
Main daemon that spins up the siege tool. Is made of a set of goroutines(threads) that 
do a set of operations. Also listens to local http port for incoming commands for
administration
2. **Admin Web UI**  
Administration UI running on NodeJS that provides interface to Add, Remove or Update
Siege sessions
3. **gosiege** command line tool   
A command line tool that can do all that the Admin UI can do

#### Design Decisions
Initial versions will be provided as **Docker** containers. Each docker container will
act as a siege node.

#### Architecture Pattern
* No component talks to each other directly. Uses channels to communicate.
* **GoSiegeState** is the king here. Every component reacts to and updates state.
    - State is stored in a distributed key-value store(e.g. etcd)
    - Components don't talk to each other, but change the state.
    - Components react to changes in the state
    - Each component subscribes to topics in the state and listen for state change 
      over a channel notifications.
* Each component is a go routine that subscribes to messages to the state engine and listen on a channel
TODO: Using etcd has issues with having to get a token from internet for each new 
cluster instantiation. This may not fly for a non-internet connected machines. 
i.e. isolated networks etc.,

#### Administrative Operations
Admin operations can be done using either the Admin UI or gosiege command line tool.
* GoSiege Cluster
    - Add a node
    - Remove a node
* GoSiege Session
    - Create a new session
    - Stop a session
    - Update a session with
        - New requests
        - New requests per second
        - New target information

#### Data Structures
1. Configuration 
2. Session State Key Value Pairs
2. Commands issued from the web UI

#### Runtime Routine Design
gosiege, the main program, spins up a 'siege' command for each session configured for the session. 

Main spins up these go routines that does the following: 
- *GoSiegeState Handler* 
    - Provides interfaces for reading or changing GoSiegeState
    - Monitors GoSiegeState for updates and notifies components that subscribed to topics
    - Accepts subscriptions for topics    
2. Session handler

### Load distribution
When there are more than on node the go cluster, the load needs to be distributed across equally. There
are a couple of options we can take. 
1. A leader elected within the group periodically decides how much each node should generate.
2. Each node gets around the network table, every 5 seconds, and decide how much each 
should generate to achieve the target.

Choosing option two since there is no need for electing a leader maintaining the same. 
Whoever talks first is the leader each iteration.

#### Load consensus protocol
Each node gets around the network table, every 5 seconds, and agrees upon how much each 
should generate to achieve the target. This is how the dialogue goes, if you evesdrop on them.

Lets say there are 3 nodes NodeA, NodeB and NodeC in a GoSiege Cluster and target requests/sec
this 5000. The GoSiege Session is configured by the manager and is just starting. 

They enter a room with a white board and one marker, figure out we are evesdropping 
so they don't talk but write on the board to communicate.  Each gets one chance with
the marker.

*(Node A gets marker first and writes)*  
*Node A* : OK, target is 5000. Am big and strong, I can do 4000 of it, writes
```
Node A: 4000, Can do 4000
Total: 4000
```
*(Node C gets the marker next)*  
*Node C* : I'm bigger and stronger. I can do all 5000, but since Node B is playing this round, 
lets distribute. I will do 5k\*5k/(4k+5k) = 2778, Node B does 5k*4k/(4k+5k) = 2222. Total: 5k   
```
Node A: 2778, Can do 4000
Node B: 2222, Can do 5000
Total: 5000
```
*(Node B gets the marker next)*  
*Node B* : I'm not as big as you folks, can do only 3000. But I'm old and wise, and I got the last 
turn and what I write goes. I will distribute equally per our strengths. So here is what I write.  
```
Node A: 1666
Node B: 2084
Node C: 1250
Total: 5000
```

##### Partition Issues and Quorum
Issue arise if there is a network disconnection and only Node A and B make it to the table, 
and Node C couldn't. In this scenario A and B will be on one table and C will get 
it's own table. This will end up sending 5k traffic from A and B and C will send 3k, 
making the server smoke and burn. 

So we introduce Quorum in each GoSiege Cluster. Quorum, as you would know, is a minimum
number of votes that should be obtained to be operational. 

In the above GoSiege cluster, the quorum should be set to 2. In this instance of 
Node C will not proceed with sending traffic since it has only one vote on it's decision.

As you would have guessed, whiteboard here is the Session State, that gets updated by 
each node.

Dependencies
------------
**Docker** - Docker is an open-source project that automates the deployment of applications inside software containers, by providing an additional layer of abstraction and automation of operating-system-level virtualization on Linux.  
**Siege** - *Siege is an http load testing and benchmarking utility. It was designed to let web developers measure their code under duress, to see how it will stand up to load on the internet. Siege supports basic authentication, cookies, HTTP, HTTPS and FTP protocols. It lets its user hit a server with a configurable number of simulated clients. Those clients place the server “under siege.”*. Taken from siege [*home page*](http://www.joedog.org/siege-home).   

**etcd** - A distributed key-value store daemon from CoreOS written in Go.

**NodeJS** - Runs admin ui.  

Code Organization
----------
### main
Initializes all the components and starts them up. Has common components like
config, logging etc.,

### state
Distributed state maintainer. Has a set of plugins to use different distributed
key-value store providers: etcd, local file(if not for cluster) etc.,

### manager/session
Session instances that gossip across the cluster and share load generation. Each 
GoSiege session will have a goroutine associated with it on each node maintaining the 
same.
- Each session has a session ID 

### manager/cluster
Administers the GoSiege cluster. Is a set of go routines that listens to commands from 
Admin UI or command line and does CRUD operations on the cluster.  

## REST API
