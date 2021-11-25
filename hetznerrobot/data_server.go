package hetznerrobot

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server ID",
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
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of server cancellation",
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

func dataSourceServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)

	server, err := c.getServer(serverID)
	if err != nil {
		return fmt.Errorf("Unable to find Server with ID %d:\n\t %q", serverID, err)
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
	d.SetId(strconv.Itoa(serverID))

	return nil
}
