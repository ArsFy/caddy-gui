package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const maxLogLines = 500

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Caddy GUI")

	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)

	configEntry := widget.NewMultiLineEntry()
	configEntry.SetPlaceHolder("Enter Caddy configuration")
	configEntry.SetMinRowsVisible(15)

	content, err := os.ReadFile("CaddyFile")
	if err != nil {
		configEntry.SetText(`:80 {
    log
    
    reverse_proxy /* {
        to http://localhost:3000
    }

    reverse_proxy /api/* {
        to http://localhost:8081
    }
}`)
	} else {
		configEntry.SetText(string(content))
	}

	outputLabel := widget.NewLabel("")
	outputLabel.Wrapping = fyne.TextWrapWord
	outputScroll := container.NewVScroll(outputLabel)
	outputScroll.SetMinSize(fyne.NewSize(600, 250))

	var cmd *exec.Cmd
	var running bool
	var mu sync.Mutex

	buttonTextChan := make(chan string)

	runButton := widget.NewButton("Run", func() {
		mu.Lock()
		defer mu.Unlock()

		if running {
			if cmd != nil && cmd.Process != nil {
				cmd.Process.Kill()
			}
			buttonTextChan <- "Run"
			running = false
			return
		}

		err := os.WriteFile("CaddyFile", []byte(configEntry.Text), 0644)
		if err != nil {
			outputLabel.SetText(fmt.Sprintf("Error writing CaddyFile: %v", err))
			return
		}

		execPath, err := os.Executable()
		if err != nil {
			outputLabel.SetText(fmt.Sprintf("Error finding executable path: %v", err))
			return
		}
		execDir := filepath.Dir(execPath)
		caddyPath := filepath.Join(execDir, "caddy")
		cmd = exec.Command(caddyPath, "run", "-c", "CaddyFile")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			outputLabel.SetText(fmt.Sprintf("Error creating StdoutPipe: %v", err))
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			outputLabel.SetText(fmt.Sprintf("Error creating StderrPipe: %v", err))
			return
		}

		if err := cmd.Start(); err != nil {
			outputLabel.SetText(fmt.Sprintf("Error starting command: %v", err))
			return
		}

		buttonTextChan <- "Stop"
		running = true

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				appendLog(outputLabel, scanner.Text())
				outputScroll.ScrollToBottom()
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				appendLog(outputLabel, scanner.Text())
				outputScroll.ScrollToBottom()
			}
		}()

		go func() {
			if err := cmd.Wait(); err != nil {
				appendLog(outputLabel, fmt.Sprintf("Command finished with error: %v", err))
			}
			mu.Lock()
			buttonTextChan <- "Run"
			running = false
			mu.Unlock()
		}()
	})

	go func() {
		for text := range buttonTextChan {
			runButton.SetText(text)
		}
	}()

	myWindow.SetContent(container.NewVBox(
		configEntry,
		runButton,
		container.NewStack(outputScroll),
	))

	myWindow.ShowAndRun()
}

func appendLog(outputLabel *widget.Label, text string) {
	lines := strings.Split(outputLabel.Text, "\n")
	if len(lines) >= maxLogLines {
		lines = lines[len(lines)-maxLogLines+1:]
	}
	lines = append(lines, text)
	outputLabel.SetText(strings.Join(lines, "\n"))
}
