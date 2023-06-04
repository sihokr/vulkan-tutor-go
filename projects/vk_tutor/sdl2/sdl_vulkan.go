package sdl2

// #cgo CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/include/SDL2 -DSDL_MAIN_HANDLED
// #cgo LDFLAGS: -L${SRCDIR}/lib -static -lmingw32 -lSDL2main -lSDL2 -mwindows
// #cgo LDFLAGS: -Wl,--dynamicbase -Wl,--nxcompat -Wl,--high-entropy-va -lm -ldinput8 -ldxguid -ldxerr8 -luser32 -lgdi32 -lwinmm -limm32 -lole32 -loleaut32 -lshell32 -lsetupapi -lversion -luuid
// #include "SDL.h"
// #include "SDL_vulkan.h"
import "C"

import (
	"example.com/vk_tutor/sdl2/internal"
	"example.com/vk_tutor/vulkan"
	"unsafe"
)

/**
 *  \file SDL_vulkan.h
 *
 *  Header file for functions to creating Vulkan surfaces on SDL windows.
 */

/* Avoid including vulkan.h, don't define VkInstance if it's already included */
// #ifdef VULKAN_H_
// #define NO_SDL_VULKAN_TYPEDEFS
// #endif
// #ifndef NO_SDL_VULKAN_TYPEDEFS
// #define VK_DEFINE_HANDLE(object) typedef struct object##_T* object;

// #if defined(__LP64__) || defined(_WIN64) || defined(__x86_64__) || defined(_M_X64) || defined(__ia64) || defined (_M_IA64) || defined(__aarch64__) || defined(__powerpc64__)
// #define VK_DEFINE_NON_DISPATCHABLE_HANDLE(object) typedef struct object##_T *object;
// #else
// #define VK_DEFINE_NON_DISPATCHABLE_HANDLE(object) typedef uint64_t object;
// #endif

type cHandleWrapper[T any] struct {
	cHandle T
}

func wrap[T any](w unsafe.Pointer, p *T) {
	var w1 = (*cHandleWrapper[T])(w)
	w1.cHandle = *p
}

func unwrap[T any](w unsafe.Pointer) *T {
	var w1 = (*cHandleWrapper[T])(w)
	return &w1.cHandle
}

// VK_DEFINE_HANDLE(VkInstance)
type VkInstance cHandleWrapper[C.VkInstance]

// VK_DEFINE_NON_DISPATCHABLE_HANDLE(VkSurfaceKHR)
type VkSurfaceKHR cHandleWrapper[C.VkSurfaceKHR]

// #endif /* !NO_SDL_VULKAN_TYPEDEFS */

// typedef VkInstance SDL_vulkanInstance;
// typedef VkSurfaceKHR SDL_vulkanSurface; /* for compatibility with Tizen */

/**
 *  \name Vulkan support functions
 *
 *  \note SDL_Vulkan_GetInstanceExtensions & SDL_Vulkan_CreateSurface API
 *        is compatable with Tizen's implementation of Vulkan in SDL.
 */
/* @{ */

/**
 * Dynamically load the Vulkan loader library.
 *
 * This should be called after initializing the video driver, but before
 * creating any Vulkan windows. If no Vulkan loader library is loaded, the
 * default library will be loaded upon creation of the first Vulkan window.
 *
 * It is fairly common for Vulkan applications to link with libvulkan instead
 * of explicitly loading it at run time. This will work with SDL provided the
 * application links to a dynamic library and both it and SDL use the same
 * search path.
 *
 * If you specify a non-NULL `path`, an application should retrieve all of the
 * Vulkan functions it uses from the dynamic library using
 * SDL_Vulkan_GetVkGetInstanceProcAddr unless you can guarantee `path` points
 * to the same vulkan loader library the application linked to.
 *
 * On Apple devices, if `path` is NULL, SDL will attempt to find the
 * `vkGetInstanceProcAddr` address within all the Mach-O images of the current
 * process. This is because it is fairly common for Vulkan applications to
 * link with libvulkan (and historically MoltenVK was provided as a static
 * library). If it is not found, on macOS, SDL will attempt to load
 * `vulkan.framework/vulkan`, `libvulkan.1.dylib`,
 * `MoltenVK.framework/MoltenVK`, and `libMoltenVK.dylib`, in that order. On
 * iOS, SDL will attempt to load `libMoltenVK.dylib`. Applications using a
 * dynamic framework or .dylib must ensure it is included in its application
 * bundle.
 *
 * On non-Apple devices, application linking with a static libvulkan is not
 * supported. Either do not link to the Vulkan loader or link to a dynamic
 * library version.
 *
 * \param path The platform dependent Vulkan loader library name or NULL
 * \returns 0 on success or -1 if the library couldn't be loaded; call
 *          SDL_GetError() for more information.
 *
 * \since This function is available since SDL 2.0.6.
 *
 * \sa SDL_Vulkan_GetVkInstanceProcAddr
 * \sa SDL_Vulkan_UnloadLibrary
 */
// extern DECLSPEC int SDLCALL SDL_Vulkan_LoadLibrary(const char *path);

/**
 * Get the address of the `vkGetInstanceProcAddr` function.
 *
 * This should be called after either calling SDL_Vulkan_LoadLibrary() or
 * creating an SDL_Window with the `SDL_WINDOW_VULKAN` flag.
 *
 * \returns the function pointer for `vkGetInstanceProcAddr` or NULL on error.
 *
 * \since This function is available since SDL 2.0.6.
 */
// extern DECLSPEC void *SDLCALL SDL_Vulkan_GetVkGetInstanceProcAddr(void);

/**
 * Unload the Vulkan library previously loaded by SDL_Vulkan_LoadLibrary()
 *
 * \since This function is available since SDL 2.0.6.
 *
 * \sa SDL_Vulkan_LoadLibrary
 */
// extern DECLSPEC void SDLCALL SDL_Vulkan_UnloadLibrary(void);

/**
 * Get the names of the Vulkan instance extensions needed to create a surface
 * with SDL_Vulkan_CreateSurface.
 *
 * If `pNames` is NULL, then the number of required Vulkan instance extensions
 * is returned in `pCount`. Otherwise, `pCount` must point to a variable set
 * to the number of elements in the `pNames` array, and on return the variable
 * is overwritten with the number of names actually written to `pNames`. If
 * `pCount` is less than the number of required extensions, at most `pCount`
 * structures will be written. If `pCount` is smaller than the number of
 * required extensions, SDL_FALSE will be returned instead of SDL_TRUE, to
 * indicate that not all the required extensions were returned.
 *
 * The `window` parameter is currently needed to be valid as of SDL 2.0.8,
 * however, this parameter will likely be removed in future releases
 *
 * \param window A window for which the required Vulkan instance extensions
 *               should be retrieved (will be deprecated in a future release)
 * \param pCount A pointer to an unsigned int corresponding to the number of
 *               extensions to be returned
 * \param pNames NULL or a pointer to an array to be filled with required
 *               Vulkan instance extensions
 * \returns SDL_TRUE on success, SDL_FALSE on error.
 *
 * \since This function is available since SDL 2.0.6.
 *
 * \sa SDL_Vulkan_CreateSurface
 */
// extern DECLSPEC SDL_bool SDLCALL SDL_Vulkan_GetInstanceExtensions(SDL_Window *window,
//                                                                   unsigned int *pCount,
//                                                                   const char **pNames);
func SDL_Vulkan_GetInstanceExtensions(window *SDL_Window, pCount *uint, pNames []string) bool {

	var window1 *C.SDL_Window
	if nil != window {
		window1 = (*C.SDL_Window)(internal.Unwrap[SDL_Window](window))
	}

	var count1 = C.uint(*pCount)

	// Assume pointer size = sizeof(uint)
	// const sizeof_p = C.sizeof_uint
	// VkInstance is defined as pointer
	const sizeof_p = C.sizeof_VkInstance

	var pNames1 **C.char
	if nil != pNames && *pCount > 0 {

		var buf1 = C.malloc(C.ulonglong((*pCount) * uint(vulkan.VK_MAX_EXTENSION_NAME_SIZE)))
		var buf2 = C.malloc(C.ulonglong((*pCount) * sizeof_p))

		defer func() {
			C.free(buf1)
			C.free(buf2)
		}()

		pNames1 = (**C.char)(buf2)

		var p1 = uintptr(buf1)
		var p2 = uintptr(buf2)

		for i := 0; i < int(*pCount); i++ {

			var p1_1 = (*C.char)(unsafe.Pointer(p1))
			var p2_1 = (**C.char)(unsafe.Pointer(p2))
			*p2_1 = p1_1

			p1 += uintptr(vulkan.VK_MAX_EXTENSION_NAME_SIZE)
			p2 += uintptr(sizeof_p)
		}
	}

	if C.SDL_TRUE != C.SDL_Vulkan_GetInstanceExtensions(
		window1,
		&count1,
		pNames1,
	) {
		return false
	}

	*pCount = uint(count1)

	if nil != pNames1 {

		var p = uintptr(unsafe.Pointer(pNames1))

		for i := 0; i < int(*pCount); i++ {
			var p1 = (**C.char)(unsafe.Pointer(p))
			pNames[i] = C.GoString(*p1)
			p += uintptr(sizeof_p)
		}
	}

	return true
}

/**
 * Create a Vulkan rendering surface for a window.
 *
 * The `window` must have been created with the `SDL_WINDOW_VULKAN` flag and
 * `instance` must have been created with extensions returned by
 * SDL_Vulkan_GetInstanceExtensions() enabled.
 *
 * \param window The window to which to attach the Vulkan surface
 * \param instance The Vulkan instance handle
 * \param surface A pointer to a VkSurfaceKHR handle to output the newly
 *                created surface
 * \returns SDL_TRUE on success, SDL_FALSE on error.
 *
 * \since This function is available since SDL 2.0.6.
 *
 * \sa SDL_Vulkan_GetInstanceExtensions
 * \sa SDL_Vulkan_GetDrawableSize
 */
// extern DECLSPEC SDL_bool SDLCALL SDL_Vulkan_CreateSurface(SDL_Window *window,
//                                                           VkInstance instance,
//                                                           VkSurfaceKHR* surface);
func SDL_Vulkan_CreateSurface(window *SDL_Window, instance vulkan.VkInstance, surface *vulkan.VkSurfaceKHR) bool {

	var instance_sdl VkInstance
	{
		var p_instance_sdl_c = unwrap[C.VkInstance](unsafe.Pointer(&instance_sdl))
		var p_instance_vk_c = instance.Unwrap()
		*p_instance_sdl_c = *((*C.VkInstance)(p_instance_vk_c))
	}

	var surface_sdl VkSurfaceKHR

	if !_SDL_Vulkan_CreateSurface(window, instance_sdl, &surface_sdl) {
		return false
	}

	//
	{
		var p_surface_sdl_c = unwrap[C.VkSurfaceKHR](unsafe.Pointer(&surface_sdl))
		var p_surface_vk_c = surface.Unwrap()
		*((*C.VkSurfaceKHR)(p_surface_vk_c)) = *p_surface_sdl_c
	}

	return true
}

func _SDL_Vulkan_CreateSurface(window *SDL_Window, instance VkInstance, surface *VkSurfaceKHR) bool {

	var window1 = internal.Unwrap[SDL_Window](window)
	var instance1 = unwrap[C.VkInstance](unsafe.Pointer(&instance))
	var surface1 C.VkSurfaceKHR

	if C.SDL_TRUE != C.SDL_Vulkan_CreateSurface(
		(*C.SDL_Window)(window1),
		*instance1,
		&surface1,
	) {
		return false
	}

	wrap[C.VkSurfaceKHR](unsafe.Pointer(surface), &surface1)

	return true
}

/**
 * Get the size of the window's underlying drawable dimensions in pixels.
 *
 * This may differ from SDL_GetWindowSize() if we're rendering to a high-DPI
 * drawable, i.e. the window was created with `SDL_WINDOW_ALLOW_HIGHDPI` on a
 * platform with high-DPI support (Apple calls this "Retina"), and not
 * disabled by the `SDL_HINT_VIDEO_HIGHDPI_DISABLED` hint.
 *
 * \param window an SDL_Window for which the size is to be queried
 * \param w Pointer to the variable to write the width to or NULL
 * \param h Pointer to the variable to write the height to or NULL
 *
 * \since This function is available since SDL 2.0.6.
 *
 * \sa SDL_GetWindowSize
 * \sa SDL_CreateWindow
 * \sa SDL_Vulkan_CreateSurface
 */
// extern DECLSPEC void SDLCALL SDL_Vulkan_GetDrawableSize(SDL_Window * window,
//                                                         int *w, int *h);
func SDL_Vulkan_GetDrawableSize(window *SDL_Window, w, h *int) {

	var window1 = internal.Unwrap[SDL_Window](window)
	var w1, h1 C.int
	C.SDL_Vulkan_GetDrawableSize((*C.SDL_Window)(window1), &w1, &h1)

	if nil != w {
		*w = int(w1)
	}

	if nil != h {
		*h = int(h1)
	}
}
