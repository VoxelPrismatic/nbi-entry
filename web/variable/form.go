package variable

import "nbientry/web"

type var_form struct {
	NbiId      int
	VariableId int
	Variable   Variable
	Index      int
	Value      string
}

func (v var_form) GetVarId() int {
	return v.VariableId
}

func (v var_form) GetNbiId() int {
	return v.NbiId
}

func (v var_form) GetIndex() int {
	v.Index = 1
	return v.Index
}

func (v *var_form) SetIndex(_ int) {
	v.Index = 1
}

func (v var_form) GetVariable() Variable {
	if v.Variable.VariableId == 0 {
		v.Variable = web.GetFirst(Variable{VariableId: v.VariableId})
	}
	return v.Variable
}

func (v var_form) Get() string {
	return v.Value
}

func (v *var_form) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_form) ToSqlEntry() VariableEntry {
	return VariableEntry{NbiId: v.NbiId, VariableId: v.VariableId, Value: v.Value}
}

func (v var_form) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}

func (v var_form) children() []VarEntry {
	children := web.GetSorted(Variable{ParentId: v.VariableId}, "VariableId ASC")
	var ret []VarEntry
	for _, c := range children {
		val := c.GetEntry(v.NbiId)
		ret = append(ret, &val)
	}
	return ret
}
