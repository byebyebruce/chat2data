build:
	export CGO_ENABLED=0 && go build -o ./chat2data ./cmd/chat2data