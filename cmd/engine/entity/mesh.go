package entity

import (
	"math"

	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/time"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	sealMesh "github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	sealUniform "github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/vulkan-go/vulkan"
)

type EntityComponentMeshData struct {
	Mesh           sealMesh.VulkanMesh
	DescriptorPool descriptor.VulkanDescriptorPool
	DescriptorSets []descriptor.VulkanDescriptorSet
}

func NewEntityComponentMesh(entity *Entity, mesh sealMesh.VulkanMesh) (EntityComponent, error) {
	var err error
	component := EntityComponent{
		Entity: entity,
		Data: EntityComponentMeshData{
			Mesh:           mesh,
			DescriptorSets: make([]descriptor.VulkanDescriptorSet, commonPipeline.MaxFramesInFlight),
		},
		Render: RenderEntityComponentMesh,
	}
	if component.Data.DescriptorPool, err = descriptor.NewVulkanDescriptorPool(mesh.Device); err != nil {
		return component, err
	}
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if component.Data.DescriptorSets[i], err = descriptor.NewVulkanDescriptorSet(mesh.Device, &component.Data.DescriptorPool, &renderer.RendererInstance.Pipeline.Layout.DescriptorSetLayout); err != nil {
			return component, err
		}
	}
	if err := UpdateEntityComponentMesh(&component, 0, commonPipeline.MaxFramesInFlight); err != nil {
		return component, err
	}

	return component, nil
}

func UpdateEntityComponentMesh(component *EntityComponent, startFrame uint32, endFrame uint32) error {
	if err := buffer.CopyVulkanMeshUniformBuffer(&component.Data.Mesh.Buffer); err != nil {
		return err
	}

	writeDescriptorSets := []vulkan.WriteDescriptorSet{}
	for i := startFrame; i < endFrame; i++ {
		writeDescriptorSets = append(writeDescriptorSets, vulkan.WriteDescriptorSet{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          component.Data.DescriptorSets[i].Handle,
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			PBufferInfo: []vulkan.DescriptorBufferInfo{
				{
					Buffer: component.Data.Mesh.Buffer.StagingBuffer.Handle,
					Offset: buffer.GetVulkanMeshBufferOptionsUniformsOffset(&component.Data.Mesh.Buffer.Options) + vulkan.DeviceSize(i)*vulkan.DeviceSize(uniform.VulkanUniformSize),
					Range:  vulkan.DeviceSize(uniform.VulkanUniformSize),
				},
			},
		})
	}
	vulkan.UpdateDescriptorSets(component.Data.Mesh.Device.Handle, uint32(len(writeDescriptorSets)), writeDescriptorSets, 0, nil)

	return nil
}

func RenderEntityComponentMesh(component *EntityComponent) error {
	data := (EntityComponentMeshData)(component.Data)

	component.Entity.Position[0] = float32(math.Mod(float64(component.Entity.Position[0]+time.DeltaTime*100), 1))
	component.Entity.Rotation += time.DeltaTime * 100
	data.Mesh.Buffer.Options.Uniforms[renderer.RendererInstance.Pipeline.CurrentFrame] = sealUniform.NewVulkanUniform(renderer.RendererInstance.Window.Data.Extent, component.Entity.Position, component.Entity.Rotation)
	if err := UpdateEntityComponentMesh(component, renderer.RendererInstance.Pipeline.CurrentFrame, renderer.RendererInstance.Pipeline.CurrentFrame); err != nil {
		return err
	}

	vulkan.CmdBindDescriptorSets(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, vulkan.PipelineBindPointGraphics, renderer.RendererInstance.Pipeline.Layout.Handle, 0, 1, []vulkan.DescriptorSet{data.DescriptorSets[renderer.RendererInstance.Pipeline.CurrentFrame].Handle}, 0, nil)
	vulkan.CmdBindVertexBuffers(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Buffer{data.Mesh.Buffer.DeviceBuffer.Handle}, []vulkan.DeviceSize{buffer.GetVulkanMeshBufferOptionsVerticesOffset(&data.Mesh.Buffer.Options)})
	vulkan.CmdBindIndexBuffer(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, data.Mesh.Buffer.DeviceBuffer.Handle, buffer.GetVulkanMeshBufferOptionsIndicesOffset(&data.Mesh.Buffer.Options), vulkan.IndexTypeUint16)
	vulkan.CmdDrawIndexed(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, uint32(len(data.Mesh.Buffer.Options.Indices)), 1, 0, 0, 0)

	return nil
}
