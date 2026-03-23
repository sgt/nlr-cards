ext := if os() == "windows" { ".exe" } else { "" }

build:
    @go build -o bin/nlr{{ext}} ./cmd/nlr

run-dl: build
    bin/nlr{{ext}} dl

run-count: build
    bin/nlr{{ext}} count

count-downloaded:
    find . -type f | wc -l

default: build
