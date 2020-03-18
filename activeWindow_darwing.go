//+build darwin

package activeWindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework Foundation
#import <Foundation/Foundation.h>
#import <CoreFoundation/CoreFoundation.h>
#import <CoreGraphics/CoreGraphics.h>
void getOwnerAndName() {
    CFArrayRef windowArray = CGWindowListCopyWindowInfo(kCGWindowListOptionOnScreenOnly, kCGNullWindowID);
    NSArray*  windows = (NSArray*)windowArray;
    NSString *result = nil;
    for (NSDictionary *window in windows) {
        NSString *windowLayer = [window objectForKey:@"kCGWindowLayer"];
        if (windowLayer.intValue == 0) {
            NSString *owner = [window objectForKey:@"kCGWindowOwnerName"];
            NSString *name = [window objectForKey:@"kCGWindowName"];
            result = [[NSString alloc] initWithFormat: @"%@+%@", owner, name];
            break;
        }
    }
    if (result == nil) {
        result = @"empty";
    }
    [windows release];
    NSLog(@"%@", result);
}

bool permissionCheck() {
    if (@available(macOS 10.15, *)) {
        CFArrayRef windowList = CGWindowListCopyWindowInfo(kCGWindowListOptionOnScreenOnly, kCGNullWindowID);
        NSUInteger numberOfWindows = CFArrayGetCount(windowList);
        NSUInteger numberOfWindowsWithInfoGet = 0;
        for (int idx = 0; idx < numberOfWindows; idx++) {

            NSDictionary *windowInfo = (NSDictionary *)CFArrayGetValueAtIndex(windowList, idx);
            NSString *windowName = windowInfo[(id)kCGWindowName];
            NSNumber* sharingType = windowInfo[(id)kCGWindowSharingState];

            if (windowName || kCGWindowSharingNone != sharingType.intValue) {
                numberOfWindowsWithInfoGet++;
            } else {
                NSNumber* pid = windowInfo[(id)kCGWindowOwnerPID];
                NSString* appName = windowInfo[(id)kCGWindowOwnerName];
                NSLog(@"windowInfo get Fail pid:%lu appName:%@", pid.integerValue, appName);
            }
        }
        CFRelease(windowList);
        if (numberOfWindows == numberOfWindowsWithInfoGet) {
            return YES;
        } else {
            return NO;
        }
    }
    return YES;
}
*/
import "C"

var title string
var owner string

func (a *ActiveWindow) getActiveWindowTitle() (string, string) {
	C.getOwnerAndName()
	return owner, title
}