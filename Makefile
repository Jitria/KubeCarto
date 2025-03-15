# SPDX-License-Identifier: Apache-2.0
PREFIX = boanlab
AGENT_NAME = sentryflow-agent
OPERATOR_NAME = sentryflow-operator
AGENT_IMAGE_NAME = $(PREFIX)/$(AGENT_NAME)
OPERATOR_IMAGE_NAME = $(PREFIX)/$(OPERATOR_NAME)
TAG = v0.1

.PHONY: create-sentryflow
apply: 
	kubectl apply -f ./deployments/sentryflow.yaml
	sleep 1
	kubectl apply -f ./deployments/$(AGENT_NAME).yaml
	kubectl apply -f ./deployments/$(OPERATOR_NAME).yaml

# client build image랑 밑에 것 해두기
.PHONY: create-client
create-client:
	kubectl apply -f ./deployments/client.yaml
	

.PHONY: delete-sentryflow
delete:
	kubectl delete all --all -n sentryflow
	kubectl delete namespace sentryflow

.PHONY: build-image
build-image:
	docker build -t $(AGENT_IMAGE_NAME):$(TAG) -f sentryflow/agent/Dockerfile .
	docker build -t $(OPERATOR_IMAGE_NAME):$(TAG) -f sentryflow/operator/Dockerfile .
	docker save -o $(AGENT_NAME)-$(TAG).tar $(AGENT_IMAGE_NAME):$(TAG)
	docker save -o $(OPERATOR_NAME)-$(TAG).tar $(OPERATOR_IMAGE_NAME):$(TAG)
	docker rmi $(AGENT_IMAGE_NAME):$(TAG)
	docker rmi $(OPERATOR_IMAGE_NAME):$(TAG)
	ctr -n=k8s.io image import $(AGENT_NAME)-$(TAG).tar
	ctr -n=k8s.io image import $(OPERATOR_NAME)-$(TAG).tar
	rm -f $(AGENT_NAME)-$(TAG).tar
	rm -f $(OPERATOR_NAME)-$(TAG).tar

