package hetznerrobot

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataVSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSwitchRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VSwitch ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VLAN ID",
			},
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

func dataSourceVSwitchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Get("id").(int)
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.name)
	d.Set("vlan", vSwitch.vlan)
	d.Set("is_cancelled", vSwitch.isCancelled)
	d.Set("servers", vSwitch.server)
	d.Set("subnets", vSwitch.subnet)
	d.Set("cloud_networks", vSwitch.cloudNetwork)
	d.SetId(strconv.Itoa(vSwitchID))

	return nil
}
