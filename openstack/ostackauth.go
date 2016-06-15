package openstack

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

func CreateClient() (*gophercloud.ProviderClient, error) {

	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println("Unable to retrieve auth options from environment: %s", err)
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		fmt.Println("Unable to retrieve openstack client: %s", err)
	}

	return provider, nil
}
