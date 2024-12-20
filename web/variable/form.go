package variable

import (
	"fmt"
	"nbientry/web"
	"slices"
)

type var_form struct {
	Real  VariableEntry
	T     Variable
	Value string
}

func (v var_form) GetVarId() int {
	return v.T.Id
}

func (v var_form) GetNbiId() int {
	return v.Real.NbiId
}

func (v var_form) GetIndex() int {
	v.Real.Index = 1
	return v.Real.Index
}

func (v *var_form) SetIndex(_ int) {
	v.Real.Index = 1
}

func (v var_form) Type() Variable {
	if v.T.Id == 0 {
		v.T = web.GetFirst(Variable{Id: v.T.Id})
	}
	return v.T
}

func (v var_form) Get() string {
	return v.Value
}

func (v *var_form) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_form) ToSqlEntry() VariableEntry {
	return v.Real
}

func (v var_form) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}

func (v *var_form) Reset() {
	v.Value = ""
}

func (v var_form) children() ([]VarEntry, []VarEntry) {
	fmt.Println(v.T.Id)
	if v.T.Id == 0 {
		return nil, nil
	}
	children := web.GetSorted(Variable{ParentId: v.T.Id}, "Id ASC")
	fmt.Println(children)
	var ret []VarEntry
	var last []VarEntry
	for _, c := range children {
		val := c.GetEntry(v.Real.NbiId)
		if slices.Contains(expandable, c.Type) {
			last = append(last, &val)
		} else {
			ret = append(ret, &val)
		}
	}

	if len(ret) == 1 {
		fmt.Println(ret[0])
	}
	fmt.Println(ret)
	fmt.Println(last)
	return ret, last
}
