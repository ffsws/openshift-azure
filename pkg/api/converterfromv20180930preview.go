package api

import (
	v20180930preview "github.com/openshift/openshift-azure/pkg/api/2018-09-30-preview/api"
)

// ConvertFromV20180930preview converts from a
// v20180930preview.OpenShiftManagedCluster to an OpenShiftManagedCluster.
func ConvertFromV20180930preview(oc *v20180930preview.OpenShiftManagedCluster) *OpenShiftManagedCluster {
	cs := &OpenShiftManagedCluster{
		ID:       oc.ID,
		Location: oc.Location,
		Name:     oc.Name,
		Tags:     oc.Tags,
		Type:     oc.Type,
	}

	if oc.Plan != nil {
		cs.Plan = &ResourcePurchasePlan{
			Name:          oc.Plan.Name,
			Product:       oc.Plan.Product,
			PromotionCode: oc.Plan.PromotionCode,
			Publisher:     oc.Plan.Publisher,
		}
	}

	if oc.Properties != nil {
		cs.Properties = &Properties{
			ProvisioningState: ProvisioningState(oc.Properties.ProvisioningState),
			OpenShiftVersion:  oc.Properties.OpenShiftVersion,
			PublicHostname:    oc.Properties.PublicHostname,
			FQDN:              oc.Properties.FQDN,
		}

		if oc.Properties.NetworkProfile != nil {
			cs.Properties.NetworkProfile = &NetworkProfile{
				VnetCIDR:   oc.Properties.NetworkProfile.VnetCIDR,
				PeerVnetID: oc.Properties.NetworkProfile.PeerVnetID,
			}
		}

		cs.Properties.RouterProfiles = make([]RouterProfile, len(oc.Properties.RouterProfiles))
		for i, rp := range oc.Properties.RouterProfiles {
			cs.Properties.RouterProfiles[i] = RouterProfile{
				Name:            rp.Name,
				PublicSubdomain: rp.PublicSubdomain,
				FQDN:            rp.FQDN,
			}
		}

		cs.Properties.AgentPoolProfiles = make([]AgentPoolProfile, 0, len(oc.Properties.AgentPoolProfiles)+1)
		for _, app := range oc.Properties.AgentPoolProfiles {
			cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, AgentPoolProfile{
				Name:       app.Name,
				Count:      app.Count,
				VMSize:     VMSize(app.VMSize),
				SubnetCIDR: app.SubnetCIDR,
				OSType:     OSType(app.OSType),
				Role:       AgentPoolProfileRole(app.Role),
			})
		}

		if oc.Properties.MasterPoolProfile != nil {
			cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, AgentPoolProfile{
				Name:       string(AgentPoolProfileRoleMaster),
				Count:      oc.Properties.MasterPoolProfile.Count,
				VMSize:     VMSize(oc.Properties.MasterPoolProfile.VMSize),
				SubnetCIDR: oc.Properties.MasterPoolProfile.SubnetCIDR,
				OSType:     OSTypeLinux,
				Role:       AgentPoolProfileRoleMaster,
			})
		}

		if oc.Properties.AuthProfile != nil {
			cs.Properties.AuthProfile = &AuthProfile{}

			cs.Properties.AuthProfile.IdentityProviders = make([]IdentityProvider, len(oc.Properties.AuthProfile.IdentityProviders))
			for i, ip := range oc.Properties.AuthProfile.IdentityProviders {
				cs.Properties.AuthProfile.IdentityProviders[i].Name = ip.Name
				switch provider := ip.Provider.(type) {
				case (*v20180930preview.AADIdentityProvider):
					cs.Properties.AuthProfile.IdentityProviders[i].Provider = &AADIdentityProvider{
						Kind:     provider.Kind,
						ClientID: provider.ClientID,
						Secret:   provider.Secret,
						TenantID: provider.TenantID,
					}

				default:
					panic("authProfile.identityProviders conversion failed")
				}
			}
		}
	}

	return cs
}
