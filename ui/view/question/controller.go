package question

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/io"
	"time"
)

type Controller struct {
	window *Window

	backendBuilder backend.Builder
	backend        backend.Handle

	rawAnswer string
	printer   io.ResponsePrinter
}

func NewController(bb backend.Builder, printer io.ResponsePrinter) *Controller {
	return &Controller{
		backendBuilder: bb,
		printer:        printer,
	}
}

func (c *Controller) SetWindow(window *Window) {
	c.window = window
}

func (c *Controller) OnInputSubmitted(content string) {
	b, err := c.getBackend()
	if err != nil {
		c.setOutput(content, fmt.Sprintf("Error: %s", err))
		return
	}

	c.window.inText.Disable()
	c.window.btnStop.Show()

	go func() {
		defer c.window.inText.Enable()
		defer c.window.btnStop.Hide()

		result, err := b.AskSomething(content)
		if err != nil {
			c.setOutput(content, fmt.Sprintf("Error: %s", err))
		} else {
			c.setOutput(content, result)
		}
	}()
}

func (c *Controller) getBackend() (backend.Handle, error) {
	if c.backend == nil {
		var err error
		c.backend, err = c.backendBuilder.Build()
		if err != nil {
			return nil, err
		}
	}

	return c.backend, nil
}

func (c *Controller) setOutput(input, output string) {
	c.rawAnswer = output

	c.printer.Print(input, output)

	if c.backendBuilder.Type == backend.TypeSingleShot {
		c.window.mdOutput.ParseMarkdown(output)
		c.window.rawOutput.SetText(output)
	} else if c.backendBuilder.Type == backend.TypeMultiShot {
		c.window.mdOutput.AppendMarkdown("# " + input)
		c.window.mdOutput.AppendMarkdown(output)
		c.window.rawOutput.SetText(c.window.rawOutput.Text + "#" + input + "\n" + output + "\n")
		c.window.inText.SetText("")
	} else {
		panic("Unknown backend type!")
	}

	if output == "" {
		c.window.mdOutputScroll.Hide()
		c.window.rawOutputScroll.Hide()
		c.window.btnClipboard.Hide()
		c.window.tabs.Hide()
	} else {
		c.window.mdOutputScroll.Show()
		c.window.rawOutputScroll.Show()
		c.window.btnClipboard.Show()
		c.window.tabs.Show()

		time.AfterFunc(250*time.Millisecond, func() {
			c.window.mdOutputScroll.ScrollToBottom()
			c.window.rawOutputScroll.ScrollToBottom()
		})
	}
}

func (c *Controller) OnClickStop() {
	if c.backend != nil {
		c.backend.Close()
	}
}

func (c *Controller) OnClickClipboard() {
	fyneClipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	fyneClipboard.SetContent(c.rawAnswer)
}

func (c *Controller) Close() error {
	if c.backend != nil {
		return c.backend.Close()
	}
	return nil
}
