package variable

import (
	"nbientry/web"
	"strings"
)

var _ = web.Migrate(&Autocomplete{})

type Autocomplete struct {
	VariableId int
	Value      string
}

type var_str_autofill struct {
	NbiId      int
	VariableId int
	Index      int
	Variable   Variable
	Value      string
}

func (v var_str_autofill) GetNbiId() int {
	return v.NbiId
}

func (v var_str_autofill) GetVarId() int {
	return v.VariableId
}

func (v var_str_autofill) GetIndex() int {
	return v.Index
}

func (v *var_str_autofill) SetIndex(idx int) {
	v.Index = idx
}

func (v var_str_autofill) GetVariable() Variable {
	if v.Variable.VariableId == 0 {
		v.Variable = web.GetFirst(Variable{VariableId: v.VariableId})
	}
	return v.Variable
}

func (v *var_str_autofill) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_str_autofill) Get() string {
	return v.Value
}

func (v var_str_autofill) ToSqlEntry() VariableEntry {
	return VariableEntry{NbiId: v.NbiId, VariableId: v.VariableId, Value: v.Value}
}

func (v var_str_autofill) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}

func (v var_str_autofill) exists() bool {
	autocompletes := web.GetSorted(Autocomplete{VariableId: v.VariableId}, "Value ASC")
	lower := strings.ToLower(v.Value)
	for _, a := range autocompletes {
		if strings.ToLower(a.Value) == lower {
			return true
		}
	}
	return false
}

func (v var_str_autofill) list() []string {
	autocompletes := web.GetSorted(Autocomplete{VariableId: v.VariableId}, "Value ASC")
	var ret []string
	for _, a := range autocompletes {
		ret = append(ret, a.Value)
	}
	return ret
}
