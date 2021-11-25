package hetznerrobot

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataServerVSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerVSwitchRead,
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

func dataSourceServerVSwitchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)

	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return err
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
		return err
	}

	ID := strconv.Itoa(serverID) + "-" + strconv.Itoa(vSwitchID)
	d.SetId(ID)

	return nil
}
