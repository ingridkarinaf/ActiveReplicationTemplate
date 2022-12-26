# ActiveReplicationTemplate
Template for active replication client-server structure resilient to one node crashing (crash-failure)

You have to implement the following distributed system: A set of nodes tries to implement a distributed Increment function. We want the following properties to hold:

- Every call to Increment returns a value. Value of Increment is always a natural number.
- First call to Increment returns 0.
- Monotonicity: For every node, if it has subsequent calls to Increment, with values Inc_1 and Inc_2 respectively, we have that Inc_1 < Inc_2
- Uniqueness: No call to Increment should ever return a previously returned value.
- Liveness: every call to Increment returns a value
- Reliability: your system can tolerate a crash-failure of 1 node.
- Partial submissions are accepted, e.g., a not fully working implementation.

## Set-up
Run the following in separate terminals:

`go run rmserver.go 5000`
`go run rmserver.go 5001`
`go run rmserver.go 5002`
`go run feserver.go 4000`
`go run feserver.go 4001`
`go run client.go 4000`
`go run client.go 4001`

See activeRep.png for visual represntation of client-server archiceture.