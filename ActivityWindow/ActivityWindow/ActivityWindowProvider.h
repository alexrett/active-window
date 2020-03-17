//
//  ActivityWindow.h
//  ActivityWindow
//
//  Created by Eugene Kulinich on 17.03.2020.
//  Copyright © 2020 Eugene Kulinich. All rights reserved.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@protocol ActivityWindowProviderInput
 
- (BOOL)permissionCheck;
- (NSString*)getOwnerAndName;
 
@end

@interface ActivityWindowProvider : NSObject <ActivityWindowProviderInput>

@end

NS_ASSUME_NONNULL_END
