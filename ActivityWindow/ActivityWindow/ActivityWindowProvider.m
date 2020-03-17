//
//  ActivityWindow.m
//  ActivityWindow
//
//  Created by Eugene Kulinich on 17.03.2020.
//  Copyright © 2020 Eugene Kulinich. All rights reserved.
//

#import "ActivityWindowProvider.h"

@implementation ActivityWindowProvider

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

@end
