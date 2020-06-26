# gitlab-service
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/gitlab-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/gitlab-service)](https://goreportcard.com/report/github.com/keptn-sandbox/gitlab-service)

This implements a gitlab-service for Keptn. If you want to learn more about Keptn visit us on [keptn.sh](https://keptn.sh)

## Compatibility Matrix

| Keptn Version    | [gitlab-service Docker Image](https://hub.docker.com/r/checkelmann/gitlab-service/tags) |
|:----------------:|:----------------------------------------:|
|       0.6.2      | checkelmann/gitlab-service:develop|

## Installation

The *gitlab-service* can be installed as a part of [Keptn's uniform](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the *gitlab-service* in your Keptn Kubernetes cluster, adjust and apply the [`deploy/service.yaml`](deploy/service.yaml) file:

_You need to set your GitLab Instance URL and API Token in the service.yaml - This will be moved in a K8s-Secret in the official release!_

```yaml
- name: GITLAB_URL
  value: gitlab.com
- name: GITLAB_TOKEN
  value: YOURGITLABAPITOKEN
```

```console
kubectl apply -f deploy/service.yaml
```

This should install the `gitlab-service` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment gitlab-service -o wide
kubectl -n keptn get pods -l run=gitlab-service
```

### Configuration

You need to configure your gitlab.conf.yaml and add it as project resource

```yaml
---
spec_version: '0.1.0'
gitlabinstance:
  - name: gitlab.com
    url: $ENV.GITLAB_URL
    token: $ENV.GITLAB_TOKEN
pipelines:
  - name: MyGitLabPipeline
    instance: gitlab.com
    projectid: 19092773
    ref: master
    variables:
      Message: My message from Keptn Project $PROJECT-$STAGE-$SERVICE
      WaitTime: 2
      Result: SUCCESS
  - name: MyGitLabPipelineWithFailure
    instance: gitlab.com
    projectid: 19092773
    ref: error
    variables:
      Message: My message from Keptn Project $PROJECT-$STAGE-$SERVICE
      WaitTime: 2
      Result: FAILURE
events:
  - event: configuration.change
    action: MyGitLabPipeline 
    timeout: 60
    onsuccess:
      event: deployment.finished
      deploymentURIPublic: http://yourdeployedapp.yourdomain
      result: pass
    onfailure:
      event: deployment.finished
      result: failed
```

Upload the configuration file as keptn resource:

```console
keptn add-resource --project=sockshop --stage=hardening --service=carts --resource=gitlab/gitlab.conf.yaml --resourceUri=gitlab/gitlab.conf.yaml
```

### Up- or Downgrading

Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):

```console
kubectl -n keptn set image deployment/gitlab-service gitlab-service=checkelmann/gitlab-service:$VERSION --record
```

### Uninstall

To delete a deployed *gitlab-service*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## License

Please find more information in the [LICENSE](LICENSE) file.
