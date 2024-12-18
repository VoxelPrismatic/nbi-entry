package variable

import (
	"fmt"
	"nbientry/web"

	"github.com/a-h/templ"
)

var _ = web.Migrate(&Variable{})
var _ = web.Migrate(&VariableEntry{})

type Variable struct {
	VariableId  int `gorm:"primaryKey"`
	ParentId    int
	Name        string
	Type        string
	Description string
	Suffix      string
}

type VariableEntry struct {
	NbiId      int      `gorm:"primaryKey;auto_increment=false"`
	VariableId int      `gorm:"primaryKey;auto_increment=false"`
	Index      int      `gorm:"primaryKey;auto_increment=false"`
	Variable   Variable `gorm:"foreignKey:VariableId"`
	Value      string
}

func (v Variable) GetParent() Variable {
	return web.GetFirst(Variable{VariableId: v.ParentId})
}

func (v Variable) GetChildren() []Variable {
	return web.GetSorted(Variable{ParentId: v.VariableId}, "VariableId ASC")
}

func (v Variable) GetEntry(id int) VariableEntry {
	return web.GetFirst(VariableEntry{NbiId: id, VariableId: v.VariableId})
}

type VarEntry interface {
	GetVarId() int
	GetNbiId() int
	GetIndex() int
	SetIndex(idx int)
	GetVariable() Variable
	Get() string
	Set(string) error
	RenderInNbi() templ.Component
	// RenderInViewer() templ.Component
	// RenderInEditor() templ.Component
	ToSqlEntry() VariableEntry
	ToTypedEntry() (VarEntry, error)
}

func (v VariableEntry) GetNbiId() int {
	return v.NbiId
}

func (v VariableEntry) GetVarId() int {
	return v.NbiId
}

func (v VariableEntry) GetIndex() int {
	return v.Index
}

func (v *VariableEntry) SetIndex(idx int) {
	v.Index = idx
}

func (v *VariableEntry) GetVariable() Variable {
	if v.Variable.VariableId == 0 {
		v.Variable = web.GetFirst(Variable{VariableId: v.VariableId})
	}
	return v.Variable
}

func (v VariableEntry) Get() string {
	return v.Value
}

func (v *VariableEntry) Set(value string) error {
	typed, err := v.ToTypedEntry()
	if err != nil {
		return err
	}
	err = typed.Set(value)
	if err != nil {
		return err
	}
	v.Value = value
	return nil
}

func (v VariableEntry) ToSqlEntry() VariableEntry {
	return v
}

func (v VariableEntry) ToTypedEntry() (VarEntry, error) {
	var ret VarEntry
	switch v.GetVariable().Type {
	case "string_plain":
		ret = &var_str_plain{NbiId: v.NbiId, VariableId: v.VariableId}
	case "string_autofill":
		ret = &var_str_autofill{NbiId: v.NbiId, VariableId: v.VariableId}
	case "form":
		ret = &var_form{NbiId: v.NbiId, VariableId: v.VariableId}
	default:
		return nil, fmt.Errorf("unknown variable type: %s", v.GetVariable().Type)
	}
	err := ret.Set(v.Value)
	return ret, err
}

func (v VariableEntry) RenderInNbi() templ.Component {
	entry, err := v.ToTypedEntry()
	if err != nil {
		return nil
	}
	return entry.RenderInNbi()
}
