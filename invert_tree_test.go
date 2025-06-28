package main

import (
	"testing"
)

// Helper function to check if two trees are equal
func areTreesEqual(t1 *TreeNode, t2 *TreeNode) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Val == t2.Val &&
		areTreesEqual(t1.Left, t2.Left) &&
		areTreesEqual(t1.Right, t2.Right)
}

func TestInvertTree(t *testing.T) {
	tests := []struct {
		name     string
		input    *TreeNode
		expected *TreeNode
	}{
		{
			name:     "nil tree",
			input:    nil,
			expected: nil,
		},
		{
			name: "single node",
			input: &TreeNode{
				Val: 1,
			},
			expected: &TreeNode{
				Val: 1,
			},
		},
		{
			name: "perfect binary tree",
			input: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 1,
					},
					Right: &TreeNode{
						Val: 3,
					},
				},
				Right: &TreeNode{
					Val: 7,
					Left: &TreeNode{
						Val: 6,
					},
					Right: &TreeNode{
						Val: 9,
					},
				},
			},
			expected: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 7,
					Left: &TreeNode{
						Val: 9,
					},
					Right: &TreeNode{
						Val: 6,
					},
				},
				Right: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 3,
					},
					Right: &TreeNode{
						Val: 1,
					},
				},
			},
		},
		{
			name: "asymmetric tree",
			input: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 4,
					},
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			expected: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 3,
				},
				Right: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val: 4,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := invertTree(tt.input)
			if !areTreesEqual(result, tt.expected) {
				t.Errorf("Test case %s failed: trees are not equal", tt.name)
			}
		})
	}
}
