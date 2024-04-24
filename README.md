# jazz records
Is a project name from [a golang tutorial](https://go.dev/doc/tutorial/database-access).

This project is for my own learning purpose of kubernetes and the tools.

# Setup
### Prerequisite
- [Install Docker](https://docs.docker.com/engine/install/)
- [Install & setup kubernetes with tools](https://kubernetes.io/docs/setup/)
- [Install & setup skaffold](https://skaffold.dev/docs/quickstart/)

# Run local app
Start minikube first. (Set profile if necessary.)

### start minikube
```bash
minikube -p jazz-records start
eval $(minikube -p jazz-records docker-env)
```

### create volume and secrets

Existing files are
```bash
$ tree k8s/local
k8s/local
├── secrets.yaml
└── volume.yaml

1 drectory, 2 files
```

create local specific resources
```bash
kubectl apply -f k8s/local
```

Finally, run app by skaffold command
```bash
skaffold dev -p local
```


Then, in another terminal pane / tab / window,

```bash
minikube -p jazz-records service jazz-records --url
```
