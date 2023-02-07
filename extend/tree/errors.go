package tree

import "errors"

var (
	ErrNodeNotFound        = errors.New("node is not found")
	ErrNodeExisting        = errors.New("node is already existing")
	ErrTreeBalancingFailed = errors.New("fail to balance tree")
	ErrNodeDeletingFailed  = errors.New("fail to balance tree")
)
