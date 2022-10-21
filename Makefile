app: parser
	./coai
	
parser:
	./config/cfg

.PHONY: proto
proto:
	@./scripts/proto.sh users

generate:
	go generate -v ./...

test:
	@./scripts/test.sh users