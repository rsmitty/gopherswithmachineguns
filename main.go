package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"

	"github.com/rsmitty/gopherswithmachineguns/openstack"
	"github.com/rsmitty/gopherswithmachineguns/ssh"
)

type ServerFlags struct {
	Count           int
	ImageName       string
	FlavorName      string
	Network         string
	FloatingNetwork string
	KeyName         string
	SSHUser         string
	SSHKeyPath      string
	Endpoint        string
	SimulRequests   int
	TotalRequests   int
}

func ParseFlags(serverflags *ServerFlags) {

	numGophers := flag.Int("count", 1, "# of OpenStack VMs to create")
	imgName := flag.String("image", "", "Name of image to use")
	imgFlavor := flag.String("flavor", "", "Name of flavor to use")
	network := flag.String("network", "", "Network to use")
	floatingNetwork := flag.String("floating-network", "", "Floating network to use")
	keyName := flag.String("keyname", "", "Openstack keypair to attach to servers")
	sshUser := flag.String("sshuser", "", "Name of ssh user for image")
	sshKeyPath := flag.String("sshkey", "", "Path to ssh key")
	endpoint := flag.String("endpoint", "", "Endpoint to attack")
	simulRequests := flag.Int("sim-reqs", 1, "# of simultaneous requests")
	totalRequests := flag.Int("tot-reqs", 1, "# of total requests")
	flag.Parse()

	serverflags.ImageName = *imgName
	serverflags.FlavorName = *imgFlavor
	serverflags.Network = *network
	serverflags.FloatingNetwork = *floatingNetwork
	serverflags.KeyName = *keyName
	serverflags.Count = *numGophers
	serverflags.SSHUser = *sshUser
	serverflags.SSHKeyPath = *sshKeyPath
	serverflags.Endpoint = *endpoint
	serverflags.SimulRequests = *simulRequests
	serverflags.TotalRequests = *totalRequests
}

func attack(host string, sshuser string, sshkey string, endpoint string, simreq int, totreq int) {
	sshclient := ssh.ConnectSSH(host, sshuser, sshkey)
	abCommand := "ab -r -n " + strconv.Itoa(totreq) + " -c " + strconv.Itoa(simreq) + " " + endpoint
	ssh.IssueCommand(sshclient, abCommand)
	ssh.CloseSSH(sshclient)

}
func main() {

	//Parse command line input
	serverStruct := new(ServerFlags)
	ParseFlags(serverStruct)

	//Create servers and floating IPs as necessary
	provider, _ := openstack.CreateClient()
	serverIDs, floatingIDs := openstack.CreateServers(provider, serverStruct.Count, serverStruct.ImageName, serverStruct.FlavorName, serverStruct.Network, serverStruct.FloatingNetwork, serverStruct.KeyName)

	fmt.Println(serverIDs)
	fmt.Println(floatingIDs)

	//Retrieve network info
	var ipList []string
	if len(floatingIDs) > 0 {
		ipList = openstack.GetFloatingIPs(provider, floatingIDs)
	} else {
		ipList = openstack.GetPrivateIPs(provider, serverIDs)
	}
	fmt.Println(ipList)

	var wg sync.WaitGroup
	wg.Add(len(ipList))

	for i := 0; i < len(ipList); i++ {
		go func(i int) {
			defer wg.Done()
			attack(ipList[i], serverStruct.SSHUser, serverStruct.SSHKeyPath, serverStruct.Endpoint, serverStruct.SimulRequests, serverStruct.TotalRequests)
		}(i)
	}

	wg.Wait()

	//Cleanup servers after attack
	openstack.DeleteServers(provider, serverIDs)
	if len(floatingIDs) > 0 {
		openstack.DeleteFloatingIPs(provider, floatingIDs)
	}
}
