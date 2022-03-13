package vrapi

type OVRModeFlags int32

const ( // OVRModeFlags
	// When an application moves backwards on the activity stack
	// the activity window it returns to is no longer flagged as fullscreen.
	// As a result Android will also render the decor view which wastes a
	// significant amount of bandwidth.
	// By setting this flag the fullscreen flag is reset on the window.
	// Unfortunately this causes Android life cycle events that mess up
	// several NativeActivity codebases like Stratum and UE4 so this
	// flag should only be set for specific applications.
	// Use "adb shell dumpsys SurfaceFlinger" to verify
	// that there is only one HWC next to the FB_TARGET.
	MODE_FLAG_RESET_WINDOW_FULLSCREEN OVRModeFlags = 0x0000FF00

	// The WindowSurface passed in is an ANativeWindow.
	MODE_FLAG_NATIVE_WINDOW OVRModeFlags = 0x00010000

	// Create the front buffer in TrustZone memory to allow protected DRM
	// content to be rendered to the front buffer. This functionality
	// requires the WindowSurface to be allocated from TimeWarp via
	// specifying the nativeWindow via MODE_FLAG_NATIVE_WINDOW.
	MODE_FLAG_FRONT_BUFFER_PROTECTED OVRModeFlags = 0x00020000
	// Create a front buffer using the sRGB color space.
	MODE_FLAG_FRONT_BUFFER_SRGB OVRModeFlags = 0x00080000

	// If set indicates the OpenGL ES Context was created with EGL_CONTEXT_OPENGL_NO_ERROR_KHR
	// attribute. The same attribute would be applied when TimeWrap creates the shared context.
	// More information could be found at:
	// https://www.khronos.org/registry/EGL/extensions/KHR/EGL_KHR_create_context_no_error.txt
	MODE_FLAG_CREATE_CONTEXT_NO_ERROR OVRModeFlags = 0x00100000
)

type OVRStructureType int32

const ( // OVRInitParms
	STRUCTURE_TYPE_INIT_PARMS        OVRStructureType = 1
	STRUCTURE_TYPE_MODE_PARMS        OVRStructureType = 2
	STRUCTURE_TYPE_FRAME_PARMS       OVRStructureType = 3
	STRUCTURE_TYPE_MODE_PARMS_VULKAN OVRStructureType = 5
)

type OVRSystemProperty int32

const ( //ovrSystemProperty
	SYS_PROP_DEVICE_TYPE                       OVRSystemProperty = 0
	SYS_PROP_MAX_FULLSPEED_FRAMEBUFFER_SAMPLES OVRSystemProperty = 1
	// Physical width and height of the display in pixels.
	SYS_PROP_DISPLAY_PIXELS_WIDE OVRSystemProperty = 2
	SYS_PROP_DISPLAY_PIXELS_HIGH OVRSystemProperty = 3
	// Returns the refresh rate of the display in cycles per second.
	SYS_PROP_DISPLAY_REFRESH_RATE OVRSystemProperty = 4
	// With a display resolution of 2560x1440 the pixels at the center
	// of each eye cover about 0.06 degrees of visual arc. To wrap a
	// full 360 degrees about 6000 pixels would be needed and about one
	// quarter of that would be needed for ~90 degrees FOV. As such Eye
	// images with a resolution of 1536x1536 result in a good 1:1 mapping
	// in the center but they need mip-maps for off center pixels. To
	// avoid the need for mip-maps and for significantly improved rendering
	// performance this currently returns a conservative 1024x1024.
	SYS_PROP_SUGGESTED_EYE_TEXTURE_WIDTH  OVRSystemProperty = 5
	SYS_PROP_SUGGESTED_EYE_TEXTURE_HEIGHT OVRSystemProperty = 6
	// This is a product of the lens distortion and the screen size
	// but there is no truly correct answer.
	// There is a tradeoff in resolution and coverage.
	// Too small of an FOV will leave unrendered pixels visible but too
	// large wastes resolution or fill rate.  It is unreasonable to
	// increase it until the corners are completely covered but we do
	// want most of the outside edges completely covered.
	// Applications might choose to render a larger FOV when angular
	// acceleration is high to reduce black pull in at the edges by
	// the time warp.
	// Currently symmetric 90.0 degrees.
	SYS_PROP_SUGGESTED_EYE_FOV_DEGREES_X OVRSystemProperty = 7
	SYS_PROP_SUGGESTED_EYE_FOV_DEGREES_Y OVRSystemProperty = 8
	SYS_PROP_DEVICE_REGION               OVRSystemProperty = 10
	// Returns an ovrHandedness enum indicating left or right hand.
	SYS_PROP_DOMINANT_HAND OVRSystemProperty = 15

	// Returns TRUE if the system supports orientation tracking.
	SYS_PROP_HAS_ORIENTATION_TRACKING OVRSystemProperty = 16
	// Returns TRUE if the system supports positional tracking.
	SYS_PROP_HAS_POSITION_TRACKING OVRSystemProperty = 17

	// Returns the number of display refresh rates supported by the system.
	SYS_PROP_NUM_SUPPORTED_DISPLAY_REFRESH_RATES OVRSystemProperty = 64
	// Returns an array of the supported display refresh rates.
	SYS_PROP_SUPPORTED_DISPLAY_REFRESH_RATES OVRSystemProperty = 65

	// Returns the number of swapchain texture formats supported by the system.
	SYS_PROP_NUM_SUPPORTED_SWAPCHAIN_FORMATS OVRSystemProperty = 66
	// Returns an array of the supported swapchain formats.
	// Formats are platform specific. For GLES this is an array of
	// GL internal formats.
	SYS_PROP_SUPPORTED_SWAPCHAIN_FORMATS OVRSystemProperty = 67
	// Returns TRUE if on-chip foveated rendering of swapchains is supported
	// for this system otherwise FALSE.
	SYS_PROP_FOVEATION_AVAILABLE OVRSystemProperty = 130
)

type OVRControllerType uint32

const ( // OVRControllerType
	OVRControllerType_None          OVRControllerType = 0
	OVRControllerType_Reserved0     OVRControllerType = (1 << 0)
	OVRControllerType_Reserved1     OVRControllerType = (1 << 1)
	OVRControllerType_TrackedRemote OVRControllerType = (1 << 2)
	OVRControllerType_Gamepad       OVRControllerType = (1 << 4) // Deprecated, will be removed in a future release
	OVRControllerType_Hand          OVRControllerType = (1 << 5)

	OVRControllerType_StandardPointer OVRControllerType = (1 << 7)
)

type OVRControllerCapabilities uint32

const ( // OVRControllerCapabilities
	OVRControllerCaps_HasOrientationTracking     OVRControllerCapabilities = 0x00000001
	OVRControllerCaps_HasPositionTracking        OVRControllerCapabilities = 0x00000002
	OVRControllerCaps_LeftHand                   OVRControllerCapabilities = 0x00000004 //< Controller is configured for left hand
	OVRControllerCaps_RightHand                  OVRControllerCapabilities = 0x00000008 //< Controller is configured for right hand
	OVRControllerCaps_ModelOculusGo              OVRControllerCapabilities = 0x00000010 //< Controller for Oculus Go devices
	OVRControllerCaps_HasAnalogIndexTrigger      OVRControllerCapabilities = 0x00000040 //< Controller has an analog index trigger vs. a binary one
	OVRControllerCaps_HasAnalogGripTrigger       OVRControllerCapabilities = 0x00000080 //< Controller has an analog grip trigger vs. a binary one
	OVRControllerCaps_HasSimpleHapticVibration   OVRControllerCapabilities = 0x00000200 //< Controller supports simple haptic vibration
	OVRControllerCaps_HasBufferedHapticVibration OVRControllerCapabilities = 0x00000400 //< Controller supports buffered haptic vibration
	OVRControllerCaps_ModelGearVR                OVRControllerCapabilities = 0x00000800 //< Controller is the Gear VR Controller
	OVRControllerCaps_HasTrackpad                OVRControllerCapabilities = 0x00001000 //< Controller has a trackpad
	OVRControllerCaps_HasJoystick                OVRControllerCapabilities = 0x00002000 //< Controller has a joystick.
	OVRControllerCaps_ModelOculusTouch           OVRControllerCapabilities = 0x00004000 //< Oculus Touch Controller For Oculus Quest
	OVRControllerCaps_EnumSize                   OVRControllerCapabilities = 0x7fffffff
)

type OVRLayerType2 uint32

const ( // OVRLayerType2
	LAYER_TYPE_PROJECTION2   OVRLayerType2 = 1
	LAYER_TYPE_CYLINDER2     OVRLayerType2 = 3
	LAYER_TYPE_CUBE2         OVRLayerType2 = 4
	LAYER_TYPE_EQUIRECT2     OVRLayerType2 = 5
	LAYER_TYPE_LOADING_ICON2 OVRLayerType2 = 6
	LAYER_TYPE_FISHEYE2      OVRLayerType2 = 7
	LAYER_TYPE_EQUIRECT3     OVRLayerType2 = 10
)

type OVRFrameLayerBlend uint32

const ( // OVRFrameLayerBlend
	FRAME_LAYER_BLEND_ZERO                OVRFrameLayerBlend = 0
	FRAME_LAYER_BLEND_ONE                 OVRFrameLayerBlend = 1
	FRAME_LAYER_BLEND_SRC_ALPHA           OVRFrameLayerBlend = 2
	FRAME_LAYER_BLEND_ONE_MINUS_SRC_ALPHA OVRFrameLayerBlend = 5
)

type OVRFrameLayerFlags uint32

const (
	/// NOTE: On Oculus standalone devices chromatic aberration correction is enabled
	/// by default.
	/// For non Oculus standalone devices this must be explicitly enabled by specifying the layer
	/// flag as it is a quality / performance trade off.
	FRAME_LAYER_FLAG_CHROMATIC_ABERRATION_CORRECTION OVRFrameLayerFlags = 1 << 1
	/// Used for some HUDs but generally considered bad practice.
	FRAME_LAYER_FLAG_FIXED_TO_VIEW OVRFrameLayerFlags = 1 << 2
	/// Spin the layer - for loading icons
	FRAME_LAYER_FLAG_SPIN OVRFrameLayerFlags = 1 << 3
	/// Clip fragments outside the layer's TextureRect
	FRAME_LAYER_FLAG_CLIP_TO_TEXTURE_RECT OVRFrameLayerFlags = 1 << 4

	/// To get gamma correct sRGB filtering of the eye textures the textures must be
	/// allocated with GL_SRGB8_ALPHA8 format and the window surface must be allocated
	/// with these attributes:
	/// EGL_GL_COLORSPACE_KHR  EGL_GL_COLORSPACE_SRGB_KHR
	///
	/// While we can reallocate textures easily enough we can't change the window
	/// colorspace without relaunching the entire application so if you want to
	/// be able to toggle between gamma correct and incorrect you must allocate
	/// the framebuffer as sRGB then inhibit that processing when using normal
	/// textures.
	///
	/// If the texture being read isn't an sRGB texture the conversion
	/// on write must be inhibited or the colors are washed out.
	/// This is necessary for using external images on an sRGB framebuffer.
	FRAME_LAYER_FLAG_INHIBIT_SRGB_FRAMEBUFFER OVRFrameLayerFlags = 1 << 8

	/// Allow Layer to use an expensive filtering mode. Only useful for 2D layers that are high
	/// resolution (e.g. a remote desktop layer) typically double or more the target resolution.
	FRAME_LAYER_FLAG_FILTER_EXPENSIVE OVRFrameLayerFlags = 1 << 19
)

type OVRTextureType uint32

const ( // OVRTextureType
	TEXTURE_TYPE_2D       OVRTextureType = 0 //< 2D textures.
	TEXTURE_TYPE_2D_ARRAY OVRTextureType = 2 //< Texture array.
	TEXTURE_TYPE_CUBE     OVRTextureType = 3 //< Cube maps.
	TEXTURE_TYPE_MAX      OVRTextureType = 4
)
