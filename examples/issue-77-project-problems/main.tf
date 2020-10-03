terraform {
  required_providers {
    #vault = {
    #  source  = "hashicorp/vault"
    #  version = "2.11.0"
    #}

    harbor = {
      source  = "nolte/harbor"
      version = "0.1.9"
    }
  }
  required_version = ">= 0.13"
}

resource "harbor_registry" "dockerhub" {
  name        = "dockerhub-acc-classic"
  url         = "https://hub.docker.com"
  type        = "docker-hub"
  description = "Docker Hub Registry"
  insecure    = false
}

resource "harbor_registry" "gcr_io" {
  name        = "k8s-gcr-io-classic"
  url         = "https://k8s.gcr.io"
  type        = "docker-hub"
  description = "Google K8S Container Registry"
  insecure    = false
}

#provider "harbor" {
#  host     = data.vault_generic_secret.harbor_credentials.data["host"]
#  username = data.vault_generic_secret.harbor_credentials.data["user"]
#  password = data.vault_generic_secret.harbor_credentials.data["password"]
#  schema   = "https"
#  insecure = false
#  basepath = "/api/v2.0"
#}
provider "harbor" {

}

resource "harbor_project" "rancher" {
  name                   = "rancher"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}

resource "harbor_project" "metallb" {
  name                   = "metallb"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}
resource "harbor_project" "ingress_nginx" {
  name                   = "ingress-nginx"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}

resource "harbor_replication" "pull_nginx" {
  name                        = "ingress-nginx-controller"
  description                 = "ingress-nginx-controller"
  source_registry_id          = harbor_registry.gcr_io.id
  source_registry_filter_name = "ingress-nginx/controller"
  source_registry_filter_tag  = "v0.35.0"
  destination_namespace       = harbor_project.ingress_nginx.name
}

resource "harbor_replication" "rancher_rancher" {
  name                        = "rancher-rancher"
  description                 = "rancher-rancher"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/rancher"
  source_registry_filter_tag  = "v2.4.8"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_rancher_agent" {
  name                        = "rancher-rancher-agent"
  description                 = "rancher-rancher-agent"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/rancher-agent"
  source_registry_filter_tag  = "v2.4.8"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_metrics_server" {
  name                        = "rancher-metrics-server"
  description                 = "rancher-metrics-server"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/metrics-server"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_hyperkube" {
  name                        = "rancher-hyperkube"
  description                 = "rancher-hyperkube"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/hyperkube"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_coreos_flannel" {
  name                        = "rancher-coreos-flannel"
  description                 = "rancher-coreos-flannel"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/coreos-flannel"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_coredns_coredns" {
  name                        = "rancher-coredns-coredns"
  description                 = "rancher-coredns-coredns"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/coredns-coredns"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_cluster_proportional_autoscaler" {
  name                        = "rancher-cluster-proportional-autoscaler"
  description                 = "rancher-cluster-proportional-autoscaler"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/cluster-proportional-autoscaler"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "rancher_calico_pod2daemon_flexvol" {
  name                        = "rancher-calico-pod2daemon-flexvol"
  description                 = "rancher-calico-pod2daemon-flexvol"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/calico-pod2daemon-flexvol"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}

resource "harbor_replication" "metallb_speaker" {
  name                        = "metallb-speaker"
  description                 = "metallb-speaker"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "metallb/speaker"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.metallb.name
}

resource "harbor_replication" "metallb_controller" {
  name                        = "metallb-controller"
  description                 = "metallb-controller"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "metallb/controller"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.metallb.name
}

resource "harbor_replication" "rancher_calico_node" {
  name                        = "rancher-calico-node"
  description                 = "rancher-calico-node"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/calico-node"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}
resource "harbor_replication" "rancher_calico_cni" {
  name                        = "rancher-calico-cni"
  description                 = "rancher-calico-cni"
  source_registry_id          = harbor_registry.dockerhub.id
  source_registry_filter_name = "rancher/calico-cni"
  source_registry_filter_tag  = "*"
  destination_namespace       = harbor_project.rancher.name
}
