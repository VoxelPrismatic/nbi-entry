package variable

import "nbientry/web"

type var_str_plain struct {
	Real  VariableEntry
	T     Variable
	Value string
}

func (v var_str_plain) GetNbiId() int {
	return v.Real.NbiId
}

func (v var_str_plain) GetVarId() int {
	return v.Real.Id
}

func (v var_str_plain) GetIndex() int {
	return v.Real.Index
}

func (v *var_str_plain) SetIndex(idx int) {
	v.Real.Index = idx
}

func (v var_str_plain) Type() Variable {
	if v.T.Id == 0 {
		v.T = web.GetFirst(Variable{Id: v.Real.Id})
	}
	return v.T
}

func (v *var_str_plain) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_str_plain) Get() string {
	return v.Value
}

func (v *var_str_plain) Reset() {
	v.Value = ""
}

func (v var_str_plain) ToSqlEntry() VariableEntry {
	return v.Real
}

func (v var_str_plain) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}
