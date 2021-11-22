package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  schema.EnvDefaultFunc("HETZNERROBOT_URL", "https://robot-ws.your-server.de"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hetznerrobot_boot":     resourceBoot(),
			"hetznerrobot_firewall": resourceFirewall(),
			"hetznerrobot_key":      resourceSSHKey(),
			"hetznerrobot_server":   resourceServer(),
			"hetznerrobot_vswitch":  resourceVSwitch(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hetznerrobot_boot":    dataBoot(),
			"hetznerrobot_key":     dataSSHKey(),
			"hetznerrobot_server":  dataServer(),
			"hetznerrobot_vswitch": dataVSwitch(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	url := d.Get("url").(string)

	var diags diag.Diagnostics

	return NewHetznerRobotClient(username, password, url), diags
}
