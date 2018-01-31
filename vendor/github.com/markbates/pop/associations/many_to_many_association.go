package associations

import (
	"fmt"
	"reflect"

	"github.com/markbates/inflect"
)

type manyToManyAssociation struct {
	fieldType           reflect.Type
	fieldValue          reflect.Value
	model               reflect.Value
	manyToManyTableName string
	owner               interface{}
	fkID                string
	orderBy             string
}

func init() {
	associationBuilders["many_to_many"] = func(p associationParams) (Association, error) {
		return &manyToManyAssociation{
			fieldType:           p.modelValue.FieldByName(p.field.Name).Type(),
			fieldValue:          p.modelValue.FieldByName(p.field.Name),
			owner:               p.model,
			model:               p.modelValue,
			manyToManyTableName: p.popTags.Find("many_to_many").Value,
			fkID:                p.popTags.Find("fk_id").Value,
			orderBy:             p.popTags.Find("order_by").Value,
		}, nil
	}
}

func (m *manyToManyAssociation) TableName() string {
	i := m.fieldValue.Interface()
	if m, ok := i.(tableNameable); ok {
		return m.TableName()
	}
	return inflect.Tableize(m.fieldType.Name())
}

func (m *manyToManyAssociation) Kind() reflect.Kind {
	return m.fieldType.Kind()
}

func (m *manyToManyAssociation) Interface() interface{} {
	if m.fieldValue.Kind() == reflect.Ptr {
		val := reflect.New(m.fieldType.Elem())
		m.fieldValue.Set(val)
		return m.fieldValue.Interface()
	}
	return m.fieldValue.Addr().Interface()
}

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (m *manyToManyAssociation) Constraint() (string, []interface{}) {
	modelColumnID := fmt.Sprintf("%s%s", inflect.Underscore(m.model.Type().Name()), "_id")

	var columnFieldID string
	i := reflect.Indirect(m.fieldValue)
	if i.Kind() == reflect.Slice || i.Kind() == reflect.Array {
		t := i.Type().Elem()
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(t.Name()), "_id")
	} else {
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(i.Type().Name()), "_id")
	}

	if m.fkID != "" {
		columnFieldID = m.fkID
	}

	subQuery := fmt.Sprintf("select %s from %s where %s = ?", columnFieldID, m.manyToManyTableName, modelColumnID)
	modelIDValue := m.model.FieldByName("ID").Interface()

	return fmt.Sprintf("id in (%s)", subQuery), []interface{}{modelIDValue}
}

func (m *manyToManyAssociation) OrderBy() string {
	return m.orderBy
}