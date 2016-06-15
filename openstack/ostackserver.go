package openstack

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

func ListServers(provider *gophercloud.ProviderClient) {

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create a compute client: %s\n", err)
	}
	opts := servers.ListOpts{}

	pager := servers.List(client, opts)

	pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, _ := servers.ExtractServers(page)

		for _, s := range serverList {
			log.Fatalf(s.Name)
		}
		return true, nil
	})

}

func CreateServers(provider *gophercloud.ProviderClient, count int, image string, flavor string, networkname string, floatingnetwork string, keyname string) ([]string, []string) {

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create a compute client: %s\n", err)
	}

	netid := GetNetworkId(provider, networkname)

	servernetwork := []servers.Network{
		{
			UUID: netid,
		},
	}

	var serverIDs []string
	var floatingIDs []string

	for i := 0; i < count; i++ {
		serveropts := servers.CreateOpts{
			Name:       "gwmg-" + strconv.Itoa(i),
			ImageName:  image,
			FlavorName: flavor,
			Networks:   servernetwork,
		}

		server, err := servers.Create(client, keypairs.CreateOptsExt{
			serveropts,
			keyname,
		}).Extract()
		if err != nil {
			log.Fatalf("Unable to create a compute client: %s\n", err)
		}
		fmt.Printf("Server %v created\n", i)

		if floatingnetwork != "" {
			floatingID := SetFloatingIP(provider, networkname, floatingnetwork, server.ID)
			floatingIDs = append(floatingIDs, floatingID)
		}

		serverIDs = append(serverIDs, server.ID)
	}

	return serverIDs, floatingIDs
}

func DeleteServers(provider *gophercloud.ProviderClient, serverIDs []string) {

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create a compute client: %s\n", err)
	}

	for i := 0; i < len(serverIDs); i++ {

		fmt.Printf("Deleting server %v\n", i)
		servers.Delete(client, serverIDs[i])
	}
}

func GetPrivateIPs(provider *gophercloud.ProviderClient, serverIDs []string) []string {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create a compute client: %s\n", err)
	}

	var privateIPs []string
	for i := 0; i < len(serverIDs); i++ {
		server, err := servers.Get(client, serverIDs[i]).Extract()

		if err != nil {
			log.Fatalf("Unable to retrieve server: %s\n", err)
		}
		privateIPs = append(privateIPs, server.AccessIPv4)

	}

	return privateIPs

}
