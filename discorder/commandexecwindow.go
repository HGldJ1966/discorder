package discorder

import (
	"github.com/jonas747/discorder/common"
	"github.com/jonas747/discorder/ui"
	"strconv"
	"strings"
)

type CommandExecWindow struct {
	*ui.BaseEntity
	app        *App
	layer      int
	menuWindow *ui.MenuWindow
	command    Command
}

type CustomMenuType int

const (
	CustomMenuExecute CustomMenuType = iota
)

func NewCommandExecWindow(layer int, app *App, command Command) *CommandExecWindow {
	execWindow := &CommandExecWindow{
		BaseEntity: &ui.BaseEntity{},
		app:        app,
		menuWindow: ui.NewMenuWindow(layer, app.ViewManager.UIManager, false),
		command:    command,
		layer:      layer,
	}

	execWindow.menuWindow.Transform.AnchorMax = common.NewVector2F(1, 1)
	execWindow.menuWindow.Transform.Top = 1
	execWindow.menuWindow.Transform.Bottom = 2

	execWindow.menuWindow.Window.Title = "Execute command"
	execWindow.menuWindow.Window.Footer = ":)"

	app.ApplyThemeToMenu(execWindow.menuWindow)

	execWindow.Transform.AddChildren(execWindow.menuWindow)

	execWindow.Transform.AnchorMax = common.NewVector2F(1, 1)

	execWindow.Transform.Right = 2
	execWindow.Transform.Left = 1

	app.ViewManager.UIManager.AddWindow(execWindow)

	execWindow.GenMenu()

	return execWindow
}

func (cew *CommandExecWindow) Destroy() {
	cew.app.ViewManager.UIManager.RemoveWindow(cew)
	cew.DestroyChildren()
}

func (cew *CommandExecWindow) GenMenu() {
	items := make([]*ui.MenuItem, 0)
	for _, arg := range cew.command.GetArgs() {
		helper := &ui.MenuItem{
			Name:       arg.Name,
			Info:       arg.Description,
			Decorative: true,
		}
		input := &ui.MenuItem{
			Name:      arg.Name,
			Info:      arg.Description,
			IsInput:   true,
			InputType: arg.Datatype,
			UserData:  arg,
		}
		if arg.CurVal != nil {
			input.InputDefaultText = arg.CurVal(cew.app)
		}

		items = append(items, helper, input)
	}

	exec := &ui.MenuItem{
		Name:     "Execute",
		Info:     "Execute the commadn with specified args",
		UserData: CustomMenuExecute,
	}
	items = append(items, exec)
	cew.menuWindow.SetOptions(items)
}

func (cew *CommandExecWindow) Select() {
	element := cew.menuWindow.GetHighlighted()
	if element == nil {
		return
	}

	if element.IsCategory {
		cew.menuWindow.Select()
		return
	}

	if element.UserData == nil {
		return
	}

	switch t := element.UserData.(type) {
	case CustomMenuType:
		switch t {
		case CustomMenuExecute:
			cew.Execute()
		}
	// Run a argument helper if any
	case *ArgumentDef:
		if t.Helper != nil {
			t.Helper.Run(cew.app, cew.layer+2, func(result string) {
				if element.Input != nil {
					element.Input.TextBuffer = result
					element.Input.CursorLocation = 0
				}
			})
		}
	}
}

func (cew *CommandExecWindow) Execute() {
	args := make(map[string]interface{})
	for _, item := range cew.menuWindow.Options {
		if !item.IsInput {
			continue
		}
		buf := item.Input.TextBuffer
		args[item.Name] = ParseArgumentString(buf, item.InputType)
	}

	cew.app.RunCommand(cew.command, Arguments(args))
	cew.Transform.Parent.RemoveChild(cew, true)
}

func ParseArgumentString(arg string, dataType ui.DataType) interface{} {
	switch dataType {
	case ui.DataTypeBool:
		lowerBuf := strings.ToLower(arg)
		b, _ := strconv.ParseBool(lowerBuf)
		return b
	case ui.DataTypeInt:
		i, _ := strconv.ParseInt(arg, 10, 64)
		return i
	case ui.DataTypeFloat:
		f, _ := strconv.ParseFloat(arg, 64)
		return f
	}

	return arg
}
