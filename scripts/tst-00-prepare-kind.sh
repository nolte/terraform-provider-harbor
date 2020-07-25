#!/usr/bin/env bash
set -o errexit
set -o pipefail
set -o nounset

INGRESS_NGINX_CHART_VERSION=${INGRESS_NGINX_CHART_VERSION:-2.1.0}

echo "==> Start Local Kind k8s Cluster for testing"

cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

echo "==> Install Ingress for Harbor ApiAccess"
set +x
kubectl apply -f "https://raw.githubusercontent.com/kubernetes/ingress-nginx/ingress-nginx-${INGRESS_NGINX_CHART_VERSION}/deploy/static/provider/kind/deploy.yaml" --wait=true
set -x
sleep 30

echo "==> Waiting Ingress is successfull started"

kubectl wait --namespace ingress-nginx \
    --for=condition=ready pod \
    --selector=app.kubernetes.io/component=controller \
    --timeout=240s

echo "==> Kind Test Env Ready for use"
