package pipeline

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
)

type VulkanPipelineCommander struct {
	Device         *logical.VulkanLogicalDevice
	Pool           command.VulkanCommandPool
	StagingBuffer  buffer.VulkanCommandBuffer
	CommandBuffers []buffer.VulkanCommandBuffer
}

func NewVulkanPipelineCommander(device *logical.VulkanLogicalDevice) (VulkanPipelineCommander, error) {
	var err error
	commander := VulkanPipelineCommander{
		Device:         device,
		CommandBuffers: make([]buffer.VulkanCommandBuffer, commonPipeline.MaxFramesInFlight),
	}

	if commander.Pool, err = command.NewVulkanCommandPool(device, uint32(device.Physical.Capabilities.Queue.GraphicsIndex)); err != nil {
		return commander, err
	}
	if commander.StagingBuffer, err = buffer.NewVulkanCommandBuffer(device, &commander.Pool); err != nil {
		return commander, err
	}
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if commander.CommandBuffers[i], err = buffer.NewVulkanCommandBuffer(device, &commander.Pool); err != nil {
			return commander, err
		}
	}

	return commander, nil
}

func FreeVulkanPipelineCommander(commander *VulkanPipelineCommander) error {
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if err := buffer.FreeVulkanCommandBuffer(&commander.CommandBuffers[i]); err != nil {
			return err
		}
	}
	if err := buffer.FreeVulkanCommandBuffer(&commander.StagingBuffer); err != nil {
		return err
	}

	return command.FreeVulkanCommandPool(&commander.Pool)
}
