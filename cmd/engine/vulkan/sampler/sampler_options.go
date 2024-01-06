package sampler

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/vulkan-go/vulkan"
)

type VulkanSamplerOptions struct {
	CreateInfo vulkan.SamplerCreateInfo
}

func NewVulkanSamplerOptions(device *logical.VulkanLogicalDevice) VulkanSamplerOptions {
	options := VulkanSamplerOptions{
		CreateInfo: vulkan.SamplerCreateInfo{
			SType:                   vulkan.StructureTypeSamplerCreateInfo,
			MagFilter:               vulkan.FilterLinear,
			MinFilter:               vulkan.FilterLinear,
			AddressModeU:            vulkan.SamplerAddressModeRepeat,
			AddressModeV:            vulkan.SamplerAddressModeRepeat,
			AddressModeW:            vulkan.SamplerAddressModeRepeat,
			AnisotropyEnable:        vulkan.True,
			MaxAnisotropy:           device.Physical.Properties.Limits.MaxSamplerAnisotropy,
			BorderColor:             vulkan.BorderColorIntOpaqueBlack,
			UnnormalizedCoordinates: vulkan.False,
			CompareEnable:           vulkan.False,
			CompareOp:               vulkan.CompareOpAlways,
			MipmapMode:              vulkan.True,
			MipLodBias:              0,
			MinLod:                  0,
			MaxLod:                  0,
		},
	}

	return options
}
