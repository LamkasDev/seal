package buffer

type VulkanAbstractBufferAction func()

type VulkanAbstractBuffer struct {
	Actions []VulkanAbstractBufferAction
}

func NewVulkanAbstractBuffer() VulkanAbstractBuffer {
	return VulkanAbstractBuffer{
		Actions: []VulkanAbstractBufferAction{},
	}
}
