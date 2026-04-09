package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Language struct {
	ent.Schema
}

func (Language) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("native_name").
			NotEmpty(),
		field.Bool("is_default").Default(false),
		field.Bool("active").Default(true),
		field.Bool("deleted").Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Language) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("resources", LanguageResource.Type),
	}
}

func (Language) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted"),
		index.Fields("code", "deleted").Unique(),
	}
}
