// Generated by PMS #337
package vpcep

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceVpcepQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the VPC endpoint resource quotas.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of used quotas.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of available quotas.`,
						},
					},
				},
			},
		},
	}
}

type QuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *QuotasDSWrapper {
	return &QuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceVpcepQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newQuotasDSWrapper(d, meta)
	listQuotaDetailsRst, err := wrapper.ListQuotaDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listQuotaDetailsToSchema(listQuotaDetailsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API VPCEP GET /v1/{project_id}/quotas
func (w *QuotasDSWrapper) ListQuotaDetails() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpcep")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/quotas"
	params := map[string]any{
		"type": w.Get("type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *QuotasDSWrapper) listQuotaDetailsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("quotas", schemas.SliceToList(body.Get("quotas.resources"),
			func(quotas gjson.Result) any {
				return map[string]any{
					"type":  quotas.Get("type").Value(),
					"used":  quotas.Get("used").Value(),
					"quota": quotas.Get("quota").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
