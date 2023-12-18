package swapchain

import (
	sealFramebuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/framebuffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchain struct {
	Handle       vulkan.Swapchain
	Pipeline     *pipeline.VulkanPipeline
	Options      VulkanSwapchainOptions
	Images       []vulkan.Image
	Framebuffers []sealFramebuffer.VulkanFramebuffer
}

func NewVulkanSwapchain(pipeline *pipeline.VulkanPipeline, surface *vulkan.Surface, old *VulkanSwapchain) (VulkanSwapchain, error) {
	swapchain := VulkanSwapchain{
		Pipeline:     pipeline,
		Options:      NewVulkanSwapchainOptions(pipeline, surface, old),
		Framebuffers: []sealFramebuffer.VulkanFramebuffer{},
	}

	var vulkanSwapchain vulkan.Swapchain
	if res := vulkan.CreateSwapchain(pipeline.Device.Handle, &swapchain.Options.CreateInfo, nil, &vulkanSwapchain); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	swapchain.Handle = vulkanSwapchain
	logger.DefaultLogger.Debug("created new vulkan swapchain")

	var imagesCount uint32
	if res := vulkan.GetSwapchainImages(pipeline.Device.Handle, swapchain.Handle, &imagesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	swapchain.Images = make([]vulkan.Image, imagesCount)
	swapchain.Framebuffers = make([]sealFramebuffer.VulkanFramebuffer, imagesCount)
	if res := vulkan.GetSwapchainImages(pipeline.Device.Handle, swapchain.Handle, &imagesCount, swapchain.Images); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}

	for i := 0; i < len(swapchain.Images); i++ {
		imageview, err := image.NewVulkanImageView(pipeline.Device, &swapchain.Images[i])
		if err != nil {
			return swapchain, err
		}
		framebuffer, err := sealFramebuffer.NewVulkanFramebuffer(pipeline.Device, &swapchain.Pipeline.RenderPass, &imageview, swapchain.Options.CreateInfo.ImageExtent)
		if err != nil {
			return swapchain, err
		}
		swapchain.Framebuffers[i] = framebuffer
	}

	return swapchain, nil
}

func FreeVulkanSwapchain(swapchain *VulkanSwapchain) error {
	vulkan.DeviceWaitIdle(swapchain.Pipeline.Device.Handle)
	for _, framebuffer := range swapchain.Framebuffers {
		if err := sealFramebuffer.FreeVulkanFramebuffer(&framebuffer); err != nil {
			return err
		}
	}
	vulkan.DestroySwapchain(swapchain.Pipeline.Device.Handle, swapchain.Handle, nil)

	return nil
}
