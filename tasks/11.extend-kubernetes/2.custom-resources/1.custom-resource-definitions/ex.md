# Extend the Kubernetes API with CustomResourceDefinitions

This page shows how to install a [custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) into the Kubernetes API by creating a [CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#customresourcedefinition-v1-apiextensions-k8s-io).

# TOC
[Create a CustomResourceDefinition]("##create-a-custom-resource-definition")

---

## Create a CustomResourceDefinition

`When you create a new CustomResourceDefinition (CRD), the Kubernetes API Server creates a new RESTful resource path for each version you specify.` The custom resource created from a CRD object can be either namespaced or cluster-scoped, as specified in the CRD's spec.scope field. As with existing built-in objects, deleting a namespace deletes all custom objects in that namespace. CustomResourceDefinitions themselves are non-namespaced and are available to all namespaces.


For example, if you have the following CustomResourceDefinition

```yaml
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    # CRD as in Model
    metadata:
        # name must match the spec fields below, and be in the form: <plural>.<group>
        # name as in SomethingModel e.g BookCollectionModel
        name: crontabs.stable.example.com
    spec:
        # group name to use for REST API: /apis/<group>/<version>
        group: stable.example.com
        # list of versions supported by this CustomResourceDefinition
        versions:
            - name: v1
              # Each version can be enabled/disabled by Served flag.
              served: true
              # One and only one version must be marked as the storage version.
              storage: true
              schema:
                openAPIV3Schema:
                    type: object
                    properties:
                        spec:
                            type: object
                            properties:
                                cronSpec:
                                    type: string
                                image:
                                    type: string
                                replicas:
                                    type: integer
        # either Namespaced or Cluster
        scope: Namespaced
        names:
            # plural name to be used in the URL: /apis/<group>/<version>/<plural>
            plural: crontabs
            # singular name to be used as an alias on the CLI and for display
            singular: crontab
            # kind is normally the CamelCased singular type. Your resource manifests use this.
            kind: CronTab
            # shortNames allow shorter string to match your resource on the CLI
            shortNames:
            - ct

```

And create it:

```sh
kubectl apply -f resourcedefinition.yaml
```

Then a new namespaced RESTful API endpoint is created at:

```
/apis/stable.example.com/v1/namespaces/*/crontabs/...
```

---

This endpoint URL can then be used to create and manage custom objects. The `kind` of these objects will be `CronTab` from the spec of the CustomResourceDefinition object you created above.

It might take a few seconds for the endpoint to be created. You can watch the `Established` condition of your CustomResourceDefinition to be true or watch the discovery information of the API server for your resource to show up.


## Create custom objects
After the CustomResourceDefinition object has been created, you can create custom objects. 