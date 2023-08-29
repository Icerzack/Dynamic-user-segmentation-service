ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=admin password=admin dbname=app host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
SQL_FOLDER=$(INTERNAL_PKG_PATH)/db/sql

.PHONY: compose-db-up
compose-db-up:
	docker-compose build
	docker-compose up -d postgres

.PHONY: compose-db-rm
compose-db-rm:
	docker-compose down

.PHONY: compose-app-up
compose-app-up:
	docker build -f build/Dockerfile -t avito-app .
	docker-compose up -d avito-app

.PHONY: compose-app-rm
compose-app-rm:
	docker-compose down

.PHONY: compose-all-up
compose-all-up: compose-db-up compose-app-up

.PHONY: compose-all-rm
compose-all-rm: compose-app-rm compose-db-rm