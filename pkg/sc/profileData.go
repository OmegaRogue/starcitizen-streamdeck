package sc

import (
	"encoding/xml"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ActivationMode struct {
	Name                    string  `xml:"name,attr"`
	OnPress                 bool    `xml:"onPress,attr"`
	OnHold                  bool    `xml:"onHold,attr"`
	OnRelease               bool    `xml:"onRelease,attr"`
	MultiTap                int     `xml:"multiTap,attr"`
	MultiTapBlock           bool    `xml:"multiTapBlock,attr"`
	PressTriggerThreshold   float64 `xml:"pressTriggerThreshold,attr"`
	ReleaseTriggerThreshold float64 `xml:"releaseTriggerThreshold,attr"`
	ReleaseTriggerDelay     float64 `xml:"releaseTriggerDelay,attr"`
	Retriggerable           bool    `xml:"retriggerable,attr"`
	Always                  bool    `xml:"always,attr,omitempty"`
	NoModifiers             bool    `xml:"noModifiers,attr,omitempty"`
	HoldTriggerDelay        float64 `xml:"holdTriggerDelay,attr,omitempty"`
}

type Gamepad struct {
	ActivationMode *ActivationMode
	Input          string
	OnHold         bool
	OnPress        bool
	OnRelease      bool
	NoModifiers    bool
	Inputdata      []string
}

type Action struct {
	Name             string
	ActivationMode   *ActivationMode
	AttrKeyboard     string
	AttrJoystick     string
	UILabel          string
	UIDescription    string
	AttrGamepad      string
	Category         string
	OnPress          bool
	OnHold           bool
	OnRelease        bool
	OptionGroup      string
	UICategory       string
	AttrMouse        string
	Always           bool
	Retriggerable    bool
	HoldTriggerDelay float64
	HoldRepeatDelay  float64
	UseAnalogCompare bool
	AnalogCompareVal float64
	AnalogCompareOp  CompareOperation
	NoModifiers      bool
	Gamepad          Gamepad
	States           map[string]string
	Joystick         Joystick
	Keyboard         Keyboard
	Mouse            []string
}

func (a Action) HasKeyboard() bool {
	return strings.TrimSpace(a.AttrKeyboard) != "" || strings.TrimSpace(a.Keyboard.Input) != "" || len(a.Keyboard.Inputdata) > 0
}
func (a Action) HasGamepad() bool {
	return strings.TrimSpace(a.AttrGamepad) != "" || strings.TrimSpace(a.Gamepad.Input) != "" || len(a.Gamepad.Inputdata) > 0
}
func (a Action) HasJoystick() bool {
	return strings.TrimSpace(a.AttrJoystick) != "" || strings.TrimSpace(a.Joystick.Input) != ""
}
func (a Action) HasMouse() bool {
	return strings.TrimSpace(a.AttrMouse) != "" || len(a.Mouse) > 0
}

type Keyboard struct {
	ActivationMode *ActivationMode
	Input          string
	Inputdata      []string
}

type Joystick struct {
	ActivationMode *ActivationMode `xml:"activationMode,attr"`
	Input          string          `xml:"input,attr"`
}

type ActionMap struct {
	Name          string             `xml:"name,attr"`
	UILabel       string             `xml:"UILabel,attr"`
	UICategory    string             `xml:"UICategory,attr"`
	UIDescription string             `xml:"UIDescription,attr"`
	Action        map[string]*Action `xml:"action"`
}

type Profile struct {
	ActivationModes map[string]*ActivationMode
	Maps            map[string]ActionMap
	Actions         map[string]*Action
}

func FromXml(data string) (*Profile, error) {
	data = strings.ReplaceAll(data, "activationmode", "activationMode")
	data = strings.ReplaceAll(data, "UILable", "UILabel")
	var profile profileReaderTemp
	err := xml.Unmarshal([]byte(data), &profile)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal XML")
	}

	return NewProfile(profile), nil
}

func NewProfile(r profileReaderTemp) *Profile {
	p := new(Profile)
	p.ActivationModes = map[string]*ActivationMode{}
	p.Maps = map[string]ActionMap{}
	p.Actions = map[string]*Action{}
	for _, mode := range r.ActivationModes {
		p.ActivationModes[mode.Name] = &mode
	}
	for _, actionMap := range r.ActionMaps {
		m := ActionMap{
			Name:          actionMap.Name,
			UILabel:       actionMap.UILabel,
			UICategory:    actionMap.UICategory,
			UIDescription: actionMap.UIDescription,
			Action:        map[string]*Action{},
		}
		for _, action2 := range actionMap.Action {
			act := Action{
				Name:             action2.Name,
				ActivationMode:   p.ActivationModes[action2.ActivationMode],
				AttrKeyboard:     action2.AttrKeyboard,
				AttrJoystick:     action2.AttrJoystick,
				UILabel:          action2.UILabel,
				UIDescription:    action2.UIDescription,
				AttrGamepad:      action2.AttrGamepad,
				Category:         action2.Category,
				OnPress:          action2.OnPress,
				OnHold:           action2.OnHold,
				OnRelease:        action2.OnRelease,
				OptionGroup:      action2.OptionGroup,
				UICategory:       action2.UICategory,
				AttrMouse:        action2.AttrMouse,
				Always:           action2.Always,
				Retriggerable:    action2.Retriggerable,
				HoldTriggerDelay: action2.HoldTriggerDelay,
				HoldRepeatDelay:  action2.HoldRepeatDelay,
				UseAnalogCompare: action2.UseAnalogCompare,
				AnalogCompareVal: action2.AnalogCompareVal,
				AnalogCompareOp:  action2.AnalogCompareOp,
				NoModifiers:      action2.NoModifiers,
				Gamepad: Gamepad{
					ActivationMode: p.ActivationModes[action2.Gamepad.ActivationMode],
					Input:          action2.Gamepad.Input,
					OnHold:         action2.Gamepad.OnHold,
					OnPress:        action2.Gamepad.OnPress,
					OnRelease:      action2.Gamepad.OnRelease,
					NoModifiers:    action2.Gamepad.NoModifiers,
					Inputdata:      []string{},
				},
				States: map[string]string{},
				Joystick: Joystick{
					ActivationMode: p.ActivationModes[action2.Joystick.ActivationMode],
					Input:          action2.Joystick.Input,
				},
				Keyboard: Keyboard{
					ActivationMode: p.ActivationModes[action2.Keyboard.ActivationMode],
					Input:          action2.Keyboard.Input,
					Inputdata:      []string{},
				},
				Mouse: []string{},
			}
			for _, inputdatum := range action2.Gamepad.Inputdata {
				act.Gamepad.Inputdata = append(act.Gamepad.Inputdata, inputdatum.Input)
			}
			for _, inputdatum := range action2.Keyboard.Inputdata {
				act.Keyboard.Inputdata = append(act.Keyboard.Inputdata, inputdatum.Input)
			}
			for _, inputdatum := range action2.Mouse {
				act.Mouse = append(act.Mouse, inputdatum.Input)
			}
			for _, state := range action2.States {
				act.States[state.Name] = state.UILabel
			}

			m.Action[action2.Name] = &act
			p.Actions[action2.Name] = &act
			p.Maps[actionMap.Name] = m

		}
	}
	return p
}

type Data struct {
	Profile *Profile
	Rebinds *ActionMapActionMaps
	Locale  map[string]map[string]string
}

func (c Data) LookupBind(name string) string {
	action := c.Profile.Actions[name]
	reActions := c.Rebinds.Actions[name]
	var possibilities []string
	possibilities = append(possibilities, action.AttrKeyboard)
	possibilities = append(possibilities, action.Keyboard.Input)
	possibilities = append(possibilities, action.Keyboard.Inputdata...)
	if reActions != nil {
		for _, rebind := range reActions.Keyboard {
			possibilities = append(possibilities, strings.TrimPrefix(rebind.DefaultInput, "kb1_"))
			possibilities = append(possibilities, strings.TrimPrefix(rebind.Input, "kb1_"))
		}
	}

	possibilities = lo.Filter(possibilities, func(item string, index int) bool {
		return strings.TrimSpace(item) != "" && !strings.HasSuffix(strings.TrimSpace(item), "_")
	})
	bind, err := lo.Last(possibilities)
	if err != nil {
		return ""
	}
	return bind
}
