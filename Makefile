REGION := asia-northeast1
SIDECAR := $(REGION)-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/sidecar
shellapp := $(REGION)-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/shellapp
SERVICE_NAME := $(SERVICE_NAME)
REPO_NAME := my-app

.PHONY: client
client:
	( cd client/ ; go build -o shell . )

.PHONY: deploy
deploy:
	REGION=$(REGION) DEPLOY_START=$(shell date '+%Y%m%d%H%M%S') envsubst < cloudrun.yaml | gcloud run services replace - --region=$(REGION)

.PHONY: bench
bench:
	( cd bench/ ; go test -timeout 3600m -count 1 -bench . )

.PHONY: sidecar
sidecar:
	( cd sidecar/ ; docker build --platform=linux/amd64 -t $(SIDECAR) . )
	docker push $(SIDECAR)

.PHONY: shellapp
shellapp:
	( cd shellapp/ ; docker build --platform=linux/amd64 -t $(shellapp) . )
	docker push $(shellapp)

.PHONY: all
all: sidecar shellapp deploy

.PHONY: service
service:
	gcloud services enable compute.googleapis.com run.googleapis.com artifactregistry.googleapis.com

.PHONY: repo
repo: service
	gcloud artifacts repositories describe --location=$(REGION) $(REPO_NAME) > /dev/null 2>&1 || \
		gcloud artifacts repositories create --location=$(REGION) --repository-format=docker $(REPO_NAME)

.PHONY: expose
expose:
	gcloud run services add-iam-policy-binding --member=allUsers $(SERVICE_NAME) --region=$(REGION) --role=roles/run.invoker
