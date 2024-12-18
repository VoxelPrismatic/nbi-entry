package variable

import "nbientry/web"

type var_str_plain struct {
	NbiId      int
	VariableId int
	Index      int
	Variable   Variable
	Value      string
}

func (v var_str_plain) GetNbiId() int {
	return v.NbiId
}

func (v var_str_plain) GetVarId() int {
	return v.VariableId
}

func (v var_str_plain) GetIndex() int {
	return v.Index
}

func (v *var_str_plain) SetIndex(idx int) {
	v.Index = idx
}

func (v var_str_plain) GetVariable() Variable {
	if v.Variable.VariableId == 0 {
		v.Variable = web.GetFirst(Variable{VariableId: v.VariableId})
	}
	return v.Variable
}

func (v *var_str_plain) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_str_plain) Get() string {
	return v.Value
}

func (v var_str_plain) ToSqlEntry() VariableEntry {
	return VariableEntry{NbiId: v.NbiId, VariableId: v.VariableId, Value: v.Value}
}

func (v var_str_plain) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}
