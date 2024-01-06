package entity

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/time"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	sealMesh "github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/transform_buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	sealUniform "github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/vulkan-go/vulkan"
)

type EntityComponentMeshData struct {
	Mesh           *sealMesh.VulkanMesh
	Buffer         transform_buffer.VulkanTransformBuffer
	DescriptorPool descriptor.VulkanDescriptorPool
	DescriptorSets []descriptor.VulkanDescriptorSet
}

func NewEntityComponentMesh(entity *Entity, mesh *sealMesh.VulkanMesh) (EntityComponent, error) {
	var err error
	component := EntityComponent{
		Entity: entity,
		Data: EntityComponentMeshData{
			Mesh:           mesh,
			DescriptorSets: make([]descriptor.VulkanDescriptorSet, commonPipeline.MaxFramesInFlight),
		},
		Render: RenderEntityComponentMesh,
		Update: func(component *EntityComponent) error { return nil },
	}
	data := component.Data.(EntityComponentMeshData)
	if data.Buffer, err = transform_buffer.NewVulkanTransformBuffer(&renderer.RendererInstance.Pipeline, transform_buffer.NewVulkanTransformBufferOptions(uniform.NewVulkanUniform(renderer.RendererInstance.Window.Data.Extent, renderer.RendererInstance.Pipeline.Camera.Position, component.Entity.Transform.Position, component.Entity.Transform.Rotation))); err != nil {
		return component, err
	}
	if data.DescriptorPool, err = descriptor.CreateVulkanDescriptorPoolWithContainer(&renderer.RendererInstance.Pipeline.DescriptorPoolContainer); err != nil {
		return component, err
	}
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if data.DescriptorSets[i], err = descriptor.NewVulkanDescriptorSet(mesh.Device, &data.DescriptorPool, &renderer.RendererInstance.Pipeline.Layout.DescriptorSetLayout); err != nil {
			return component, err
		}
	}
	component.Data = data
	if err := UpdateEntityComponentMesh(&component, 0, commonPipeline.MaxFramesInFlight); err != nil {
		return component, err
	}

	return component, nil
}

func UpdateEntityComponentMesh(component *EntityComponent, startFrame uint32, endFrame uint32) error {
	data := component.Data.(EntityComponentMeshData)
	if err := transform_buffer.CopyVulkanTransformBuffer(&data.Buffer); err != nil {
		return err
	}

	writeDescriptorSets := []vulkan.WriteDescriptorSet{}
	for i := startFrame; i < endFrame; i++ {
		uniformDescriptorSet := vulkan.WriteDescriptorSet{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          data.DescriptorSets[i].Handle,
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			PBufferInfo: []vulkan.DescriptorBufferInfo{
				{
					Buffer: data.Buffer.StagingBuffer.Handle,
					Offset: transform_buffer.GetVulkanTransformBufferOptionsUniformsOffset(&data.Buffer.Options) + vulkan.DeviceSize(i)*vulkan.DeviceSize(uniform.VulkanUniformSize),
					Range:  vulkan.DeviceSize(uniform.VulkanUniformSize),
				},
			},
		}
		writeDescriptorSets = append(writeDescriptorSets, uniformDescriptorSet)

		textureDescriptorSet := vulkan.WriteDescriptorSet{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          data.DescriptorSets[i].Handle,
			DstBinding:      1,
			DstArrayElement: 0,
			DescriptorType:  vulkan.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 1,
			PImageInfo: []vulkan.DescriptorImageInfo{
				{
					ImageLayout: vulkan.ImageLayoutShaderReadOnlyOptimal,
					ImageView:   renderer.RendererInstance.Pipeline.TextureContainer.Textures[texture.TEXTURE_BASIC].ImageView.Handle,
					Sampler:     renderer.RendererInstance.Pipeline.Sampler.Handle,
				},
			},
		}
		writeDescriptorSets = append(writeDescriptorSets, textureDescriptorSet)
	}
	vulkan.UpdateDescriptorSets(data.Mesh.Device.Handle, uint32(len(writeDescriptorSets)), writeDescriptorSets, 0, nil)

	return nil
}

func RenderEntityComponentMesh(component *EntityComponent) error {
	data := component.Data.(EntityComponentMeshData)
	data.Buffer.Options.Uniforms[renderer.RendererInstance.Pipeline.CurrentFrame] = sealUniform.NewVulkanUniform(renderer.RendererInstance.Window.Data.Extent, renderer.RendererInstance.Pipeline.Camera.Position, component.Entity.Transform.Position, component.Entity.Transform.Rotation)
	component.Data = data

	component.Entity.Transform.Rotation += time.DeltaTime * 100
	if err := UpdateEntityComponentMesh(component, renderer.RendererInstance.Pipeline.CurrentFrame, renderer.RendererInstance.Pipeline.CurrentFrame); err != nil {
		return err
	}

	vulkan.CmdBindDescriptorSets(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, vulkan.PipelineBindPointGraphics, renderer.RendererInstance.Pipeline.Layout.Handle, 0, 1, []vulkan.DescriptorSet{data.DescriptorSets[renderer.RendererInstance.Pipeline.CurrentFrame].Handle}, 0, nil)
	vulkan.CmdBindVertexBuffers(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Buffer{data.Mesh.Buffer.DeviceBuffer.Handle}, []vulkan.DeviceSize{buffer.GetVulkanMeshBufferOptionsVerticesOffset(&data.Mesh.Buffer.Options)})
	vulkan.CmdBindIndexBuffer(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, data.Mesh.Buffer.DeviceBuffer.Handle, buffer.GetVulkanMeshBufferOptionsIndicesOffset(&data.Mesh.Buffer.Options), vulkan.IndexTypeUint16)
	vulkan.CmdDrawIndexed(renderer.RendererInstance.Pipeline.CommandBuffer.Handle, uint32(len(data.Mesh.Buffer.Options.Indices)), 1, 0, 0, 0)

	return nil
}
