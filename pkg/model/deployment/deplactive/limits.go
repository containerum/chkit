package deplactive

import "github.com/ninedraft/ranger/intranger"

const (
	MIN_REPLICAS = 1
	MAX_REPLICAS = 15

	MAX_CPU = 3000
	MIN_CPU = 10

	MAX_MEM = 8000
	MIN_MEM = 10
)

var (
	MemLimit      = intranger.IntRanger(MIN_MEM, MAX_MEM)
	CPULimit      = intranger.IntRanger(MIN_CPU, MAX_CPU)
	ReplicasLimit = intranger.IntRanger(MIN_REPLICAS, MAX_REPLICAS)
)
