# SPDX-License-Identifier: Apache-2.0
PREFIX = boanlab
AGENT_NAME = sentryflow-agent
OPERATOR_NAME = sentryflow-operator
LOG_CLIENT_NAME = sentryflow-log-client
MONGO_CLIENT_NAME = sentryflow-mongo-client

AGENT_IMAGE_NAME = $(PREFIX)/$(AGENT_NAME)
OPERATOR_IMAGE_NAME = $(PREFIX)/$(OPERATOR_NAME)
LOG_CLIENT_IMAGE_NAME = $(PREFIX)/$(LOG_CLIENT_NAME)
MONGO_CLIENT_IMAGE_NAME = $(PREFIX)/$(MONGO_CLIENT_NAME)

TAG = v0.1

.PHONY: create-sentryflow
create-sentryflow: build-image
	kubectl label namespace default sentryflow istio-injection=enabled
	kubectl apply -f ./deployments/sentryflow.yaml
	sleep 1
	kubectl apply -f ./deployments/$(AGENT_NAME).yaml
	kubectl apply -f ./deployments/$(OPERATOR_NAME).yaml

.PHONY: create-sentryflow-operator
create-sentryflow-operator:
	kubectl label namespace default sentryflow istio-injection=enabled
	docker build -t $(OPERATOR_IMAGE_NAME):$(TAG) -f sentryflow/operator/Dockerfile .
	docker save -o $(OPERATOR_NAME)-$(TAG).tar $(OPERATOR_IMAGE_NAME):$(TAG)
	docker rmi $(OPERATOR_IMAGE_NAME):$(TAG)
	ctr -n=k8s.io image import $(OPERATOR_NAME)-$(TAG).tar
	rm -f $(OPERATOR_NAME)-$(TAG).tar
	kubectl apply -f ./deployments/sentryflow.yaml
	kubectl delete -f ./deployments/$(OPERATOR_NAME).yaml --ignore-not-found
	kubectl apply -f ./deployments/$(OPERATOR_NAME).yaml

.PHONY: create-sentryflow-agent
create-sentryflow-agent:
	kubectl label namespace default sentryflow istio-injection=enabled
	docker build -t $(AGENT_IMAGE_NAME):$(TAG) -f sentryflow/agent/Dockerfile .
	docker save -o $(AGENT_NAME)-$(TAG).tar $(AGENT_IMAGE_NAME):$(TAG)
	docker rmi $(AGENT_IMAGE_NAME):$(TAG)
	ctr -n=k8s.io image import $(AGENT_NAME)-$(TAG).tar
	rm -f $(AGENT_NAME)-$(TAG).tar
	kubectl apply -f ./deployments/sentryflow.yaml
	kubectl delete -f ./deployments/$(AGENT_NAME).yaml --ignore-not-found
	kubectl apply -f ./deployments/$(AGENT_NAME).yaml

.PHONY: create-client
create-client: delete-client
	kubectl apply -f ./deployments/log-client.yaml
	kubectl apply -f ./deployments/mongo-client.yaml
	
.PHONY: create-example
create-example: delete-example
	kubectl label namespace default sentryflow istio-injection=enabled
	kubectl apply -f examples/httpbin/httpbin.yaml -f examples/httpbin/sleep.yaml

.PHONY: delete-sentryflow
delete-sentryflow:
	kubectl delete all --all -n sentryflow --ignore-not-found
	kubectl delete namespace sentryflow --ignore-not-found

.PHONY: delete-sentryflow-operator
delete-sentryflow-operator:
	kubectl delete -f ./deployments/$(OPERATOR_NAME).yaml --ignore-not-found

.PHONY: delete-sentryflow-agent
delete-sentryflow-agent:
	kubectl delete -f ./deployments/$(AGENT_NAME).yaml --ignore-not-found

.PHONY: delete-client
delete-client:
	kubectl delete -f ./deployments/log-client.yaml --ignore-not-found
	kubectl delete -f ./deployments/mongo-client.yaml --ignore-not-found

.PHONY: delete-example
delete-example:
	kubectl delete -f examples/httpbin/httpbin.yaml -f examples/httpbin/sleep.yaml --ignore-not-found

.PHONY: build-image
build-image:
	docker build -t $(AGENT_IMAGE_NAME):$(TAG) -f sentryflow/agent/Dockerfile .
	docker build -t $(OPERATOR_IMAGE_NAME):$(TAG) -f sentryflow/operator/Dockerfile .
	docker build -t $(LOG_CLIENT_IMAGE_NAME):$(TAG) -f clients/log-client/Dockerfile .
	docker build -t $(MONGO_CLIENT_IMAGE_NAME):$(TAG) -f clients/mongo-client/Dockerfile .
	docker save -o $(AGENT_NAME)-$(TAG).tar $(AGENT_IMAGE_NAME):$(TAG)
	docker save -o $(OPERATOR_NAME)-$(TAG).tar $(OPERATOR_IMAGE_NAME):$(TAG)
	docker save -o $(LOG_CLIENT_NAME)-$(TAG).tar $(LOG_CLIENT_IMAGE_NAME):$(TAG)
	docker save -o $(MONGO_CLIENT_NAME)-$(TAG).tar $(MONGO_CLIENT_IMAGE_NAME):$(TAG)
	docker rmi $(AGENT_IMAGE_NAME):$(TAG)
	docker rmi $(OPERATOR_IMAGE_NAME):$(TAG)
	docker rmi $(LOG_CLIENT_IMAGE_NAME):$(TAG)
	docker rmi $(MONGO_CLIENT_IMAGE_NAME):$(TAG)
	ctr -n=k8s.io image import $(AGENT_NAME)-$(TAG).tar
	ctr -n=k8s.io image import $(OPERATOR_NAME)-$(TAG).tar
	ctr -n=k8s.io image import $(LOG_CLIENT_NAME)-$(TAG).tar
	ctr -n=k8s.io image import $(MONGO_CLIENT_NAME)-$(TAG).tar
	rm -f $(AGENT_NAME)-$(TAG).tar
	rm -f $(OPERATOR_NAME)-$(TAG).tar
	rm -f $(LOG_CLIENT_NAME)-$(TAG).tar
	rm -f $(MONGO_CLIENT_NAME)-$(TAG).tar

.PHONY: watch
watch:
	watch kubectl get all -n sentryflow