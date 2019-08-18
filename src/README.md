Requirement:
==========
Distributed Load-balancer:

Design and implement a load-balancer algorithm that distributes a set of objects say [O1, O2, ….., On] across a set of servers [S1, S2, ….., Sm]. 
For implementation, assume that the objects are integer values in a bounded range say [0…..4095], and servers are unique IPv4 addresses.

Design Considerations:
This algorithm runs in a distributed manner on each server. Each server receives the full set of objects [O1, O2, ….., On] and the full set of servers [S1, S2, ….., Sm]. Each server runs the same algorithm and should deterministically arrive at the same result, with respect to object assignments to each server (including itself). The reason is to minimize the inter-server communication to a minimum.

Algorithm should handle server failures and recovery and re-assign objects accordingly. An important design consideration is to try and minimize disruption on a server failure, aka, limit re-assignment to a minimum set of objects as far as possible. Again, keep in mind that the result needs be deterministic and consistent across all the servers.

Protocol between the servers to detect a server failure or recovery is beyond the scope of this problem. Assume a server failure OR recovery is magically relayed to each servers. From the point of view of your implementation, this is one of the events that the algorithm should handle. How this event is produced is not relevant.

Solution:
=========

Alogrithm used: Objects are redistributed across servers based on Round Robin plus server status, CPU and Memory usage.

Assumption: Server IPs and Server IDs are unique which are linearly created/updated during initialization, server down or server up, Acutally this suppose to happen via IPC events.
Approach 1:
==========
1) Servers are mainted in the map, indexed with server IP which are dynamically created based on number of number of servers
2) Each server has map of objects, cpu and memory status
3) Every server allocates objects from pool of objects based on Object_ID modulo by number of active servers 
4) When server goes down/up, number of active servers, server ID are updated, then objects are redistributed.


Language: Golang 1.11, its a modern language has lots of built-in echo system to write these type of application

DataStructure Used : Map and slices, Its very similar to Hast table with list of linked list in c


Test cases:

Test cases are automated using testing package, included in test sub directory

Test case 1:
Configure number of server as 5, check objects are redistributed.
Log:
$ go test -run "TestObjectDistributionAcrossServer"
Number of Total servers:5
Number of Active Servers:5
Number of Total objects:4095
Server IP : 10.1.0.1, Server ID: 0 Number of objects: 819
Server IP : 10.1.0.2, Server ID: 1 Number of objects: 819
Server IP : 10.1.0.3, Server ID: 2 Number of objects: 819
Server IP : 10.1.0.4, Server ID: 3 Number of objects: 819
Server IP : 10.1.0.5, Server ID: 4 Number of objects: 819
PASS
ok      test    0.481s



Test case 2:
Bring down a server, check number of active server is reduced and objects are redistbuted to all other active servers.
$ go test -run "TestServerDown"
Servers status before a server Down

Number of Total servers:5
Number of Active Servers:5
Number of Total objects:4095
Server IP : 10.1.0.5, Server ID: 1 Number of objects: 819
Server IP : 10.1.0.1, Server ID: 2 Number of objects: 819
Server IP : 10.1.0.2, Server ID: 3 Number of objects: 819
Server IP : 10.1.0.3, Server ID: 4 Number of objects: 819
Server IP : 10.1.0.4, Server ID: 0 Number of objects: 819
Redistribute_obj_Server_down:Server 10.1.0.1 down
Servers status after a server Down

Number of Total servers:5
Number of Active Servers:4
Number of Total objects:4095
Server IP : 10.1.0.2, Server ID: 0 Number of objects: 1023
Server IP : 10.1.0.3, Server ID: 1 Number of objects: 1024
Server IP : 10.1.0.4, Server ID: 2 Number of objects: 1024
Server IP : 10.1.0.5, Server ID: 3 Number of objects: 1024
Server IP : 10.1.0.1, Server ID: -1 Number of objects: 0
PASS
ok      test    0.482s

SMUTHIAH@SMUTHIAH-PC MINGW64 ~/go/awesomeProject1/loadbalancer/src/test (master)
$

Test case 3:
Bring up the server , check objects are rebalanced again.
$ go test -run "TestServerUp"
Redistribute_obj_Server_down:Server 10.1.0.1 down
Servers status before a server Up

Number of Total servers:5
Number of Active Servers:4
Number of Total objects:4095
Server IP : 10.1.0.3, Server ID: 1 Number of objects: 1024
Server IP : 10.1.0.4, Server ID: 2 Number of objects: 1024
Server IP : 10.1.0.5, Server ID: 3 Number of objects: 1024
Server IP : 10.1.0.1, Server ID: -1 Number of objects: 0
Server IP : 10.1.0.2, Server ID: 0 Number of objects: 1023
Redistribute_obj_Server_up:Server 10.1.0.1 upServers status before a server Up

Number of Total servers:5
Number of Active Servers:5
Number of Total objects:4095
Server IP : 10.1.0.1, Server ID: 2 Number of objects: 819
Server IP : 10.1.0.2, Server ID: 3 Number of objects: 819
Server IP : 10.1.0.3, Server ID: 4 Number of objects: 819
Server IP : 10.1.0.4, Server ID: 0 Number of objects: 819
Server IP : 10.1.0.5, Server ID: 1 Number of objects: 819
PASS
ok      test    0.449s


Test case 4:
Bring down 10 servers, check number of active servers and its object.

Test case 5:
Bring up 10 servers, check number of active servers and its object.

Test case 6:
configure number of server as out of bound, 0 or 4096.(nagative test case).


Approach 2:
==========
same algorimthm, but implemention is different,
Additionaly 
1) maintain array of all objects in the global strucutre similar to Servers which has slices of objects, these slices can be increased/decreased from global array based on number of active servers, this way we can avoid looping objects.

2) use "goroutine" to multithread the application to handle uses cases like multiple servers goes down/up at the same time.
(It is not yet implemented)


Pending activity:
================
Testing, found vet error, needs to be fixed
