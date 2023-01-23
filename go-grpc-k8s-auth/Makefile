.PHONY: all deploy

all: deploy

deploy:
	kubectl kustomize example/server/config/default | ko apply --sbom none -f -
