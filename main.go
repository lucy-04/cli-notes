package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/textarea"
)

var (
	vaultDir string
)

type model struct {
		NoteHeadingInput textinput.Model
		NoteHeadingInputValue bool
		NoteTextArea textarea.Model
		FilePointer *os.File
}

func init() {
	
	homeDir,err := os.UserHomeDir()
	if(err != nil){
		print(err)
	}
	vaultDir = fmt.Sprintf("%s/Documents/.notes",homeDir)
}
func initialiseModel() model{

	err := os.MkdirAll(vaultDir,0750)
	if err != nil{
		log.Fatal(err)
	}
	ti := textinput.New()
	ti.Placeholder = "Write you Note heading here"
	ti.SetVirtualCursor(true)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(50)

	ta := textarea.New()
	ta.Placeholder = "Once upon a time..."
	ta.SetVirtualCursor(true)
	ta.SetStyles(textarea.DefaultStyles(true)) // default to dark styles.
	ta.Focus()

	
	return model{
			NoteHeadingInput : ti,
			NoteHeadingInputValue : false,
			NoteTextArea: ta,
		}
}

func (m model) Init() tea.Cmd{
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model,tea.Cmd){

	var cmd tea.Cmd
	
	if(m.NoteHeadingInputValue) {
		m.NoteHeadingInput,cmd = m.NoteHeadingInput.Update(msg)
	}
	switch msg := msg.(type){
	case tea.KeyPressMsg:
		switch msg.String(){
		case "super+c","q":
			return m,tea.Quit
		case "n":
			m.NoteHeadingInputValue = true
		case "enter":
			fileName := m.NoteHeadingInput.Value()
			filePath := fmt.Sprintf("%s/%s.md",vaultDir,fileName)
			
			if _,err := os.Stat(filePath) ; err == nil{
				log.Fatalf("File already exists: %v",err)
			}
			if fileName != ""{
				f,err := os.Create(filePath)
				if err != nil {
					log.Fatal("Error in creating file")
				}
			
			
			
			m.FilePointer = f
			m.NoteHeadingInput.SetValue("")
			m.NoteHeadingInputValue = false
			}

		}
		
	}
	
	return m, cmd
}

func (m model) View() tea.View{

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(4).
		Width(22)

	welcome := style.Render("Welcome to Notes")
	help := "Ctrl+N: New File, Ctrl+O: Open Notes, Ctrl+h: help, Ctrl+q/q: Quit"
	view := ""
	if m.NoteHeadingInputValue{
		view = m.NoteHeadingInput.View()
	}
	if(m.FilePointer != nil){
		view = m.NoteTextArea.View()
	}
	
	return tea.NewView(fmt.Sprintf("\n\n%s\n\n%s\n\n%s",welcome,view,help))
}

func main() {
    p := tea.NewProgram(initialiseModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}