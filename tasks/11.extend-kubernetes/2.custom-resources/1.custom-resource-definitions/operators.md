## Kubernetes Operators

Custom resources are neat but useless in isolation. You need some custom code to interact with them:
On their own, custom resources simply let you store and retrieve structured data. When you combine a custom resource with a custom controller, custom resources provide a true declarative API.

`Logic in Kubernetes is organized in form of controllers. There is an extensive set of built-in controllers watching the built-in resources and applying changes to the clusters.` But how can we get some custom controllers?

Well, apparently, it's a simple as just starting a Pod. I.e. everything you need is to program a control loop logic in a language of your choice (it'd be wise to use one of the official API clients though), pack this program into a [Docker] image, and deploy it to your cluster.


Actually, you can run such code wherever you like. I.e. it's not mandatory to run it inside of a Kubernetes cluster. For example, code can run on a stand-alone virtual or bare metal machine, assuming it has sufficient permissions to call Kubernetes API. But to be honest, I don't see many good reasons for doing so.


So, what is a Kubernetes Operator? Citing the official docs one more time, "an operator is a client of the Kubernetes API that acts as a controller for a Custom Resource." As simple as just that. Of course, you can deviate from this definition a bit by, say, adding multiple Custom Resources, but one of the best practices of writing operators states that you'd need to introduce multiple control loops (i.e. a controller per resource) to keep the implementation clear.

What kind of logic is a good fit for an operator? Guessing on the name of the pattern, it should have something to do with operating [an application]. Originally, operators were meant to automate human beings out of the application management. The Operator Pattern page mentions some example use cases:

- deploying, backuping, restoring, or upgrading an application;
- exposing a Kubernetes service to non-Kubernetes applications;
- chaos engineering (i.e. simulating failures);
- automatically choosing a leader for a distributed application.

You can also find lots of real-world operators on OperatorHub.io. However, in my opinion, your imagination is the only thing that bounds what an operator could do.


## CRD annotations

While Custom Resources are great, Custom Resources are also a challenge today because there is no good way for YAML developers to understand their usage or their semantics. For instance, when a Custom Resource is created, what all Kubernetes resources get created behind the scene? Or, what labels/annotations are important from the perspective of a Custom Resource and its Operator? Or, what is the CPU/Memory/Storage consumption of a Custom Resource instance?

In this post we present a generic technique that we have developed to address such questions about Kubernetes Custom Resources. You can use this technique on all required Operators on your cluster and significantly simplify usage of Custom Resources for your YAML developers.

The basis for our technique is the observation that inter-resource relationships form one of the fundamental aspects for Kubernetes YAML automation. Some examples of inter-resource relationships between Kubernetes built-in resources are — a Pod is related to a PersistentVolumeClaim through Spec property; a Service is related to a Pod through labels; a Pod is related to a ReplicaSet through ownerReference.

![Building Kubernetes YAML automation][def]

[def]: https://miro.medium.com/max/1400/0*_5t3G86mpACcvofl