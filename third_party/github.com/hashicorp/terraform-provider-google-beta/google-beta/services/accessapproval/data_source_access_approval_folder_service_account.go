// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accessapproval

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func DataSourceAccessApprovalFolderServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessApprovalFolderServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"folder_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessApprovalFolderServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessApprovalBasePath}}folders/{{folder_id}}/serviceAccount")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("AccessApprovalFolderServiceAccount %q", d.Id()))
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("account_email", res["accountEmail"]); err != nil {
		return fmt.Errorf("Error setting account_email: %s", err)
	}
	d.SetId(res["name"].(string))

	return nil
}