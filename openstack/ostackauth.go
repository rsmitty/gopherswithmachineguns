package openstack

import (
	"log"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

func CreateClient() *gophercloud.ProviderClient {

	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("Unable to retrieve auth options from environment: %s", err)
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		log.Fatalf("Unable to retrieve openstack client: %s", err)
	}

	return provider
}
