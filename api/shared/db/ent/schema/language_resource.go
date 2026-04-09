package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type LanguageResource struct {
	ent.Schema
}

func (LanguageResource) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_id"),
		field.String("key").
			NotEmpty(),
		field.Text("value").
			Default(""),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (LanguageResource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("language", Language.Type).
			Ref("resources").
			Field("language_id").
			Unique().
			Required(),
	}
}

func (LanguageResource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("language_id"),
		index.Fields("language_id", "key").Unique(),
	}
}
