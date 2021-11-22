package hetznerrobot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceServerVSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerVSwitchCreate,
		ReadContext:   resourceServerVSwitchRead,
		UpdateContext: resourceServerVSwitchUpdate,
		DeleteContext: resourceServerVSwitchDelete,

		Importer: &schema.ResourceImporter{
			State: resourceServerVSwitchImportState,
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server ID",
			},
			"vswitch_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VSwitch ID",
			},
		},
	}
}

func resourceServerVSwitchImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)
	ID := strconv.Itoa(serverID) + "-" + strconv.Itoa(vSwitchID)
	d.SetId(ID)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceServerVSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)

	err := c.addServerToVSwitch(serverID, vSwitchID)
	if err != nil {
		return diag.FromErr(err)
	}

	ID := strconv.Itoa(serverID) + "-" + strconv.Itoa(vSwitchID)
	d.SetId(ID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerVSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)

	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for lServerID := range vSwitch.server {
		if lServerID == serverID {
			found = true
		}
	}

	if !found {
		d.SetId("")
		err := fmt.Errorf("server %d not associated with vswitch %d", serverID, vSwitchID)
		return diag.FromErr(err)
	}

	ID := strconv.Itoa(serverID) + "-" + strconv.Itoa(vSwitchID)
	d.SetId(ID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerVSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerVSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)

	err := c.removeServerFromVSwitch(serverID, vSwitchID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
