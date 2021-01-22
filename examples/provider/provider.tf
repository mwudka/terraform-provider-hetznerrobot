terraform {
  required_providers {
    hetznerrobot = {
      version = "0.0.1"
      source  = "github.com/mwudka/hetznerrobot"
    }
  }

}

resource "hetznerrobot_firewall" "firewall" {
  server_ip     = "95.216.6.55"
  active        = true
  whitelist_hos = true

  rule {
    name     = "Allow ssh"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "22"
    protocol = "tcp"
    action   = "accept"
  }

  rule {
    name     = "Allow inbound"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "32768-65535"
    action   = "accept"
  }

  rule {
    name     = "Allow ICMP"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "0-65535"
    protocol = "icmp"
    action   = "accept"
  }


  rule {
    name     = "Deny others"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "0-65535"
    action   = "discard"
  }
}
