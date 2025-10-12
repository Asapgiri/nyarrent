name=nyarrent

build:
	go build -o ${name}

release:
	go build -o ${name} -ldflags "-s -w"

run:
	echo -e "\e[32mStarting Server!\e[0m"
	./${name}

reflex:
	reflex -R '\.git' -r '\.go' -s -- make build run

restart:
	sudo systemctl restart nyarrent.service
