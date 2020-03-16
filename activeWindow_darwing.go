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
	"reflect"
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

func CFStringToString(cfStr C.CFStringRef) string {
	// NB: don't use CFStringGetCStringPtr() because it will stop at the first NUL
	length := C.CFStringGetLength(cfStr)
	if length == 0 {
		// short-cut for empty strings
		return ""
	}
	cfRange := C.CFRange{0, length}
	enc := C.CFStringEncoding(C.kCFStringEncodingUTF8)
	// first find the buffer size necessary
	var usedBufLen C.CFIndex
	if C.CFStringGetBytes(cfStr, cfRange, enc, 0, C.false, nil, 0, &usedBufLen) > 0 {
		bytes := make([]byte, usedBufLen)
		buffer := (*C.UInt8)(unsafe.Pointer(&bytes[0]))
		if C.CFStringGetBytes(cfStr, cfRange, enc, 0, C.false, buffer, usedBufLen, nil) > 0 {
			// bytes is now filled up
			// convert it to a string
			header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
			strHeader := &reflect.StringHeader{
				Data: header.Data,
				Len:  header.Len,
			}
			return *(*string)(unsafe.Pointer(strHeader))
		}
	}

	// we failed to convert, for some reason. Too bad there's no nil string
	return ""
}
