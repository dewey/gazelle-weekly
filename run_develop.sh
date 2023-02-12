source develop.env

function cleanup() {
    rm -f gazelle-weekly
}
trap cleanup EXIT

# Compile Go
GO111MODULE=on GOGC=off go build -mod=vendor -v -o gazelle-weekly .
./gazelle-weekly