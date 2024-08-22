// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *queueResource) convert(d *client.Queue, plan *QueueResourceModel) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
}
