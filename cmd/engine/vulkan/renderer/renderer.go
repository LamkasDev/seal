package renderer

import (
	"github.com/EngoEngine/glm"
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/device"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/font"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	sealPipeline "github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline_layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/sampler"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/swapchain"
	sealTexture "github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/transform"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/samber/lo"
	"github.com/vulkan-go/vulkan"
)

var RendererInstance *VulkanRenderer

type VulkanRenderer struct {
	VulkanInstance sealVulkan.VulkanInstance
	Window         window.Window
	Surface        vulkan.Surface

	Camera   transform.VulkanTransform
	Layout   pipeline_layout.VulkanPipelineLayout
	Viewport viewport.VulkanViewport

	ShaderContainer         shader.VulkanShaderContainer
	TextureContainer        sealTexture.VulkanTextureContainer
	MeshContainer           mesh.VulkanMeshContainer
	FontContainer           font.VulkanFontContainer
	DescriptorPoolContainer descriptor.VulkanDescriptorPoolContainer
	BufferContainer         sealBuffer.VulkanBufferContainer

	Sampler           sampler.VulkanSampler
	RendererSyncer    VulkanRendererSyncer
	RendererCommander VulkanRendererCommander
	ShaderPipelines   []map[string]sealPipeline.VulkanPipeline

	CurrentImageIndex uint32
	CurrentFrame      uint32

	Swapchain swapchain.VulkanSwapchain
}

func NewRenderer() (VulkanRenderer, error) {
	var err error
	renderer := VulkanRenderer{
		Camera:          transform.VulkanTransform{Position: glm.Vec3{0, 0, 2}},
		ShaderPipelines: []map[string]sealPipeline.VulkanPipeline{},
	}

	progress.AdvanceLoading()
	if renderer.VulkanInstance, err = sealVulkan.NewVulkanInstance(); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Window, err = window.NewWindow(); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	var surfaceRaw uintptr
	if surfaceRaw, err = renderer.Window.Handle.CreateWindowSurface(renderer.VulkanInstance.Handle, nil); err != nil {
		return renderer, err
	}
	renderer.Surface = vulkan.Surface(vulkan.SurfaceFromPointer(surfaceRaw))

	progress.AdvanceLoading()
	if err := sealVulkan.InitializeVulkanInstanceDevices(&renderer.VulkanInstance, &renderer.Window, &renderer.Surface); err != nil {
		return renderer, err
	}

	if renderer.Layout, err = pipeline_layout.NewVulkanPipelineLayout(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	renderer.Viewport = viewport.NewVulkanViewport(renderer.Window.Data.Extent)

	progress.AdvanceLoading()
	if renderer.TextureContainer, err = sealTexture.NewVulkanTextureContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	if renderer.ShaderContainer, err = shader.NewVulkanShaderContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	if renderer.MeshContainer, err = mesh.NewVulkanMeshContainer(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.TextureContainer); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	if renderer.FontContainer, err = font.NewVulkanFontContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	if renderer.DescriptorPoolContainer, err = descriptor.NewVulkanDescriptorPoolContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	if renderer.BufferContainer, err = sealBuffer.NewVulkanBufferContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Sampler, err = sampler.NewVulkanSampler(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	progress.AdvanceLoading()
	if renderer.RendererSyncer, err = NewVulkanRendererSyncer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}
	progress.AdvanceLoading()
	a := renderer.VulkanInstance.Devices.LogicalDevice.Physical.Capabilities.Surface.ImageFormats[renderer.VulkanInstance.Devices.LogicalDevice.Physical.Capabilities.Surface.ImageFormatIndex].Format
	if renderer.RendererCommander, err = NewVulkanRendererCommander(&renderer.VulkanInstance.Devices.LogicalDevice, a, lo.Keys(renderer.ShaderContainer.Shaders)); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	renderer.ShaderPipelines = make([]map[string]sealPipeline.VulkanPipeline, len(renderer.RendererCommander.RenderPassContainer.Passes))
	for layer, pass := range renderer.RendererCommander.RenderPassContainer.Passes {
		renderer.ShaderPipelines[layer] = map[string]sealPipeline.VulkanPipeline{}
		for key, shader := range renderer.ShaderContainer.Shaders {
			if renderer.ShaderPipelines[layer][key], err = sealPipeline.NewVulkanPipeline(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.Window, &renderer.Layout, &renderer.Viewport, pass, &shader); err != nil {
				return renderer, err
			}
		}
	}
	progress.AdvanceLoading()
	if renderer.Swapchain, err = swapchain.NewVulkanSwapchain(renderer.Layout.Device, &renderer.Window, renderer.RendererCommander.RenderPassContainer.Passes[0], &renderer.Surface, nil); err != nil {
		return renderer, err
	}

	if err := AdvanceVulkanRendererFrame(&renderer); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func ResizeVulkanRenderer(renderer *VulkanRenderer) error {
	var err error
	if err = swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err = device.ResizeVulkanInstanceDevices(&renderer.VulkanInstance.Devices); err != nil {
		return err
	}
	renderer.Viewport = viewport.NewVulkanViewport(renderer.Window.Data.Extent)
	renderer.Swapchain, err = swapchain.NewVulkanSwapchain(renderer.Layout.Device, &renderer.Window, renderer.RendererCommander.RenderPassContainer.Passes[0], &renderer.Surface, nil)
	if err != nil {
		return err
	}

	return nil
}

func AcquireNextImageRenderer(renderer *VulkanRenderer) error {
	if res := vulkan.AcquireNextImage(renderer.Layout.Device.Handle, renderer.Swapchain.Handle, vulkan.MaxUint64, renderer.RendererSyncer.ImageAvailableSemaphores[renderer.CurrentFrame].Handle, nil, &renderer.CurrentImageIndex); res != vulkan.Success && res != vulkan.Suboptimal {
		switch res {
		case vulkan.ErrorOutOfDate:
			if err := ResizeVulkanRenderer(renderer); err != nil {
				return err
			}
			return AcquireNextImageRenderer(renderer)
		default:
			return vulkan.Error(res)
		}
	}

	return nil
}

func BeginRendererFrame(renderer *VulkanRenderer) error {
	if renderer.Window.Data.Extent.Width == 0 || renderer.Window.Data.Extent.Height == 0 {
		return nil
	}

	var err error
	if res := vulkan.WaitForFences(renderer.Layout.Device.Handle, 1, []vulkan.Fence{renderer.RendererSyncer.InFlightFences[renderer.CurrentFrame].Handle}, vulkan.Bool32(1), vulkan.MaxUint64); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.ResetFences(renderer.Layout.Device.Handle, 1, []vulkan.Fence{renderer.RendererSyncer.InFlightFences[renderer.CurrentFrame].Handle}); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if err = AcquireNextImageRenderer(renderer); err != nil {
		logger.DefaultLogger.Error(err)
		return err
	}

	renderer.RendererCommander.CurrentCommandBuffer = &renderer.RendererCommander.CommandBuffers[renderer.CurrentFrame]
	if err := sealBuffer.BeginVulkanCommandBuffer(renderer.RendererCommander.CurrentCommandBuffer); err != nil {
		return err
	}
	vulkan.CmdSetViewport(renderer.RendererCommander.CurrentCommandBuffer.Handle, 0, 1, []vulkan.Viewport{
		{
			X:        0,
			Y:        0,
			Width:    float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Width),
			Height:   float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Height),
			MinDepth: 0,
			MaxDepth: 0,
		},
	})
	vulkan.CmdSetScissor(renderer.RendererCommander.CurrentCommandBuffer.Handle, 0, 1, []vulkan.Rect2D{{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: renderer.Swapchain.Options.CreateInfo.ImageExtent}})

	return nil
}

func EndRendererFrame(renderer *VulkanRenderer) error {
	for layer, pass := range renderer.RendererCommander.RenderPassContainer.Passes {
		if err := BeginVulkanRenderPass(renderer, pass, renderer.RendererCommander.CurrentCommandBuffer, renderer.CurrentImageIndex); err != nil {
			return err
		}
		for shader, buffer := range pass.AbstractCommandBuffers {
			if len(buffer.Actions) < 1 {
				continue
			}
			vulkan.CmdBindPipeline(renderer.RendererCommander.CurrentCommandBuffer.Handle, vulkan.PipelineBindPointGraphics, renderer.ShaderPipelines[layer][shader].Handle)
			RunVulkanRendererCommanderCommands(&renderer.RendererCommander, uint8(layer), shader)
			ResetVulkanRendererCommanderCommands(&renderer.RendererCommander, uint8(layer), shader)
		}
		vulkan.CmdEndRenderPass(renderer.RendererCommander.CurrentCommandBuffer.Handle)
	}
	if err := sealBuffer.EndVulkanCommandBuffer(renderer.RendererCommander.CurrentCommandBuffer); err != nil {
		return err
	}
	if err := QueueSubmitVulkanRenderer(renderer); err != nil {
		return err
	}
	if err := QueuePresentRenderer(renderer, renderer.CurrentImageIndex); err != nil {
		return err
	}
	if err := AdvanceVulkanRendererFrame(renderer); err != nil {
		return err
	}

	return nil
}

func AdvanceVulkanRendererFrame(renderer *VulkanRenderer) error {
	renderer.CurrentFrame = (renderer.CurrentFrame + 1) % commonPipeline.MaxFramesInFlight
	return nil
}

func PushVulkanRendererBuffers(renderer *VulkanRenderer) error {
	if err := sealBuffer.BeginVulkanCommandBuffer(&renderer.RendererCommander.StagingBuffer); err != nil {
		return err
	}

	for _, texture := range renderer.TextureContainer.Textures {
		PrepareVulkanTexture(renderer, texture)
	}
	for _, font := range renderer.FontContainer.Fonts {
		PrepareVulkanTexture(renderer, &font.Texture)
	}

	for _, mesh := range renderer.MeshContainer.Meshes {
		bufferCopy := vulkan.BufferCopy{
			Size: mesh.Buffer.StagingBuffer.Options.CreateInfo.Size,
		}
		vulkan.CmdCopyBuffer(renderer.RendererCommander.StagingBuffer.Handle, mesh.Buffer.StagingBuffer.Handle, mesh.Buffer.DeviceBuffer.Handle, 1, []vulkan.BufferCopy{bufferCopy})
	}

	if err := sealBuffer.EndVulkanCommandBuffer(&renderer.RendererCommander.StagingBuffer); err != nil {
		return err
	}

	submitInfo := vulkan.SubmitInfo{
		SType:              vulkan.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    []vulkan.CommandBuffer{renderer.RendererCommander.StagingBuffer.Handle},
	}
	if res := vulkan.QueueSubmit(renderer.VulkanInstance.Devices.LogicalDevice.Queues[uint32(renderer.VulkanInstance.Devices.LogicalDevice.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.QueueWaitIdle(renderer.VulkanInstance.Devices.LogicalDevice.Queues[uint32(renderer.VulkanInstance.Devices.LogicalDevice.Physical.Capabilities.Queue.GraphicsIndex)]); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func PrepareVulkanTexture(renderer *VulkanRenderer, texture *sealTexture.VulkanTexture) {
	ApplyVulkanTextureBarrier(renderer, texture, vulkan.ImageLayoutUndefined, vulkan.ImageLayoutTransferDstOptimal)
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
	vulkan.CmdCopyBufferToImage(renderer.RendererCommander.StagingBuffer.Handle, texture.Buffer.StagingBuffer.Handle, texture.Image.Handle, vulkan.ImageLayoutTransferDstOptimal, 1, []vulkan.BufferImageCopy{imageCopy})
	ApplyVulkanTextureBarrier(renderer, texture, vulkan.ImageLayoutTransferDstOptimal, vulkan.ImageLayoutShaderReadOnlyOptimal)
}

func ApplyVulkanTextureBarrier(renderer *VulkanRenderer, texture *sealTexture.VulkanTexture, oldLayout vulkan.ImageLayout, newLayout vulkan.ImageLayout) {
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

	vulkan.CmdPipelineBarrier(renderer.RendererCommander.StagingBuffer.Handle, srcStage, dstStage, 0, 0, nil, 0, nil, 1, []vulkan.ImageMemoryBarrier{barrier})
}

func FreeRenderer(renderer *VulkanRenderer) error {
	if err := swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err := FreeVulkanRendererCommander(&renderer.RendererCommander); err != nil {
		return err
	}
	if err := FreeVulkanRendererSyncer(&renderer.RendererSyncer); err != nil {
		return err
	}
	if err := sampler.FreeVulkanSampler(&renderer.Sampler); err != nil {
		return err
	}
	if err := font.FreeVulkanFontContainer(&renderer.FontContainer); err != nil {
		return err
	}
	if err := mesh.FreeVulkanMeshContainer(&renderer.MeshContainer); err != nil {
		return err
	}
	if err := sealBuffer.FreeVulkanBufferContainer(&renderer.BufferContainer); err != nil {
		return err
	}
	if err := descriptor.FreeVulkanDescriptorPoolContainer(&renderer.DescriptorPoolContainer); err != nil {
		return err
	}
	if err := shader.FreeVulkanShaderContainer(&renderer.ShaderContainer); err != nil {
		return err
	}
	if err := pipeline_layout.FreeVulkanPipelineLayout(&renderer.Layout); err != nil {
		return err
	}
	for _, pipelines := range renderer.ShaderPipelines {
		for _, pipeline := range pipelines {
			if err := sealPipeline.FreeVulkanPipeline(&pipeline); err != nil {
				return err
			}
		}
	}
	if err := sealTexture.FreeVulkanTextureContainer(&renderer.TextureContainer); err != nil {
		return err
	}
	vulkan.DestroySurface(renderer.VulkanInstance.Handle, renderer.Surface, nil)
	if err := window.FreeWindow(&renderer.Window); err != nil {
		return err
	}

	return sealVulkan.FreeVulkanInstance(&renderer.VulkanInstance)
}
