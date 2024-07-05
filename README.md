# Port forwarding issue
It looks like Kubernetes port-forwarding cannot deal with a cancelled TCP request. It will show the following error and the port forwading is dead:
```
E0705 12:50:42.897359   76107 portforward.go:394] error copying from local connection to remote stream: read tcp6 [::1]:8080->[::1]:44368: read: connection reset by peer
E0705 12:50:42.897360   76107 portforward.go:381] error copying from remote stream to local connection: readfrom tcp6 [::1]:8080->[::1]:44368: write tcp6 [::1]:8080->[::1]:44368: write: broken pipe
```

## How to reproduce
I have created a very simple HTTP server that responds with a 30MiB empty file to every HTTP request.

### How it should work (local)
The server can be run locally by invoking `run server.go` and it will expose the service on port 8080. Now run `curl -o /dev/null http://localhost:8080` to fetch the 30MiB file from the server. Curl will download the file and also the server shows the following output:
```plain
2024/07/05 12:53:08 5: got request from [::1]:50194
2024/07/05 12:53:08 5: completed request
```

When invoking `curl http://localhost:8080` it will invoke the request, but terminate with `Warning: Binary output can mess up your terminal.` and cancel the HTTP request. You can run a cancelled request multiple times without any issue, but the serer will shows the following output:
```
2024/07/05 12:52:12 1: got request from [::1]:57244
2024/07/05 12:52:12 1: error writing data (32702464 bytes left to write): write tcp [::1]:8080->[::1]:57244: write: connection reset by peer
2024/07/05 12:52:13 2: got request from [::1]:57246
2024/07/05 12:52:13 2: error writing data (32374784 bytes left to write): write tcp [::1]:8080->[::1]:57246: write: connection reset by peer
2024/07/05 12:52:14 3: got request from [::1]:57262
2024/07/05 12:52:14 3: error writing data (32505856 bytes left to write): write tcp [::1]:8080->[::1]:57262: write: connection reset by peer
```

### How it fails (Kubernetes)
For this test, I use [kind](https://kind.sigs.k8s.io/) to create a Kubernetes cluster, so make sure it has been installed. Now run the server inside Kubernetes by invoking `make` that will do the following:
- Build the `server` executable for use in a Docker image.
- Create the Docker image.
- Create the Kind cluster.
- Load the Docker image into the Kubernetes cluster and spin up the `server` pod
- Start a port-forwarding to the `server` pod

Now run `curl -o /dev/null http://localhost:8080` to fetch the 30MiB file from the server. Curl will download the file and this can be done multiple times. Now run `curl http://localhost:8080` and it will be cancelled, because of the binary output. The first time it seems to work fine, but it's impossible to invoke the command again (or connect any other way to this port). The port-forwarding needs to be restarted to deal with this.

The script also deploys a `curl` pod, that can be accessed using `kubectl exec -it curl -- /bin/sh`. Even when the port-forwarding has crashed, you can still invoke `curl http://server:8080`, so the actual server and pod is perfectly fine. This illustrates that the port-forwarding is broken.

The service that exposes the port as a `NodePort` on the Kubernetes node also works fine. This can be tested by invoking `docker exec -it kind-control-plane /bin/bash` and run `curl http://localhost:30303` several times.