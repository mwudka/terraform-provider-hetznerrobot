package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"whitelist_hos": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dst_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(HetznerRobotClient)

	serverIP := d.Get("server_ip").(string)

	status := "disabled"
	if d.Get("active").(bool) {
		status = "active"
	}

	rules := make([]HetznerRobotFirewallRule, 0)
	for _, ruleMap := range d.Get("rule").([]interface{}) {
		ruleProperties := ruleMap.(map[string]interface{})
		rules = append(rules, HetznerRobotFirewallRule{
			Name:     ruleProperties["name"].(string),
			SrcIP:    ruleProperties["src_ip"].(string),
			SrcPort:  ruleProperties["src_port"].(string),
			DstIP:    ruleProperties["dst_ip"].(string),
			DstPort:  ruleProperties["dst_port"].(string),
			Protocol: ruleProperties["protocol"].(string),
			Action:   ruleProperties["action"].(string),
		})
	}

	if err := c.setFirewall(HetznerRobotFirewall{
		IP:                       serverIP,
		WhitelistHetznerServices: d.Get("whitelist_hos").(bool),
		Status:                   status,
		Rules:                    HetznerRobotFirewallRules{Input: rules},
	}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverIP)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(HetznerRobotClient)

	serverIP := d.Id()

	_, err := c.getFirewall(serverIP)
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFirewallRead(ctx, d, m)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
