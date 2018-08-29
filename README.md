# Ambassador Operator (Operator-SDK version)

**Project Status: Experimental**

> [Ambassador](https://www.getambassador.io/), the "Kubernetes-native API gateway for microservices built on Envoy", currently pulls its configuration (mainly Mappings of URL prefixes to Kubernetes Services) [from annotations on Services](https://www.getambassador.io/reference/configuration); I wanted Ambassador to be even more Kubernetes native, so I've created a Mapping CRD and a controller that maintains a dummy Service for each Mapping, annotated according to the Mapping's Spec.

From ["Kubernetes Custom Resource, Controller and Operator Development Tools"](https://admiralty.io/kubernetes-custom-resource-controller-and-operator-development-tools.html) (Admiralty blog); the post discusses the Operator SDK, Kubebuilder, and Metacontroller. As a use case, three different yet functionally equivalent Ambassador operators were created, one with each tool. This repository was built with the Operator SDK v0.0.5.

## Getting Started

1. [Install Ambassador](https://www.getambassador.io/user-guide/install) (steps 1 and 2). *N.B.*: you may have to wait a few minutes before the Ambassador Service's external IP is provisioned, but keep going, it's only blocking for step 5 below.
1. Install ambassador-shim-operator-sdk, which consists of:
	- a CustomResourceDefinition: Mapping,
	- a one-replica Deployment: the operator,
	- RBAC for the operator to control Mappings and Services.
	```sh
	kubectl apply -f https://raw.githubusercontent.com/admiraltyio/ambassador-shim-operator-sdk/master/deploy/operator.yaml
	kubectl apply -f https://raw.githubusercontent.com/admiraltyio/ambassador-shim-operator-sdk/master/deploy/rbac.yaml
	```
1. Deploy the sample: a stock NGINX Deployment and corresponding Service and Mapping:
	```sh
	kubectl apply -f https://raw.githubusercontent.com/admiraltyio/ambassador-shim-operator-sdk/master/examples/getting_started/ambassadorshim_v1alpha1_mapping.yaml
	```
1. Verify that the Mapping was configured and is up-to-date, _i.e._, that a corresponding dummy Service was created and properly annotated:
	```sh
	kubectl get mapping foo -o yaml
	kubectl get service foo-ambassadorshim -o yaml
	```
1. Once the external IP from step 1 is provisioned, you can access the foo service at /foo/:
	```sh
	EXTERNAL_IP=$(kubectl get service ambassador -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
	curl http://$EXTERNAL_IP/foo/
	```

## Next Steps

This shim is a proof of concept that only supports the most basic feature of Ambassador: mapping a prefix to a Service. It could be expanded (full Mapping Spec and Status, other CRDs, etc.) into a fully-featured side-car of Ambassador, or merge with it. See the full discussion in [the blog post](https://admiralty.io/kubernetes-custom-resource-controller-and-operator-development-tools.html).

## See Also

- [ambassador-shim-kubebuilder](https://github.com/admiraltyio/ambassador-shim-kubebuilder)
- [ambassador-shim-metacontroller](https://github.com/admiraltyio/ambassador-shim-metacontroller)
