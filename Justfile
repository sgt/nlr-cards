ext := if os() == "windows" { ".exe" } else { "" }

build:
    @go build -o bin/nlr-find-max-id{{ext}} ./cmd/nlr-find-max-id
    @go build -o bin/nlr-dl{{ext}} ./cmd/nlr-dl

run-find-max-id: build
    bin/nlr-find-max-id{{ext}}

run-dl: build
    bin/nlr-dl{{ext}}

count-downloaded:
    find . -type f | wc -l

default: build
