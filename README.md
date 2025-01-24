This package provides the API for handling Breakdowns

### Local Development set up

go build - build package
go run main.go - start the server

### MongoDB connection

TODO

- Currently using a whitelisted IP address 0.0.0.0/0 (Global access restricted by time)

# API

## Endpoints

/GET health - health check

### Breakdown

/GET breakdown/$user_id - retrieve all breakdowns for user
/POST breakdown - create breakdown
