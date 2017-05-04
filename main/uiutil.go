package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// IDs to access the tree view columns by
const (
	COLUMN_VERSION = iota
	COLUMN_FEATURE
	COLUMN_FEATURE1
)

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

// Creates a tree view and the list store that holds its data
func setupTreeView3(c1 string, c2 string) ( *gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	treeView.AppendColumn(createColumn(c1, COLUMN_VERSION))
	treeView.AppendColumn(createColumn(c2, COLUMN_FEATURE))

	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)

	return treeView, listStore
}
func setupTreeView(c1 string) ( *gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	treeView.AppendColumn(createColumn(c1, COLUMN_VERSION))


	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)

	return treeView, listStore
}

func setupTreeView2(c1 string, c2 string) ( *gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	treeView.AppendColumn(createColumn(c1, COLUMN_VERSION))
	treeView.AppendColumn(createColumn(c2, COLUMN_FEATURE))


	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)

	return treeView, listStore
}

// Append a row to the list store for the tree view
func addRow2(listStore *gtk.ListStore, version, feature string) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()


	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_VERSION, COLUMN_FEATURE},
		[]interface{}{version, feature})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}
func addRow3(listStore *gtk.ListStore, version, feature,feature1 string) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_VERSION, COLUMN_FEATURE,COLUMN_FEATURE1},
		[]interface{}{version, feature, feature1})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}
func addRow1(listStore *gtk.ListStore, version string) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_VERSION},
		[]interface{}{version})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}
