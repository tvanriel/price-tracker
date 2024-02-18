.PHONY: docker encoder dev-discord dev-web web
docker:
	docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/price-tracker:latest .
dev:
	CompileDaemon -build "go build -o /tmp/price-tracker ." -command "/tmp/price-tracker"

