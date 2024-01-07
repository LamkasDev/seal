package mesh

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
)

type VulkanMesh struct {
	Buffer buffer.VulkanMeshBuffer
	Device *logical.VulkanLogicalDevice
	Shader string
}

func NewVulkanMesh(device *logical.VulkanLogicalDevice, shader string, options buffer.VulkanMeshBufferOptions) (VulkanMesh, error) {
	var err error
	mesh := VulkanMesh{
		Device: device,
		Shader: shader,
	}

	if mesh.Buffer, err = buffer.NewVulkanMeshBuffer(device, options); err != nil {
		return mesh, err
	}

	return mesh, nil
}

func FreeVulkanMesh(mesh *VulkanMesh) error {
	return buffer.FreeVulkanMeshBuffer(&mesh.Buffer)
}
