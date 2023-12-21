package renderer

import (
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/device"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/swapchain"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

var RendererInstance Renderer

type Renderer struct {
	VulkanInstance  sealVulkan.VulkanInstance
	Window          window.Window
	Surface         vulkan.Surface
	ShaderContainer shader.VulkanShaderContainer
	Pipeline        pipeline.VulkanPipeline
	Swapchain       swapchain.VulkanSwapchain
}

func NewRenderer() (Renderer, error) {
	var err error
	renderer := Renderer{}

	progress.AdvanceLoading()
	renderer.VulkanInstance, err = sealVulkan.NewVulkanInstance()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	renderer.Window, err = window.NewWindow()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
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

	progress.AdvanceLoading()
	if renderer.ShaderContainer, err = shader.NewVulkanShaderContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Pipeline, err = pipeline.NewVulkanPipeline(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.ShaderContainer, &renderer.Window); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.Pipeline, &renderer.Surface, nil); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func ResizeVulkanRenderer(renderer *Renderer) error {
	var err error
	if err = swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err = device.ResizeVulkanInstanceDevices(&renderer.VulkanInstance.Devices); err != nil {
		return err
	}
	if err = pipeline.ResizeVulkanPipeline(&renderer.Pipeline); err != nil {
		return err
	}
	renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.Pipeline, &renderer.Surface, nil)
	if err != nil {
		return err
	}

	return nil
}

func AcquireNextImageRenderer(renderer *Renderer) error {
	if res := vulkan.AcquireNextImage(renderer.Pipeline.Device.Handle, renderer.Swapchain.Handle, vulkan.MaxUint64, renderer.Pipeline.Syncer.ImageAvailableSemaphores[renderer.Pipeline.CurrentFrame].Handle, nil, &renderer.Pipeline.ImageIndex); res != vulkan.Success && res != vulkan.Suboptimal {
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

func BeginRendererFrame(renderer *Renderer) error {
	if renderer.Window.Data.Extent.Width == 0 || renderer.Window.Data.Extent.Height == 0 {
		return nil
	}

	var err error
	if res := vulkan.WaitForFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.Syncer.InFlightFences[renderer.Pipeline.CurrentFrame].Handle}, vulkan.Bool32(1), vulkan.MaxUint64); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.ResetFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.Syncer.InFlightFences[renderer.Pipeline.CurrentFrame].Handle}); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	if err = AcquireNextImageRenderer(renderer); err != nil {
		logger.DefaultLogger.Error(err)
		return err
	}

	if err := buffer.BeginVulkanCommandBuffer(renderer.Pipeline.CommandBuffer); err != nil {
		return err
	}

	if err := BeginVulkanRenderPass(renderer, renderer.Pipeline.ImageIndex); err != nil {
		return err
	}
	vulkan.CmdBindVertexBuffers(renderer.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Buffer{renderer.Pipeline.Mesh.Buffer.DeviceBuffer.Handle}, []vulkan.DeviceSize{buffer.GetVulkanMeshBufferOptionsVerticesOffset(&renderer.Pipeline.Mesh.Buffer.Options)})
	vulkan.CmdBindIndexBuffer(renderer.Pipeline.CommandBuffer.Handle, renderer.Pipeline.Mesh.Buffer.DeviceBuffer.Handle, buffer.GetVulkanMeshBufferOptionsIndicesOffset(&renderer.Pipeline.Mesh.Buffer.Options), vulkan.IndexTypeUint16)
	vulkan.CmdBindPipeline(renderer.Pipeline.CommandBuffer.Handle, vulkan.PipelineBindPointGraphics, renderer.Pipeline.Handle)
	vulkan.CmdSetViewport(renderer.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Viewport{
		{
			X:        0,
			Y:        0,
			Width:    float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Width),
			Height:   float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Height),
			MinDepth: 0,
			MaxDepth: 0,
		},
	})
	vulkan.CmdSetScissor(renderer.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Rect2D{{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: renderer.Swapchain.Options.CreateInfo.ImageExtent}})

	return nil
}

func EndRendererFrame(renderer *Renderer) error {
	vulkan.CmdEndRenderPass(renderer.Pipeline.CommandBuffer.Handle)
	if err := buffer.EndVulkanCommandBuffer(renderer.Pipeline.CommandBuffer); err != nil {
		return err
	}

	if err := QueueSubmitVulkanRenderer(renderer); err != nil {
		return err
	}
	if err := QueuePresentRenderer(renderer, renderer.Pipeline.ImageIndex); err != nil {
		return err
	}
	if err := pipeline.AdvanceVulkanPipelineFrame(&renderer.Pipeline); err != nil {
		return err
	}

	return nil
}

func FreeRenderer(renderer *Renderer) error {
	if err := swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err := pipeline.FreeVulkanPipeline(&renderer.Pipeline); err != nil {
		return err
	}
	if err := shader.FreeVulkanShaderContainer(&renderer.ShaderContainer); err != nil {
		return err
	}
	vulkan.DestroySurface(renderer.VulkanInstance.Handle, renderer.Surface, nil)
	if err := window.FreeWindow(&renderer.Window); err != nil {
		return err
	}

	return sealVulkan.FreeVulkanInstance(&renderer.VulkanInstance)
}
