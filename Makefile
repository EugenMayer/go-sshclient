init:
	go mod tidy
	go mod verify
	go mod vendor

update:
	go get -u
	go mod tidy

run-tests:
	docker-compose up -d sshserver
	sleep 5s
	docker-compose up test
