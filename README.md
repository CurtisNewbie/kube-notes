# kube-notes

Notes for kubernetes

- Minikube Tutorial: https://minikube.sigs.k8s.io/docs/start/
- Kubernetes Doc: 
- Kubectl Cheat Sheet: https://kubernetes.io/docs/reference/kubectl/cheatsheet/#bash
- Nginx Ingress Controller: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/ 
- Ingress nginx for TCP and UDP services (Minikube): https://minikube.sigs.k8s.io/docs/tutorials/nginx_tcp_udp_ingress/


## Getting Started

### 1. Install Minikube 

- Minikube Getting Started https://minikube.sigs.k8s.io/docs/start/

To install.

```sh
brew install minikube
```

Or download it directly :D

Minikube needs a driver, say, Docker. But, Dockerd is not available on MacOS, it only ships with a cli tool: the docker and docker-compose. We either install Docker Desktop or Colima. 

- Colima Github https://github.com/abiosoft/colima

```sh
# install colima using homebrew 
# if you are using a very old OS version, and there is no bottle for it, it may compile and it's extremely slow (at least for me)
brew install docker-compose 
brew install docker 
brew install colima

# or you may download and install it directly
curl -LO https://github.com/abiosoft/colima/releases/download/v0.5.4/colima-Darwin-x86_64
install colima-Darwin-x86_64 /usr/local/bin/colima

# and then start the colima
colima start
```

Unfortunately, my Mac is stil using Catalina OS, which is pretty old. Colima and its dependencies must be compiled rather than automatically downloaded and installed (for me by Homebrew, e.g., LLVM). I chose to use an old version of Docker Desktop (4.15.0).

Source: [Stackoverflow install-docker-on-macos-catalina](https://stackoverflow.com/questions/68373008/install-docker-on-macos-catalina)

```sh
# download Cask code for Docker Desktop 4.15.0,93002
curl https://raw.githubusercontent.com/Homebrew/homebrew-cask/1a83f3469ab57b01c0312aa70503058f7a27bd1d/Casks/docker.rb -O

# install Docker Desktop from Cask Code
brew install --cask docker.rb
```

### 2. Run Minikube

Start Minikube

```sh
minikube start
```

Start Minikube Dashboard, the dashboard is just what it's named. When it's ready, it opens up a new tab in your browser, and shows the 'Dashboard'.

```sh
minikube dashboard
```

The opened tab's url is like the following:

```
http://127.0.0.1:55901/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/workloads?namespace=default
```

This command also serves as a proxy for us to access, by default the kubernetes network is not accessible externally. We can kill it by ctrl-c.

### 3. The Basic 

In k8s, Pod is a group of one or more containers. K8s deployment monitors and controls the pod, and restart the containers in pod if necessary. The concept is very similar to the docker-compose and the docker world, but with a larger scale. Deployment is also reponsible for scaling pods.

To create a deployment:

```sh
# from the tutorial
kubectl create deployment hello-minikube --image=kicbase/echo-server:1.0
```

To list all deployments:

```sh
kubectl get deployments

# NAME         READY   UP-TO-DATE   AVAILABLE   AGE
# hello-node   0/1     1            0           7s
```

To list the pods:

```sh
kubectl get pods

# NAME                          READY   STATUS             RESTARTS   AGE
# hello-node-7b87cd5f68-76s74   0/1     ImagePullBackOff   0          4m25s
```

To list all pods:

```sh
kubectl get pods -A

# NAMESPACE              NAME                                        READY   STATUS             RESTARTS      AGE
# default                hello-node-779bd496d-shsmz                  0/1     ImagePullBackOff   0             6m49s
# kube-system            coredns-787d4945fb-pkdkt                    1/1     Running            0             38m
# kube-system            etcd-minikube                               1/1     Running            0             38m
# kube-system            kube-apiserver-minikube                     1/1     Running            0             38m
# kube-system            kube-controller-manager-minikube            1/1     Running            0             38m
# kube-system            kube-proxy-sfg57                            1/1     Running            0             38m
# kube-system            kube-scheduler-minikube                     1/1     Running            0             38m
# kube-system            storage-provisioner                         1/1     Running            1 (37m ago)   38m
# kubernetes-dashboard   dashboard-metrics-scraper-5c6664855-fnbxb   1/1     Running            0             37m
# kubernetes-dashboard   kubernetes-dashboard-55c4cbbc7c-hptc9       1/1     Running            0             37m
```

If you have followed the outdated tutorial, the pod may not start. For example, the `https://kubernetes.io/docs/tutorials/hello-minikube/` one. 

But the experience is useful tho, the *ImagePullBackOff* error seemed to be a problem with pulling images from registry.
Another similar error that I found is "ErrImagePull". This can also be found on the Dashboard.

We can check the logs of the pod:

```sh
kubectl logs -p hello-node-7b87cd5f68-76s74

# Error from server (BadRequest): previous terminated container "agnhost" in pod "hello-node-7b87cd5f68-76s74" not found
```

This error msg states that the contain "agnhost" doesn't exist.

We can delete the deployment using the following command

```sh
kubectl delete deployment hello-node

# deployment.apps "hello-node" deleted
```

If everything goes right, then we have:

```sh
kubectl create deployment hello-minikube --image=kicbase/echo-server:1.0
# deployment.apps/hello-minikube created

kubectl get deployments
# NAME             READY   UP-TO-DATE   AVAILABLE   AGE
# hello-minikube   1/1     1            1           81s

kubectl get pods
# NAME                              READY   STATUS    RESTARTS   AGE
# hello-minikube-77b6f68484-dcfkn   1/1     Running   0          2m6s
```

We can expose the deployment by binding it's port to a specified port using `NodePort`. The `NodePort` according to minikube's handboook, is *"NodePort, as the name implies, opens a specific port, and any traffic that is sent to this port is forwarded to the service."* But the service is not externally accessible yet.

```sh
kubectl expose deployment hello-minikube --type=NodePort --port=8080

# service/hello-minikube exposed
```

We check services for this deployment:

```sh
kubectl get services hello-minikube

# NAME             TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
# hello-minikube   NodePort   10.97.53.50   <none>        8080:30856/TCP   117s
```

We then create a port-forward to 8080 for this specific service. The service exposes 8080 port, and we forward any traffic from 7080 to 8080. Then we have access to the service, using `http://localhost:7080`. *"Kubectl port-forward is a method to access, interact and manage internal Kubernetes clusters directly from your local network."*

```sh
kubectl port-forward service/hello-minikube 7080:8080

# Forwarding from 127.0.0.1:7080 -> 8080
# Forwarding from [::1]:7080 -> 8080
# Handling connection for 7080
# Handling connection for 7080
```

To stop the cluster:

```sh
minikube pause
```

To unpause the cluster:

```sh
minikube unpause
```

To halt the cluster:

```sh
minikube stop
```

To change configuration, e.g., default mem limit, require restart:

```sh
minikube config set memory 1234
```

To view configuration:

```sh
minikube config view
```

To list kubernetes addons (services that can be installed and used, such as minikube dashboard):

```sh
minikube addons list
```

Delete all minikube clusters:

```sh
minikube delete --all
```

## More Stuff

To access applications inside Kubernetes, we use services. There are two major cagegories:

- NodePort
- LoadBalancer

### NodePort

NodePort is very straightforward, it opens a specific port, and the traffic sent to this port is forwarded to the service.

In minikube, we have the following shortcut to get the minikube's IP and the service's NodePort, it's not a kubernetes thing.

```sh
minikube service hello-minikube  --url

# http://127.0.0.1:61322
```

With this open, we can access the service using the url returned. Without the `--url` flag, we have:

```sh
minikube service hello-minikube

# |-----------|----------------|-------------|---------------------------|
# | NAMESPACE |      NAME      | TARGET PORT |            URL            |
# |-----------|----------------|-------------|---------------------------|
# | default   | hello-minikube |        8080 | http://192.168.49.2:30856 |
# |-----------|----------------|-------------|---------------------------|
# 🏃  Starting tunnel for service hello-minikube.
# |-----------|----------------|-------------|------------------------|
# | NAMESPACE |      NAME      | TARGET PORT |          URL           |
# |-----------|----------------|-------------|------------------------|
# | default   | hello-minikube |             | http://127.0.0.1:61342 |
# |-----------|----------------|-------------|------------------------|
```

These two commands open the tunnel for the service (it seems like it's a Drawin/Windows/WSL thing, and it's not needed on Linux. On Linux, no tunnel is created).

So **NodePort** is a type of service, we export the deployment using the NodePort service type, and then we use the port to access the service inside kubenetes. `minikube service` command is merely for certain type of OS.

```sh
# deployment created, not accessible yet
kubectl create deployment hello-minikube1 --image=kicbase/echo-server:1.0

# deployment exposed on port 8080
kubectl expose deployment hello-minikube1 --type=NodePort --port=8080
```

Using kubectl, we can check the service port binding. 

```sh
kubectl get service hello-minikube

# NAME             TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
# hello-minikube   NodePort   10.97.53.50   <none>        8080:30856/TCP   167m
```

### LoadBalancer

According to documentation: *"A LoadBalancer service is the standard way to expose a service to the internet. With this method, each service gets its own IP address."*

As usual, create a new deployment, but it's not exposed by any means.

```sh
# create deployment without exposing it
kubectl create deployment hello-minikube1 --image=kicbase/echo-server:1.0
```

Expose the deployment using the service type LoadBalancer.

```sh
kubectl expose deployment hello-minikube1 --type=LoadBalancer --port=8080
```

Then we get service to check the ip and port assigned for this service.

```sh
kubectl get service

# NAME              TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
# hello-minikube1   LoadBalancer   10.97.163.13   <pending>     8080:31170/TCP   11s
# kubernetes        ClusterIP      10.96.0.1      <none>        443/TCP          4h43m
```

Since we are running Minikube on Darwin (Mac OS), we always need to open a tunnel in order to connect to our exposed services.

```sh
minikube tunnel

# ✅  Tunnel successfully started
# 📌  NOTE: Please do not close this terminal as this process must stay alive for the tunnel to be accessible ...
# 🏃  Starting tunnel for service hello-minikube1.
```

Then the external IP for our exposed services should be assigned.

```sh
kubectl get service

# NAME              TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
# hello-minikube1   LoadBalancer   10.97.163.13   127.0.0.1     8080:31170/TCP   2m33s
# kubernetes        ClusterIP      10.96.0.1      <none>        443/TCP          4h46m
```

With this setup, we can access `hello-minikube1` through `http://127.0.0.1:8080`.

## Pushing Images To Minikube And Deploy Them

Say we have a Golang app with the following Dockerfile.

```dockerfile
FROM golang:1.18-alpine3.17 as build
LABEL author="Yongjie.Zhuang"

WORKDIR /go/src/build/

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# build executable
COPY . .
RUN go build -o main


FROM alpine:3.17
WORKDIR /usr/src/
COPY --from=build /go/src/build/main ./main
COPY --from=build /go/src/build/app-conf-dev.yml ./app-conf-dev.yml
EXPOSE 8080

CMD ["./main"]
```

Then we build it with following command:

```sh
docker build . -t empty-head:latest
```

We can list the images using:

```sh
docker images

# REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
# empty-head       latest    9770f9061f21   59 minutes ago   27.1MB
```

Verify the image actually runs:

```sh
docker run -p 8080:8080 9770f9061f21
```

Then we push the image to Minikube:

- About pushing images to minikube: https://minikube.sigs.k8s.io/docs/handbook/pushing/

```sh
minikube image load empty-head:latest
```

We can check that the image is actually inside Minikube, which is the "docker.io/library/empty-head:latest"

```sh
minikube image list

# registry.k8s.io/pause:3.9
# registry.k8s.io/kube-scheduler:v1.26.3
# registry.k8s.io/kube-proxy:v1.26.3
# registry.k8s.io/kube-controller-manager:v1.26.3
# registry.k8s.io/kube-apiserver:v1.26.3
# registry.k8s.io/ingress-nginx/kube-webhook-certgen:<none>
# registry.k8s.io/ingress-nginx/controller:<none>
# registry.k8s.io/etcd:3.5.6-0
# registry.k8s.io/e2e-test-images/agnhost:2.39
# registry.k8s.io/coredns/coredns:v1.9.3
# gcr.io/k8s-minikube/storage-provisioner:v5
# gcr.io/k8s-minikube/minikube-ingress-dns:<none>
# docker.io/library/empty-head:latest
# docker.io/kubernetesui/metrics-scraper:<none>
# docker.io/kubernetesui/dashboard:<none>
# docker.io/kicbase/echo-server:1.0
```

Get the deployment file for our deployment: 

```sh
kubectl create deployment empty-head --image=docker.io/library/empty-head:latest -o yaml --dry-run=client > empty-head.yaml
```

which has the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: empty-head
  name: empty-head
spec:
  replicas: 1
  selector:
    matchLabels:
      app: empty-head
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: empty-head
    spec:
      containers:
      - image: docker.io/library/empty-head:latest
        name: empty-head
        resources: {}
status: {}
```

Then we add the imagePullPolicy, setting it to Never, so that it will actually use our cached docker image.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: empty-head
  name: empty-head
spec:
  replicas: 1
  selector:
    matchLabels:
      app: empty-head
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: empty-head
    spec:
      containers:
      - image: docker.io/library/empty-head:latest
        name: empty-head
        resources: {}
        imagePullPolicy: Never # CHANGED HERE!!!
status: {}
```

Lets create a deployment for it.

```sh
kubectl create -f empty-head.yaml

# deployment.apps/empty-head created
```

Once we have deployment created, we expose it as service.

```sh 
kubectl expose deployment empty-head --type=NodePort --port=8080
```

This app handles "https://0.0.0.0:8080/ping" endpoint. So we can do the following to request it:

```sh
kubectl expose deployment empty-head --type=NodePort --port=8080
# service/empty-head exposed

kubectl get service
# NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
# empty-head   NodePort    10.100.248.90   <none>        8080:32685/TCP   5s
# kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP          8h
# photon@Yongjie ~$ mk service empty-head
# |-----------|------------|-------------|---------------------------|
# | NAMESPACE |    NAME    | TARGET PORT |            URL            |
# |-----------|------------|-------------|---------------------------|
# | default   | empty-head |        8080 | http://192.168.49.2:32685 |
# |-----------|------------|-------------|---------------------------|
# 🏃  Starting tunnel for service empty-head.
# |-----------|------------|-------------|------------------------|
# | NAMESPACE |    NAME    | TARGET PORT |          URL           |
# |-----------|------------|-------------|------------------------|
# | default   | empty-head |             | http://127.0.0.1:53441 |
# |-----------|------------|-------------|------------------------|
# 🎉  Opening service default/empty-head in default browser...
# ❗  Because you are using a Docker driver on darwin, the terminal needs to be open to run it.

curl http://127.0.0.1:53441/ping
# pong at 2023-04-27 11:29:39.135640338 +0000 UTC m=+240.537862673
```



