# Extend the Kubernetes API with CustomResourceDefinitions

And indeed, in Kubernetes, one can easily register custom resources. The procedure is fully dynamic and doesn't require restarting or updating the API server.

How such a custom resource can be added? Well, again, it's Kubernetes! Of course, by interacting with another, already existing resource! There is a special API resource called CustomResourceDefinition (CRD):

`The CustomResourceDefinition API resource allows you to define custom resources. Defining a CRD object creates a new custom resource with a name and schema that you specify.

Custom Resource Definition is a resource (with it's own endpoint) already added to kubernetes
You use it to define a new resource which will then also be given an endpoint
`

This page shows how to install a [custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) into the Kubernetes API by creating a [CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#customresourcedefinition-v1-apiextensions-k8s-io).

# TOC
[Create a CustomResourceDefinition]("##create-a-custom-resource-definition")

---

## Create a CustomResourceDefinition

Let's try to create a custom resource. Remember, a resource specifies a certain kind of Kubernetes object. Canonically, objects possess some attributes. So, our CustomResourceDefinition should be mostly concerned with describing the attributes of our future resource. Additionally, it's good to know that custom resources can be either namespaced or cluster-scoped. This is specified in the CRD's scope field.


`When you create a new CustomResourceDefinition (CRD), the Kubernetes API Server creates a new RESTful resource path for each version you specify.` The custom resource created from a CRD object can be either namespaced or cluster-scoped, as specified in the CRD's spec.scope field. As with existing built-in objects, deleting a namespace deletes all custom objects in that namespace. CustomResourceDefinitions themselves are non-namespaced and are available to all namespaces.


For example, if you have the following CustomResourceDefinition

```yaml
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    # CRD as in Model
    metadata:
        # name must match the spec fields below, and be in the form: <plural>.<group> i.e plural-name.group
        # name as in SomethingModel e.g blogpost, crontab
        # group is mostly like an organization name/domain e.g fluidcoins/v1/blogpost
        # name.group togther is expected to be a valid DNS domain
        name: crontabs.stable.example.com
    spec:
        # group name to use for REST API: /apis/<group>/<version>
        group: stable.example.com
        # list of versions supported by this CustomResourceDefinition
        names:
            # NAME IS JUST WHAT THE MODEL(RESOURCE) IS CALLED
            # OBJECTS ARE INSTANCES OF A RESOURCE
            # plural name to be used in the URL: /apis/<group>/<version>/<plural>
            plural: crontabs
            # singular name to be used as an alias on the CLI and for display
            singular: crontab
            # kind is normally the CamelCased singular type. Your resource manifests use this.
            kind: CronTab
            # shortNames allow shorter string to match your resource on the CLI
            shortNames:
            - ct
            listKind: CronTabList
        # either Namespaced or Cluster
        scope: Namespaced
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

```

Example 2

```yaml
kubectl apply -f - <<EOF
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: blogposts.iximiuz.com
spec:
  group: iximiuz.com
  names:
    kind: BlogPost
    listKind: BlogPostList
    plural: blogposts
    singular: blogpost
  scope: Namespaced
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        description: BlogPost is a custom resource exemplar
        type: object
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: BlogPostSpec is the spec for a BlogPost resource
            type: object
            properties:
              title:
                type: string
              author:
                type: string
          status:
            description: BlogPostStatus is the status for a BlogPost resource
            type: object
            properties:
              publishedAt:
                type: string
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```

To create a custom resource, feed the YAML definition to `kubectl apply`

```sh
kubectl apply -f resourcedefinition.yaml
```

Then a new namespaced RESTful API endpoint is created at:

```
/apis/stable.example.com/v1/namespaces/*/crontabs/...
```

For example 2, You can easily validate that the resource has been created:
```sh
kubectl api-resources --api-group=iximiuz.com
```

---

This endpoint URL can then be used to create and manage custom objects. The `kind` of these objects will be `CronTab` from the spec of the CustomResourceDefinition object you created above.

It might take a few seconds for the endpoint to be created. You can watch the `Established` condition of your CustomResourceDefinition to be true or watch the discovery information of the API server for your resource to show up.


## Create custom objects (Instances of a resource)
After the CustomResourceDefinition object has been created, you can create custom objects. Custom objects can contain custom fields. These fields can contain arbitrary JSON. In the following example, the cronSpec and image custom fields are set in a custom object of kind CronTab. The kind CronTab comes from the spec of the CustomResourceDefinition object you created above.

If you save the following YAML to my-crontab.yaml:

```yaml
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
spec:
  cronSpec: "* * * * */5"
  image: my-awesome-cron-image
```

and create it:

```sh
kubectl apply -f my-crontab.yaml
```

You can then manage your CronTab objects using kubectl. For example:

```sh
kubectl get crontab
```

Should print a list like this:

```sh
NAME                 AGE
my-new-cron-object   6s
```

Resource names are not case-sensitive when using kubectl, and you can use either the singular or plural forms defined in the CRD, as well as any short names.

You can also view the raw YAML data:

```sh
kubectl get ct -o yaml
```

You should see that it contains the custom cronSpec and image fields from the YAML you used to create it:

```yaml
apiVersion: v1
items:
- apiVersion: stable.example.com/v1
  kind: CronTab
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
                {"apiVersion":"stable.example.com/v1","kind":"CronTab","metadata":{"annotations":{},"name":"my-new-cron-object","namespace":"default"},"spec":{"cronSpec":"* * * * */5","image":"my-awesome-cron-image"}}
    creationTimestamp: "2021-06-20T07:35:27Z"
    generation: 1
    name: my-new-cron-object
    namespace: default
    resourceVersion: "1326"
    uid: 9aab1d66-628e-41bb-a422-57b8b3b1f5a9
  spec:
    cronSpec: '* * * * */5'
    image: my-awesome-cron-image
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
```

## Delete a CustomResourceDefinition
When you delete a CustomResourceDefinition, the server will uninstall the RESTful API endpoint and delete all custom objects stored in it.

```sh
kubectl delete -f resourcedefinition.yaml
kubectl get crontabs
```

The commmand returns error. If you later recreate the same CustomResourceDefinition, it will start out empty.

## Specifying a structural schema
CustomResources store structured data in custom fields (alongside the built-in fields `apiVersion`, `kind` and `metadata`, which the API server validates implicitly).`With OpenAPI v3.0 validation a schema can be specified, which is validated during creation and updates`.

With apiextensions.k8s.io/v1 the definition of a structural schema is mandatory for CustomResourceDefinitions. In the beta version of CustomResourceDefinition, the structural schema was optional.

A structural schema is an [OpenAPI v3.0 validation](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation) schema which:

- specifies a non-empty type (via type in OpenAPI) for the root, for each specified field of an object node (via properties or additionalProperties in OpenAPI) and for each item in an array node (via items in OpenAPI), with the exception of:

    a node with `x-kubernetes-int-or-string: true`
    a node with `x-kubernetes-preserve-unknown-fields: true`

- for each field in an object and each item in an array which is specified within any of allOf, anyOf, oneOf or not, the schema also specifies the field/item outside of those logical junctors (compare example 1 and 2).

- does not set description, type, default, additionalProperties, nullable within an allOf, anyOf, oneOf or not, with the exception of the two pattern for x-kubernetes-int-or-string: true (see below).

- if metadata is specified, then only restrictions on metadata.name and metadata.generateName are allowed.


Violations of the structural schema rules are reported in the NonStructural condition in the CustomResourceDefinition.

Examples of structural and non-structureal schemas:

Non-structural example 1:
```yaml
allOf:
- properties:
    foo:
      ...
```

conflicts with rule 2. The following would be correct:
```yaml
properties:
  foo:
    ...
allOf:
- properties:
    foo:
      ...
```

Non-structural example 2:
```yaml
allOf:
- items:
    properties:
      foo:
        ...
```

conflicts with rule 2. The following would be correct:
```yaml
items:
  properties:
    foo:
      ...
allOf:
- items:
    properties:
      foo:
        ...

```

Non-structural example 3:
```yaml
properties:
  foo:
    pattern: "abc"
  metadata:
    type: object
    properties:
      name:
        type: string
        pattern: "^a"
      finalizers:
        type: array
        items:
          type: string
          pattern: "my-finalizer"
anyOf:
- properties:
    bar:
      type: integer
      minimum: 42
  required: ["bar"]
  description: "foo bar object"
```

is not a structural schema because of the following violations:

- the type at the root is missing (rule 1).
- the type of foo is missing (rule 1).
- bar inside of anyOf is not specified outside (rule 2).
- bar's type is within anyOf (rule 3).
- the description is set within anyOf (rule 3).
- metadata.finalizers might not be restricted (rule 4).

`In contrast, the following, corresponding schema is structural:`

```yaml
type: object
description: "foo bar object"
properties:
  foo:
    type: string
    pattern: "abc"
  bar:
    type: integer
  metadata:
    type: object
    properties:
      name:
        type: string
        pattern: "^a"
anyOf:
- properties:
    bar:
      minimum: 42
  required: ["bar"]
```

## Field pruning
CustomResourceDefinitions store validated resource data in the cluster's persistence store, etcd. As with native Kubernetes resources such as ConfigMap, if you specify a field that the API server does not recognize, the unknown field is pruned (removed) before being persisted.