package main

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>

#define MAX(a, b) ((a) > (b) ? (a) : (b))
*/
import "C"

func main() {
	var dpy *C.Display
	var attr C.XWindowAttributes
	var start C.XButtonEvent
	var ev C.XEvent
	var ch *C.char

	if dpy != C.XOpenDisplay(ch) {
		return
	}

	C.XGrabKey(
		dpy,
		C.XKeysymToKeycode(dpy, C.XStringToKeysym(C.CString("F1"))),
		C.Mod1Mask,
		C.XDefaultRootWindow(dpy),
		1,
		C.GrabModeAsync,
		C.GrabModeAsync,
	)

	C.XGrabButton(
		dpy,
		1,
		C.Mod1Mask,
		C.XDefaultRootWindow(dpy),
		1,
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync,
		C.GrabModeAsync,
		C.None,
		C.None,
	)

	C.XGrabButton(
		dpy,
		3,
		C.Mod1Mask,
		C.XDefaultRootWindow(dpy),
		1,
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync,
		C.GrabModeAsync,
		C.None,
		C.None,
	)

	start.subwindow = C.None

	for {
		C.XNextEvent(dpy, &ev)
		if ev.Type == C.KeyPress && ev.xkey.subwindow != C.None {
			C.XRaiseWindow(dpy, ev.xkey.subwindow)
		} else if ev.Type == C.ButtonPress && ev.xbutton.subwindow != C.None {
			C.XGetWindowAttributes(dpy, ev.xbutton.subwindow, &attr)
			start = ev.xbutton
		} else if ev.Type == C.MotionNotify && start.subwindow != C.None {
			xdiff := ev.xbutton.x_root - start.x_root
			ydiff := ev.xbutton.y_root - start.y_root

			toDiffX := 0
			toDiffY := 0

			if start.button == 1 {
				toDiffX = xdiff
				toDiffY = ydiff
			}

			toWidth := 0
			toHeight := 0

			if start.button == 3 {
				toWidth = xdiff
				toHeight = ydiff
			}

			C.XMoveResizeWindow(
				dpy,
				start.subwindow,
				attr.x+toDiffX,
				attr.y+toDiffY,
				MAX(1, attr.width+toWidth),
				MAX(1, attr.height+toHeight))
		} else if ev.Type == ButtonRelease {
			start.subwindow = C.None
		}
	}
}
