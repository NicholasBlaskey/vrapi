package vrapi

import "C"

type OVRInitParams C.ovrInitParms // HMMM alias this type?

func DefaultInitParms() OVRInitParms {
	return C.vrapi_DefaultInitParms()
}
