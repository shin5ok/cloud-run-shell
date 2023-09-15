.PHONY: client
client:
	( cd client/ ; go build -o shell . )

.PHONY: deploy
deploy:
	( cd shellapp/ ; bash ./deploy.sh )

.PHONY: bench
bench:
	( cd bench/ ; go test -timeout 3600m -count 1 -bench . )
