# ActiveReplicationTemplate
Template for active replication client-server structure resilient to one node crashing (crash-failure)

You have to implement the following distributed system: A set of nodes tries to implement a distributed Increment function. We want the following properties to hold:

## Set-up
Run the following in separate terminals:

`go run rmserver.go 5000`
`go run rmserver.go 5001`
`go run rmserver.go 5002`
`go run feserver.go 4000`
`go run feserver.go 4001`
`go run client.go 4000`
`go run client.go 4001`
