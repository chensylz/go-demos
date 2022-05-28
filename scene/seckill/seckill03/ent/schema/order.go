package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("merchandise_id").Comment("商品ID"),
		field.Int("user_id").Comment("用户ID"),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return nil
}
