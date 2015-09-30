package main

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>

#define MAX(a, b) ((a) > (b) ? (a) : (b))
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func unionToInt(cbytes [192]byte) (result int) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return *(*int)(unsafe.Pointer(uptr))
	}
	return 0
}

func unionToXKeyEvent(cbytes [192]byte) (result *C.XKeyEvent) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return (*C.XKeyEvent)(unsafe.Pointer(uptr))
	}
	return nil
}

func unionToXButtonEvent(cbytes [192]byte) (result *C.XButtonEvent) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return (*C.XButtonEvent)(unsafe.Pointer(uptr))
	}
	return nil
}

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
		C.int(C.XKeysymToKeycode(dpy, C.XStringToKeysym(C.CString("F1")))),
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
		if unionToInt(ev) == C.KeyPress && unionToXKeyEvent(ev).subwindow != C.None {
			C.XRaiseWindow(dpy, unionToXKeyEvent(ev).subwindow)
		} else if unionToInt(ev) == C.ButtonPress && unionToXButtonEvent(ev).subwindow != C.None {
			C.XGetWindowAttributes(dpy, unionToXButtonEvent(ev).subwindow, &attr)
			start = *unionToXButtonEvent(ev)
		} else if unionToInt(ev) == C.MotionNotify && start.subwindow != C.None {
			xdiff := unionToXButtonEvent(ev).x_root - start.x_root
			ydiff := unionToXButtonEvent(ev).y_root - start.y_root

			var toDiffX C.int
			var toDiffY C.int

			if start.button == 1 {
				toDiffX = xdiff
				toDiffY = ydiff
			}

			var toWidth C.int
			var toHeight C.int

			if start.button == 3 {
				toWidth = xdiff
				toHeight = ydiff
			}

			C.XMoveResizeWindow(
				dpy,
				start.subwindow,
				attr.x+toDiffX,
				attr.y+toDiffY,
				max(1, attr.width+toWidth),
				max(1, attr.height+toHeight))
		} else if unionToInt(ev) == C.ButtonRelease {
			start.subwindow = C.None
		}
	}
}

func max(a, b C.int) C.uint {
	if a > b {
		return C.uint(a)
	}

	return C.uint(b)
}
