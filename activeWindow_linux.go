//+build linux

package activeWindow

/*
#cgo LDFLAGS: -lXss -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdio.h>
#include <stdlib.h>

#define MAXSTR 1000

Display *display;
unsigned long window;
unsigned char *prop;

void check_status(int status, unsigned long window)
{
    if (status == BadWindow) {
        printf("window id # 0x%lx does not exists!", window);
        exit(1);
    }

    if (status != Success) {
        printf("XGetWindowProperty failed!");
        exit(2);
    }
}

unsigned char* get_string_property(char* property_name)
{
    Atom actual_type, filter_atom;
    int actual_format, status;
    unsigned long nitems, bytes_after;

    filter_atom = XInternAtom(display, property_name, True);
    status = XGetWindowProperty(display, window, filter_atom, 0, MAXSTR, False, AnyPropertyType,
                                &actual_type, &actual_format, &nitems, &bytes_after, &prop);
    check_status(status, window);
    return prop;
}

unsigned long get_long_property(char* property_name)
{
    get_string_property(property_name);
    unsigned long long_property = prop[0] + (prop[1]<<8) + (prop[2]<<16) + (prop[3]<<24);
    return long_property;
}

char *owner = NULL;
char *title = NULL;

void getActiveWindowTitle()
{
    char *display_name = NULL;  // could be the value of $DISPLAY

    display = XOpenDisplay(display_name);
    if (display == NULL) {
        fprintf (stderr, "unable to open display '%s'\n", XDisplayName (display_name));
    }
    int screen = XDefaultScreen(display);
    window = RootWindow(display, screen);

    window = get_long_property("_NET_ACTIVE_WINDOW");

    title = get_string_property("WM_CLASS");
    owner = get_string_property("_NET_WM_NAME");

    XCloseDisplay(display);
}
*/
import "C"

// in the code above seems like has memory leaks and not properly closed XCloseDisplay
// i tested it in virtualbox... maybe there is no problem on real device
// one of possible solution is transfer some code from xdotool, but i feel it like not really needed
func (a *ActiveWindow) getActiveWindowTitle() (string, string) {
	C.getActiveWindowTitle(&title, &owner)

	return C.GoString(C.owner), C.GoString(C.title)
}
