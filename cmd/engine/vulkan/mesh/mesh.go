package mesh

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
)

type VulkanMesh struct {
	Id      string
	Shader  string
	Texture string
	Buffer  buffer.VulkanMeshBuffer
	Device  *logical.VulkanLogicalDevice
}

func NewVulkanMesh(device *logical.VulkanLogicalDevice, id string, shader string, texture string, options buffer.VulkanMeshBufferOptions) (VulkanMesh, error) {
	var err error
	mesh := VulkanMesh{
		Id:      id,
		Device:  device,
		Shader:  shader,
		Texture: texture,
	}

	if mesh.Buffer, err = buffer.NewVulkanMeshBuffer(device, options); err != nil {
		return mesh, err
	}

	return mesh, nil
}

func FreeVulkanMesh(mesh *VulkanMesh) error {
	return buffer.FreeVulkanMeshBuffer(&mesh.Buffer)
}
