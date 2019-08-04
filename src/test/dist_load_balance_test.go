package test

import "testing"
import "lb"
//TestObjectDistributionAcrossServer test object distribution across servers
func TestObjectDistributionAcrossServer(t *testing.T) {
	distLB := new (lb.Distributed_lb)
	distLB.Initialize_lb(100)
	distLB.Redistribute_objects()
	distLB.GetServerDetails("")
}

//TestServerDown test object distribution across servers when server is down
func TestServerDown(t *testing.T) {
	distLB := new (lb.Distributed_lb)
	distLB.Initialize_lb(100)
	distLB.Redistribute_objects()
	distLB.GetServerDetails("")
	distLB.Redistribute_obj_Server_down("10.1.0.1")
	distLB.GetServerDetails("")
}

//TestObjectDistributionAcrossServer test object distribution across servers
func TestServerUp(t *testing.T) {
	distLB := new (lb.Distributed_lb)
	distLB.Initialize_lb(100)
	distLB.Redistribute_objects()
	distLB.GetServerDetails("")
	distLB.Redistribute_obj_Server_up("10.1.0.1")
	distLB.GetServerDetails("")
}

