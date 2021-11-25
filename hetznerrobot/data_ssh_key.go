package hetznerrobot

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key data in OpenSSH format",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key name/desc",
			},

			// read-only / computed
			"fingerprint": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Key fingerprint",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    false,
				Computed:    true,
				Description: "Key size in bits",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Key algorithm type",
			},
		},
	}
}

func dataSourceSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(HetznerRobotClient)

	fingerprint := d.Get("fingerprint").(string)

	sshKey, err := c.getSSHKey(fingerprint)
	if err != nil {
		return fmt.Errorf("Unable to find SSH key with fingerprint %s:\n\t %q", fingerprint, err)
	}

	d.Set("data", sshKey.Data)
	d.Set("fingerprint", fingerprint)
	d.Set("name", sshKey.Name)
	d.Set("size", sshKey.Size)
	d.Set("type", sshKey.Type)
	d.SetId(fingerprint)

	return nil
}
