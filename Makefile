build:
	export CGO_ENABLED=1 && go build -o ./chat2data ./cmd/chat2data