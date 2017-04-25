package main

import (
	"fmt"
	"log"
	"github.com/gotk3/gotk3/gtk"
	"strconv"
	"github.com/gotk3/gotk3/glib"
	"bufio"
	"os"
)

func setup_window(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func setup_box(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

func setup_tview() (*gtk.TextView, *gtk.ScrolledWindow) {
	tv, err := gtk.TextViewNew()
	//tv.SetName("Results")
	sw, err := gtk.ScrolledWindowNew(nil, nil)
	sw.Add(tv)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}

	tv.SetHExpand(true)
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv,sw
}

func setup_btn(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.Connect("clicked", onClick)
	return btn
}

func get_buffer_from_tview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func get_text_from_tview(tv *gtk.TextView) string {
	buffer := get_buffer_from_tview(tv)
	start, end := buffer.GetBounds()

	text, err := buffer.GetText(start, end, true)
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	return text
}


func set_text_in_tview(tv *gtk.TextView, text string) {
	buffer := get_buffer_from_tview(tv)
	buffer.SetText(text)
}

func newStackFull() gtk.IWidget {
	// get a stack and its switcher.
	stack, err := gtk.StackNew()
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}

	sw, err := gtk.StackSwitcherNew()
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	sw.SetStack(stack)

	// Fill the stack with 3 pages.
	boxText1 := isPrime("isPrime", "Candidate Int64")
	boxText2 := preImage("Candidate Hash")
	boxText3 := repTable("Rank","Node", "Primality Score","Pre-Image Score")

	stack.AddTitled(boxText1, "key1", "Primality Test")
	stack.AddTitled(boxText2, "key2", "Pre-Image Test")
	stack.AddTitled(boxText3, "key3", "Reputation Board")

	// You can use icons for a switcher page (the page title will be visible as tooltip).
	//stack.ChildSetProperty(boxRadio, "", "")

	// Pack in a box.
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	box.PackStart(sw, false, false, 0)
	box.PackStart(stack, true, true, 0)
	return box
}


func preImage(c1 string) gtk.IWidget {
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	treeView, listStore := setupTreeView(c1)
	sw, _ := gtk.ScrolledWindowNew(nil, nil)
	sw.Add(treeView)

	box.PackStart(sw, true, true, 10)
	result, hor:=setup_tview()
	inFile, ioErr := os.Open("tdict.txt")

	if ioErr != nil{
		fmt.Println(ioErr)
	}

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var fileTextLine string

	for scanner.Scan() {
		fileTextLine = scanner.Text()
		addRow1(listStore,fileTextLine)

	}

	btn := setup_btn("Submit Pre-Image Test to Network", func() {
		selection, err := treeView.GetSelection()
		if err != nil {
			log.Fatal("Could not get tree selection object.")
		}
		selection.SetMode(gtk.SELECTION_SINGLE)

		_, iter, _ := selection.GetSelected()

		x, err := listStore.GetValue(iter, 0)
		if err != nil {
			log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
			return
		}
		hash,_:=x.GetString()
		fmt.Printf(hash)
		//nSets:=0
		go func() {
				//time.Sleep(time.Second*2)
				_, err := glib.IdleAdd(LabelSetTextIdle, "")
				if err != nil {
					log.Fatal("IdleAdd() failed:", err)
				}
				var dummy metric
				msg := encode(nodeinf.NodeName, "leader", "Hash",hash,"Request",nodeinf.NodeGroup,"",false,"","",dummy,hash)
				nodeSend(string(msg), nodesoc)

		}()

		//fmt.Println(text)
	})
	btn1 := setup_btn("Check Result", func() {

		go func() {
			str=get_text_from_tview(result)
			//time.Sleep(time.Second*2)
			_, err := glib.IdleAdd(LabelSetTextIdle, "")
			if err != nil {
				log.Fatal("IdleAdd() failed:", err)
			}
			ml := MQpopAll(nodesoc.appq)
			if ml.Front() == nil{
				str+=""
			}
			for n :=  ml.Front(); n != nil ; n = n.Next(){
				test := n.Value.(Message)

				str+="Input: "
				str+=test.Input
				str+="\n"
				str+=test.Value
				str+="\n"
				str+="Processed By: "
				str+=test.Sender
				str+="\n\n"
			}

			set_text_in_tview(result,str)

		}()

		//fmt.Println(text)
	})

	btn2 := setup_btn("Clear Results", func() {

		set_text_in_tview(result,"")

		//fmt.Println(text)
	})

	box.Add(btn)
	box.Add(btn1)
	box.Add(btn2)
	box.PackEnd(hor,true,true,10)

	return box

}


func isPrime(c1 string, c2 string) gtk.IWidget {
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	treeView, listStore := setupTreeView2(c1,c2)
	treeView.SetBorderWidth(23)
	sw, _ := gtk.ScrolledWindowNew(nil, nil)
	sw.Add(treeView)

	box.PackStart(sw, true, true, 10)
	result, hor:=setup_tview()

	i:=0
	for i<100 {
		x:= generateCandidate()
		isPrime:= verifyPrime(*x)
		addRow2(listStore, strconv.FormatBool(isPrime), x.String())
		i=i+1

	}
	// Add some rows to the list store

	btn := setup_btn("Submit Primality Test to Network", func() {
		selection, err := treeView.GetSelection()
		if err != nil {
			log.Fatal("Could not get tree selection object.")
		}
		selection.SetMode(gtk.SELECTION_SINGLE)

		_, iter, _ := selection.GetSelected()

		x, err := listStore.GetValue(iter, 1)
		if err != nil {
			log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
			return
		}
		isPrime,_:=x.GetString()
		//fmt.Printf(isPrime)
		//nSets:=0
		go func() {
			//time.Sleep(time.Second*2)
			_, err := glib.IdleAdd(LabelSetTextIdle, "")
			if err != nil {
				log.Fatal("IdleAdd() failed:", err)
			}
			var dummy metric
			msg := encode(nodeinf.NodeName, "leader", "Prime",isPrime,"Request",nodeinf.NodeGroup,"",false,"","", dummy,isPrime)
			nodeSend(string(msg), nodesoc)

		}()

		//fmt.Println(text)
	})
	btn1 := setup_btn("Check Result", func() {

		go func() {
			str=get_text_from_tview(result)
			//time.Sleep(time.Second*2)
			_, err := glib.IdleAdd(LabelSetTextIdle, "")
			if err != nil {
				log.Fatal("IdleAdd() failed:", err)
			}
			ml := MQpopAll(nodesoc.appq)
			if ml.Front() == nil{
				str+=""
			}
			for n :=  ml.Front(); n != nil ; n = n.Next(){
				test := n.Value.(Message)

				//fmt.Println("====Results====")

				str+="Input: "
				str+=test.Input
				str+="\n"
				str+=test.Value
				str+="\n"
				str+="Processed By: "
				str+=test.Sender
				str+="\n\n"
			}

			set_text_in_tview(result,str)

		}()

		//fmt.Println(text)
	})

	btn2 := setup_btn("Clear Results", func() {

			set_text_in_tview(result,"")

		//fmt.Println(text)
	})

	box.Add(btn)
	box.Add(btn1)
	box.Add(btn2)
	box.PackEnd(hor,true,true,10)

	return box

}

func repTable(c1 string, c2 string, c3 string, c4 string) gtk.IWidget {
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	treeView, listStore := setupTreeView3(c1,c2,c3)
	sw, _ := gtk.ScrolledWindowNew(nil, nil)
	sw.Add(treeView)
	box.PackStart(sw, true, true, 10)

	treeView1, listStore1 := setupTreeView3(c1,c2,c4)
	sw1, _ := gtk.ScrolledWindowNew(nil, nil)
	sw1.Add(treeView1)
	box.PackStart(sw1, true, true, 10)


	// Add some rows to the list store
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")
	addRow3(listStore1, "r57", "Gofix command added for rewriting code for new APIs","r57")

	btn := setup_btn("Refresh", func() {
		//text := get_text_from_tview(treeView)
		//fmt.Println(text)
	})
	btn1 := setup_btn("Train Network with test data", func() {

		//text := get_text_from_tview(treeView)
		//fmt.Println(text)
	})
	box.Add(btn)
	box.Add(btn1)
	return box
}
var nodesoc NodeSocket
var nodeinf NodeInfo
var str string
var (
	winTitle = "Distributed Load Prioritization"
)


func startUI(self NodeSocket, nodeinfo NodeInfo) {
	nodesoc=self
	nodeinf=nodeinfo

	gtk.Init(nil)

	win := setup_window(winTitle)

	box := newStackFull()
	win.Add(box)
	win.SetBorderWidth(20)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()

}

func LabelSetTextIdle( text string) bool {
	//label.SetText(text)

	// Returning false here is unnecessary, as anything but returning true
	// will remove the function from being called by the GTK main loop.
	return false
}


