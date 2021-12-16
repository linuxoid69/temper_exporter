.PHONY: all  build clean

IMAGE="alpine:3.14"

all:
	@printf '%*s' "$${COLUMNS:-$$(tput cols)}" '' | tr ' ' -
	@printf '%s\n' "$$(tput setaf 2)" 'Run `make build` for build'
	@printf '%s\n' "$$(tput setaf 2)" 'Run `make clean` for clean'
	@printf '%s' "$$(tput init)"
	@printf '%*s\n' "$${COLUMNS:-$$(tput cols)}" '' | tr ' ' -


prepare_dir_for_build:
	mkdir -p target
	mkdir -p rootfs/usr/bin/
	mkdir -p rootfs/usr/lib/systemd/system/
	mkdir -p rootfs/etc/default/

clean:
	rm -rf rootfs
	rm -rf target

docker: docker_check
	docker run -it --rm -v $(PWD):/mnt $(IMAGE) \
	sh -c "apk update && apk add go ruby-etc ruby ruby-dev make upx \
	&& gem install fpm && cd /mnt && GOARCH=arm go build \
	&& install -c temper_exporter rootfs/usr/bin/ \
	&& install -c temper_exporter.service rootfs/usr/lib/systemd/system/"

docker_check:
	@which docker  || printf '%s\n' "$$(tput setaf 1)" 'Please install docker'

build: prepare_dir_for_build docker
