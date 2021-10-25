package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lenstra/hetzner"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HETZNERROBOT_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HETZNERROBOT_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hetznerrobot_firewall": resourceFirewall(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := hetzner.DefaultConfig()
	config.User = d.Get("username").(string)
	config.Password = d.Get("password").(string)

	return hetzner.NewClient(config), nil
}
