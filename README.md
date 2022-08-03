# k8s-adventures

## Local tools

e.g. `kubens`, `kubectx`, `k9s`, `lens`

<https://www.youtube.com/watch?v=CB79eTFbR0w> and <https://martinheinz.dev/blog/75>

## Ingress

Nginx ingress controller: <https://github.com/kubernetes/ingress-nginx/blob/main/docs/deploy/index.md#quick-start>

## Crossplane

### Compositions

Manifest are [here](./crossplane-compositions/)

An entire custom monitoring stack (prometheus, grafana, loki) abstracted to a bit of YAML:

```yaml
apiVersion: devopstoolkitseries.com/v1alpha1
kind: MonitoringClaim
metadata:
  name: monitoring
spec:
  compositionSelector:
    matchLabels:
      monitor: prometheus
      alert: alert-manager
      dashboard: grafana
      log: loki
  parameters:
    monitor:
      host: monitor.127.0.0.1.nip.io
    alert:
      host: alert.127.0.0.1.nip.io
    dashboard:
      host: dashboard.127.0.0.1.nip.io
```

- <https://gist.github.com/vfarcic/e27c3e62438479efc3b676edbe57aacf>
- <https://www.youtube.com/watch?v=yFLV_mOSiYI>
- <https://github.com/crossplane/crossplane/blob/2842bcf3b9bcc7e66bcd8bca1d07a8646e42324d/docs/getting-started/create-configuration.md>
- <https://crossplane.io/docs/v1.9/getting-started/install-configure.html#select-a-getting-started-configuration>

## prometheus-operator

<https://getbetterdevops.io/setup-prometheus-and-grafana-on-kubernetes/>

`helm upgrade --namespace monitoring --install kube-stack-prometheus prometheus-community/kube-prometheus-stack --set prometheus-node-exporter.hostRootFsMount.enabled=false`

Relevant commands for using Nginx ingress controller:

```sh
kubectl create ingress prom --class=nginx \
  --rule="prom.localdev.me/*=kube-stack-prometheus-kube-prometheus:80"

kubectl create ingress graf --class=nginx \
  --rule="graf.localdev.me/*=kube-stack-prometheus-grafana:80"

kubectl create ingress alert --class=nginx \
  --rule="alert.localdev.me/*=kube-stack-prometheus-kube-alertmanager:9093"
```
