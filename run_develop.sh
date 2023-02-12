source develop.env

function cleanup() {
    rm -f redacted-weekly
}
trap cleanup EXIT

# Compile Go
GO111MODULE=on GOGC=off go build -mod=vendor -v -o redacted-weekly .
./redacted-weekly