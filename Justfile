ext := if os() == "windows" { ".exe" } else { "" }

build:
    @go build -o bin/findmax{{ext}} ./cmd/findmax

run: build
    bin/findmax{{ext}}

default: build
