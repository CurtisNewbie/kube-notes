# kube-notes

Notes for kubernetes

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

Unfortunately, my Mac is stil using Catalina OS, which is pretty old. Colima and its dependencies must be compiled rather than automatically downloaded and installed (by Homebrew, e.g., LLVM). I chose to use an old version of Docker Desktop (4.15.0).

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


### 3. Create Deployment

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

We can notice that the EXTERNAL-IP is still none. We create a port-forward to 8080 for this specific service. The service exposes 8080 port, and we forward any traffic from 7080 to 8080. Then we have access to the service, using `http://localhost:7080`. *"Kubectl port-forward is a method to access, interact and manage internal Kubernetes clusters directly from your local network."*

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



























