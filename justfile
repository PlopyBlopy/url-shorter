set shell := ["powershell.exe", "-NoProfile", "-c"]

# default go run
run:
    go run .\cmd\app\main.go 

# go run with env development
run_dev:
    go run .\cmd\app\main.go  -env dev

# go run with env production
run_prod:
    go run .\cmd\app\main.go  -env prod

# test
test:
    go test .\...

# test all
test_all:
    go test -tags="integration,e2e" .\...

# test with escape analyzis
test_ea:
    go test -gcflags='-m -l' .\...

# integration tests
test_integration:
    go test -tags=integration .\...

# e2e tests
test_e2e:
    go test -tags=e2e .\...
