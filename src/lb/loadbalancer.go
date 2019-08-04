//Author: Selvam Muthiah
//Distributed load balance algorithm
package lb

import (
	"errors"
	"fmt"
	"strconv"
)

//Server structure contain server details
type Server struct {
	Ip_address   string
	Oper_status  bool     // Server operational status, true is up and false is down
	Cpu_usage    int      // represented in percentage 100 is maximum, 80% is threshold, this is just a place holder for this program
	Memory_usage int      // represented in percentage 100 is maximum, 80% is threshold, this is just a place holder for this program
	Objects      []Object // list of running objects in the Server
}

//Object details
type Object struct {
	Object_id int // object identifier
	signature string
}

//Distributed_lb Redistributed Load Balance table
type Distributed_lb struct {
	Num_servers        int // number of servers
	Num_active_servers int // number of actively running server, it is primary used when server goes down/recover
	Num_object         int //4095 for this project but it can be changed
	Server_map         map[string]Server
}

//Initialize_lb Initialze distributed load balancer
//This function builds a Map table of servers, similar to Hashtable in C
func (lb *Distributed_lb) Initialize_lb(list int) error {
	if list < 1 || list > 4095 {
		fmt.Errorf("Number of servers is out of bound %d, it should be in the range of 1-4095", list)
		return errors.New("Number of servers is out of bound")
	}
	host_id_max_per_byte := 254
	lb.Num_servers = list
	lb.Num_object = 4095 // This is as per requirement
	lb.Server_map = make(map[string]Server, 0)
	host_msb := list / host_id_max_per_byte // thi is used to created unique server IP address for each server
	host_lsb := list % host_id_max_per_byte
	//create server list
	//update default values , assumption is default up, ideally it should be update by IPC
	for i := 1; i <= host_lsb; i++ {
		ip_address := "10.1.0." + strconv.Itoa(i)
		serv := Server{}
		serv.Ip_address = ip_address
		serv.Oper_status = true
		serv.Cpu_usage = 0
		serv.Memory_usage = 0
		lb.Server_map[ip_address] = serv
	}

	for i := 1; i <= host_msb; i++ {
		for j := 1; j <= host_id_max_per_byte; j++ {
			ip_address := "10.1." + strconv.Itoa(i) + "." + strconv.Itoa(i)
			serv := Server{}
			serv.Ip_address = ip_address
			serv.Oper_status = true
			serv.Cpu_usage = 0
			serv.Memory_usage = 0
			lb.Server_map[ip_address] = serv
		}
	}
	return nil
}

//Redistribute_objects Load balance objects across server
func (lb *Distributed_lb) Redistribute_objects() {
	num_obj_per_server := lb.Num_object / lb.Num_servers
	num_remain_obj := lb.Num_object % lb.Num_servers
	object_id := 1
	extra_object_offset := num_obj_per_server * lb.Num_servers
	for server_ip, server := range lb.Server_map {
		//Distribute object only for the servers online and cpu, memory threshold less than 75%
		if server.Oper_status == true && server.Cpu_usage <= 75 && server.Memory_usage <= 75 {
			for count := 1; count <= num_obj_per_server; count++ {

				var obj Object
				lb.Num_active_servers++
				obj.Object_id = object_id
				object_id++
				obj.signature = server_ip
				server.Objects = append(server.Objects, obj)
				//calculate CPU and memory usage percentage, actually this is just place holder, we should get it from actual server
				server.Memory_usage = (len(server.Objects) / 4095) * 100
				server.Cpu_usage = (len(server.Objects) / 4095) * 100
			}
			//Add an object from extra pool of objects
			if num_remain_obj > 0 {

				var obj Object
				extra_object_offset++
				obj.Object_id = extra_object_offset
				obj.signature = server_ip
				server.Objects = append(server.Objects, obj)
				//calculate CPU and memory usage percentage, actually this is just place holder, we should get it from actual server
				server.Memory_usage = (len(server.Objects) / 4095) * 100
				server.Cpu_usage = (len(server.Objects) / 4095) * 100
			}
		}
	}
}

//Redistribute_obj_Server_down redistribute objects from down server
func (lb *Distributed_lb) Redistribute_obj_Server_down(ip_address string) error {

	fmt.Printf("Redistribute_obj_Server_down:Server %s down", ip_address)
	if lb.Server_map[ip_address].Oper_status == false {
		fmt.Printf("Server %s is already down", ip_address)
		return nil
	}
	lb.Num_active_servers--
	server_down := lb.Server_map[ip_address]
	server_down.Oper_status = false

	num_obj_per_server := len(lb.Server_map[ip_address].Objects) / lb.Num_active_servers
	num_remain_obj := len(lb.Server_map[ip_address].Objects) % lb.Num_active_servers
	Index := 1
	extra_object_offset := num_obj_per_server * lb.Num_servers
	for _, server := range lb.Server_map {
		//Distribute object only for the servers online and cpu, memory threshold less than 75%
		if server.Oper_status == true && server.Cpu_usage <= 75 && server.Memory_usage <= 75 {
			for count := 1; count <= num_obj_per_server; count++ {

				server.Objects = append(server.Objects, server_down.Objects[Index])
				Index++
				//calculate CPU and memory usage percentage, actually this is just place holder, we should get it from actual server
				server.Memory_usage = (len(server.Objects) / 4095) * 100
				server.Cpu_usage = (len(server.Objects) / 4095) * 100
			}
			//Add an object from extra pool of objects
			if num_remain_obj > 0 {

				extra_object_offset++
				server.Objects = append(server.Objects, server_down.Objects[extra_object_offset])
				//calculate CPU and memory usage percentage, actually this is just place holder, we should get it from actual server
				server.Memory_usage = (len(server.Objects) / 4095) * 100
				server.Cpu_usage = (len(server.Objects) / 4095) * 100
			}
		}
	}
	server_down.Objects = nil
	server_down.Cpu_usage = 0
	server_down.Memory_usage = 0

	return nil

}

//Redistribute_obj_Server_up redistribute objects from running servers to new server
func (lb *Distributed_lb) Redistribute_obj_Server_up(ip_address string) error {

	fmt.Printf("Redistribute_obj_Server_up:Server %s up", ip_address)
	if lb.Server_map[ip_address].Oper_status == true {
		fmt.Printf("Server %s is already up", ip_address)
		return nil
	}
	lb.Num_active_servers++
	server_up := lb.Server_map[ip_address]
	server_up.Oper_status = true
	//copy objects from already running servers
	num_obj_per_server := lb.Num_object / lb.Num_active_servers
	num_remain_obj := lb.Num_object % lb.Num_active_servers
	start_index := 0

	for _, server := range lb.Server_map {
		//Get objects from online servers
		if server.Oper_status == true {

			offset_index := num_obj_per_server
			//add one more object if there is a un even number of objects
			if num_remain_obj > 0 {
				offset_index++
				num_remain_obj--
			}
			original := server.Objects
			if len(server.Objects) > offset_index {
				max_num_obj := len(server.Objects) - offset_index
				server_up.Objects[start_index:] = original[max_num_obj:]
				//change the start index to copy from next server
				start_index = len(server_up.Objects)
				server.Objects[:offset_index] = original[:offset_index]
			}

			//Re calculate CPU and memory usage percentage,
			server.Memory_usage = (len(server.Objects) / 4095) * 100
			server.Cpu_usage = (len(server.Objects) / 4095) * 100

			server_up.Memory_usage = (len(server_up.Objects) / 4095) * 100
			server_up.Cpu_usage = (len(server_up.Objects) / 4095) * 100
		}
	}
	return nil
}

//GetServerDetails Get all server or per server details
func (lb *Distributed_lb) GetServerDetails(ip_address string) {
	if ip_address != "" {
		fmt.Printf("Server IP : %s", ip_address)
		fmt.Printf("Number of objects: %d", len(lb.Server_map[ip_address].Objects))
		fmt.Printf("List of objects: %v", lb.Server_map[ip_address].Objects)
	} else {
		fmt.Printf("Number of running Servers:%d", lb.Num_active_servers)
		fmt.Printf("Number of Total servers:%d", lb.Num_servers)
		fmt.Printf("Number of Total objects:%d", lb.Num_object)
		for server_ip, server := range lb.Server_map {
			fmt.Printf("Server IP : %s", server_ip)
			fmt.Printf("Number of objects: %d", len(server.Objects))
			fmt.Printf("List of objects: %v", server.Objects)
		}
	}
}
