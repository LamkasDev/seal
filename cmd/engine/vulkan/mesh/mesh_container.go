package mesh

import (
	"github.com/EngoEngine/glm"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline_layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/LamkasDev/seal/cmd/logger"
)

const MESH_BASIC = "basic"

type VulkanMeshContainer struct {
	Device *logical.VulkanLogicalDevice
	Layout *pipeline_layout.VulkanPipelineLayout
	Meshes map[string]*VulkanMesh
}

func NewVulkanMeshContainer(device *logical.VulkanLogicalDevice, layout *pipeline_layout.VulkanPipelineLayout) (VulkanMeshContainer, error) {
	container := VulkanMeshContainer{
		Device: device,
		Layout: layout,
		Meshes: map[string]*VulkanMesh{},
	}
	if _, err := CreateVulkanMeshWithContainer(&container, layout, MESH_BASIC, buffer.NewVulkanMeshBufferOptions(vertex.DefaultVertices, vertex.DefaultVerticesIndex, uniform.NewVulkanUniform(device.Physical.Window.Data.Extent, glm.Vec3{2, 2, 2}, glm.Vec3{0, 0, 0}, 0))); err != nil {
		return container, err
	}
	logger.DefaultLogger.Debug("created new vulkan mesh container")

	return container, nil
}

func CreateVulkanMeshWithContainer(container *VulkanMeshContainer, layout *pipeline_layout.VulkanPipelineLayout, id string, options buffer.VulkanMeshBufferOptions) (VulkanMesh, error) {
	mesh, err := NewVulkanMesh(container.Device, layout, options)
	if err != nil {
		return mesh, err
	}
	container.Meshes[id] = &mesh

	return mesh, nil
}

func FreeVulkanMeshContainer(container *VulkanMeshContainer) error {
	for _, mesh := range container.Meshes {
		if err := FreeVulkanMesh(mesh); err != nil {
			return err
		}
	}

	return nil
}
