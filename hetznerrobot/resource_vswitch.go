package hetznerrobot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceVSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVSwitchCreate,
		ReadContext:   resourceVSwitchRead,
		UpdateContext: resourceVSwitchUpdate,
		DeleteContext: resourceVSwitchDelete,

		Importer: &schema.ResourceImporter{
			State: resourceVSwitchImportState,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VSwitch ID",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN ID",
			},
			// computed / read-only fields
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Cancellation status",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached server list",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached subnet list",
			},
			"cloud_networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached cloud network list",
			},
		},
	}
}

func resourceVSwitchImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := strconv.Atoi(d.Id())
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return nil, fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.name)
	d.Set("vlan", vSwitch.vlan)
	d.Set("is_cancelled", vSwitch.isCancelled)
	d.Set("servers", vSwitch.server)
	d.Set("subnets", vSwitch.subnet)
	d.Set("cloud_networks", vSwitch.cloudNetwork)
	d.SetId(strconv.Itoa(vSwitchID))

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	vSwitch, err := c.createVSwitch(name, vlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to create VSwitch :\n\t %q", err))
	}

	d.Set("is_cancelled", vSwitch.isCancelled)
	d.Set("servers", vSwitch.server)
	d.Set("subnets", vSwitch.subnet)
	d.Set("cloud_networks", vSwitch.cloudNetwork)
	d.SetId(strconv.Itoa(vSwitch.id))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := strconv.Atoi(d.Id())
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err))
	}

	d.Set("name", vSwitch.name)
	d.Set("vlan", vSwitch.vlan)
	d.Set("cancelled", vSwitch.isCancelled)
	d.Set("servers", vSwitch.server)
	d.Set("subnets", vSwitch.subnet)
	d.Set("cloud_networks", vSwitch.cloudNetwork)
	d.SetId(strconv.Itoa(vSwitchID))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Get("id").(int)
	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	vSwitch, err := c.updateVSwitch(vSwitchID, name, vlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to update VSwitch :\n\t %q", err))
	}

	d.Set("is_cancelled", vSwitch.isCancelled)
	d.Set("servers", vSwitch.server)
	d.Set("subnets", vSwitch.subnet)
	d.Set("cloud_networks", vSwitch.cloudNetwork)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := strconv.Atoi(d.Id())
	err := c.deleteVSwitch(vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
