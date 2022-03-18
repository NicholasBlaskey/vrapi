//go:build darwin || linux || windows
// +build android, gldebug
package vrapi

/*
#cgo CPPFLAGS: -I./Include -I/usr/local/include
#cgo LDFLAGS: -v -march=armv8-a -shared -L./lib/arm64-v8a/ -lvrapi -landroid

#include <VrApi.h>
#include <VrApi_Helpers.h>
#include <VrApi_Input.h>

ovrResult submit(ovrMobile* ovr, ovrSubmitFrameDescription2* frameDesc, ovrLayerProjection2 layer) {
//ovrResult submit(ovrMobile* ovr, ovrSubmitFrameDescription2* frameDesc, ovrLayerHeader2* layer) {

	// Sets a C pointer to a C pointer?
	const ovrLayerHeader2* layers[] = { &layer.Header };

	//const ovrLayerHeader2* layers[] = { layer };
	(*frameDesc).Layers = layers;
    return vrapi_SubmitFrame2(ovr, frameDesc);
}
*/
import "C"

import (
	"fmt"
	"unsafe"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func jplToHamiltonQuats(quat mgl.Quat) mgl.Quat {
	// https://fzheng.me/2017/11/12/quaternion_conventions_en/
	//
	// https://naif.jpl.nasa.gov/pub/naif/toolkit_docs/C/cspice/q2m_c.html
	//   Relationship between SPICE and Engineering Quaternions
	// Not sure why we don't need the negatives as described in this relationship above?
	return mgl.Quat{quat.V[2], mgl.Vec3{quat.W, quat.V[0], quat.V[1]}}
}

// Constants here since we get these from runtime?
// Should we not get these from runtime?
// Should we get all constants from runtime?
const (
	INITIALIZE_SUCCESS = C.VRAPI_INITIALIZE_SUCCESS

	OVRSuccess = C.ovrSuccess

	FRAME_LAYER_EYE_MAX = C.VRAPI_FRAME_LAYER_EYE_MAX
)

func GetSystemPropertyInt(java *OVRJava, parm OVRSystemProperty) int { // int or int32?
	cJava := (*C.ovrJava)(java)
	return int(C.vrapi_GetSystemPropertyInt(cJava, C.ovrSystemProperty(parm)))
}

type OVRInitParms C.ovrInitParms // HMMM alias this type?

//type OVRModeParms C.ovrModeParms
// Just experiment with this?
type OVRModeParms struct {
	Type  OVRStructureType
	Flags OVRModeFlags
	Java  OVRJava
	//Padding       int32 // ??? // Add in build constraint for padding here?
	Display       uint64
	WindowSurface uint64
	ShareContext  uint64
}

type OVRJava C.ovrJava
type OVRMobile C.ovrMobile

type OVRTracking2 struct {
	Status  uint32 // Sensor status described by ovrTrackingStatus flags.
	Padding [4]byte

	// Predicted head configuration at the requested absolute time.
	// The pose describes the head orientation and center eye position.
	HeadPose OVRRigidBodyPosef // TODO
	Eye      [2]Tracking2Matrices
}

type Tracking2Matrices struct {
	ProjectionMatrix mgl.Mat4
	ViewMatrix       mgl.Mat4
}

func DefaultInitParms(java *OVRJava) OVRInitParms {
	cParms := C.vrapi_DefaultInitParms((*C.ovrJava)(java))
	return OVRInitParms(cParms)
}

func DefaultModeParms(java *OVRJava) OVRModeParms {
	cParms := C.vrapi_DefaultModeParms((*C.ovrJava)(java))
	return *(*OVRModeParms)(unsafe.Pointer(&cParms))
}

type OVRLayerHeader2 struct {
	Type       OVRLayerType2
	Flags      OVRFrameLayerFlags
	ColorScale mgl.Vec4
	SrcBlend   OVRFrameLayerBlend
	DstBlend   OVRFrameLayerBlend
	Reserved   unsafe.Pointer // VOID*?
}

type OVRLayerProjection2 struct {
	Header OVRLayerHeader2
	// TODO padding for 32 bit
	//Padding       int32 // ??? // Add in build constraint for padding here?

	HeadPose OVRRigidBodyPosef // TODO
	// We have to conver to, then convert back upon submission supposedly?

	Textures [FRAME_LAYER_EYE_MAX]EyeInformation
}

type OVRSubmitFrameDescription2 struct {
	Flags        uint32
	SwapInterval uint32
	FrameIndex   uint64
	DisplayTime  float64
	Pad          [8]byte // Unused
	LayerCount   uint32
	Layers       []*OVRLayerHeader2 // Only supports single layers for now

	//Layers       [1]*OVRLayerHeader2
	//Layers       []OVRLayerHeader2 // TODO when calling stuff pass a pointer to first element
}

type Context struct {
	workAvailable chan<- struct{}
	work          chan<- func()
	workDone      <-chan struct{}
}

func NewContext() (Context, Worker) {
	workAvailable := make(chan struct{}, 1)
	work := make(chan func(), 1)
	workDone := make(chan struct{})

	c := Context{
		workAvailable: workAvailable,
		work:          work,
		workDone:      workDone,
	}
	w := Worker{
		workAvailable: workAvailable,
		work:          work,
		workDone:      workDone,
	}

	return c, w
}

type Worker struct {
	workAvailable <-chan struct{}
	work          <-chan func()
	workDone      chan<- struct{}
}

func (w *Worker) WorkAvailable() <-chan struct{} {
	return w.workAvailable
}

func (w *Worker) DoWork() {
	// (<-w.work)() would be so much more confusing...
	fun := <-w.work
	fun()
	w.workDone <- struct{}{}
}

func (c *Context) EnterVrMode(modeParms *OVRModeParms) *OVRMobile {
	var ovr *OVRMobile
	c.work <- func() {
		cParms := (*C.ovrModeParms)(unsafe.Pointer(modeParms))
		cOVR := C.vrapi_EnterVrMode(cParms)
		ovr = (*OVRMobile)(cOVR)
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return ovr
}

func (c *Context) Initialize(parms *OVRInitParms) error {
	var err error
	c.work <- func() {
		status := C.vrapi_Initialize((*C.ovrInitParms)(parms))
		if status != INITIALIZE_SUCCESS {
			err = fmt.Errorf("vrapi_Initialize status %d not equal to sucess %d",
				status, INITIALIZE_SUCCESS)
		}
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return err
}

func (c *Context) CreateTextureSwapChain3(texType OVRTextureType, format int64,
	width, height, levels, bufferCount int) *OVRTextureSwapChain {

	var swapChain *OVRTextureSwapChain
	c.work <- func() {
		cSwapChain := C.vrapi_CreateTextureSwapChain3(
			C.ovrTextureType(texType), C.long(format),
			C.int(width), C.int(height), C.int(levels), C.int(bufferCount))
		swapChain = (*OVRTextureSwapChain)(unsafe.Pointer(cSwapChain))
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return swapChain
}

func (c *Context) GetTextureSwapChainLength(swapChain *OVRTextureSwapChain) int {
	var length int
	c.work <- func() {
		cSwapChain := (*C.ovrTextureSwapChain)(unsafe.Pointer(swapChain))
		length = int(C.vrapi_GetTextureSwapChainLength(cSwapChain))
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return length
}

func (c *Context) GetTextureSwapChainHandle(swapChain *OVRTextureSwapChain, i int) uint32 {
	var handle uint32
	c.work <- func() {
		cSwapChain := (*C.ovrTextureSwapChain)(unsafe.Pointer(swapChain))
		handle = uint32(C.vrapi_GetTextureSwapChainHandle(cSwapChain, C.int(i)))
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return handle
}

func (c *Context) SubmitFrame2(vrApp *OVRMobile, frameDesc *OVRSubmitFrameDescription2) error {
	var err error
	c.work <- func() {
		// TODO fix this constaint (allow multiple layers possibly using "varadic" C functions)
		if len(frameDesc.Layers) != 1 {
			err = fmt.Errorf("TODO not implmeneted layers must be size 1 for now passed in %+v",
				frameDesc.Layers)
			return
		}

		// Save layers and clear out layer field of frame description
		layers := frameDesc.Layers
		frameDesc.Layers = nil

		// Think we should be good with casting this memory? Even though Layers is
		// a slice in Go? We can just overwrite the memory with the Go pointer, and
		// the slice will always be larger than the pointer?
		cFrameDesc := (*C.ovrSubmitFrameDescription2)(unsafe.Pointer(frameDesc))

		cApp := (*C.ovrMobile)(unsafe.Pointer(vrApp))
		cLayer := *(*C.ovrLayerProjection2)(unsafe.Pointer(layers[0]))

		res := C.submit(cApp, cFrameDesc, cLayer)
		frameDesc.Layers = layers

		if res != OVRSuccess {
			err = fmt.Errorf("get current input state expected sucess (%d) got %d",
				OVRSuccess, res)
			return
		}
	}
	c.workAvailable <- struct{}{}
	<-c.workDone

	return err
}

type OVRTextureSwapChain C.ovrTextureSwapChain // TODO what is this type???

type EyeInformation struct {
	ColorSwapChain         *OVRTextureSwapChain
	SwapChainIndex         int32
	TexCoordsFromTanAngles mgl.Mat4
	TextureRect            OVRRectf
}

type OVRRigidBodyPosef struct {
	Pose                OVRPosef // TODO
	AngularVelocity     mgl.Vec3
	LinearVelocity      mgl.Vec3
	AngularAcceleration mgl.Vec3
	LinearAcceleration  mgl.Vec3

	TimeInSeconds       float64 //< Absolute time of this pose.
	PredictionInSeconds float64 //< Seconds this pose was predicted ahead.
}

type OVRRectf struct { // Make this a vec4? Or img rect???
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func DefaultLayerProjection2() OVRLayerProjection2 {
	cLayer := C.vrapi_DefaultLayerProjection2()
	layer := *(*OVRLayerProjection2)(unsafe.Pointer(&cLayer))

	return layer
}

func GetPredictedDisplayTime(vrApp *OVRMobile, frameIndex int64) float64 {
	return float64(C.vrapi_GetPredictedDisplayTime((*C.ovrMobile)(vrApp),
		C.longlong(frameIndex)))
}

func GetPredictedTracking2(vrApp *OVRMobile, displayTime float64) OVRTracking2 {
	cOVR := (*C.ovrMobile)(unsafe.Pointer(vrApp))
	cTracking := C.vrapi_GetPredictedTracking2(cOVR, C.double(displayTime))

	// TODO
	return *(*OVRTracking2)(unsafe.Pointer(&cTracking))
}

// Input (move to seperate file)

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

type OVRInputStateTrackedRemote struct {
	Header OVRInputStateHeader

	Buttons uint32 // Values for buttons described by ovrButton.
	// Finger contact status for trackpad
	// true = finger is on trackpad, false = finger is off trackpad
	TrackpadStatus uint32

	TrackpadPosition        mgl.Vec2 // X and Y coordinates of the Trackpad
	BatteryPercentRemaining uint8    // The percentage of max battery charge remaining.

	// Increments every time the remote is recentered. If this changes, the application may need
	// to adjust its arm model accordingly.
	RecenterCount uint8
	Reserved      uint16 // Reserved for future use.

	// Analog values from 0.0 - 1.0 of the pull of the triggers
	// added in API version 1.1.13.0
	IndexTrigger float32
	GripTrigger  float32

	// added in API version 1.1.15.0
	Touches    uint32
	Reserved5a uint32

	// Analog values from -1.0 - 1.0
	// The value is set to 0.0 on Joystick, if the magnitude of the vector is < 0.1f
	Joystick mgl.Vec2
	// JoystickNoDeadZone does change the raw values of the data.
	JoystickNoDeadZone mgl.Vec2
}

type OVRInputStateStandardPointer struct {
	Header           OVRInputStateHeader
	PointerPose      OVRPosef // to hamiltoned
	PointerStrength  float32
	GripPose         OVRPosef // to hamiltoned
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
	ControllerCapabilities OVRControllerCapabilities // Mask of controller capabilities described by ovrControllerCapabilities
	HapticSamplesMax       uint32                    // Maximum submittable samples for the haptics buffer
	HapticSampleDurationMS uint32                    // length in milliseconds of a sample in the haptics buffer.
	Reserved               [20]uint64                // Reserved for future use
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

	if inputState.ControllerType == OVRControllerType_StandardPointer {
		cPointer := (*OVRInputStateStandardPointer)(unsafe.Pointer(inputState))
		cPointer.GripPose.Orientation = jplToHamiltonQuats(cPointer.GripPose.Orientation)
		cPointer.PointerPose.Orientation = jplToHamiltonQuats(cPointer.PointerPose.Orientation)
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
