## UAPI-UI Operator deploys API, UI and MongoDB service. 

### Build & Run Operator locally 
```bash
# CRD
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_crd.yaml
# Run operator locally 
operator-sdk up local --namespace=uapi
# CR
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_cr.yaml
# Cleanup CR
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_cr.yaml
# Cleanup CRD
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_crd.yaml
```

### Debug Operator locally 
```bash
# With your favorite IDE export env vars   
export KUBERNETES_CONFIG=/Users/dima/.kube/config
export WATCH_NAMESPACE=uapi
export OPERATOR_NAME=UAPI-OPERATOR
# Run cmd/manager/main.go
cmd/manager/main.go
```

### Build UAPI Go operator
```bash
operator-sdk build docker.io/dimssss/uapi-operator:TAG
docker push docker.io/dimssss/uapi-operator:TAG
```

### Deploy CRD & CR 
```bash
# CRD
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_crd.yaml
# CR
oc create -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_cr.yaml
```

### Cleanup
```bash
# CR
oc delete -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_cr.yaml
# CRD
oc delete -f https://raw.githubusercontent.com/Dimss/uapi-go-operator/master/deploy/crds/uiapi_v1alpha1_uapi_crd.yaml
```

Customize CR
```bash
apiVersion: uiapi.com/v1alpha1
kind: Uapi
metadata:
  name: uapi
spec:
  namespace: "uapi"
  ui:
    size: 1
    name: "ui"
    serviceNodePort: 30080
    apiUrl: "http://127.0.0.1:30081"
    image: "docker-registry.default.svc:5000/uapi/ui:latest"
  api:
    size: 1
    name: "api"
    serviceNodePort: 30081
    confSecretName: "uapisecret"
    image: "docker-registry.default.svc:5000/uapi/uapi:latest"
  db:
    image: "registry.redhat.io/rhscl/mongodb-36-rhel7:latest"
    port: 27017
    host: "mongo"
    name: "uapi"
```