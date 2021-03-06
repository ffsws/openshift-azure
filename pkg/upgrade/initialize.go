package upgrade

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/util/azureclient"
	"github.com/openshift/openshift-azure/pkg/util/azureclient/storage"
)

func (u *simpleUpgrader) InitializeCluster(ctx context.Context, cs *api.OpenShiftManagedCluster) error {
	if u.accountsClient == nil {
		authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
		if err != nil {
			return err
		}

		u.accountsClient = azureclient.NewAccountsClient(cs.Properties.AzProfile.SubscriptionID, authorizer, u.pluginConfig.AcceptLanguages)
	}
	if u.storageClient == nil {
		keys, err := u.accountsClient.ListKeys(ctx, cs.Properties.AzProfile.ResourceGroup, cs.Config.ConfigStorageAccount)
		if err != nil {
			return err
		}

		u.storageClient, err = storage.NewClient(cs.Config.ConfigStorageAccount, *(*keys.Keys)[0].Value, storage.DefaultBaseURL, storage.DefaultAPIVersion, true)
		if err != nil {
			return err
		}
	}
	bsc := u.storageClient.GetBlobService()

	// etcd data container
	c := bsc.GetContainerReference("etcd")
	_, err := c.CreateIfNotExists(nil)
	if err != nil {
		return err
	}

	// cluster config container
	c = bsc.GetContainerReference("config")
	_, err = c.CreateIfNotExists(nil)
	if err != nil {
		return err
	}

	b := c.GetBlobReference("config")

	csj, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	return b.CreateBlockBlobFromReader(bytes.NewReader(csj), nil)
}
