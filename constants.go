package vrapi

type OVRStructureType int32

const ( // OVRInitParms
	STRUCTURE_TYPE_INIT_PARMS        = 1
	STRUCTURE_TYPE_MODE_PARMS        = 2
	STRUCTURE_TYPE_FRAME_PARMS       = 3
	STRUCTURE_TYPE_MODE_PARMS_VULKAN = 5
)

type OVRSystemProperty int32

const ( //ovrSystemProperty
	SYS_PROP_DEVICE_TYPE                       = 0
	SYS_PROP_MAX_FULLSPEED_FRAMEBUFFER_SAMPLES = 1
	// Physical width and height of the display in pixels.
	SYS_PROP_DISPLAY_PIXELS_WIDE = 2
	SYS_PROP_DISPLAY_PIXELS_HIGH = 3
	// Returns the refresh rate of the display in cycles per second.
	SYS_PROP_DISPLAY_REFRESH_RATE = 4
	// With a display resolution of 2560x1440 the pixels at the center
	// of each eye cover about 0.06 degrees of visual arc. To wrap a
	// full 360 degrees about 6000 pixels would be needed and about one
	// quarter of that would be needed for ~90 degrees FOV. As such Eye
	// images with a resolution of 1536x1536 result in a good 1:1 mapping
	// in the center but they need mip-maps for off center pixels. To
	// avoid the need for mip-maps and for significantly improved rendering
	// performance this currently returns a conservative 1024x1024.
	SYS_PROP_SUGGESTED_EYE_TEXTURE_WIDTH  = 5
	SYS_PROP_SUGGESTED_EYE_TEXTURE_HEIGHT = 6
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
	SYS_PROP_SUGGESTED_EYE_FOV_DEGREES_X = 7
	SYS_PROP_SUGGESTED_EYE_FOV_DEGREES_Y = 8
	SYS_PROP_DEVICE_REGION               = 10
	// Returns an ovrHandedness enum indicating left or right hand.
	SYS_PROP_DOMINANT_HAND = 15

	// Returns TRUE if the system supports orientation tracking.
	SYS_PROP_HAS_ORIENTATION_TRACKING = 16
	// Returns TRUE if the system supports positional tracking.
	SYS_PROP_HAS_POSITION_TRACKING = 17

	// Returns the number of display refresh rates supported by the system.
	SYS_PROP_NUM_SUPPORTED_DISPLAY_REFRESH_RATES = 64
	// Returns an array of the supported display refresh rates.
	SYS_PROP_SUPPORTED_DISPLAY_REFRESH_RATES = 65

	// Returns the number of swapchain texture formats supported by the system.
	SYS_PROP_NUM_SUPPORTED_SWAPCHAIN_FORMATS = 66
	// Returns an array of the supported swapchain formats.
	// Formats are platform specific. For GLES this is an array of
	// GL internal formats.
	SYS_PROP_SUPPORTED_SWAPCHAIN_FORMATS = 67
	// Returns TRUE if on-chip foveated rendering of swapchains is supported
	// for this system otherwise FALSE.
	SYS_PROP_FOVEATION_AVAILABLE = 130
)

type OVRControllerType uint32

const ( // OVRControllerType
	OVRControllerType_None          = 0
	OVRControllerType_Reserved0     = (1 << 0)
	OVRControllerType_Reserved1     = (1 << 1)
	OVRControllerType_TrackedRemote = (1 << 2)
	OVRControllerType_Gamepad       = (1 << 4) // Deprecated, will be removed in a future release
	OVRControllerType_Hand          = (1 << 5)

	OVRControllerType_StandardPointer = (1 << 7)
)

type OVRControllerCapabilities uint32

const ( // OVRControllerCapabilities
	OVRControllerCaps_HasOrientationTracking     = 0x00000001
	OVRControllerCaps_HasPositionTracking        = 0x00000002
	OVRControllerCaps_LeftHand                   = 0x00000004 //< Controller is configured for left hand
	OVRControllerCaps_RightHand                  = 0x00000008 //< Controller is configured for right hand
	OVRControllerCaps_ModelOculusGo              = 0x00000010 //< Controller for Oculus Go devices
	OVRControllerCaps_HasAnalogIndexTrigger      = 0x00000040 //< Controller has an analog index trigger vs. a binary one
	OVRControllerCaps_HasAnalogGripTrigger       = 0x00000080 //< Controller has an analog grip trigger vs. a binary one
	OVRControllerCaps_HasSimpleHapticVibration   = 0x00000200 //< Controller supports simple haptic vibration
	OVRControllerCaps_HasBufferedHapticVibration = 0x00000400 //< Controller supports buffered haptic vibration
	OVRControllerCaps_ModelGearVR                = 0x00000800 //< Controller is the Gear VR Controller
	OVRControllerCaps_HasTrackpad                = 0x00001000 //< Controller has a trackpad
	OVRControllerCaps_HasJoystick                = 0x00002000 //< Controller has a joystick.
	OVRControllerCaps_ModelOculusTouch           = 0x00004000 //< Oculus Touch Controller For Oculus Quest
	OVRControllerCaps_EnumSize                   = 0x7fffffff
)

type OVRLayerType2 uint32

const ( // OVRLayerType2
	LAYER_TYPE_PROJECTION2   = 1
	LAYER_TYPE_CYLINDER2     = 3
	LAYER_TYPE_CUBE2         = 4
	LAYER_TYPE_EQUIRECT2     = 5
	LAYER_TYPE_LOADING_ICON2 = 6
	LAYER_TYPE_FISHEYE2      = 7
	LAYER_TYPE_EQUIRECT3     = 10
)

type OVRFrameLayerBlend uint32

const ( // OVRFrameLayerBlend
	FRAME_LAYER_BLEND_ZERO                = 0
	FRAME_LAYER_BLEND_ONE                 = 1
	FRAME_LAYER_BLEND_SRC_ALPHA           = 2
	FRAME_LAYER_BLEND_ONE_MINUS_SRC_ALPHA = 5
)

type OVRFrameLayerFlags uint32

const (
	/// NOTE: On Oculus standalone devices chromatic aberration correction is enabled
	/// by default.
	/// For non Oculus standalone devices this must be explicitly enabled by specifying the layer
	/// flag as it is a quality / performance trade off.
	FRAME_LAYER_FLAG_CHROMATIC_ABERRATION_CORRECTION = 1 << 1
	/// Used for some HUDs but generally considered bad practice.
	FRAME_LAYER_FLAG_FIXED_TO_VIEW = 1 << 2
	/// Spin the layer - for loading icons
	FRAME_LAYER_FLAG_SPIN = 1 << 3
	/// Clip fragments outside the layer's TextureRect
	FRAME_LAYER_FLAG_CLIP_TO_TEXTURE_RECT = 1 << 4

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
	FRAME_LAYER_FLAG_INHIBIT_SRGB_FRAMEBUFFER = 1 << 8

	/// Allow Layer to use an expensive filtering mode. Only useful for 2D layers that are high
	/// resolution (e.g. a remote desktop layer) typically double or more the target resolution.
	FRAME_LAYER_FLAG_FILTER_EXPENSIVE = 1 << 19
)

type OVRTextureType uint32

const ( // OVRTextureType
	TEXTURE_TYPE_2D       = 0 //< 2D textures.
	TEXTURE_TYPE_2D_ARRAY = 2 //< Texture array.
	TEXTURE_TYPE_CUBE     = 3 //< Cube maps.
	TEXTURE_TYPE_MAX      = 4
)
