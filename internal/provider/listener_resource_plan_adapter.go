// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *listenerResource) convertFromApiModel(d *client.ListenerResponse, plan *ListenerResourceModel) (*ListenerResourceModel, error) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.LivenessUrl = types.StringValue(d.LivenessUrl)
	plan.NotifyUrl = types.StringValue(d.NotifyUrl)

	return plan, nil
}

func (r *listenerResource) convertToApiModel(_ context.Context, plan *ListenerResourceModel) (client.ListenerRequest, error) {
	req := client.ListenerRequest{
		Alias:       plan.Alias.ValueString(),
		LivenessUrl: plan.LivenessUrl.ValueString(),
		NotifyUrl:   plan.NotifyUrl.ValueString(),
	}

	return req, nil
}
