package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.

func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.UUID("author_id", uuid.UUID{}),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Field("author_id").
			Unique().
			Required().
			Ref("todos"),
	}
}
