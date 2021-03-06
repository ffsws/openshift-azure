package apimanagement

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// UserSubscriptionsClient is the apiManagement Client
type UserSubscriptionsClient struct {
	BaseClient
}

// NewUserSubscriptionsClient creates an instance of the UserSubscriptionsClient client.
func NewUserSubscriptionsClient(subscriptionID string) UserSubscriptionsClient {
	return NewUserSubscriptionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewUserSubscriptionsClientWithBaseURI creates an instance of the UserSubscriptionsClient client.
func NewUserSubscriptionsClientWithBaseURI(baseURI string, subscriptionID string) UserSubscriptionsClient {
	return UserSubscriptionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByUser lists the collection of subscriptions of the specified user.
// Parameters:
// resourceGroupName - the name of the resource group.
// serviceName - the name of the API Management service.
// UID - user identifier. Must be unique in the current API Management service instance.
// filter - | Field        | Supported operators    | Supported functions                         |
// |--------------|------------------------|---------------------------------------------|
// | id           | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
// | name         | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
// | stateComment | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
// | userId       | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
// | productId    | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
// | state        | eq                     |                                             |
// top - number of records to return.
// skip - number of records to skip.
func (client UserSubscriptionsClient) ListByUser(ctx context.Context, resourceGroupName string, serviceName string, UID string, filter string, top *int32, skip *int32) (result SubscriptionCollectionPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: serviceName,
			Constraints: []validation.Constraint{{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil}}},
		{TargetValue: UID,
			Constraints: []validation.Constraint{{Target: "UID", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "UID", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "UID", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil}}},
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("apimanagement.UserSubscriptionsClient", "ListByUser", err.Error())
	}

	result.fn = client.listByUserNextResults
	req, err := client.ListByUserPreparer(ctx, resourceGroupName, serviceName, UID, filter, top, skip)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "ListByUser", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByUserSender(req)
	if err != nil {
		result.sc.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "ListByUser", resp, "Failure sending request")
		return
	}

	result.sc, err = client.ListByUserResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "ListByUser", resp, "Failure responding to request")
	}

	return
}

// ListByUserPreparer prepares the ListByUser request.
func (client UserSubscriptionsClient) ListByUserPreparer(ctx context.Context, resourceGroupName string, serviceName string, UID string, filter string, top *int32, skip *int32) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serviceName":       autorest.Encode("path", serviceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"uid":               autorest.Encode("path", UID),
	}

	const APIVersion = "2016-07-07"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/users/{uid}/subscriptions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByUserSender sends the ListByUser request. The method will close the
// http.Response Body if it receives an error.
func (client UserSubscriptionsClient) ListByUserSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByUserResponder handles the response to the ListByUser request. The method always
// closes the http.Response Body.
func (client UserSubscriptionsClient) ListByUserResponder(resp *http.Response) (result SubscriptionCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByUserNextResults retrieves the next set of results, if any.
func (client UserSubscriptionsClient) listByUserNextResults(lastResults SubscriptionCollection) (result SubscriptionCollection, err error) {
	req, err := lastResults.subscriptionCollectionPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "listByUserNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByUserSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "listByUserNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByUserResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.UserSubscriptionsClient", "listByUserNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByUserComplete enumerates all values, automatically crossing page boundaries as required.
func (client UserSubscriptionsClient) ListByUserComplete(ctx context.Context, resourceGroupName string, serviceName string, UID string, filter string, top *int32, skip *int32) (result SubscriptionCollectionIterator, err error) {
	result.page, err = client.ListByUser(ctx, resourceGroupName, serviceName, UID, filter, top, skip)
	return
}
