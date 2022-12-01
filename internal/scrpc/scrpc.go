package scrpc

import "github.com/pja237/slurmcommander-dev/internal/slurm"

// Main object that is passed around: fetcher->bank->listener->RPC Client
type CachedSqueue struct {
	Counter    uint
	SqueueJSON slurm.SqueueJSON
}

// RPC Request arguments comming in from client
type ReqArgs struct {
	Cid  uint   // client id
	Cstr string // client string
}

// RPC Response sent back to client
type ReplyArgs struct {
	N int
	CachedSqueue
}
