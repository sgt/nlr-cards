ext := if os() == "windows" { ".exe" } else { "" }

build:
    @go build -o bin/nlr-find-max{{ext}} ./cmd/nlr-find-max
    @go build -o bin/nlr-dl{{ext}} ./cmd/nlr-dl

run-find-max: build
    bin/nlr-find-max{{ext}}

run-dl: build
    bin/nlr-dl{{ext}}

count-downloaded:
    find . -type f | wc -l

default: build
