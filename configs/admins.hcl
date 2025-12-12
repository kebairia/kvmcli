// ==================================================
// Network Definition
// ==================================================

# network "homelab2" {
#   namespace  = "infra"
#   netaddress = "10.40.40.0"
#   netmask    = "255.255.255.0"
#   bridge     = "virbr3"
#   mode       = "nat"
#   autostart  = true
#
#   dhcp {
#     # range = "10.10.10.2-10.10.10.254"
#     start = "10.40.40.2"
#     end   = "10.40.40.254"
#   }
#
#   labels = {
#     role        = "work,study"
#     environment = "homelab"
#   }
# }

data "network" "homelab2" {}
data "store" "homelab" {}
// ==================================================
// DNS Servers (admins namespace)
// ==================================================

vm "zakaria" {
  image     = "rocky-9.5"
  namespace = "infra"
  cpu       = 1
  memory    = 2048
  disk      = "20G"
  store     = data.store.homelab
  network   = data.network.homelab2
  # network   = network.homelab
  mac = "02:A3:10:00:00:04"
  ip  = "10.30.30.4"

  labels = {
    role        = "dns"
    environment = "homelab"
  }
}

vm "test2" {
  image     = "rocky-10.1"
  namespace = "infra"
  cpu       = 1
  memory    = 2048
  disk      = "20G"
  store     = data.store.homelab
  network   = data.network.homelab2
  # net = "homelab"
  mac = "02:A3:10:00:00:05"
  ip  = "10.30.30.5"

  labels = {
    role        = "dns"
    environment = "homelab"
  }
}
