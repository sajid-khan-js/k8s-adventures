# k8s-adventures

## Local tools

e.g. `kubens`, `kubectx`, `k9s`, `lens`

<https://www.youtube.com/watch?v=CB79eTFbR0w> and
<https://martinheinz.dev/blog/75>

## Ingress

Nginx ingress controller:
<https://github.com/kubernetes/ingress-nginx/blob/main/docs/deploy/index.md#quick-start>

## Crossplane

### Compositions

Manifest are [here](./crossplane-compositions/)

An entire custom monitoring stack (prometheus, grafana, loki) abstracted to a
bit of YAML:

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

`helm upgrade --namespace monitoring --install kube-stack-prometheus
prometheus-community/kube-prometheus-stack --set
prometheus-node-exporter.hostRootFsMount.enabled=false`

Relevant commands for using Nginx ingress controller:

```sh
kubectl create ingress prom --class=nginx \
  --rule="prom.localdev.me/*=kube-stack-prometheus-kube-prometheus:80"

kubectl create ingress graf --class=nginx \
  --rule="graf.localdev.me/*=kube-stack-prometheus-grafana:80"

kubectl create ingress alert --class=nginx \
  --rule="alert.localdev.me/*=kube-stack-prometheus-kube-alertmanager:9093"
```

## ArgoCD

- <https://www.youtube.com/watch?v=MeU5_k9ssrs> and
  <https://gitlab.com/nanuchi/argocd-app-config>
- <https://www.youtube.com/watch?v=vpWQeoaiRM4>
- Official docs:
  - <https://argo-cd.readthedocs.io/en/stable/>
  - <https://argo-cd.readthedocs.io/en/stable/core_concepts/>
  - <https://argo-cd.readthedocs.io/en/stable/operator-manual/architecture/>
  
### Using ArgoCD in a PaaS

- https://github.com/argoproj/argo-cd/discussions/5667

#### Structure

##### Shared resources

- Consumable pipeline libraries/GitHub actions
  - use these to setup your golden pipeline (e.g. auto pushes to team's
    container registry in the platform after build, deploy step that updates config repo image tag), but let people pick and mix
    and create their own pipelines
- Helm chart/base kustomize to deploy an app. Used by teams so they don't have to
  define basic k8s building blocks, and also lets you bake in things like ingress
  controllers annotations

##### User repo

- CI pipeline e.g. how to build, test the app. Leave this flexible to devs but
  GitHub actions is a good choice
- Helm chart/kustomization overlays to set k8s deployment options e.g. replicas, ports, cpu limits, env vars, healthz probes
  - Per env overrides required for this config too

##### Config repo

Basically where the team defines their "product" on the platform

- Hydrated deployment.yamls (i.e. where image tag gets updated), values taken from app repo
- Onboard metadata: e.g. team email/github group, cost centre, app git repos
  - allow teams to have multiple apps under one "product"
- CD pipeline switches e.g. promotion between envs, tests before deploy, canary deployment, manual approval required for prod
- Backing infra e.g. Crossplane resources
- Notifications of build/deploy events e.g. to slack
- Who can access my namespace/logs/secrets i.e. anything todo with my product on the platform

:memo: It makes sense to have the config repo split in two. A user facing facade repo which allows platform users to define simple YAML. Then a repo which is actually watched by ArgoCD/Flux which takes the simple YAML and converts it to real K8s objects e.g. a Crossplane resource, an ArgoCD application, FluxCD helm release etc.
