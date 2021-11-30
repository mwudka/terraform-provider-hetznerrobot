package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSSHKeyImportState,
		},

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

func resourceSSHKeyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	fingerprint := d.Id()

	sshKey, err := c.getSSHKey(fingerprint)
	if err != nil {
		return nil, err
	}

	d.Set("data", sshKey.Data)
	d.Set("fingerprint", sshKey.FingerPrint)
	d.Set("name", sshKey.Name)
	d.Set("size", sshKey.Size)
	d.Set("type", sshKey.Type)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	data := d.Get("data").(string)
	name := d.Get("name").(string)

	sshKey, err := c.addSSHKey(name, data)
	if err != nil {
		return diag.FromErr(err)
	}

	fingerprint := sshKey.FingerPrint
	d.Set("fingerprint", fingerprint)
	d.Set("size", sshKey.Size)
	d.Set("type", sshKey.Type)
	d.SetId(fingerprint)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	fingerprint := d.Id()

	sshKey, err := c.getSSHKey(fingerprint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("data", sshKey.Data)
	d.Set("fingerprint", fingerprint)
	d.Set("name", sshKey.Name)
	d.Set("size", sshKey.Size)
	d.Set("type", sshKey.Type)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	fingerprint := d.Id()
	name := d.Get("name").(string)

	if err := c.updSSHKey(name, fingerprint); err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	fingerprint := d.Id()

	err := c.delSSHKey(fingerprint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
