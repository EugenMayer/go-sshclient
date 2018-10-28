test: init
	CGO_ENABLED=0 go test -tags netgo test/*.go

init:
	dep ensure

run-tests:
	docker-compose up -d sshserver
	sleep 5s
	docker-compose up test
