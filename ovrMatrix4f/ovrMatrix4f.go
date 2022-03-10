package ovrMatrix4f

/*
#cgo CPPFLAGS: -I../Include -I../usr/local/include
#cgo LDFLAGS: -v -march=armv8-a -shared -L../lib/arm64-v8a/ -lvrapi -landroid

#include <VrApi_Helpers.h>
*/
import "C"

import (
	"unsafe"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// TODO figure out how we want to do matrices...?
// We could just return
// https://pkg.go.dev/golang.org/x/mobile/exp/f32#Mat4.Translate
// or we could use mgl32?
// f32 might be good to test out the package
// but it has less features.
// and we have to consider if we should wrap / replace / or not include the
// matrix helpers provided.
func CreateTranslation(x, y, z float32) {
}

func TanAngleMatrixFromProjection(projection *mgl.Mat4) mgl.Mat4 {
	cProj := (*C.ovrMatrix4f)(unsafe.Pointer(projection))
	cMat := C.ovrMatrix4f_TanAngleMatrixFromProjection(cProj)
	return *(*mgl.Mat4)(unsafe.Pointer(&cMat))
}
