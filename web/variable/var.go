package variable

import (
	"fmt"
	"nbientry/web"

	"github.com/a-h/templ"
)

var _ = web.Migrate(&Variable{})
var _ = web.Migrate(&VariableEntry{})

var expandable = []string{"form", "table", "reason", "paragraph"}

type Variable struct {
	Id          int `gorm:"primaryKey"`
	ParentId    int
	Name        string
	Type        string
	Description string
	Suffix      string
}

type VariableEntry struct {
	NbiId int      `gorm:"primaryKey;auto_increment=false"`
	Id    int      `gorm:"primaryKey;auto_increment=false"`
	Index int      `gorm:"primaryKey;auto_increment=false"`
	T     Variable `gorm:"foreignKey:Id"`
	Value string
}

func (v Variable) GetParent() Variable {
	return web.GetFirst(Variable{Id: v.ParentId})
}

func (v Variable) GetChildren() []Variable {
	return web.GetSorted(Variable{ParentId: v.Id}, "Id ASC")
}

func (v Variable) GetEntry(id int) VariableEntry {
	ret := web.GetFirst(VariableEntry{NbiId: id, Id: v.Id})
	if ret.Id == 0 {
		ret = v.New()
	}
	return ret
}

func (v *Variable) New() VariableEntry {
	if v.Id == 0 {
		v.Id = -1
	}

	return VariableEntry{
		Id: v.Id,
		T:  *v,
	}
}

func (v *Variable) Delete() {
	for _, c := range v.GetChildren() {
		c.Delete()
	}

	for _, e := range web.GetSorted(VariableEntry{Id: v.Id}, "nbi_id ASC") {
		web.Db().Delete(e)
	}

	web.Db().Delete(v)
}

type VarEntry interface {
	GetVarId() int
	GetNbiId() int
	GetIndex() int
	SetIndex(idx int)
	Type() Variable
	Get() string
	Set(string) error
	Reset()
	RenderInNbi() templ.Component
	RenderInViewer() templ.Component
	RenderInEditor() templ.Component
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

func (v *VariableEntry) Type() Variable {
	if v.T.Id == 0 {
		v.T = web.GetFirst(Variable{Id: v.Id})
	}
	return v.T
}

func (v VariableEntry) Get() string {
	return v.Value
}

func (v *VariableEntry) Reset() {
	v.Value = ""
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
	var err error

	switch v.Type().Type {
	case "text-plain":
		ret = &var_str_plain{Real: v, T: v.T}
	case "text-suggest":
		ret = &var_str_suggest{Real: v, T: v.T}
	case "form":
		ret = &var_form{Real: v, T: v.T}
	case "num":
		ret = &var_num{Real: v, T: v.T}
	default:
		fmt.Printf("\x1b[91;1munknown variable type:\x1b[0m %s\n", v.Type().Type)
		return &v, fmt.Errorf("unknown variable type")
	}

	if v.Value == "" {
		ret.Reset()
	} else {
		err = ret.Set(v.Value)
	}
	return ret, err
}

func (v VariableEntry) RenderInNbi() templ.Component {
	entry, err := v.ToTypedEntry()
	if err != nil {
		return nil
	}
	return entry.RenderInNbi()
}

func (v VariableEntry) RenderInEditor() templ.Component {
	entry, err := v.ToTypedEntry()
	if err != nil {
		return nil
	}
	return entry.RenderInEditor()
}

func (v VariableEntry) RenderInViewer() templ.Component {
	entry, err := v.ToTypedEntry()
	if err != nil {
		return nil
	}
	return entry.RenderInViewer()
}
