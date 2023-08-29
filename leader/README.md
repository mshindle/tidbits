# leader

Note on configuration, instead of implementing service discovery - which should be done for a production application - the app makes the following assumptions:

 * Nodes live on ports starting at 6001 and increase sequentially based on rank. So rank 2 lives at 6002, rank 4 lives at 6004, etc.
 * Nodes all live at the same IP address (i.e. localhost).
 * Node IDs are formed automatically from the rank value (e.g. rank 1 node is "node-01"). In production, a UUID should be used here.
 * The total number of nodes available for leader election should be passed as the NODES environment variable or command line.

See [Implement a basic leader election using Go and RPC](https://itnext.io/lets-implement-a-basic-leader-election-algorithm-using-go-with-rpc-6cd012515358) for more information about how this algorithm was implemented.

