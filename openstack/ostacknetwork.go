package openstack

import (
	"log"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
)

func GetNetworkId(provider *gophercloud.ProviderClient, networkname string) string {

	networkclient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create network client:%s\n", err)
	}

	networkid, err := networks.IDFromName(networkclient, networkname)
	if err != nil {
		log.Fatalf("Unable to retrieve network id: %s\n", err)
	}

	networkobj, err := networks.Get(networkclient, networkid).Extract()
	if err != nil {
		log.Fatalf("Unable to retrieve network: %s\n", err)
	}

	return networkobj.ID
}

func SetFloatingIP(provider *gophercloud.ProviderClient, networkname string, floatingname string, serverid string) string {

	networkclient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create network client:%s\n", err)
	}

	networkid := GetNetworkId(provider, networkname)
	floatingnetworkid := GetNetworkId(provider, floatingname)

	var portID string
	var fixedIP string

	timer := 0
	timeout := 60
	for timer < timeout {

		portpages := ports.List(networkclient, ports.ListOpts{
			DeviceID:  serverid,
			NetworkID: networkid,
		})

		err = portpages.EachPage(func(page pagination.Page) (bool, error) {
			portList, err := ports.ExtractPorts(page)
			if err != nil {
				log.Fatalf("Unable to extract ports: %s", err)
			}
			for _, port := range portList {
				portID = port.ID
			}
			return true, nil
		})

		if portID != "" {
			break
		}

		time.Sleep(10 * time.Second)
		timer += 10
	}

	flip, err := floatingips.Create(networkclient, floatingips.CreateOpts{
		FloatingNetworkID: floatingnetworkid,
		PortID:            portID,
		FixedIP:           fixedIP,
	}).Extract()
	if err != nil {
		log.Fatalf("Unable to create a floating IP: %s\n", err)
	}

	return flip.ID
}

func DeleteFloatingIPs(provider *gophercloud.ProviderClient, floatingIDs []string) {

	networkclient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create network client:%s\n", err)
	}

	for i := 0; i < len(floatingIDs); i++ {
		floatingips.Delete(networkclient, floatingIDs[i])
	}
}

func GetFloatingIPs(provider *gophercloud.ProviderClient, floatingIDs []string) []string {

	networkclient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Unable to create network client:%s\n", err)
	}

	var floatingIPs []string
	for i := 0; i < len(floatingIDs); i++ {
		flip, err := floatingips.Get(networkclient, floatingIDs[i]).Extract()
		if err != nil {
			log.Fatalf("Unable to retrieve floating ip: %s\n", err)
		}
		floatingIPs = append(floatingIPs, flip.FloatingIP)

	}

	return floatingIPs

}
