ext := if os() == "windows" { ".exe" } else { "" }

test:
    go test ./...

build:
    go build -o bin/nlr{{ext}} ./cmd/nlr

run-dl: build
    bin/nlr{{ext}} dl

run-count: build
    bin/nlr{{ext}} count

count-downloaded:
    find . -type f | wc -l

top-10-ids:
    jq "to_entries | sort_by(-.value) | from_entries" cards.json | head -11

last-ids:
    jq -r 'keys | map(tonumber) | sort | reverse' cards.json | head

total-cards:
    jq '[.[]] | add' cards.json

default: build
