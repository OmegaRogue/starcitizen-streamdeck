package sc

type ActionCategory struct {
	Name       string           `xml:"name,attr"`
	Action     string           `xml:"action,attr"`
	UILabel    string           `xml:"UILabel,attr"`
	UIIcon     string           `xml:"UIIcon,attr"`
	ShowEmpty  bool             `xml:"showEmpty,attr,omitempty"`
	Categories []ActionCategory `xml:"Category"`
}
type profileReaderTemp struct {
	ActivationModes  []ActivationMode `xml:"ActivationModes>ActivationMode"`
	ActionCategories []ActionCategory `xml:"ActionCategories>Category"`
	ActionMaps       []actionMapRaw   `xml:"actionmap"`
}

type gamepadRaw struct {
	ActivationMode string      `xml:"activationMode,attr"`
	Input          string      `xml:"input,attr"`
	OnHold         bool        `xml:"onHold,attr"`
	OnPress        bool        `xml:"onPress,attr"`
	OnRelease      bool        `xml:"onRelease,attr"`
	NoModifiers    bool        `xml:"noModifiers,attr"`
	Inputdata      []inputData `xml:"inputdata"`
}

type joystickRaw struct {
	ActivationMode string `xml:"activationMode,attr"`
	Input          string `xml:"input,attr"`
}

type inputData struct {
	Input string `xml:"input,attr"`
}

type keyboardRaw struct {
	ActivationMode string      `xml:"activationMode,attr"`
	Input          string      `xml:"input,attr"`
	Inputdata      []inputData `xml:"inputdata"`
}
type actionRaw struct {
	Name             string           `xml:"name,attr"`
	ActivationMode   string           `xml:"activationMode,attr"`
	AttrKeyboard     string           `xml:"keyboard,attr"`
	AttrJoystick     string           `xml:"joystick,attr"`
	UILabel          string           `xml:"UILabel,attr"`
	UIDescription    string           `xml:"UIDescription,attr"`
	AttrGamepad      string           `xml:"gamepad,attr"`
	Category         string           `xml:"Category,attr"`
	OnPress          bool             `xml:"onPress,attr"`
	OnHold           bool             `xml:"onHold,attr"`
	OnRelease        bool             `xml:"onRelease,attr"`
	OptionGroup      string           `xml:"optionGroup,attr"`
	UICategory       string           `xml:"UICategory,attr"`
	AttrMouse        string           `xml:"mouse,attr"`
	Always           bool             `xml:"always,attr"`
	Retriggerable    bool             `xml:"retriggerable,attr"`
	HoldTriggerDelay float64          `xml:"holdTriggerDelay,attr"`
	HoldRepeatDelay  float64          `xml:"holdRepeatDelay,attr"`
	UseAnalogCompare bool             `xml:"useAnalogCompare,attr"`
	AnalogCompareVal float64          `xml:"analogCompareVal,attr"`
	AnalogCompareOp  CompareOperation `xml:"analogCompareOp,attr"`
	Activationmode   string           `xml:"activationmode,attr"`
	NoModifiers      bool             `xml:"noModifiers,attr"`
	Gamepad          gamepadRaw       `xml:"gamepad"`
	States           []state          `xml:"states>state"`
	Joystick         joystickRaw      `xml:"joystick"`
	Keyboard         keyboardRaw      `xml:"keyboard"`
	Mouse            []inputData      `xml:"mouse>inputdata"`
}

type state struct {
	Name    string `xml:"name,attr"`
	UILabel string `xml:"UILabel,attr"`
}

type actionMapRaw struct {
	Name          string      `xml:"name,attr"`
	UILabel       string      `xml:"UILabel,attr"`
	UICategory    string      `xml:"UICategory,attr"`
	UILable       string      `xml:"UILable,attr"`
	UIDescription string      `xml:"UIDescription,attr"`
	Action        []actionRaw `xml:"action"`
}
