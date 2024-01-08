package mesh

import (
	"fmt"

	"github.com/EngoEngine/glm"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	sealModel "github.com/LamkasDev/seal/cmd/engine/vulkan/model"
	sealShader "github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/LamkasDev/seal/cmd/logger"
)

type VulkanMeshTemplate struct {
	Id      string
	Mesh    string
	Shader  string
	Texture string
}

const MESH_BASIC = "basic"
const MESH_UI = "ui"

var DefaultMeshes = []VulkanMeshTemplate{
	{
		Id:      MESH_BASIC,
		Mesh:    MESH_BASIC,
		Shader:  sealShader.SHADER_BASIC,
		Texture: texture.TEXTURE_BASIC,
	},
	{
		Id:      MESH_UI,
		Mesh:    MESH_BASIC,
		Shader:  sealShader.SHADER_BASIC,
		Texture: texture.TEXTURE_SKY,
	},
}

type VulkanMeshContainer struct {
	Device *logical.VulkanLogicalDevice
	Meshes map[string]*VulkanMesh
}

func NewVulkanMeshContainer(device *logical.VulkanLogicalDevice, textureContainer *texture.VulkanTextureContainer) (VulkanMeshContainer, error) {
	container := VulkanMeshContainer{
		Device: device,
		Meshes: map[string]*VulkanMesh{},
	}
	for _, template := range DefaultMeshes {
		model := sealModel.ConvertModel(sealModel.NewModel(fmt.Sprintf("../../resources/models/%s.obj", template.Mesh)))
		if _, err := CreateVulkanMeshWithContainer(&container, template.Id, template.Shader, textureContainer.Textures[template.Texture], buffer.NewVulkanMeshBufferOptions(model.Vertices, model.Indices, uniform.NewVulkanUniform3D(device.Physical.Window.Data.Extent, glm.Vec3{2, 2, 2}, glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, 0}))); err != nil {
			return container, err
		}
	}
	logger.DefaultLogger.Debug("created new vulkan mesh container")

	return container, nil
}

func CreateVulkanMeshWithContainer(container *VulkanMeshContainer, id string, shader string, texture *texture.VulkanTexture, options buffer.VulkanMeshBufferOptions) (VulkanMesh, error) {
	mesh, err := NewVulkanMesh(container.Device, id, shader, texture, options)
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
