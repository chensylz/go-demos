package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/zaihui/go-hutils"
)

// Merchandise holds the schema definition for the Merchandise entity.
type Merchandise struct {
	hutils.BaseSchema
}

// Fields of the Merchandise.
func (Merchandise) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("stock").Default(0).Comment("库存"),
	}
}

// Edges of the Merchandise.
func (Merchandise) Edges() []ent.Edge {
	return nil
}
