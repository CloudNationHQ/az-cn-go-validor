package tests

import (
	"testing"

	"github.com/cloudnationhq/az-cn-go-validor"
)

// TestApplyNoError tests a single Terraform module
func TestApplyNoError(t *testing.T) {
	validor.TestApplyNoError(t)
}

// TestApplyAllParallel tests all Terraform modules in parallel
func TestApplyAllParallel(t *testing.T) {
	validor.TestApplyAllParallel(t)
}
