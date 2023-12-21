package component

import (
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	sealMesh "github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/vulkan-go/vulkan"
)

type EntityComponentMeshData struct {
	Mesh sealMesh.VulkanMesh
}

func NewEntityComponent(mesh sealMesh.VulkanMesh) (EntityComponent, error) {
	component := EntityComponent{
		Data: EntityComponentMeshData{
			Mesh: mesh,
		},
		Render: RenderEntityComponentMesh,
	}

	return component, nil
}

func RenderEntityComponentMesh(component *EntityComponent) error {
	mesh := (EntityComponentMeshData)(component.Data).Mesh
	vulkan.CmdBindVertexBuffers(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Buffer{mesh.Buffer.DeviceBuffer.Handle}, []vulkan.DeviceSize{buffer.GetVulkanMeshBufferOptionsVerticesOffset(&mesh.Buffer.Options)})
	vulkan.CmdBindIndexBuffer(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, mesh.Buffer.DeviceBuffer.Handle, buffer.GetVulkanMeshBufferOptionsIndicesOffset(&mesh.Buffer.Options), vulkan.IndexTypeUint16)
	vulkan.CmdDrawIndexed(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, uint32(len(mesh.Buffer.Options.Indices)), 1, 0, 0, 0)

	return nil
}
