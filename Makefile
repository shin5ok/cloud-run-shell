SIDECAR := asia-northeast1-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/sidecar
shellapp := asia-northeast1-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/shellapp
REGION := asia-northeast1
SERVICE_NAME := $(SERVICE_NAME)

.PHONY: client
client:
	( cd client/ ; go build -o shell . )

.PHONY: deploy
deploy:
	DEPLOY_START=$(shell date '+%Y%m%d%H%M%S') envsubst < cloudrun.yaml | gcloud run services replace - --region=$(REGION)
	gcloud run services add-iam-policy-binding --member=allUsers $(SERVICE_NAME) --region=$(REGION) --role=roles/run.invoker

.PHONY: bench
bench:
	( cd bench/ ; go test -timeout 3600m -count 1 -bench . )

.PHONY: sidecar
sidecar:
	( cd sidecar/ ; docker build -t $(SIDECAR) . )
	docker push $(SIDECAR)

.PHONY: shellapp
shellapp:
	( cd shellapp/ ; docker build -t $(shellapp) . )
	docker push $(shellapp)

.PHONY: all
all: sidecar shellapp deploy
