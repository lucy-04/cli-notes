package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/list"
	
)

var (
	vaultDir string
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type item struct{
	title,desc string
}

type model struct {
		NoteHeadingInput textinput.Model
		NoteHeadingInputValue bool
		NoteTextArea textarea.Model
		NoteTextAreaValue bool
		FilePointer *os.File
		list list.Model
		showList bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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

	listInit := list.New(listData(), list.NewDefaultDelegate(), 0, 0)

	
	return model{
			NoteHeadingInput : ti,
			NoteHeadingInputValue : false,
			NoteTextArea: ta,
			NoteTextAreaValue: false,
			list: listInit,
			showList: false,
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
	if(m.NoteTextAreaValue){
		m.NoteTextArea , cmd = m.NoteTextArea.Update(msg)
	}
	if m.showList{
		m.list ,cmd  = m.list.Update(msg)
	}
	switch msg := msg.(type){
	case tea.KeyPressMsg:
		switch msg.String(){
		case "ctrl+c":
			return m,tea.Quit
		case "ctrl+n":
			m.NoteHeadingInputValue = true
		case "ctrl+s": // make it autoSave in future

			if m.FilePointer == nil{
				break;
			}

			if err := m.FilePointer.Truncate(0); err != nil{
				fmt.Printf("Cannot save the file")
				return m,nil
			}

			if _,err := m.FilePointer.Seek(0,0); err != nil{
				fmt.Printf("Cannot save the file")
				return m,nil
			}

			if _, err := m.FilePointer.WriteString(m.NoteTextArea.Value()); err != nil{
				fmt.Printf("Cannot save the file")
				return m,nil
			}


			if err := m.FilePointer.Close(); err != nil {
				fmt.Printf("Unable to close the file")
			}
			
			m.NoteTextArea.SetValue("")
			m.NoteTextAreaValue = false
			m.FilePointer = nil

			return m,nil

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
			m.NoteTextAreaValue = true
			}
		
		case "ctrl+l":
			m.showList = true
		
	}
	case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)
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
	if m.showList{
		view = m.list.View()
	}
	
	return tea.NewView(fmt.Sprintf("\n\n%s\n\n%s\n\n%s",welcome,view,help))
}

func listData() []list.Item{
	items := make([]list.Item,0)
	DirList, err := os.ReadDir(vaultDir)
	if err != nil{
		log.Fatal("file reading error")
	}
	
	for _,entry := range DirList{
		if !entry.IsDir(){
			info , err := entry.Info()
			if err != nil{
				continue
			}

			modTime := info.ModTime().Format("2020-02-13 15:32")

			items = append(items, item{
				title : entry.Name(),
				desc: fmt.Sprintf("Last Modified: %s",modTime),
			})
		}
	}
	return items
	
}

func main() {
    p := tea.NewProgram(initialiseModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}