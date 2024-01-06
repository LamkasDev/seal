package pipeline

import (
	"github.com/EngoEngine/glm"
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline_layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/sampler"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	sealTexture "github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/transform"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipeline struct {
	Handle vulkan.Pipeline
	Device *logical.VulkanLogicalDevice
	Window *window.Window

	Camera                  transform.VulkanTransform
	Layout                  pipeline_layout.VulkanPipelineLayout
	Viewport                viewport.VulkanViewport
	TextureContainer        sealTexture.VulkanTextureContainer
	ShaderContainer         shader.VulkanShaderContainer
	DescriptorPoolContainer descriptor.VulkanDescriptorPoolContainer
	BufferContainer         buffer.VulkanBufferContainer
	MeshContainer           mesh.VulkanMeshContainer

	Sampler    sampler.VulkanSampler
	RenderPass pass.VulkanRenderPass
	Syncer     VulkanPipelineSyncer
	Commander  VulkanPipelineCommander
	Options    VulkanPipelineOptions

	ImageIndex    uint32
	CommandBuffer *buffer.VulkanCommandBuffer
	CurrentFrame  uint32
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, cwindow *window.Window) (VulkanPipeline, error) {
	var err error
	pipeline := VulkanPipeline{
		Device:   device,
		Window:   cwindow,
		Camera:   transform.VulkanTransform{Position: glm.Vec3{0, 0, 2}},
		Viewport: viewport.NewVulkanViewport(cwindow.Data.Extent),
	}

	if pipeline.Layout, err = pipeline_layout.NewVulkanPipelineLayout(device); err != nil {
		return pipeline, err
	}
	if pipeline.ShaderContainer, err = shader.NewVulkanShaderContainer(device); err != nil {
		return pipeline, err
	}
	if pipeline.DescriptorPoolContainer, err = descriptor.NewVulkanDescriptorPoolContainer(device); err != nil {
		return pipeline, err
	}
	if pipeline.BufferContainer, err = buffer.NewVulkanBufferContainer(device); err != nil {
		return pipeline, err
	}
	if pipeline.TextureContainer, err = sealTexture.NewVulkanTextureContainer(device); err != nil {
		return pipeline, err
	}
	if pipeline.MeshContainer, err = mesh.NewVulkanMeshContainer(device, &pipeline.Layout); err != nil {
		return pipeline, err
	}
	if pipeline.Sampler, err = sampler.NewVulkanSampler(device); err != nil {
		return pipeline, err
	}
	if pipeline.RenderPass, err = pass.NewVulkanRenderPass(device, device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format); err != nil {
		return pipeline, err
	}
	if pipeline.Syncer, err = NewVulkanPipelineSyncer(device); err != nil {
		return pipeline, err
	}
	if pipeline.Commander, err = NewVulkanPipelineCommander(device); err != nil {
		return pipeline, err
	}
	pipeline.Options = NewVulkanPipelineOptions(&pipeline.Layout, &pipeline.Viewport, &pipeline.RenderPass, &pipeline.ShaderContainer)

	vulkanPipelines := make([]vulkan.Pipeline, 1)
	if res := vulkan.CreateGraphicsPipelines(device.Handle, nil, 1, []vulkan.GraphicsPipelineCreateInfo{pipeline.Options.CreateInfo}, nil, vulkanPipelines); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pipeline, vulkan.Error(res)
	}
	pipeline.Handle = vulkanPipelines[0]
	logger.DefaultLogger.Debug("created new vulkan pipeline")

	if err := PushVulkanPipelineBuffers(&pipeline); err != nil {
		return pipeline, err
	}
	AdvanceVulkanPipelineFrame(&pipeline)

	return pipeline, nil
}

func PushVulkanPipelineBuffers(pipeline *VulkanPipeline) error {
	if err := buffer.BeginVulkanCommandBuffer(&pipeline.Commander.StagingBuffer); err != nil {
		return err
	}

	for _, texture := range pipeline.TextureContainer.Textures {
		ApplyVulkanTextureBarrier(pipeline, &texture, vulkan.ImageLayoutUndefined, vulkan.ImageLayoutTransferDstOptimal)

		imageCopy := vulkan.BufferImageCopy{
			BufferOffset:      0,
			BufferRowLength:   0,
			BufferImageHeight: 0,
			ImageSubresource: vulkan.ImageSubresourceLayers{
				AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			ImageOffset: vulkan.Offset3D{X: 0, Y: 0, Z: 0},
			ImageExtent: texture.Image.Options.CreateInfo.Extent,
		}
		vulkan.CmdCopyBufferToImage(pipeline.Commander.StagingBuffer.Handle, texture.Buffer.StagingBuffer.Handle, texture.Image.Handle, vulkan.ImageLayoutTransferDstOptimal, 1, []vulkan.BufferImageCopy{imageCopy})

		ApplyVulkanTextureBarrier(pipeline, &texture, vulkan.ImageLayoutTransferDstOptimal, vulkan.ImageLayoutShaderReadOnlyOptimal)
	}

	for _, mesh := range pipeline.MeshContainer.Meshes {
		bufferCopy := vulkan.BufferCopy{
			Size: mesh.Buffer.StagingBuffer.Options.CreateInfo.Size,
		}
		vulkan.CmdCopyBuffer(pipeline.Commander.StagingBuffer.Handle, mesh.Buffer.StagingBuffer.Handle, mesh.Buffer.DeviceBuffer.Handle, 1, []vulkan.BufferCopy{bufferCopy})
	}

	if err := buffer.EndVulkanCommandBuffer(&pipeline.Commander.StagingBuffer); err != nil {
		return err
	}

	submitInfo := vulkan.SubmitInfo{
		SType:              vulkan.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    []vulkan.CommandBuffer{pipeline.Commander.StagingBuffer.Handle},
	}
	if res := vulkan.QueueSubmit(pipeline.Device.Queues[uint32(pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.QueueWaitIdle(pipeline.Device.Queues[uint32(pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)]); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func ApplyVulkanTextureBarrier(pipeline *VulkanPipeline, texture *sealTexture.VulkanTexture, oldLayout vulkan.ImageLayout, newLayout vulkan.ImageLayout) {
	var srcStage, dstStage vulkan.PipelineStageFlags
	barrier := vulkan.ImageMemoryBarrier{
		SType:               vulkan.StructureTypeImageMemoryBarrier,
		OldLayout:           oldLayout,
		NewLayout:           newLayout,
		SrcQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		DstQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		Image:               texture.Image.Handle,
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	if oldLayout == vulkan.ImageLayoutUndefined && newLayout == vulkan.ImageLayoutTransferDstOptimal {
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessTransferWriteBit)
		srcStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTopOfPipeBit)
		dstStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
	} else if oldLayout == vulkan.ImageLayoutTransferDstOptimal && newLayout == vulkan.ImageLayoutShaderReadOnlyOptimal {
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessTransferWriteBit)
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessShaderReadBit)
		srcStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
		dstStage = vulkan.PipelineStageFlags(vulkan.PipelineStageFragmentShaderBit)
	} else {
		panic("what")
	}

	vulkan.CmdPipelineBarrier(pipeline.Commander.StagingBuffer.Handle, srcStage, dstStage, 0, 0, nil, 0, nil, 1, []vulkan.ImageMemoryBarrier{barrier})
}

func ResizeVulkanPipeline(pipeline *VulkanPipeline) error {
	pipeline.Viewport = viewport.NewVulkanViewport(pipeline.Window.Data.Extent)
	return nil
}

func AdvanceVulkanPipelineFrame(pipeline *VulkanPipeline) error {
	pipeline.CurrentFrame = (pipeline.CurrentFrame + 1) % commonPipeline.MaxFramesInFlight
	pipeline.CommandBuffer = &pipeline.Commander.CommandBuffers[pipeline.CurrentFrame]

	return nil
}

func FreeVulkanPipeline(pipeline *VulkanPipeline) error {
	if err := FreeVulkanPipelineCommander(&pipeline.Commander); err != nil {
		return err
	}
	if err := FreeVulkanPipelineSyncer(&pipeline.Syncer); err != nil {
		return err
	}
	if err := pass.FreeVulkanRenderPass(&pipeline.RenderPass); err != nil {
		return err
	}
	if err := sampler.FreeVulkanSampler(&pipeline.Sampler); err != nil {
		return err
	}
	if err := mesh.FreeVulkanMeshContainer(&pipeline.MeshContainer); err != nil {
		return err
	}
	if err := sealTexture.FreeVulkanTextureContainer(&pipeline.TextureContainer); err != nil {
		return err
	}
	if err := buffer.FreeVulkanBufferContainer(&pipeline.BufferContainer); err != nil {
		return err
	}
	if err := descriptor.FreeVulkanDescriptorPoolContainer(&pipeline.DescriptorPoolContainer); err != nil {
		return err
	}
	if err := shader.FreeVulkanShaderContainer(&pipeline.ShaderContainer); err != nil {
		return err
	}
	if err := pipeline_layout.FreeVulkanPipelineLayout(&pipeline.Layout); err != nil {
		return err
	}
	vulkan.DestroyPipeline(pipeline.Device.Handle, pipeline.Handle, nil)

	return nil
}
