package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lenstra/hetzner"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallUpdate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
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

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hetzner.Client)

	serverIP := d.Id()

	firewall, err := c.Firewall().Info(serverIP)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", firewall.Status == "active"); err != nil {
		return diag.Errorf("failed to set 'active': %s", err)
	}

	if err := d.Set("whitelist_hos", firewall.WhitelistHOS); err != nil {
		return diag.Errorf("failed to set 'whitelist_hos': %s", err)
	}

	rules := []interface{}{}
	for _, r := range firewall.Rules.Input {
		rules = append(rules, map[string]interface{}{
			"name":     r.Name,
			"dst_ip":   r.DstIP,
			"dst_port": r.DstPort,
			"src_ip":   r.SrcIP,
			"src_port": r.SrcPort,
			"protocol": r.Protocol,
			"action":   r.Action,
		})
	}
	if err := d.Set("rule", rules); err != nil {
		return diag.Errorf("failed to set 'rules': %s", err)
	}

	return nil
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hetzner.Client)

	serverIP := d.Get("server_ip").(string)

	status := "disabled"
	if d.Get("active").(bool) {
		status = "active"
	}

	rules := []*hetzner.FirewallRule{}
	for _, ruleMap := range d.Get("rule").([]interface{}) {
		ruleProperties := ruleMap.(map[string]interface{})
		rules = append(rules, &hetzner.FirewallRule{
			Name:     ruleProperties["name"].(string),
			SrcIP:    hetzner.String(ruleProperties["src_ip"].(string)),
			SrcPort:  hetzner.String(ruleProperties["src_port"].(string)),
			DstIP:    hetzner.String(ruleProperties["dst_ip"].(string)),
			DstPort:  hetzner.String(ruleProperties["dst_port"].(string)),
			Protocol: hetzner.String(ruleProperties["protocol"].(string)),
			Action:   ruleProperties["action"].(string),
		})
	}

	_, err := c.Firewall().Update(&hetzner.FirewallRequest{
		ServerIP:     serverIP,
		WhitelistHOS: hetzner.Bool(d.Get("whitelist_hos").(bool)),
		Status:       hetzner.String(status),
		Rules:        hetzner.FirewallRules{Input: rules},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverIP)

	return resourceFirewallRead(ctx, d, m)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Since all Hetzner servers have an associated firewall there is no real
	// way to delete this resource but removing all rules help by making sure we
	// don't leave a server unknowingly exposed to Internet.

	c := m.(*hetzner.Client)
	serverIP := d.Get("server_ip").(string)
	status := "disabled"
	if d.Get("active").(bool) {
		status = "active"
	}

	_, err := c.Firewall().Update(&hetzner.FirewallRequest{
		ServerIP:     serverIP,
		Status:       hetzner.String(status),
		WhitelistHOS: hetzner.Bool(false),
	})

	return diag.FromErr(err)
}
