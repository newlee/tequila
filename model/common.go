package model

type Entity struct {
	name     string
	Entities []*Entity
	VOs      []*ValueObject
	Refs     []*Entity
}

type ValueObject struct {
	name string
}

type Repository struct {
	name string
	For  string
}
type Provider struct {
	name string
}

func NewEntity(name string) *Entity {
	return &Entity{name: name}
}

func NewValueObject(name string) *ValueObject {
	return &ValueObject{name: name}
}

func NewRepository(name string) *Repository {
	return &Repository{name: name}
}

func NewProvider(name string) *Provider {
	return &Provider{name: name}
}

func (entity *Entity) AppendVO(vo *ValueObject) {
	for _, item := range entity.VOs {
		if item.name == vo.name {
			return
		}
	}
	entity.VOs = append(entity.VOs, vo)
}
func (entity *Entity) Compare(other *Entity) bool {
	if len(entity.Entities) != len(other.Entities) {
		return false
	}
	if len(entity.VOs) != len(other.VOs) {
		return false
	}
	em := make(map[string]*Entity)
	for _, childEntity := range entity.Entities {
		em[childEntity.name] = childEntity
	}
	for _, childEntity := range other.Entities {
		if !em[childEntity.name].Compare(childEntity) {
			return false
		}
	}
	vom := make(map[string]*ValueObject)
	for _, vo := range entity.VOs {
		vom[vo.name] = vo
	}
	for _, vo := range other.VOs {
		if _, ok := vom[vo.name]; !ok {
			return false
		}
	}
	return true
}

func (repo *Repository) Compare(other *Repository) bool {
	return repo.For == other.For
}
