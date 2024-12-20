package variable

import (
	"fmt"
	"nbientry/web"
	"strconv"
)

type var_num struct {
	Real  VariableEntry
	T     Variable
	Value float64
}

func (v var_num) GetNbiId() int {
	return v.Real.NbiId
}

func (v var_num) GetVarId() int {
	return v.Real.Id
}

func (v var_num) GetIndex() int {
	return v.Real.Index
}

func (v var_num) SetIndex(idx int) {
	v.Real.Index = idx
}

func (v var_num) Type() Variable {
	if v.T.Id == 0 {
		v.T = web.GetFirst(Variable{Id: v.Real.Id})
	}
	return v.T
}

func (v var_num) Get() string {
	return fmt.Sprint(v.Value)
}

func (v *var_num) Set(value string) error {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	v.Value = n
	return nil
}

func (v *var_num) Reset() {
	v.Value = 0
}

func (v var_num) ToSqlEntry() VariableEntry {
	return v.Real
}

func (v var_num) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}
