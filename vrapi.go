//go:build darwin || linux || windows
// +build android, gldebug
package vrapi

/*
#cgo CPPFLAGS: -I./Include -I/usr/local/include
#cgo LDFLAGS: -v -march=armv8-a -shared -L./lib/arm64-v8a/ -lvrapi -landroid

#include <VrApi.h>
#include <VrApi_Helpers.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	INITIALIZE_SUCCESS = C.VRAPI_INITIALIZE_SUCCESS
)

type OVRInitParms C.ovrInitParms // HMMM alias this type?

type OVRStructureType int32

const (
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

/*
func (p *OVRModeParms) fromC() {

}

func (p *OVRModeParms) toC() *C.ovrModeParms {

}
*/

type OVRJava C.ovrJava

//type OVRMobile C.ovrMobile

func DefaultInitParms(java *OVRJava) OVRInitParms {
	cParms := C.vrapi_DefaultInitParms((*C.ovrJava)(java))
	return OVRInitParms(cParms)
}

func DefaultModeParms(java *OVRJava) OVRModeParms {
	cParms := C.vrapi_DefaultModeParms((*C.ovrJava)(java))
	return *(*OVRModeParms)(unsafe.Pointer(&cParms))
}

// This should run with vrctx like glctx on a seperate thread.
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

/*
func GetPredictedDisplayTime(vrApp *OVRMobile, frameIndex int64) float32 {
	return float32(C.vrapi_GetPredictedDisplayTime((*C.ovrMobile)(vrApp),
		C.longlong(frameIndex)))
}
*/

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
