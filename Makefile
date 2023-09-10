.PHONY: client
client:
	( cd client ; go build -o shell . )

.PHONY: deploy
deploy:
	bash ./deploy.sh

.PHONY: bench
bench:
	( cd bench ; go test -timeout 3600m -bench . )
