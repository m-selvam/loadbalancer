package test

import (
	"fmt"
	"testing"
)
import "lb"
//TestObjectDistributionAcrossServer test object distribution across servers
func TestObjectDistributionAcrossServer(t *testing.T) {
	distLB := new (lb.DistributedLB)
	distLB.InitializeLoadBalancer(5)
	distLB.UpdateServers()
	distLB.GetServerDetails("")
}

//TestServerDown test object distribution across servers when server is down
func TestServerDown(t *testing.T) {
	distLB := new (lb.DistributedLB)
	distLB.InitializeLoadBalancer(5)
	distLB.UpdateServers()
	fmt.Println("Servers status before a server Down\n\r")
	distLB.GetServerDetails("")
	distLB.RedistributeObjServerdown("10.1.0.1")
	distLB.UpdateServers()

	fmt.Println("Servers status after a server Down\n\r")
	distLB.GetServerDetails("")
}

//TestObjectDistributionAcrossServer test object distribution across servers
func TestServerUp(t *testing.T) {
	distLB := new (lb.DistributedLB)
	distLB.InitializeLoadBalancer(5)
	distLB.UpdateServers()
	distLB.RedistributeObjServerdown("10.1.0.1")
	distLB.UpdateServers()
	fmt.Println("Servers status before a server Up \n\r")
	distLB.GetServerDetails("")
	distLB.RedistributeObjServerUp("10.1.0.1")
	distLB.UpdateServers()
	fmt.Println("Servers status before a server Up \n\r")
	distLB.GetServerDetails("")

}

