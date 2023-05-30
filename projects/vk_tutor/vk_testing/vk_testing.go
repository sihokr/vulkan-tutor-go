package main

import (
	"fmt"
	"runtime"

	"example.com/vk_tutor/sdl2"
	"example.com/vk_tutor/vulkan"
)

func main() {

	runtime.LockOSThread()

	sdl2.SDL_SetMainReady()
	sdl2.SDL_Init(sdl2.SDL_INIT_EVERYTHING)

	var window = sdl2.SDL_CreateWindow("Vulkan window", 100, 100, 640, 480, sdl2.SDL_WINDOW_SHOWN|sdl2.SDL_WINDOW_VULKAN)

	var ext_cnt uint32
	vulkan.VkEnumerateInstanceExtensionProperties(nil, &ext_cnt, nil)
	fmt.Printf("%v extensions supported\n", ext_cnt)

	untilQuit()

	sdl2.SDL_DestroyWindow(window)
	sdl2.SDL_Quit()
}

func untilQuit() {

	var e sdl2.SDL_Event

loop:
	for {
		if 0 == sdl2.SDL_PollEvent(&e) {
			continue
		}

		switch e.Type() {
		case sdl2.SDL_QUIT:
			break loop
		} // switch
	} // for
}
