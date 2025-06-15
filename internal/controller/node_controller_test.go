package controller_test

import (
	"context"
	"testing"

	"github.com/flohansen/nova-cloud/internal/controller"
	v1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"github.com/stretchr/testify/assert"
)

func TestNodeController_NewNodeController(t *testing.T) {
	// assign
	// act
	ctrl := controller.NewNodeController()

	// assert
	assert.NotNil(t, ctrl)
}

func TestNodeController_GetResources(t *testing.T) {
	// assign
	ctx := context.Background()
	ctrl := controller.NewNodeController()

	// act
	res, err := ctrl.GetResources(ctx, &v1.GetResourcesRequest{})

	// assert
	assert.NoError(t, err)
	assert.Greater(t, res.CpuCores, int32(0))
	assert.NotEqual(t, res.CpuArchitecture, v1.CpuArch_CPU_ARCH_UNSPECIFIED)
}
