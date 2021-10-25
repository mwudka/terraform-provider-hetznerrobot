package hetzner

type ServerSummary struct {
	ServerIP         string   `json:"server_ip"`
	ServerNumber     int      `json:"server_number"`
	ServerName       string   `json:"server_name"`
	Product          string   `json:"product"`
	DC               string   `json:"dc"`
	Traffic          string   `json:"traffic"`
	Status           string   `json:"status"`
	Cancelled        bool     `json:"cancelled"`
	PaidUntil        string   `json:"paid_until"`
	IP               []string `json:"ip"`
	Subnet           []Subnet `json:"subnet"`
	LinkedStoragebox *int     `json:"linked_storagebox"`
}

type Server struct {
	ServerIP         string   `json:"server_ip"`
	ServerNumber     int      `json:"server_number"`
	ServerName       string   `json:"server_name"`
	Product          string   `json:"product"`
	DC               string   `json:"dc"`
	Traffic          string   `json:"traffic"`
	Status           string   `json:"status"`
	Cancelled        bool     `json:"cancelled"`
	PaidUntil        string   `json:"paid_until"`
	IP               []string `json:"ip"`
	Subnet           []Subnet `json:"subnet"`
	Reset            bool     `json:"reset"`
	Rescue           bool     `json:"rescue"`
	Vnc              bool     `json:"vnc"`
	Windows          bool     `json:"windows"`
	Plesk            bool     `json:"plesk"`
	CPanel           bool     `json:"cpanel"`
	WOL              bool     `json:"wol"`
	HotSwap          bool     `json:"hot_swap"`
	LinkedStoragebox *int     `json:"linked_storagebox"`
}

type Subnet struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type ServerRequest struct {
	ServerIP   string `form:"-"`
	ServerName string `form:"server_name"`
}

type Rescue struct {
	ServerIP       string   `json:"server_ip"`
	ServerNumber   int      `json:"server_number"`
	OS             string   `json:"os"`
	Arch           int      `json:"arch"`
	Active         bool     `json:"active"`
	Password       string   `json:"password"`
	AuthorizedKeys []string `json:"authorized_key"`
	HostKeys       []string `json:"host_key"`
}

type RescueRequest struct {
	ServerIP       string   `form:"-"`
	OS             string   `form:"os"`
	Arch           int      `form:"arch"`
	AuthorizedKeys []string `form:"authorized_key"`
}

type Reset struct {
	ServerIP string `json:"server_ip"`
	Type     string `json:"type"`
}

type ResetRequest struct {
	ServerIP string `form:"-"`
	Type     string `form:"type"`
}

type StorageBoxSummary struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	Product      string `json:"product"`
	Cancelled    bool   `json:"cancelled"`
	Locked       bool   `json:"locked"`
	Location     string `json:"location"`
	LinkedServer int    `json:"linked_server"`
	PaidUntil    string `json:"paid_until"`
}

type StorageBox struct {
	ID                   int    `json:"id"`
	Login                string `json:"login"`
	Name                 string `json:"name"`
	Product              string `json:"product"`
	Cancelled            bool   `json:"cancelled"`
	Locked               bool   `json:"locked"`
	Location             string `json:"location"`
	LinkedServer         int    `json:"linked_server"`
	PaidUntil            string `json:"paid_until"`
	DiskQuota            int    `json:"disk_quota"`
	DiskUsage            int    `json:"disk_usage"`
	DiskUsageData        int    `json:"disk_usage_data"`
	DiskUsageSnapshots   int    `json:"disk_usage_snapshots"`
	Webdav               bool   `json:"webdav"`
	Samba                bool   `json:"samba"`
	SSH                  bool   `json:"ssh"`
	ExternalReachability bool   `json:"external_reachability"`
	ZFS                  bool   `json:"zfs"`
	Server               string `json:"server"`
	HostSystem           string `json:"host_system"`
}

type StorageBoxRequest struct {
	ID                   int     `form:"-"`
	StorageBoxName       *string `form:"storagebox_name,omitempty"`
	Samba                *bool   `form:"samba,omitempty"`
	Webdav               *bool   `form:"webdav,omitempty"`
	SSH                  *bool   `form:"ssh,omitempty"`
	ExternalReachability *bool   `form:"external_reachability,omitempty"`
	ZFS                  *bool   `form:"zfs,omitempty"`
}

type Firewall struct {
	ServerIP     string        `json:"server_ip"`
	ServerNumber int           `json:"server_number"`
	Status       string        `json:"status"`
	WhitelistHOS bool          `json:"whitelist_hos"`
	Port         string        `json:"port"`
	Rules        FirewallRules `json:"rules"`
}

type FirewallRules struct {
	Input []*FirewallRule `json:"input" form:"input"`
}

type FirewallRule struct {
	IPVersion string  `json:"ip_version" form:"ip_version"`
	Name      string  `json:"name" form:"name"`
	DstIP     *string `json:"dst_ip" form:"dst_ip"`
	SrcIP     *string `json:"src_ip" form:"src_ip"`
	DstPort   *string `json:"dst_port" form:"dst_port"`
	SrcPort   *string `json:"src_port" form:"src_port"`
	Protocol  *string `json:"protocol" form:"protocol"`
	TcpFlags  *string `json:"tcp_flags" form:"tcp_flags"`
	Action    string  `json:"action" form:"action"`
}

type FirewallRequest struct {
	ServerIP     string        `form:"-"`
	Status       *string       `form:"status"`
	WhitelistHOS *bool         `form:"whitelist_hos"`
	TemplateID   *string       `form:"template_id"`
	Rules        FirewallRules `form:"rules"`
}
