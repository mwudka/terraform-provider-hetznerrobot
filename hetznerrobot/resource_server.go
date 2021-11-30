package hetznerrobot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceServerImportState,
		},

		Schema: map[string]*schema.Schema{
			"is_cancelled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Status of server cancellation",
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server name",
			},

			// read-only / computed
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data center",
			},
			"paid_until": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Paid until date",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server product name",
			},
			"server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Server ID",
			},
			"server_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of assigned single IP addresses",
			},
			"server_ip_v4_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IP address",
			},
			"server_ip_v6_net": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IPv6 net address",
			},
			"server_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of assigned subnets",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server status (\"ready\" or \"in process\")",
			},
			"traffic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Free traffic quota, 'unlimited' in case of unlimited traffic",
			},

			/*
			   reset (Boolean)	Flag of reset system availability
			   rescue (Boolean)	Flag of Rescue System availability
			   vnc (Boolean)	Flag of VNC installation availability
			   windows (Boolean)	Flag of Windows installation availability
			   plesk (Boolean)	Flag of Plesk installation availability
			   cpanel (Boolean)	Flag of cPanel installation availability
			   wol (Boolean)	Flag of Wake On Lan availability
			   hot_swap (Boolean)	Flag of Hot Swap availability

			   linked_storagebox (Integer)	Linked Storage Box ID
			*/
		},
	}
}

func resourceServerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	serverID, _ := strconv.Atoi(d.Id())
	server, err := c.getServer(serverID)
	if err != nil {
		return nil, fmt.Errorf("Unable to find Server with ID %d:\n\t %q", serverID, err)
	}

	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	d.Set("server_id", serverID)
	d.Set("server_ip_addresses", server.IPList)
	d.Set("server_ip_v4_addr", server.ServerIPv4)
	d.Set("server_ip_v6_net", server.ServerIPv6)
	d.Set("server_name", server.ServerName)
	d.Set("server_subnets", server.SubnetList)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/* Server resource can't be created via API */

	c := meta.(HetznerRobotClient)

	var serverID = d.Get("'server_id").(int)
	var serverName = d.Get("server_name").(string)

	server, err := c.setServerName(serverID, serverName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	d.Set("server_ip_addresses", server.IPList)
	d.Set("server_ip_v4_addr", server.ServerIPv4)
	d.Set("server_ip_v6_net", server.ServerIPv6)
	d.Set("server_subnets", server.SubnetList)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)
	d.SetId(strconv.Itoa(serverID))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID, _ := strconv.Atoi(d.Id())
	server, err := c.getServer(serverID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find Server with ID %d:\n\t %q", serverID, err))
	}

	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	d.Set("server_id", serverID)
	d.Set("server_ip_addresses", server.IPList)
	d.Set("server_ip_v4_addr", server.ServerIPv4)
	d.Set("server_ip_v6_net", server.ServerIPv6)
	d.Set("server_name", server.ServerName)
	d.Set("server_subnets", server.SubnetList)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	var serverID = d.Get("'server_id").(int)
	var serverName = d.Get("server_name").(string)
	server, err := c.setServerName(serverID, serverName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	d.Set("server_ip_addresses", server.IPList)
	d.Set("server_ip_v4_addr", server.ServerIPv4)
	d.Set("server_ip_v6_net", server.ServerIPv6)
	d.Set("server_subnets", server.SubnetList)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)
	d.SetId(strconv.Itoa(serverID))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
