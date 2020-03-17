//
//  ActivityWindow.m
//  ActivityWindow
//
//  Created by Eugene Kulinich on 17.03.2020.
//  Copyright © 2020 Eugene Kulinich. All rights reserved.
//

#import "ActivityWindowProvider.h"

const char* getOwner() {
    ActivityWindowProvider *window = [[ActivityWindowProvider alloc] init];
    NSString* ownerName = [window getOwnerAndName];
    const char *result = [ownerName UTF8String];
    [window release];
    [ownerName release];
    return result;
}

@implementation ActivityWindowProvider

- (BOOL)permissionCheck {
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

- (NSString*)getOwnerAndName {
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
    return result;
}

char* getOwnerName(void) {
    return "YES";
}

@end
