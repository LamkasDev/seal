package viewport

import "github.com/vulkan-go/vulkan"

type VulkanViewport struct {
	Viewport vulkan.Viewport
	Scissor  vulkan.Rect2D
}

func NewVulkanViewport(extent vulkan.Extent2D) VulkanViewport {
	viewport := VulkanViewport{
		Viewport: vulkan.Viewport{
			X:        0,
			Y:        0,
			Width:    float32(extent.Width),
			Height:   float32(extent.Height),
			MinDepth: 0,
			MaxDepth: 1,
		},
		Scissor: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: extent,
		},
	}

	return viewport
}
