
.PHONY: client
client:
	( cd client ; go build -o shell . )

.PHONY: deploy
deploy:
	bash ./deploy.sh
