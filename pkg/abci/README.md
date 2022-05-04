# Bhojpur State - Application BlockChain Interface

The `Blockchain` systems replicate data as multi-master state machine. The **ABCI** is an
interface that defines the boundary between the replication engine (the blockchain), and
the state machine (the application). Using a socket protocol, a consensus engine running
in one process can manage an application state running in another.
