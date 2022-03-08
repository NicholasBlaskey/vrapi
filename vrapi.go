//go:build darwin || linux || windows
// +build android, gldebug
package vrapi

/*
#cgo CPPFLAGS: -I./Include -I/usr/local/include
#cgo LDFLAGS: -v -march=armv8-a -shared -L./lib/arm64-v8a/ -lvrapi -landroid

#include <VrApi.h>
#include <VrApi_Helpers.h>
#include <VrApi_Input.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	INITIALIZE_SUCCESS = C.VRAPI_INITIALIZE_SUCCESS

	OVRSuccess = C.ovrSuccess
)

type OVRInitParms C.ovrInitParms // HMMM alias this type?
type OVRStructureType int32

const ( // OVRInitParms
	STRUCTURE_TYPE_INIT_PARMS        = 1
	STRUCTURE_TYPE_MODE_PARMS        = 2
	STRUCTURE_TYPE_FRAME_PARMS       = 3
	STRUCTURE_TYPE_MODE_PARMS_VULKAN = 5
)

//type OVRModeParms C.ovrModeParms
// Just experiment with this?
type OVRModeParms struct {
	Type  OVRStructureType
	Flags uint32
	Java  OVRJava
	//Padding       int32 // ??? // Add in build constraint for padding here?
	Display       uint64
	WindowSurface uint64
	ShareContext  uint64
}

type OVRJava C.ovrJava
type OVRMobile C.ovrMobile

type OVRTracking2 C.ovrTracking2

func DefaultInitParms(java *OVRJava) OVRInitParms {
	cParms := C.vrapi_DefaultInitParms((*C.ovrJava)(java))
	return OVRInitParms(cParms)
}

func DefaultModeParms(java *OVRJava) OVRModeParms {
	cParms := C.vrapi_DefaultModeParms((*C.ovrJava)(java))
	return *(*OVRModeParms)(unsafe.Pointer(&cParms))
}

// glctx
func EnterVrMode(modeParms *OVRModeParms) *OVRMobile {
	cParms := (*C.ovrModeParms)(unsafe.Pointer(modeParms))
	ovr := C.vrapi_EnterVrMode(cParms)
	return (*OVRMobile)(ovr)
}

// This should run with vrctx like glctx on a seperate thread.???
// Could be intresting / helpful to seperate functions on seperate threads.
// Following same worker pattern as gl.
// Should this be a status? or an error code
func Initialize(parms *OVRInitParms) error {
	status := C.vrapi_Initialize((*C.ovrInitParms)(parms))
	if status != INITIALIZE_SUCCESS {
		return fmt.Errorf("vrapi_Initialize status %d not equal to sucess %d",
			status, INITIALIZE_SUCCESS)
	}
	return nil
}

func GetPredictedDisplayTime(vrApp *OVRMobile, frameIndex int64) float64 {
	return float64(C.vrapi_GetPredictedDisplayTime((*C.ovrMobile)(vrApp),
		C.longlong(frameIndex)))
}

func GetPredictedTracking2(vrApp *OVRMobile, displayTime float64) OVRTracking2 {
	cOVR := (*C.ovrMobile)(unsafe.Pointer(vrApp))
	cTracking := C.vrapi_GetPredictedTracking2(cOVR, C.double(displayTime))
	return OVRTracking2(cTracking)
}

// Input (move to seperate file)
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

type OVRDeviceID uint32

type OVRInputCapabilityHeader struct {
	Type     OVRControllerType
	DeviceID OVRDeviceID
}

// TODO should this return an error?
func EnumerateInputDevices(vrApp *OVRMobile, index uint32,
	capsHeader *OVRInputCapabilityHeader) int32 {

	cOVR := (*C.ovrMobile)(unsafe.Pointer(vrApp))
	cHeader := (*C.ovrInputCapabilityHeader)(unsafe.Pointer(capsHeader))
	res := C.vrapi_EnumerateInputDevices(cOVR, C.uint(index), cHeader)

	//log.Printf("cHeader %+v %p", cHeader, cHeader)
	//capsHeader = (*OVRInputCapabilityHeader)(unsafe.Pointer(cHeader)) // Convert back?
	//log.Printf("goHeader %+v %p %+v", capsHeader, capsHeader, capsHeader.Type)

	return int32(res)
}

type OVRInputStateStandardPointer struct {
	Header           OVRInputStateHeader
	PointerPose      OVRPosef
	PointerStrength  float32
	GripPose         OVRPosef
	InputStateStatus uint32
	Reserved         [20]uint64 // Reserved for future use
}

type OVRInputStateHeader struct {
	ControllerType OVRControllerType
	TimeInSeconds  float64
}

type OVRPosef struct {
	Orientation mgl.Quat
	Position    mgl.Vec3 // aka Translation (Limitation due to Go not having unions)
}

type OVRInputStandardPointerCapabilities struct {
	Header                 OVRInputCapabilityHeader
	ControllerCapabilities uint32     // Mask of controller capabilities described by ovrControllerCapabilities
	HapticSamplesMax       uint32     // Maximum submittable samples for the haptics buffer
	HapticSampleDurationMS uint32     // length in milliseconds of a sample in the haptics buffer.
	Reserved               [20]uint64 // Reserved for future use
}

func GetCurrentInputState(vrApp *OVRMobile,
	deviceID OVRDeviceID, inputState *OVRInputStateHeader) error {

	cOVR := (*C.ovrMobile)(unsafe.Pointer(vrApp))
	cInputState := (*C.ovrInputStateHeader)(unsafe.Pointer(inputState))
	res := C.vrapi_GetCurrentInputState(cOVR, C.uint(deviceID), cInputState)
	if res != OVRSuccess {
		return fmt.Errorf("get current input state expected sucess (%d) got %d",
			OVRSuccess, res)
	}
	return nil
}

func GetInputDeviceCapabilities(vrApp *OVRMobile,
	capsHeader *OVRInputCapabilityHeader) error {

	cOVR := (*C.ovrMobile)(unsafe.Pointer(vrApp))
	cCapsHeader := (*C.ovrInputCapabilityHeader)(unsafe.Pointer(capsHeader))
	res := C.vrapi_GetInputDeviceCapabilities(cOVR, cCapsHeader)
	if res != OVRSuccess {
		return fmt.Errorf("get input device capabilities expected sucess (%d) got %d",
			OVRSuccess, res)
	}
	return nil
}

// END input

// Helpers not in the original API.
// Expects the values from
// app.RunOnJVM from the "golang.org/x/mobile/app" package
// This function is behind an android build contraint.
// https://github.com/golang/mobile/blob/8a0a1e50732f652b0c7a0ef4a9f6903c5dc0ca13/app/android.go#L73
func CreateJavaObject(vm, jniEnv, ctx uintptr) OVRJava {
	var java OVRJava
	java.Vm = (*C.JavaVM)(unsafe.Pointer(vm))
	java.Env = (*C.JNIEnv)(unsafe.Pointer(jniEnv))
	java.ActivityObject = (C.jobject)(unsafe.Pointer(ctx))

	return java
}
