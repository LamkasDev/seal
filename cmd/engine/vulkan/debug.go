package vulkan

type VulkanDebugger struct{}

func NewVulkanDebugger() (VulkanDebugger, error) {
	// var err error
	debugger := VulkanDebugger{}

	return debugger, nil
}

func FreeVulkanDebugger(debugger *VulkanDebugger) error {
	return nil
}
