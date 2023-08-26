ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=admin password=admin dbname=app host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
SQL_FOLDER=$(INTERNAL_PKG_PATH)/db/sql

.PHONY: compose-up
compose-up:
	docker-compose build
	docker-compose up -d postgres

.PHONY: compose-rm
compose-rm:
	docker-compose down