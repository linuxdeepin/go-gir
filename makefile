bin/gir-generator:
	mkdir -p bin
	GOPATH=/dev/shm/gir-generator/vendor:/dev/shm/gir-generator go build -o bin/gir-generator gir-generator
