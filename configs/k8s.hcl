// ==================================================
// REFERENCE EXISTING RESOURCES
// ==================================================

// Reference an existing network (already created)
data "network" "homelab" {
  name = "homelab"
}
data "store" "homelab-store" {}

// Or create a new one
network "production" {
  cidr = "10.20.0.0/24"
}

// ==================================================
// VM DEFINITIONS WITH REFERENCES
// ==================================================

vm "k8s-master-01" {
  image  = "rocky-10"
  cpu    = 2
  memory = "2G"
  disk   = "20G"
  net    = data.network.homelab // Reference existing network
  mac    = "02:A3:10:00:01:01"
}

vm "k8s-master-02" {
  image  = "rocky-10"
  cpu    = 2
  memory = "2G"
  disk   = "20G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:01:02"
}

vm "k8s-worker-01" {
  image  = "rocky-9.5"
  cpu    = 4
  memory = "4G"
  disk   = "40G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:02:01"
  ip     = "10.10.0.201"
}

vm "k8s-worker-02" {
  image  = "rocky-9.5"
  cpu    = 4
  memory = "4G"
  disk   = "40G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:02:02"
  ip     = "10.10.0.202"
}

vm "k8s-worker-03" {
  image  = "rocky-9.5"
  cpu    = 2
  memory = "2G"
  disk   = "20G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:02:03"
}

vm "k8s-lb-01" {
  image  = "nginx-alpine"
  cpu    = 1
  memory = "1G"
  disk   = "10G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:05:01"
}

// ==================================================
// LOGICAL GROUPING: Cluster as VM collection
// ==================================================

// Cluster is just a logical grouping of existing VMs
cluster "k8s-prod" {
  masters = [
    vm.k8s-master-01,
    vm.k8s-master-02,
  ]

  workers = [
    vm.k8s-worker-01,
    vm.k8s-worker-02,
    vm.k8s-worker-03,
  ]

  lb = [
    vm.k8s-lb-01,
  ]

  // Optional: cluster-level metadata
  labels = {
    environment = "production"
    app         = "kubernetes"
  }

  // Optional: cluster-level operations
  lifecycle {
    start_order = ["masters", "workers", "lb"]
    stop_order  = ["lb", "workers", "masters"]
  }
}

// ==================================================
// ALTERNATIVE 1: Role-based grouping
// ==================================================

/*
cluster "k8s-prod" {
  masters = [
    vm.k8s-master-01,
    vm.k8s-master-02,
  ]
  
  workers = [
    vm.k8s-worker-01,
    vm.k8s-worker-02,
    vm.k8s-worker-03,
  ]
  
  loadbalancers = [
    vm.k8s-lb-01,
  ]
}
*/

// ==================================================
// ALTERNATIVE 2: Tag-based auto-grouping
// ==================================================

/*
vm "k8s-master-01" {
  image = "rocky-10"
  cpu   = 2
  memory = "2G"
  disk  = "20G"
  net   = data.network.homelab
  
  tags = ["k8s-prod", "master"]  // Automatically grouped
}

vm "k8s-worker-01" {
  image = "rocky-9.5"
  cpu   = 4
  memory = "4G"
  disk  = "40G"
  net   = data.network.homelab
  
  tags = ["k8s-prod", "worker"]
}

// Cluster auto-discovers VMs with matching tag
cluster "k8s-prod" {
  discover = "tag:k8s-prod"
  
  roles {
    master = "tag:master"
    worker = "tag:worker"
  }
}
*/

// ==================================================
// ALTERNATIVE 3: Query-based grouping
// ==================================================

/*
cluster "k8s-prod" {
  query = {
    labels = {
      app = "kubernetes"
      environment = "production"
    }
  }
}

// Or even simpler selector syntax
cluster "k8s-prod" {
  select = "app=kubernetes,environment=production"
}
*/

// ==================================================
// ALTERNATIVE 4: Inline minimal syntax
// ==================================================

/*
cluster "k8s-prod" [
  vm.k8s-master-01,
  vm.k8s-master-02,
  vm.k8s-worker-01,
  vm.k8s-worker-02,
  vm.k8s-worker-03,
  vm.k8s-lb-01
]
*/

// ==================================================
// MIXED APPROACH: Best of all worlds
// ==================================================

/*
// Reference existing network
data "network" "homelab" {}

// VMs with simple tagging
vm "k8s-master-01" {
  image = "rocky-10"
  cpu   = 2
  memory = "2G"
  disk  = "20G"
  net   = data.network.homelab
  
  cluster = "k8s-prod"
  role    = "master"
}

vm "k8s-worker-01" {
  image = "rocky-9.5"
  cpu   = 4
  memory = "4G"
  disk  = "40G"
  net   = data.network.homelab
  
  cluster = "k8s-prod"
  role    = "worker"
}

// Cluster definition (optional, for metadata/operations)
cluster "k8s-prod" {
  // VMs auto-discovered via cluster = "k8s-prod"
  
  operations {
    start_order = ["role=master", "role=worker", "role=lb"]
  }
}
*/

// ==================================================
// PRACTICAL EXAMPLE: DNS Cluster
// ==================================================

data "network" "homelab" {}

vm "dns-primary" {
  image  = "rocky-9.5"
  cpu    = 1
  memory = "2G"
  disk   = "20G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:03:01"
}

vm "dns-secondary" {
  image  = "rocky-9.5"
  cpu    = 1
  memory = "2G"
  disk   = "20G"
  net    = data.network.homelab
  mac    = "02:A3:10:00:03:02"
}

cluster "dns-ha" {
  vms = [
    vm.dns-primary,
    vm.dns-secondary,
  ]

  labels = {
    service = "dns"
    ha      = "true"
  }
}

