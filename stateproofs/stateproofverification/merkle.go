package stateproofverification

const (
	// MaxEncodedTreeDepth is the maximum tree depth (root only depth 0) for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxEncodedTreeDepth = 16

	// MaxNumLeavesOnEncodedTree is the maximum number of leaves allowed for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxNumLeavesOnEncodedTree = 1 << MaxEncodedTreeDepth
)
