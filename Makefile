clean:
	@rm -rf ./bin
	@mkdir -p ./bin/temp
	@mkdir -p ./bin/bundle
	@mkdir -p ./bin/data

dist-tools:
	go get -u github.com/pkieltyka/fresh
	go get -v github.com/robfig/glock

deps:
	@glock -v sync github.com/pressly/selfie

build: clean
	@mkdir -p ./bin
	cd cmd/selfie && GOGC=off go build -i -o ../../bin/selfie

build-all: clean
	cd cmd/selfie; \
	for GOOS in darwin linux windows; do \
		for GOARCH in 386 amd64; do \
			echo "building $$GOOS $$GOARCH ..."; \
			export GOOS=$$GOOS; \
			export GOARCH=$$GOARCH; \
			go build -o ../../bin/selfie-$$GOOS-$$GOARCH; \
		done \
	done

dev: kill
	@(export CONFIG=$$PWD/etc/selfie.conf && \
		cd ./cmd/selfie && \
		fresh -c ../../etc/fresh-runner.conf -w=../..)

##
## Database mgmt
##
reset-db:
	bash scripts/init_db.sh

kill-fresh:
	ps -ef | grep 'f[r]esh' | awk '{print $$2}' | xargs kill

kill-by-port:
	lsof -t -i:7331 | xargs kill

kill: kill-fresh kill-by-port


docker-rm-all-images: docker-rm-existing-ps
	docker images | grep "^<none>" | awk '{print $$3}' | xargs docker rmi

docker-rm-existing-ps:
	docker ps -a | awk 'NR>1' | awk '{print $$1}' | xargs docker rm

docker-dev:
	docker run -d -p 5432:5432 -v $$pwd/bin/data:/var/lib/postgresql/data --restart=always --name postgres -e POSTGRES_PASSWORD=betame postgres

docker-psql:
	docker run -it --rm --link selfiedb:selfiedb selfiedb sh -c 'exec psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U selfiedb'

start-db:
	docker run -d \
						 -p 5432:5432 \
						 -v $$pwd/bin/data:/var/lib/postgresql/data \
						 --restart=always \
						 --name postgres \
						 -e POSTGRES_PASSWORD=betame \
						 -e POSTGRES_USER=ali \
						 postgres


kill-selfie:
	docker ps | awk 'NR>1' | grep selfie | awk '{print $$1}' | xargs docker kill
	docker ps | awk 'NR>1' | grep selfiedb | awk '{print $$1}' | xargs docker kill
	docker ps -a | awk 'NR>1' | grep selfie | awk '{print $$1}' | xargs docker rm
	docker ps -a | awk 'NR>1' | grep selfiedb | awk '{print $$1}' | xargs docker rm

run-prod: kill-selfie
	docker run -d --name selfiedb selfiedb
	docker run -d -p 7331:7331 --name selfie --restart=always --link selfiedb selfie
