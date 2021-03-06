#!/bin/bash -ex

set +x
if ! az account show >/dev/null; then
    exit 1
fi

if [[ -z "$AZURE_SUBSCRIPTION_ID" ]]; then
    echo error: must set AZURE_SUBSCRIPTION_ID
    exit 1
fi

if [[ -z "$AZURE_TENANT_ID" ]]; then
    echo error: must set AZURE_TENANT_ID
    exit 1
fi

if [[ -z "$AZURE_CLIENT_ID" ]]; then
    echo error: must set AZURE_CLIENT_ID
    exit 1
fi

if [[ -z "$AZURE_CLIENT_SECRET" ]]; then
    echo error: must set AZURE_CLIENT_SECRET
    exit 1
fi

if [[ -z "$AZURE_REGION" ]]; then
    echo error: must set AZURE_REGION
    exit 1
fi
valid_regions=(eastus westeurope australiasoutheast)
match=0
for region in "${valid_regions[@]}"; do
    if [[ $region = "$valid_regions" ]]; then
        match=1
        break
    fi
done
if [[ $match = 0 ]]; then
    echo "Error invalid region: must be one of ${valid_regions[@]}"
    exit 1
fi

if [[ -z "$DNS_DOMAIN" ]]; then
    echo error: must set DNS_DOMAIN
    exit 1
fi

if [[ -z "$DNS_RESOURCEGROUP" ]]; then
    echo error: must set DNS_RESOURCEGROUP
    exit 1
fi
set -x

if [[ $# -ne 1 ]]; then
    echo usage: $0 resourcegroup
    exit 1
fi

export RESOURCEGROUP=$1

rm -rf _data
mkdir -p _data/_out

if [[ -z "$NOGROUPTAGS" ]]; then
  ttl=76h
  if [[ -n "$RESOURCEGROUP_TTL" ]]; then
    ttl=$RESOURCEGROUP_TTL
  fi
  GROUPTAGS="--tags now=$(date +%s) ttl=$ttl"
fi
az group create --subscription $AZURE_SUBSCRIPTION_ID -n $RESOURCEGROUP -l $AZURE_REGION $GROUPTAGS >/dev/null

# if AZURE_CLIENT_ID is used as AZURE_AAD_CLIENT_ID, script will reset global team account!
set +x
if [[ "$AZURE_AAD_CLIENT_ID" && "$AZURE_AAD_CLIENT_ID" != "$AZURE_CLIENT_ID" ]]; then
    . <(hack/aad.sh app-update $AZURE_AAD_CLIENT_ID https://$RESOURCEGROUP.$AZURE_REGION.cloudapp.azure.com/oauth2callback/Azure%20AD)
else
    AZURE_AAD_CLIENT_ID=$AZURE_CLIENT_ID
    AZURE_AAD_CLIENT_SECRET=$AZURE_CLIENT_SECRET
fi
set -x

cat >_data/manifest.yaml <<EOF
name: $RESOURCEGROUP
location: $AZURE_REGION
properties:
  openShiftVersion: "$DEPLOY_VERSION"
  fqdn: $RESOURCEGROUP.$AZURE_REGION.cloudapp.azure.com
  authProfile:
    identityProviders:
    - name: Azure AD
      provider:
        kind: AADIdentityProvider
        clientId: $AZURE_AAD_CLIENT_ID
        secret: $AZURE_AAD_CLIENT_SECRET
        tenantId: $AZURE_TENANT_ID
  networkProfile:
    vnetCidr: 10.0.0.0/8
  routerProfiles:
  - name: default
  masterPoolProfile:
    count: 3
    vmSize: Standard_D2s_v3
    subnetCidr: 10.0.0.0/24
  agentPoolProfiles:
  - name: infra
    role: infra
    count: 2
    vmSize: Standard_D2s_v3
    subnetCidr: 10.0.0.0/24
    osType: Linux
  - name: compute
    role: compute
    count: 1
    vmSize: Standard_D2s_v3
    subnetCidr: 10.0.0.0/24
    osType: Linux
EOF

go generate ./...
if [[ -n "$TEST_IN_PRODUCTION" ]]; then
  go test ./test/e2erp -tags e2erp -test.v -ginkgo.v -ginkgo.randomizeAllSpecs -ginkgo.focus=Real -timeout 4h
else
  go run cmd/createorupdate/createorupdate.go

  hack/dns.sh zone-create $RESOURCEGROUP
  hack/dns.sh cname-create $RESOURCEGROUP '*' $RESOURCEGROUP-router.$AZURE_REGION.cloudapp.azure.com
fi

echo
echo  Cluster available at https://$RESOURCEGROUP.$AZURE_REGION.cloudapp.azure.com/
echo
