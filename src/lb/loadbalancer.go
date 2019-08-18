/*Author: Selvam Muthiah
Distributed load balance algorithm */

package lb

import (
	"errors"
	"fmt"
	"strconv"
)

//Server structure contain server details
type Server struct {
	IPAddress   string
	ServerID    int            //Server identifier
	OperStatus  bool           // Server operational status, true is up and false is down
	CPUUsage    int            // represented in percentage 100 is maximum, 80% is threshold, this is just a place holder for this program
	MemoryUsage int            // represented in percentage 100 is maximum, 80% is threshold, this is just a place holder for this program
	Objects     map[int]Object // list of running objects in the Server
}

//Object details
type Object struct {
	OjbectID  int // object identifier
	Signature string
}

//DistributedLB Redistributed Load Balance table
type DistributedLB struct {
	NumServers       int // number of servers
	NumActiveServers int // number of actively running server, it is primary used when server goes down/recover
	NumObject        int //4095 for this project but it can be changed
	ObjectStore      []int
	ServerMap        map[string]Server
}

//InitializeLoadBalancer Initialze distributed load balancer
//This function builds a Map table of servers, similar to Hashtable in C
func (lb *DistributedLB) InitializeLoadBalancer(list int) error {
	if list < 1 || list > 4095 {
		//fmt.Errorf("Number of servers is out of bound %d, it should be in the range of 1-4095", list)
		return errors.New("Number of servers is out of bound")
	}
	hostIDMaxPerByte := 254
	lb.NumServers = list
	lb.NumActiveServers = list
	lb.NumObject = 4095 // This is as per requirement
	lb.ServerMap = make(map[string]Server, 0)
	hostMSB := list / hostIDMaxPerByte // thi is used to created unique server IP address for each server
	hostLSB := list % hostIDMaxPerByte
	serverID := 0
	//create server list
	//update default values , assumption is default up, ideally it should be update by IPC
	for i := 1; i <= hostLSB; i++ {
		IPAddress := "10.1.0." + strconv.Itoa(i)
		serv := Server{}
		serv.IPAddress = IPAddress
		serv.ServerID = serverID
		serverID++
		serv.OperStatus = true
		serv.CPUUsage = 0
		serv.MemoryUsage = 0
		lb.ServerMap[IPAddress] = serv
	}

	for i := 1; i <= hostMSB; i++ {
		for j := 1; j <= hostIDMaxPerByte; j++ {
			IPAddress := "10.1." + strconv.Itoa(i) + "." + strconv.Itoa(i)
			serv := Server{}
			serv.IPAddress = IPAddress
			serv.ServerID = serverID
			serverID++
			serv.OperStatus = true
			serv.CPUUsage = 0
			serv.MemoryUsage = 0
			lb.ServerMap[IPAddress] = serv
		}
	}

	//store object in object store
	for i := 1; i <= 4095; i++ {
		lb.ObjectStore = append(lb.ObjectStore, i)
	}

	return nil
}

//RedistributeObjects Load balance objects across server
//This is redistribute algorithm will be called on each server during following scenarios
// 1) after server boot up/init , 2) Server down notification from other server 3) Server up notification from other server
func (lb *DistributedLB) RedistributeObjects(server *Server) error {
	//Distribute object only for the servers online and cpu, memory threshold less than 75%
	if server.OperStatus != true && server.CPUUsage >= 75 && server.MemoryUsage >= 75 {
		fmt.Printf("server is not up or it reached its max threshold CPU Usage %d, Memory Usage %d", server.CPUUsage, server.MemoryUsage)
		return errors.New("server is not up or it reached its max threshold CPU Usage")
	}
	objectMap := make(map[int]Object)
	for _, ojbectID := range lb.ObjectStore {
		//This ensure that every server will have unique objects, for example if object id is 1, then it will be assigned to only server id 1,
		myObject := ojbectID % lb.NumActiveServers

		if myObject == server.ServerID {
			var obj Object
			obj.OjbectID = ojbectID
			obj.Signature = server.IPAddress
			objectMap[ojbectID] = obj
			//calculate CPU and memory usage percentage, actually this is just place holder
			server.MemoryUsage = (len(objectMap) / 4095) * 100
			server.CPUUsage = (len(objectMap) / 4095) * 100
		}
	}
	server.Objects = objectMap
	return nil
}

//UpdateServers recalulate server id and update object IDs, Actually this will be updated through IPC,
func (lb *DistributedLB) UpdateServers() {
	id := 0
	//re initialze server ID, Assumption: acutally this will by synchronized across the server by IPC
	for serverIP, server := range lb.ServerMap {
		if server.IPAddress != "" && server.OperStatus == true {
			server.ServerID = id
			id++
			lb.ServerMap[serverIP] = server
		}
	}
	for _, server := range lb.ServerMap {
		//assign unique objects per server
		lb.RedistributeObjects(&server)
		lb.ServerMap[server.IPAddress] = server
	}
}

//RedistributeObjServerdown redistribute objects from down server
func (lb *DistributedLB) RedistributeObjServerdown(IPAddress string) error {

	fmt.Printf("Redistribute_obj_Server_down:Server %s down\n\r", IPAddress)
	if lb.ServerMap[IPAddress].OperStatus == false {
		fmt.Printf("Server %s is already down", IPAddress)
		return nil
	}
	lb.NumActiveServers--

	server := lb.ServerMap[IPAddress]
	server.OperStatus = false
	server.ServerID = -1
	lb.ServerMap[IPAddress] = server

	return nil

}

//RedistributeObjServerUp redistribute objects from running servers to new server
func (lb *DistributedLB) RedistributeObjServerUp(IPAddress string) error {

	fmt.Printf("Redistribute_obj_Server_up:Server %s up", IPAddress)
	if lb.ServerMap[IPAddress].OperStatus == true {
		fmt.Printf("Server %s is already up", IPAddress)
		return nil
	}
	lb.NumActiveServers++
	server := lb.ServerMap[IPAddress]
	server.OperStatus = true
	lb.ServerMap[IPAddress] = server

	return nil
}

//GetServerDetails Get all server or per server details
func (lb *DistributedLB) GetServerDetails(IPAddress string) {
	if IPAddress != "" {
		fmt.Printf("Server IP : %s \n\r", IPAddress)
		fmt.Printf("Number of objects: %d\n\r", len(lb.ServerMap[IPAddress].Objects))
		fmt.Printf("List of objects: %v\n\r", lb.ServerMap[IPAddress].Objects)
	} else {
		fmt.Printf("Number of Total servers:%d\n\r", lb.NumServers)
		fmt.Printf("Number of Active Servers:%d\n\r", lb.NumActiveServers)
		fmt.Printf("Number of Total objects:%d\n\r", lb.NumObject)
		for serverIP, server := range lb.ServerMap {
			fmt.Printf("Server IP : %s, Server ID: %d ", serverIP, server.ServerID)
			fmt.Printf("Number of objects: %d\n\r", len(server.Objects))
			//fmt.Printf("List of objects: %v", server.Objects)
		}
	}
}
