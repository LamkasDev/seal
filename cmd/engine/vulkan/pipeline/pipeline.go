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
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
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
	ShaderContainer         shader.VulkanShaderContainer
	DescriptorPoolContainer descriptor.VulkanDescriptorPoolContainer
	BufferContainer         buffer.VulkanBufferContainer
	MeshContainer           mesh.VulkanMeshContainer

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
	if pipeline.MeshContainer, err = mesh.NewVulkanMeshContainer(device, &pipeline.Layout); err != nil {
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
	if err := mesh.FreeVulkanMeshContainer(&pipeline.MeshContainer); err != nil {
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
