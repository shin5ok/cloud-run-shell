
.PHONY: client
client:
	( cd client ; go build -o client . )

.PHONY: deploy
deploy:
	bash ./deploy.sh
