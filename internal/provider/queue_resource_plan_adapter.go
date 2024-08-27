// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *queueResource) convertToApiModel(_ctx context.Context, plan *QueueResourceModel) (client.Queue, error) {
	return client.Queue{
		Alias: plan.Alias.ValueString(),
	}, nil
}

func (r *queueResource) convertFromApiModel(d *client.Queue, plan *QueueResourceModel) error {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)

	return nil
}
