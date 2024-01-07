package entity

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/time"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	sealMesh "github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/renderer"
	sealRenderer "github.com/LamkasDev/seal/cmd/engine/vulkan/renderer"
	sealShader "github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
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
	if data.Buffer, err = transform_buffer.NewVulkanTransformBuffer(sealRenderer.RendererInstance, transform_buffer.NewVulkanTransformBufferOptions(uniform.NewVulkanUniform3D(sealRenderer.RendererInstance.Window.Data.Extent, sealRenderer.RendererInstance.Camera.Position, component.Entity.Transform.Position, component.Entity.Transform.Rotation))); err != nil {
		return component, err
	}
	if data.DescriptorPool, err = descriptor.CreateVulkanDescriptorPoolWithContainer(&sealRenderer.RendererInstance.DescriptorPoolContainer); err != nil {
		return component, err
	}
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if data.DescriptorSets[i], err = descriptor.NewVulkanDescriptorSet(mesh.Device, &data.DescriptorPool, &sealRenderer.RendererInstance.Layout.DescriptorSetLayout); err != nil {
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
					ImageView:   renderer.RendererInstance.TextureContainer.Textures[texture.TEXTURE_BASIC].ImageView.Handle,
					Sampler:     renderer.RendererInstance.Sampler.Handle,
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
	data.Buffer.Options.Uniforms[sealRenderer.RendererInstance.CurrentFrame] = sealUniform.NewVulkanUniform3D(sealRenderer.RendererInstance.Window.Data.Extent, sealRenderer.RendererInstance.Camera.Position, component.Entity.Transform.Position, component.Entity.Transform.Rotation)
	component.Data = data

	if data.Mesh.Shader == sealShader.SHADER_BASIC {
		component.Entity.Transform.Rotation += time.DeltaTime * 100
		if err := UpdateEntityComponentMesh(component, sealRenderer.RendererInstance.CurrentFrame, sealRenderer.RendererInstance.CurrentFrame); err != nil {
			return err
		}
	}

	renderer.RecordVulkanRendererCommanderCommands(&renderer.RendererInstance.RendererCommander, data.Mesh.Shader, func() {
		vulkan.CmdBindDescriptorSets(sealRenderer.RendererInstance.RendererCommander.CurrentCommandBuffer.Handle, vulkan.PipelineBindPointGraphics, sealRenderer.RendererInstance.Layout.Handle, 0, 1, []vulkan.DescriptorSet{data.DescriptorSets[sealRenderer.RendererInstance.CurrentFrame].Handle}, 0, nil)
		vulkan.CmdBindVertexBuffers(sealRenderer.RendererInstance.RendererCommander.CurrentCommandBuffer.Handle, 0, 1, []vulkan.Buffer{data.Mesh.Buffer.DeviceBuffer.Handle}, []vulkan.DeviceSize{buffer.GetVulkanMeshBufferOptionsVerticesOffset(&data.Mesh.Buffer.Options)})
		vulkan.CmdBindIndexBuffer(sealRenderer.RendererInstance.RendererCommander.CurrentCommandBuffer.Handle, data.Mesh.Buffer.DeviceBuffer.Handle, buffer.GetVulkanMeshBufferOptionsIndicesOffset(&data.Mesh.Buffer.Options), vulkan.IndexTypeUint16)
		vulkan.CmdDrawIndexed(sealRenderer.RendererInstance.RendererCommander.CurrentCommandBuffer.Handle, uint32(len(data.Mesh.Buffer.Options.Indices)), 1, 0, 0, 0)
	})

	return nil
}
