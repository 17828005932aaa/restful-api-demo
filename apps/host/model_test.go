package host_test

import (
	"restful-api-demo/apps/host"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostUpdate(t *testing.T) {
	should := assert.New(t)

	h := host.NewDefaultHost()
	patch := host.NewDefaultHost()
	patch.Resource.Name = "patch01"

	err := h.Patch(patch.Resource, patch.Describe)
	if should.NoError(err) {
		should.Equal(patch.Resource.Name, h.Resource.Name)
	}
}
