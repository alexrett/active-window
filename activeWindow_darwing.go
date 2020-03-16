//+build darwin

package activeWindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics  -framework CoreFoundation

#import <CoreFoundation/CoreFoundation.h>
#import <CoreGraphics/CoreGraphics.h>
*/
import "C"
import (
	"unsafe"
)

var title string
var owner string

// there is memory leaks... I can't find where it is...
// todo: find and fix memory leaks on Darwin method getActiveWindowTitle()
func (a *ActiveWindow) getActiveWindowTitle() (string, string) {
	var (
		dictionary  C.CFDictionaryRef
		windowList  C.CFArrayRef
		layer       C.int
		cfiI        C.CFIndex
		cfiLen      C.CFIndex
		windowLayer C.CFNumberRef
	)

	windowList = C.CGWindowListCopyWindowInfo(C.kCGWindowListOptionOnScreenOnly, C.kCGNullWindowID)
	cfiLen = C.CFArrayGetCount(windowList)
	for cfiI = 0; cfiI < cfiLen; cfiI++ {
		dictionary = C.CFDictionaryRef(C.CFArrayGetValueAtIndex(windowList, cfiI))
		windowLayer = C.CFNumberRef(C.CFDictionaryGetValue(dictionary, unsafe.Pointer(C.kCGWindowLayer)))
		C.CFNumberGetValue(windowLayer, C.kCFNumberIntType, unsafe.Pointer(&layer))
		if layer == 0 {
			var tmpTitle C.CFStringRef
			var tmpOwner C.CFStringRef
			tmpOwner = C.CFStringRef(C.CFDictionaryGetValue(dictionary, unsafe.Pointer(C.kCGWindowOwnerName)))
			tmpTitle = C.CFStringRef(C.CFDictionaryGetValue(dictionary, unsafe.Pointer(C.kCGWindowName)))
			owner = CFStringToString(tmpOwner)
			title = CFStringToString(tmpTitle)
			C.CFRelease(C.CFTypeRef(tmpOwner))
			C.CFRelease(C.CFTypeRef(tmpTitle))
			break
		}
	}

	defer C.free(unsafe.Pointer(dictionary))
	defer C.CFRetain(C.CFTypeRef(windowList))
	defer C.CFRelease(C.CFTypeRef(windowLayer))

	return owner, title
}

// CFStringToString converts a CFStringRef to a string.
func CFStringToString(s C.CFStringRef) string {
	p := C.CFStringGetCStringPtr(s, C.kCFStringEncodingUTF8)
	if p != nil {
		return C.GoString(p)
	}
	length := C.CFStringGetLength(s)
	if length == 0 {
		return ""
	}
	maxBufLen := C.CFStringGetMaximumSizeForEncoding(length, C.kCFStringEncodingUTF8)
	if maxBufLen == 0 {
		return ""
	}
	buf := make([]byte, maxBufLen)
	var usedBufLen C.CFIndex
	_ = C.CFStringGetBytes(s, C.CFRange{0, length}, C.kCFStringEncodingUTF8, C.UInt8(0), C.false, (*C.UInt8)(&buf[0]), maxBufLen, &usedBufLen)
	return string(buf[:usedBufLen])
}
