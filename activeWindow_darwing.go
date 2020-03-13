//+build darwin

package activeWindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics
#include <Foundation/Foundation.h>
#include <CoreGraphics/CoreGraphics.h>

void getActiveWindowTitle(CFStringRef *title, CFStringRef *owner) {
	CFArrayRef windowList = CGWindowListCopyWindowInfo(kCGWindowListOptionOnScreenOnly, kCGNullWindowID);
    CFIndex cfiLen = CFArrayGetCount(windowList);
	CFDictionaryRef dictionary;
	for (CFIndex cfiI = 0; cfiI < cfiLen; cfiI++){
		dictionary = (CFDictionaryRef) CFArrayGetValueAtIndex(windowList, cfiI);
		CFStringRef owner2 = (CFStringRef) CFDictionaryGetValue(dictionary,kCGWindowOwnerName);
        CFStringRef title2 = (CFStringRef) CFDictionaryGetValue(dictionary,kCGWindowName);
        CFNumberRef window_layer = (CFNumberRef) CFDictionaryGetValue(dictionary, kCGWindowLayer);
        int layer;
        CFNumberGetValue(window_layer, kCFNumberIntType, &layer);
        if (layer==0){
			*title = title2;
			*owner = owner2;
            break;
        }
	}
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

var title C.CFStringRef
var owner C.CFStringRef

func (a *ActiveWindow) getActiveWindowTitle() (string, string) {
	C.getActiveWindowTitle(&title, &owner)

	return convertCFStringToString(owner), convertCFStringToString(title)
}

// get from https://github.com/lilyball/go-osx-plist/blob/a0f875443af46a7c67f72a2fe1fd245e966a77c2/convert.go#L156
func convertCFStringToString(cfStr C.CFStringRef) string {
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
