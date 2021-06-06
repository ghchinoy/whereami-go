# whereami-go

A Go version of [whereami](https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/tree/master/whereami)

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https://github.com/ghchinoy/whereami-go&cloudshell_tutorial=README.md&cloudshell_workspace=/)




Whereami is a simple kubernetes-oriented Go app "useful for describing the location of the pod serving a request via its attributes (cluster name, cluster region, pod name, namespace, service account, etc). This is useful for a variety of demos where you just need to understand how traffic is getting to and returning from your app."

## work in progress

This version does not mimic the [whereami python version](https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/tree/master/whereami) completely. Missing are:

* gRPC