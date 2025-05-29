// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package testutil

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, file, _, _ = runtime.Caller(0)
	dir           = filepath.Dir(file)
	tpl           = template.Must(template.New("").ParseFiles(dir + "/table.gohtml"))
)

// Record represents one row of a decision table.
// This struct is to make it easy to handle data in Go templates.
type Record struct {
	// Name is the name of row.
	// A condition name or an action name.
	Name string
	// Values are the value in the cells of table.
	Values []string
}

// TableInfo has all the data of a decision table.
// Data is held in the way to make it easy to handle in Go templates.
type TableInfo struct {
	// Name is the table name.
	Name string
	// Descriptions are the description for the table.
	Descriptions []string
	// EntryNames are the names of all entries.
	// It's order is left to right in the decision table.
	EntryNames []string
	// EntryCndDescriptions are the descriptions fo the conditions.
	// The order is the same as the EntryNames field.
	EntryCndDescriptions []string
	// EntryActDescriptions are the descriptions fo the actions.
	// The order is the same as the EntryNames field.
	EntryActDescriptions []string
	// Conditions are the conditions.
	// The values are true ("T") or false ("F").
	// The order is the same as the EntryNames field.
	Conditions []*Record
	// Actions are the actions.
	// The values are true ("T") or false ("F").
	// The order is the same as the EntryNames field.
	Actions []*Record
	// CndDescriptions are the descriptions for each conditions.
	// The order is the same as the Conditions field.
	CndDescriptions []string
	// ActDescriptions are the descriptions for each actions.
	// The order is the same as the Actions field.
	ActDescriptions []string
}

// NewSaver creates a new instance of Saver.
//
// Usage:
//
//	s := NewSaver().Title("TestsDecisionTable").FileName("table.html")
//	s.Add(x.Info(), y.Info(), z.Info()) //x,y,x are the instance of Table.
//	s.Save() // output decision table html.
//
// When saving a html on the end of unit tests, TestMain function may be helpful like below.
//
//	func TestMain(m *testing.M) {
//		flag.Parse()
//		code := m.Run()
//		s.Save()
//		os.Exit(code)
//	}
func NewSaver() *Saver {
	s := &Saver{}
	return s
}

// Saver holds a HTML file info and output a decision table html using Go template.
//
// Usage:
//
//	s := NewSaver().Title("TestsDecisionTable").FileDir(../tables/).FileName("table.html")
//	s.Add(x.Info(), y.Info(), z.Info()) //x,y,x are the instance of Table.
//	s.Save() // output decision table html.
type Saver struct {
	title    string
	fileDir  string
	fileName string
	infos    []*TableInfo
}

// Title sets the title of the html content.
// The title comes to the top of the html.
// Calling this multiple times overwrites the already set one.
func (s *Saver) Title(t string) *Saver {
	s.title = t
	return s
}

// FileDir sets the directory to save html.
// Calling this multiple times overwrites the already set one.
func (s *Saver) FileDir(d string) *Saver {
	s.fileDir = d
	return s
}

// FileName sets the filename of the html.
// The file name should be end with ".html".
// Calling this multiple times overwrites the already set one.
func (s *Saver) FileName(f string) *Saver {
	s.fileName = f
	return s
}

// Add register table information for this saver.
// Tables are listed in the html as the same order as added.
func (s *Saver) Add(infos ...*TableInfo) *Saver {
	if s != nil {
		s.infos = append(s.infos, infos...)
	}
	return s
}

// Save outputs a html file.
// This panics if failed to create the file.
func (s *Saver) Save() *Saver {
	caller := ""
	if _, file, _, ok := runtime.Caller(1); ok {
		caller = filepath.Base(file)
	}
	if s.fileName == "" {
		if caller != "" {
			s.fileName = caller + ".html"
		} else {
			s.fileName = "DecisionTable.html"
		}
	}
	if s.title == "" {
		if caller != "" {
			s.title = caller
		} else {
			s.title = "DecisionTable"
		}
	}
	output := filepath.Join(s.fileDir, s.fileName)

	root := struct {
		Title string
		Infos []*TableInfo
	}{
		s.title,
		s.infos,
	}

	// Generate html from the template.
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, root); err != nil {
		panic(err)
	}
	// compress the html.
	content := buf.String()
	// content = strings.ReplaceAll(content, "\t", "")
	// content = strings.ReplaceAll(content, "\n", "")

	// Output to a file.
	if _, err := os.Stat(s.fileDir); os.IsNotExist(err) {
		if err := os.Mkdir(s.fileDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	if _, err := f.WriteString(content); err != nil {
		panic(err)
	}
	return s
}

// contains checks if the given array of the string contains the given string.
func contains(s string, list []string) bool {
	for _, val := range list {
		if s == val {
			return true
		}
	}
	return false
}

// Table holds all information of a table.
// This struct should be instantiated by using TableBuilder.
type Table[C, A any] struct {
	name            string
	descriptions    []string
	conditions      []string
	actions         []string
	cndDescriptions []string
	actDescriptions []string
	entries         []*Entry[C, A]
}

// Add adds Entry[C,A] instances to this table.
// Entries to be added should be created by a entry builder
// created by EntryBuilder() method of this table instance.
//
//	// table := <create a table here.>
//	entryBuilder := table.EntryBuilder()
//	// entry := <create an entry with entryBuilder here.>
//	table.Add(entry)
func (t *Table[C, A]) Add(es ...*Entry[C, A]) *Table[C, A] {
	t.entries = append(t.entries, es...)
	return t
}

// Entries returns all Entry[C,A] instances added to this table.
func (t *Table[C, A]) Entries() []*Entry[C, A] {
	return t.entries
}

// EntryBuilder returns a EntryBuilder[C, A] instance suited to this table.
func (t *Table[C, A]) EntryBuilder() *EntryBuilder[C, A] {
	conditions := map[string]bool{}
	for _, k := range t.conditions {
		conditions[k] = false
	}
	actions := map[string]bool{}
	for _, k := range t.actions {
		actions[k] = false
	}
	return &EntryBuilder[C, A]{
		conditions: conditions,
		actions:    actions,
	}
}

// Info returns table information that can be registered to a Saver.
//
//	// saver := <create a saver here.>
//	// table := <create a table here.>
//	saver.Add(table.Info())
func (t *Table[C, A]) Info() *TableInfo {
	entryNames := make([]string, 0, len(t.entries))
	for _, e := range t.entries {
		entryNames = append(entryNames, e.Name())
	}

	entryCndDescriptions := make([]string, 0, len(t.entries))
	entryActDescription := make([]string, 0, len(t.entries))
	for _, e := range t.entries {
		entryCndDescriptions = append(entryCndDescriptions, fmt.Sprintf("%#v", e.C()))
		entryActDescription = append(entryActDescription, fmt.Sprintf("%#v", e.C()))
	}

	conditions := make([]*Record, 0, len(t.conditions))
	for _, key := range t.conditions {
		var val []string
		for _, e := range t.entries {
			if e.Conditions()[key] {
				val = append(val, "T")
				continue
			}
			val = append(val, "F")
		}
		conditions = append(conditions, &Record{key, val})
	}

	actions := make([]*Record, 0, len(t.actions))
	for _, key := range t.actions {
		var val []string
		for _, e := range t.entries {
			if e.Actions()[key] {
				val = append(val, "T")
				continue
			}
			val = append(val, "F")
		}
		actions = append(actions, &Record{key, val})
	}

	info := &TableInfo{
		Name:                 t.name,
		Descriptions:         t.descriptions,
		EntryNames:           entryNames,
		EntryCndDescriptions: entryCndDescriptions,
		EntryActDescriptions: entryActDescription,
		Conditions:           conditions,
		Actions:              actions,
		CndDescriptions:      t.cndDescriptions,
		ActDescriptions:      t.actDescriptions,
	}
	return info
}

// NewTableBuilder returns a new empty instance of TableBuilder[C, A].
func NewTableBuilder[C, A any]() *TableBuilder[C, A] {
	return &TableBuilder[C, A]{}
}

// TableBuilder is a builder for Table.
//
// Usage:
//
//	type C struct {Hoge: string} // define the type of condition
//	type A struct {Fuga: string} // define the type of action
//
//	tb := NewTableBuilder[C, A]()
//	// <additional setup of the table builder>
//	table := tb.Build()
//
//	// saver := <create a saver.>
//	saver.Add(table.Info())
//	saver.Save()
type TableBuilder[C, A any] struct {
	// name is the name of this table.
	name string
	// descriptions are the description of this table.
	descriptions []string
	// conditions are the condition names.
	//These correspond to the name of rows for conditions.
	conditions      []string
	cndDescriptions []string

	// actions are the action names.
	//These correspond to the name of rows for actions.
	actions         []string
	actDescriptions []string
}

// Name sets the name of this table.
// Calling this multiple times overwrites the already set name.
func (b *TableBuilder[C, A]) Name(n string) *TableBuilder[C, A] {
	b.name = n
	return b
}

// Description adds description for the table.
func (b *TableBuilder[C, A]) Description(d string) *TableBuilder[C, A] {
	b.descriptions = append(b.descriptions, d)
	return b
}

// Condition adds new condition entry.
// The order of conditions shown in the table becomes the same as the order added.
func (b *TableBuilder[C, A]) Condition(condition string, description string) string {
	if contains(condition, b.conditions) {
		return condition
	}
	b.conditions = append(b.conditions, condition)
	b.cndDescriptions = append(b.cndDescriptions, description)
	return condition
}

// Action adds new action entry.
// The order of actions shown in the table becomes the same as the order added.
func (b *TableBuilder[C, A]) Action(action string, d string) string {
	if contains(action, b.actions) {
		return action
	}
	b.actions = append(b.actions, action)
	b.actDescriptions = append(b.actDescriptions, d)
	return action
}

// Build builds a new decision table instance.
func (b *TableBuilder[C, A]) Build() *Table[C, A] {
	t := &Table[C, A]{
		name:            b.name,
		descriptions:    b.descriptions,
		conditions:      b.conditions,
		actions:         b.actions,
		cndDescriptions: b.cndDescriptions,
		actDescriptions: b.actDescriptions,
	}
	return t
}

type Entry[C, A any] struct {
	conditions map[string]bool
	actions    map[string]bool
	c          C
	a          A
	name       string
}

func (e *Entry[C, A]) Name() string {
	return e.name
}

func (e *Entry[C, A]) Conditions() map[string]bool {
	return e.conditions
}

func (e *Entry[C, A]) Actions() map[string]bool {
	return e.actions
}

func (e *Entry[C, A]) C() C {
	return e.c
}

func (e *Entry[C, A]) A() A {
	return e.a
}

func (e *Entry[C, A]) Register(t *Table[C, A]) *Entry[C, A] {
	t.Add(e)
	return e
}

type EntryBuilder[C, A any] struct {
	conditions map[string]bool
	actions    map[string]bool
	c          C
	a          A
	name       string
}

func (b *EntryBuilder[C, A]) Name(n string) *EntryBuilder[C, A] {
	b.name = n
	return b
}

func (b *EntryBuilder[C, A]) CTrue(c string) *EntryBuilder[C, A] {
	if _, ok := b.conditions[c]; !ok {
		panic("no such condition: " + c)
	}
	b.conditions[c] = true
	return b
}

func (b *EntryBuilder[C, A]) CTrueAll() *EntryBuilder[C, A] {
	for k := range b.conditions {
		b.conditions[k] = true
	}
	return b
}

func (b *EntryBuilder[C, A]) ATrue(a string) *EntryBuilder[C, A] {
	if _, ok := b.actions[a]; !ok {
		panic("no such action: " + a)
	}
	b.actions[a] = true
	return b
}

func (b *EntryBuilder[C, A]) ATrueAll(a string) *EntryBuilder[C, A] {
	for k := range b.actions {
		b.actions[k] = true
	}
	return b
}

func (b *EntryBuilder[C, A]) Condition(c C) *EntryBuilder[C, A] {
	b.c = c
	return b
}

func (b *EntryBuilder[C, A]) Action(a A) *EntryBuilder[C, A] {
	b.a = a
	return b
}

func (b *EntryBuilder[C, A]) Reset() *EntryBuilder[C, A] {
	b.name = ""
	conditions := map[string]bool{}
	actions := map[string]bool{}
	for k := range b.conditions {
		conditions[k] = false
	}
	for k := range b.actions {
		actions[k] = false
	}
	b.conditions = conditions
	b.actions = actions
	var c C
	b.c = c
	var a A
	b.a = a
	return b
}

func (b *EntryBuilder[C, A]) Build() *Entry[C, A] {
	e := &Entry[C, A]{
		name:       b.name,
		conditions: b.conditions,
		actions:    b.actions,
		c:          b.c,
		a:          b.a,
	}
	b.Reset()
	return e
}
