package mesh

import (
	"github.com/EngoEngine/glm"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	sealShader "github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/LamkasDev/seal/cmd/logger"
)

type VulkanMeshTemplate struct {
	Id     string
	Shader string
}

const MESH_BASIC = "basic"
const MESH_UI = "ui"

var DefaultMeshes = []VulkanMeshTemplate{
	VulkanMeshTemplate{
		Id:     MESH_BASIC,
		Shader: sealShader.SHADER_BASIC,
	},
	VulkanMeshTemplate{
		Id:     MESH_UI,
		Shader: sealShader.SHADER_UI,
	},
}

type VulkanMeshContainer struct {
	Device *logical.VulkanLogicalDevice
	Meshes map[string]*VulkanMesh
}

func NewVulkanMeshContainer(device *logical.VulkanLogicalDevice) (VulkanMeshContainer, error) {
	container := VulkanMeshContainer{
		Device: device,
		Meshes: map[string]*VulkanMesh{},
	}
	for _, template := range DefaultMeshes {
		if _, err := CreateVulkanMeshWithContainer(&container, template.Id, template.Shader, buffer.NewVulkanMeshBufferOptions(vertex.DefaultVertices, vertex.DefaultVerticesIndex, uniform.NewVulkanUniform3D(device.Physical.Window.Data.Extent, glm.Vec3{2, 2, 2}, glm.Vec3{0, 0, 0}, 0))); err != nil {
			return container, err
		}
	}
	logger.DefaultLogger.Debug("created new vulkan mesh container")

	return container, nil
}

func CreateVulkanMeshWithContainer(container *VulkanMeshContainer, id string, shader string, options buffer.VulkanMeshBufferOptions) (VulkanMesh, error) {
	mesh, err := NewVulkanMesh(container.Device, shader, options)
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
