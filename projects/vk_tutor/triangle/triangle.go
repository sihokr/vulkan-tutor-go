package main

import (
	"fmt"
	"os"
	"runtime"

	"example.com/vk_tutor/sdl2"
	"example.com/vk_tutor/vulkan"
)

func main() {

	runtime.LockOSThread()

	var app HelloTriangleApplication
	app.Run()

	fmt.Println("Quit")
}

const (
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 600
)

const enableValidationLayers = true

var validationLayers = []string{
	// "VK_LAYER_KHRONOS_validation",
	"VK_LAYER_LUNARG_standard_validation",
}

var deviceExtensions = []string{
	vulkan.VK_KHR_SWAPCHAIN_EXTENSION_NAME,
}

type HelloTriangleApplication struct {
	Window         *sdl2.SDL_Window
	Instance       vulkan.VkInstance
	Surface        vulkan.VkSurfaceKHR
	PhysicalDevice vulkan.VkPhysicalDevice
	Device         vulkan.VkDevice
	GraphicsQueue  vulkan.VkQueue
	PresentQueue   vulkan.VkQueue
	SwapChain      vulkan.VkSwapchainKHR
	Images         []vulkan.VkImage
	ImageFormat    vulkan.VkFormat
	Extent         vulkan.VkExtent2D
	ImageViews     []vulkan.VkImageView
}

type QueueFamilyIndices struct {
	GraphicsFamily int
	PresentFamily  int
}

type SwapChainSupportDetails struct {
	Capabilities vulkan.VkSurfaceCapabilitiesKHR
	Formats      []vulkan.VkSurfaceFormatKHR
	PresentModes []vulkan.VkPresentModeKHR
}

func (o *HelloTriangleApplication) Run() {
	o.initWindow()
	o.initVulkan()
	o.mainLoop()
	o.cleanup()
}

func (o *HelloTriangleApplication) initWindow() {

	sdl2.SDL_SetMainReady()

	if 0 != sdl2.SDL_Init(sdl2.SDL_INIT_EVERYTHING) {
		// var msg = sdl2.SDL_GetError()
		// TODO
		fmt.Println("SDL_Init failed")
	}

	var window = sdl2.SDL_CreateWindow("Triangle", sdl2.SDL_WINDOWPOS_UNDEFINED, sdl2.SDL_WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT,
		sdl2.SDL_WINDOW_SHOWN|sdl2.SDL_WINDOW_VULKAN)
	if nil == window {
		// var msg = sdl2.SDL_GetError()
		// TODO
		fmt.Println("SDL_CreateWindow failed")
	}

	o.Window = window

}

func (o *HelloTriangleApplication) initVulkan() {

	o.createInstance()
	o.createSurface()

	var queue_families QueueFamilyIndices
	var swap_chain_support SwapChainSupportDetails
	o.pickPhysicalDevice(&queue_families, &swap_chain_support)

	o.createLogicalDevice(&queue_families)
	o.createSwapChain(&queue_families, &swap_chain_support)
	o.createImageViews()
	// o.createGraphicsPipeline()

}

func (o *HelloTriangleApplication) mainLoop() {

	var e sdl2.SDL_Event

loop:
	for {

		// Handle events
		if 0 != sdl2.SDL_PollEvent(&e) {
			switch e.Type() {
			case sdl2.SDL_QUIT:
				break loop
			} // switch
		}

		// More work here..

	} // for
}

func (o *HelloTriangleApplication) cleanup() {

	for _, v := range o.ImageViews {
		vulkan.VkDestroyImageView(o.Device, v, nil)
	}

	vulkan.VkDestroySwapchainKHR(o.Device, o.SwapChain, nil)
	vulkan.VkDestroyDevice(o.Device, nil)
	vulkan.VkDestroySurfaceKHR(o.Instance, o.Surface, nil)
	vulkan.VkDestroyInstance(o.Instance, nil)
	sdl2.SDL_DestroyWindow(o.Window)
	sdl2.SDL_Quit()
}

func (o *HelloTriangleApplication) createInstance() {

	var app_info vulkan.VkApplicationInfo
	var app_name = "Hello Triangle"
	var engine_name = "No Engine"
	app_info.PApplicationName = &app_name
	app_info.ApplicationVersion = vulkan.VK_MAKE_VERSION(1, 0, 0)
	app_info.PEngineName = &engine_name
	app_info.EngineVersion = vulkan.VK_MAKE_VERSION(1, 0, 0)
	app_info.ApiVersion = vulkan.VK_API_VERSION_1_0

	var create_info vulkan.VkInstanceCreateInfo
	create_info.PApplicationInfo = &app_info

	var ext_names []string
	{
		// Get avaible extension count
		var cnt uint32
		vulkan.VkEnumerateInstanceExtensionProperties(nil, &cnt, nil)

		if 0 == cnt {
			// TODO
			fmt.Println("VkEnumerateInstanceExtensionProperties returns 0")
		}

		ext_names = make([]string, cnt)

		// Get required extensions by SDL
		if cnt1 := uint(cnt); sdl2.SDL_Vulkan_GetInstanceExtensions(o.Window, &cnt1, ext_names) {
			cnt = uint32(cnt1)
			ext_names = ext_names[:cnt]
		} else {
			// TODO
			fmt.Println("SDL_Vulkan_GetInstanceExtensions failed")
		}

		fmt.Printf("Required instance extension count is %v\n", cnt)
		for i := 0; i < int(cnt); i++ {
			fmt.Printf("\t#%v: %v\n", i, ext_names[i])
		}
	}

	create_info.EnabledExtensionCount = len(ext_names)
	create_info.PpEnabledExtensionNames = ext_names

	if enableValidationLayers {

		if !checkValidationLayerSupport(validationLayers) {
			// TODO
			fmt.Println("Validation layers not supported")
		}

		create_info.EnabledLayerCount = len(validationLayers)
		create_info.PpEnabledLayerNames = validationLayers
	} else {
		create_info.EnabledLayerCount = 0
	}

	var instance vulkan.VkInstance
	var err = vulkan.VkCreateInstance(&create_info, nil, &instance)
	if vulkan.VK_SUCCESS != err {
		// TODO
		fmt.Println("VkCreateInstance failed")
	}

	o.Instance = instance
}

func checkValidationLayerSupport(a []string) bool {

	var cnt uint32
	vulkan.VkEnumerateInstanceLayerProperties(&cnt, nil)

	var a1 = make([]vulkan.VkLayerProperties, cnt)
	vulkan.VkEnumerateInstanceLayerProperties(&cnt, a1)

	fmt.Printf("%v available layers\n", cnt)
	for i := 0; i < int(cnt); i++ {
		var p = a1[i]
		fmt.Printf("\t#%v: %v\n", i, p.LayerName)
	}

	for i := 0; i < len(a); i++ {
		var name = a[i]

		var b = false
		for j := 0; j < len(a1); j++ {
			var prop = &a1[j]
			if name == prop.LayerName {
				b = true
				break
			}
		} // for j

		if !b {
			return false
		}
	} // for i

	return true
}

func (o *HelloTriangleApplication) createSurface() {

	var surface vulkan.VkSurfaceKHR

	if !sdl2.SDL_Vulkan_CreateSurface(o.Window, o.Instance, &surface) {
		// TODO
		fmt.Println("SDL_Vulkan_CreateSurface() failed")
	}

	o.Surface = surface
}

func (o *HelloTriangleApplication) pickPhysicalDevice(
	queue_families *QueueFamilyIndices,
	swap_chain_support *SwapChainSupportDetails,
) {

	var devices []vulkan.VkPhysicalDevice
	{
		var cnt uint32
		vulkan.VkEnumeratePhysicalDevices(o.Instance, &cnt, nil)

		fmt.Printf("%v physical devices\n", cnt)

		if 0 == cnt {
			// TODO
			fmt.Println("No physical device")
		}

		devices = make([]vulkan.VkPhysicalDevice, cnt)
		vulkan.VkEnumeratePhysicalDevices(o.Instance, &cnt, devices)
	}

	for _, d := range devices {
		if o.isDeviceSuitable(d, queue_families, swap_chain_support) {
			o.PhysicalDevice = d
			return
		}
	}

	// TODO
	fmt.Println("No suitable device")
	panic("no suitable device")
}

func (o *HelloTriangleApplication) isDeviceSuitable(
	device vulkan.VkPhysicalDevice,
	queue_families *QueueFamilyIndices,
	swap_chain_support *SwapChainSupportDetails,
) bool {

	var prop vulkan.VkPhysicalDeviceProperties
	var feat vulkan.VkPhysicalDeviceFeatures

	vulkan.VkGetPhysicalDeviceProperties(device, &prop)
	vulkan.VkGetPhysicalDeviceFeatures(device, &feat)

	// Dedicated graphics card
	if vulkan.VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU != prop.DeviceType {
		fmt.Println("VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU not available")
		return false
	}

	// Support geometry shader
	if !feat.GeometryShader {
		fmt.Println("GeometryShader feature not available")
		return false
	}

	if !findQueueFamilies(device, o.Surface, queue_families) {
		fmt.Println("Required queue families not avaible")
		return false
	}

	if !checkDeviceExtensionSupport(device) {
		fmt.Println("Device extensions not avaible")
		return false
	}

	querySwapChainSupport(device, o.Surface, swap_chain_support)
	if 0 == len(swap_chain_support.Formats) && 0 == len(swap_chain_support.PresentModes) {
		fmt.Println("Swap chian support not available")
		return false
	}

	return true
}

func findQueueFamilies(device vulkan.VkPhysicalDevice, surface vulkan.VkSurfaceKHR, queue_families *QueueFamilyIndices) bool {

	var cnt uint32
	vulkan.VkGetPhysicalDeviceQueueFamilyProperties(device, &cnt, nil)

	var a = make([]vulkan.VkQueueFamilyProperties, cnt)
	vulkan.VkGetPhysicalDeviceQueueFamilyProperties(device, &cnt, a)

	if 0 == cnt {
		// TODO
		fmt.Println("No queue families")
		return false
	}

	queue_families.GraphicsFamily = -1
	queue_families.PresentFamily = -1

	for i, prop := range a {

		if 0 != prop.QueueFlags&vulkan.VkQueueFlags(vulkan.VK_QUEUE_GRAPHICS_BIT) {
			queue_families.GraphicsFamily = i
		}

		var b = false
		vulkan.VkGetPhysicalDeviceSurfaceSupportKHR(device, i, surface, &b)
		if b {
			queue_families.PresentFamily = i
		}

		if -1 != queue_families.GraphicsFamily && -1 != queue_families.PresentFamily {

			// TODO
			fmt.Printf("Graphics queue family = %v, present queue faimly = %v\n",
				queue_families.GraphicsFamily, queue_families.PresentFamily)

			return true
		}
	} // for

	return false
}

func checkDeviceExtensionSupport(device vulkan.VkPhysicalDevice) bool {

	var cnt int
	vulkan.VkEnumerateDeviceExtensionProperties(device, nil, &cnt, nil)

	var a = make([]vulkan.VkExtensionProperties, cnt)
	vulkan.VkEnumerateDeviceExtensionProperties(device, nil, &cnt, a)

	fmt.Printf("%v avaible device extensions: \n", cnt)
	{
		const max = 10
		for i, p := range a {

			if i >= max {
				fmt.Println("\t...")
				break
			}

			fmt.Printf("\t#%v: %v\n", i, p.ExtensionName)
		}
	}

	for _, name := range deviceExtensions {

		var exists = false
		for _, prop := range a {
			if name == prop.ExtensionName {
				exists = true
				break
			}
		} // for

		if !exists {
			fmt.Printf("Device extension %v not avaible\n", name)
			return false
		}

	} // for

	return true
}

func querySwapChainSupport(device vulkan.VkPhysicalDevice, surface vulkan.VkSurfaceKHR, swap_chain_support *SwapChainSupportDetails) {

	// Surface capabilities ---------

	var surface_cap vulkan.VkSurfaceCapabilitiesKHR
	{
		vulkan.VkGetPhysicalDeviceSurfaceCapabilitiesKHR(device, surface, &surface_cap)
	}

	// Surface formats ------

	var formats []vulkan.VkSurfaceFormatKHR
	{
		var cnt int
		vulkan.VkGetPhysicalDeviceSurfaceFormatsKHR(device, surface, &cnt, nil)

		formats = make([]vulkan.VkSurfaceFormatKHR, cnt)
		vulkan.VkGetPhysicalDeviceSurfaceFormatsKHR(device, surface, &cnt, formats)
	}

	// Present modes -----------

	var present_modes []vulkan.VkPresentModeKHR
	{
		var cnt int
		vulkan.VkGetPhysicalDeviceSurfacePresentModesKHR(device, surface, &cnt, nil)

		present_modes = make([]vulkan.VkPresentModeKHR, cnt)
		vulkan.VkGetPhysicalDeviceSurfacePresentModesKHR(device, surface, &cnt, present_modes)
	}

	swap_chain_support.Capabilities = surface_cap
	swap_chain_support.Formats = formats
	swap_chain_support.PresentModes = present_modes
}

func (o *HelloTriangleApplication) createLogicalDevice(queue_families *QueueFamilyIndices) {

	var create_info vulkan.VkDeviceCreateInfo

	// Queue create info -----

	var queue_create_infos []vulkan.VkDeviceQueueCreateInfo
	{
		var indices []int

		if queue_families.GraphicsFamily == queue_families.PresentFamily {
			indices = make([]int, 1)
			indices[0] = queue_families.GraphicsFamily
		} else {
			indices = make([]int, 2)
			indices[0] = queue_families.GraphicsFamily
			indices[1] = queue_families.PresentFamily
		}

		var queue_create_infos = make([]vulkan.VkDeviceQueueCreateInfo, len(indices))

		const queue_prio float32 = 1.0

		for i, family := range indices {
			var p1 = &queue_create_infos[i]
			p1.QueueFamilyIndex = family
			p1.QueueCount = 1
			p1.PQueuePriorities = []float32{queue_prio}
		} // for
	}

	create_info.PQueueCreateInfos = queue_create_infos
	create_info.QueueCreateInfoCount = len(queue_create_infos)

	// Device features -----

	var features vulkan.VkPhysicalDeviceFeatures
	create_info.PEnabledFeatures = &features

	// Enable extensions ------

	create_info.PpEnabledExtensionNames = deviceExtensions // Global variable
	create_info.EnabledExtensionCount = len(deviceExtensions)

	// Create logical device -----

	var device vulkan.VkDevice
	var err = vulkan.VkCreateDevice(o.PhysicalDevice, &create_info, nil, &device)
	if vulkan.VK_SUCCESS != err {
		// TODO
		fmt.Println("VkCreateDevice() failed")
	}

	// Look up queues -----

	var graph_queue, present_queue vulkan.VkQueue
	vulkan.VkGetDeviceQueue(device, queue_families.GraphicsFamily, 0, &graph_queue)

	if queue_families.GraphicsFamily == queue_families.PresentFamily {
		present_queue = graph_queue
	} else {
		vulkan.VkGetDeviceQueue(device, queue_families.PresentFamily, 0, &present_queue)
	}

	//
	o.Device = device
	o.GraphicsQueue = graph_queue
	o.PresentQueue = present_queue
}

func (o *HelloTriangleApplication) createSwapChain(
	queue_families *QueueFamilyIndices,
	swap_chain_support *SwapChainSupportDetails,
) {

	var format = chooseSwapSurfaceFormat(swap_chain_support.Formats)
	var present_mode = chooseSwapPresentMode(swap_chain_support.PresentModes)
	var extent = chooseSwapExtent(swap_chain_support.Capabilities, o.Window)

	var image_cnt = swap_chain_support.Capabilities.MinImageCount + 1
	if max := swap_chain_support.Capabilities.MaxImageCount; max > 0 && image_cnt > max {
		image_cnt = max
	}

	var create_info vulkan.VkSwapchainCreateInfoKHR
	create_info.Surface = o.Surface
	create_info.MinImageCount = image_cnt
	create_info.ImageFormat = format.Format
	create_info.ImageColorSpace = format.ColorSpace
	create_info.ImageExtent = extent
	create_info.ImageArrayLayers = 1
	create_info.ImageUsage = vulkan.VkImageUsageFlags(vulkan.VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT)

	if queue_families.GraphicsFamily == queue_families.PresentFamily {
		create_info.ImageSharingMode = vulkan.VK_SHARING_MODE_EXCLUSIVE
		create_info.QueueFamilyIndexCount = 0
		create_info.PQueueFamilyIndices = nil
	} else {
		create_info.ImageSharingMode = vulkan.VK_SHARING_MODE_CONCURRENT
		create_info.QueueFamilyIndexCount = 2

		var a = make([]int, 2)
		create_info.PQueueFamilyIndices = a
		a[0] = queue_families.GraphicsFamily
		a[1] = queue_families.PresentFamily
	}

	create_info.PreTransform = swap_chain_support.Capabilities.CurrentTransform
	create_info.CompositeAlpha = vulkan.VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR
	create_info.PresentMode = present_mode
	create_info.Clipped = true
	create_info.OldSwapchain = vulkan.VkSwapchainKHR{} // Default zero value as NULL

	var swap_chain vulkan.VkSwapchainKHR
	if vulkan.VK_SUCCESS != vulkan.VkCreateSwapchainKHR(o.Device, &create_info, nil, &swap_chain) {
		fmt.Println("VkCreateSwapchainKHR() failed")
	}

	// ---
	var images []vulkan.VkImage
	{
		var cnt int
		vulkan.VkGetSwapchainImagesKHR(o.Device, swap_chain, &cnt, nil)

		fmt.Printf("%v swap chain images\n", cnt)

		images = make([]vulkan.VkImage, cnt)
		vulkan.VkGetSwapchainImagesKHR(o.Device, swap_chain, &cnt, images)
	}

	o.SwapChain = swap_chain
	o.Images = images
	o.ImageFormat = format.Format
	o.Extent = extent
}

func chooseSwapSurfaceFormat(a []vulkan.VkSurfaceFormatKHR) vulkan.VkSurfaceFormatKHR {

	for _, fmt := range a {
		if vulkan.VK_COLOR_SPACE_SRGB_NONLINEAR_KHR == fmt.ColorSpace && vulkan.VK_FORMAT_B8G8R8A8_SRGB == fmt.Format {
			return fmt
		}
	}

	return a[0]
}

func chooseSwapPresentMode(a []vulkan.VkPresentModeKHR) vulkan.VkPresentModeKHR {

	for _, m := range a {
		if vulkan.VK_PRESENT_MODE_MAILBOX_KHR == m {
			return m
		}
	}

	return vulkan.VK_PRESENT_MODE_FIFO_KHR
}

func chooseSwapExtent(cap vulkan.VkSurfaceCapabilitiesKHR, window *sdl2.SDL_Window) vulkan.VkExtent2D {

	const max = 0xffff_ffff

	if max != cap.CurrentExtent.Width && max != cap.CurrentExtent.Height {
		return cap.CurrentExtent
	}

	var w, h int
	sdl2.SDL_Vulkan_GetDrawableSize(window, &w, &h)

	var ext = vulkan.VkExtent2D{
		Width:  uint32(w),
		Height: uint32(h),
	}

	var clamp = func(n, min, max uint32) uint32 {
		if n < min {
			return min
		} else if n > max {
			return max
		} else {
			return n
		}
	}

	ext.Width = clamp(ext.Width, cap.MinImageExtent.Width, cap.MaxImageExtent.Width)
	ext.Height = clamp(ext.Height, cap.MinImageExtent.Height, cap.MaxImageExtent.Height)
	return ext
}

func (o *HelloTriangleApplication) createImageViews() {

	var image_views = make([]vulkan.VkImageView, len(o.Images))

	for i, image := range o.Images {
		o.createImageView(image, &image_views[i])
	} // for

	o.ImageViews = image_views
}

func (o *HelloTriangleApplication) createImageView(image vulkan.VkImage, image_view *vulkan.VkImageView) {

	var create_info vulkan.VkImageViewCreateInfo
	create_info.Image = image
	create_info.ViewType = vulkan.VK_IMAGE_VIEW_TYPE_2D
	create_info.Format = o.ImageFormat
	create_info.Components = vulkan.VkComponentMapping{
		R: vulkan.VK_COMPONENT_SWIZZLE_IDENTITY,
		G: vulkan.VK_COMPONENT_SWIZZLE_IDENTITY,
		B: vulkan.VK_COMPONENT_SWIZZLE_IDENTITY,
		A: vulkan.VK_COMPONENT_SWIZZLE_IDENTITY,
	}
	create_info.SubresourceRange = vulkan.VkImageSubresourceRange{
		AspectMask:     vulkan.VkImageAspectFlags(vulkan.VK_IMAGE_ASPECT_COLOR_BIT),
		BaseMipLevel:   1,
		LevelCount:     1,
		BaseArrayLayer: 0,
		LayerCount:     1,
	}

	if vulkan.VK_SUCCESS != vulkan.VkCreateImageView(o.Device, &create_info, nil, image_view) {
		fmt.Println("VkCreateImageView() failed")
	}

}

func (o *HelloTriangleApplication) createGraphicsPipeline() {

	var vert_shader, frag_shader = o.createShaderModule("shaders/vert.spv"),
		o.createShaderModule("shaders/frag.spv")

	defer func() {
		vulkan.VkDestroyShaderModule(o.Device, vert_shader, nil)
		vulkan.VkDestroyShaderModule(o.Device, frag_shader, nil)
	}()

	var shader_stages = make([]vulkan.VkPipelineShaderStageCreateInfo, 2)
	var name = "main"

	//
	{
		var create_info = &shader_stages[0]
		create_info.Stage = vulkan.VK_SHADER_STAGE_VERTEX_BIT
		create_info.Module = vert_shader
		create_info.PName = &name
	}

	//
	{
		var create_info = &shader_stages[1]
		create_info.Stage = vulkan.VK_SHADER_STAGE_FRAGMENT_BIT
		create_info.Module = frag_shader
		create_info.PName = &name
	}

}

func (o *HelloTriangleApplication) createShaderModule(file string) vulkan.VkShaderModule {

	var code, err = os.ReadFile(file)
	if nil != err {
		fmt.Printf("Cannot read file: %v\n", file)
	}

	var create_info vulkan.VkShaderModuleCreateInfo
	create_info.CodeSize = len(code)
	create_info.PCode = code

	var shader vulkan.VkShaderModule
	if vulkan.VK_SUCCESS != vulkan.VkCreateShaderModule(
		o.Device,
		&create_info,
		nil,
		&shader,
	) {
		fmt.Println("VkCreateShaderModule() failed")
	}

	return shader
}
