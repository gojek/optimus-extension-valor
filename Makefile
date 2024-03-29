coverage_file=coverage.html
binary_outdir=out
project_name=valor

test:
	go test ./... --cover

coverage:
	go test -coverprofile ${coverage_file} ./... && go tool cover -html=${coverage_file}

bin: test
	go build -o ${binary_outdir}/${project_name} .

dist: test
	./script/build.sh
