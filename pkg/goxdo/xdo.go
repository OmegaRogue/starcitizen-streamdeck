package goxdo

// #include <xdo.h>
// #cgo LDFLAGS: -lxdo
import "C"

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	CURRENTWINDOW  = 0
	MBUTTON_LEFT   = 1
	MBUTTON_MIDDLE = 2
	MBUTTON_RIGHT  = 3
	MWHEELUP       = 4
	MWHEELDOWN     = 5
)

type Window int

type Xdo struct {
	xdo *C.xdo_t
}

func NewXdo() *Xdo {
	x := new(Xdo)
	x.xdo = C.xdo_new(nil)
	return x
}

func (t *Xdo) MoveMouse(x, y int, window Window) {
	C.xdo_move_mouse(t.xdo, C.int(x), C.int(y), C.int(window))
}

func (t *Xdo) MoveMouseRelativeToWindow(window Window, x, y int) {
	C.xdo_move_mouse_relative_to_window(t.xdo, C.Window(window), C.int(x), C.int(y))
}

func (t *Xdo) MoveMouseRelative(x, y int) {
	C.xdo_move_mouse_relative(t.xdo, C.int(x), C.int(y))
}

func (t *Xdo) MouseDown(window Window, button int) {
	C.xdo_mouse_down(t.xdo, C.Window(window), C.int(button))
}

func (t *Xdo) MouseUp(window Window, button int) {
	C.xdo_mouse_up(t.xdo, C.Window(window), C.int(button))
}

func (t *Xdo) GetMouseLocation() (x, y, screen int) {
	cX := C.int(0)
	cY := C.int(0)
	cScreen := C.int(0)

	C.xdo_get_mouse_location(t.xdo, &cX, &cY, &cScreen)
	return int(cX), int(cY), int(cScreen)
}

func (t *Xdo) GetWindowAtMouse() Window {
	var window C.Window
	C.xdo_get_window_at_mouse(t.xdo, &window)
	return Window(window)
}

func (t *Xdo) GetMouseLocation2() (x, y, screen int, window Window) {
	var cWindow C.Window
	cX := C.int(0)
	cY := C.int(0)
	cScreen := C.int(0)

	C.xdo_get_mouse_location2(t.xdo, &cX, &cY, &cScreen, &cWindow)
	return int(cX), int(cY), int(cScreen), Window(cWindow)
}

func (t *Xdo) WaitForMouseMoveFrom(x, y int) {
	C.xdo_wait_for_mouse_move_from(t.xdo, C.int(x), C.int(y))
}

func (t *Xdo) WaitForMouseMoveTo(x, y int) {
	C.xdo_wait_for_mouse_move_to(t.xdo, C.int(x), C.int(y))
}

func (t *Xdo) ClickWindow(window Window, button int) {
	C.xdo_click_window(t.xdo, C.Window(window), C.int(button))
}

func (t *Xdo) ClickWindowMultiple(window Window, button, repeat, useconds int) {
	C.xdo_click_window_multiple(t.xdo, C.Window(window), C.int(button), C.int(repeat),
		C.useconds_t(useconds))
}

func (t *Xdo) EnterTextWindow(window Window, text string, udelay int) {
	C.xdo_enter_text_window(t.xdo, C.Window(window), C.CString(text), C.useconds_t(udelay))
}

func (t *Xdo) SendKeysequenceWindow(window Window, sequence string, udelay int) {
	C.xdo_send_keysequence_window(t.xdo, C.Window(window), C.CString(sequence),
		C.useconds_t(udelay))
}

func (t *Xdo) SendKeysequenceWindowUp(window Window, sequence string, udelay int) {
	C.xdo_send_keysequence_window_up(t.xdo, C.Window(window), C.CString(sequence),
		C.useconds_t(udelay))
}

func (t *Xdo) SendKeysequenceWindowDown(window Window, sequence string, udelay int) {
	C.xdo_send_keysequence_window_down(t.xdo, C.Window(window), C.CString(sequence),
		C.useconds_t(udelay))
}

func (t *Xdo) Free() {
	C.xdo_free(t.xdo)
}
