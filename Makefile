.PHONY: up-infra
up-infra:
	docker-compose up pgdb kafka-ui zookeeper kafka-1 --build -d

.PHONY: down-infra
down-infra:
	docker-compose down -v