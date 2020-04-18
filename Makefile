run:
	docker-compose up -d
	#./run.sh

runlogs:
	docker-compose up

restart:
	docker-compose restart app

logs:
	docker logs -f storage-api_app_1

stop:
	docker-compose down