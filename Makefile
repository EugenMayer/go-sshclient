run-tests:
	docker-compose up -d sshserver
	sleep 5s
	docker-compose up test
