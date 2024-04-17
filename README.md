# jazz records
Is a project name from [a golang tutorial](https://go.dev/doc/tutorial/database-access).

This project is for my own learning purpose of kubernetes and the tools.

# Setup
### Prerequisite
- [Install Docker](https://docs.docker.com/engine/install/)
- [Install & setup kubernetes with tools](https://kubernetes.io/docs/setup/)
- [Install & setup skaffold](https://skaffold.dev/docs/quickstart/)

### Secrets
All the secrets are defined in [secrets.yaml.example](k8s/secrets.yaml.example).

Copy it as `k8s/secrets.yaml`

```bash
cp k8s/secrets.yaml.example k8s/secrets.yaml
```

# Run
Start minikube first. (Set profile if necessary.)

```bash
minikube -p jazz-records start
eval $(minikube -p jazz-records docker-env)
```

```bash
skaffold dev
```

Then, in another terminal pane / tab / window,

```bash
minikube -p jazz-records service jazz-records --url
```
