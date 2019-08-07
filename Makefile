# These are the values we want to pass for Version and BuildTime
GITTAG=`git describe --tags`
COMMIT=`git rev-parse HEAD`
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
FLAGS=-ldflags "-X main.Commit=${COMMIT} -X main.Build=${BUILD_TIME}"

init:
	[ -d apidoc ] || git submodule add --force ssh://git@git.yixindev.net:2201/general/apidoc.git apidoc;
	[ -d codes ] || git submodule add --force ssh://git@git.yixindev.net:2201/general/codes.git codes;
	[ -d manifest ] || git submodule add --force ssh://git@git.yixindev.net:2201/devops/manifest.git manifest;
clean:
	@go clean
	@rm -rf pb
pb: $(shell which protoc)
	@rm -rf pb
	@mkdir -p pb
#	protoc -I apidoc apidoc/milter/v1/*.proto --go_out=plugins=xgrpc:./pb/
code:
	go generate codes/codes.go

submodule: init
	git submodule init
	git submodule update

update: submodule
	git submodule update --recursive --remote; 

gen: pb code

build:
	kubectl config use-context develop-admin@develop
	go build -a -installsuffix cgo ${FLAGS}

image: build
	docker build -t registry.yixindev.net:5000/milter-test:${COMMIT} .
	docker push registry.yixindev.net:5000/milter-test:${COMMIT}
	docker tag registry.yixindev.net:5000/milter-test:${COMMIT} registry.yixindev.net:5000/milter-test:latest
	docker push registry.yixindev.net:5000/milter-test:latest
	docker image remove registry.yixindev.net:5000/milter-test:${COMMIT}

