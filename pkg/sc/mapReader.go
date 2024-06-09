package sc

import "strings"

type Rebind struct {
	Input        string `xml:"input,attr"`
	MultiTap     int    `xml:"multiTap,attr"`
	DefaultInput string `xml:"defaultInput,attr"`
}

type action struct {
	Name   string   `xml:"name,attr"`
	Rebind []Rebind `xml:"rebind"`
}

type actionMap struct {
	Name   string   `xml:"name,attr"`
	Action []action `xml:"action"`
}

type MappedAction struct {
	Name     string
	Keyboard []Rebind
	Mouse    []Rebind
	Joystick []Rebind
}

type DeviceOption struct {
	Input        string `xml:"input,attr"`
	Acceleration string `xml:"acceleration,attr"`
	Saturation   string `xml:"saturation,attr"`
	Deadzone     string `xml:"deadzone,attr"`
}

type deviceOptions struct {
	Name   string         `xml:"name,attr"`
	Option []DeviceOption `xml:"option"`
}

type Option struct {
	Type     string `xml:"type,attr"`
	Instance string `xml:"instance,attr"`
	Product  string `xml:"Product,attr"`
}

type ActionProfiles struct {
	Version        string          `xml:"Version,attr"`
	OptionsVersion string          `xml:"optionsVersion,attr"`
	RebindVersion  string          `xml:"rebindVersion,attr"`
	ProfileName    string          `xml:"profileName,attr"`
	DeviceOptions  []deviceOptions `xml:"deviceoptions"`
	Options        []Option        `xml:"options"`
	Modifiers      string          `xml:"modifiers"`
	Actionmap      []actionMap     `xml:"actionmap"`
}
type ActionMapActionMaps struct {
	ActionProfiles  ActionProfiles `xml:"ActionProfiles"`
	ActionMaps      map[string]map[string]*MappedAction
	DeviceOptionMap map[string]*[]DeviceOption
	Actions         map[string]*MappedAction
}

func (a *ActionMapActionMaps) Prepare() {
	a.ActionMaps = make(map[string]map[string]*MappedAction)
	a.DeviceOptionMap = make(map[string]*[]DeviceOption)
	a.Actions = make(map[string]*MappedAction)
	for _, aMap := range a.ActionProfiles.Actionmap {
		act := make(map[string]*MappedAction)
		for _, acti := range aMap.Action {
			reAct := MappedAction{
				Name:     acti.Name,
				Keyboard: []Rebind{},
				Mouse:    []Rebind{},
				Joystick: []Rebind{},
			}
			for _, rebind := range acti.Rebind {
				if strings.HasPrefix(rebind.Input, "kb") {
					reAct.Keyboard = append(reAct.Keyboard, rebind)
				} else if strings.HasPrefix(rebind.Input, "js") {
					reAct.Joystick = append(reAct.Joystick, rebind)
				} else if strings.HasPrefix(rebind.Input, "mo") {
					reAct.Mouse = append(reAct.Mouse, rebind)
				}
			}
			act[acti.Name] = &reAct
			a.Actions[acti.Name] = &reAct
		}
		a.ActionMaps[aMap.Name] = act
	}
	for _, option := range a.ActionProfiles.DeviceOptions {
		a.DeviceOptionMap[option.Name] = &option.Option
	}
}
