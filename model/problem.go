package model

type SubDomain struct {
	ARs       map[string]*Entity
	Repos     map[string]*Repository
	Providers map[string]*Provider
	es        map[string]*Entity
	vos       map[string]*ValueObject
}
type ProblemModel struct {
	SubDomains map[string]*SubDomain
}

func (model *ProblemModel) Validate() bool {
	for key := range model.SubDomains {
		if !model.SubDomains[key].Validate() {
			return false
		}
	}
	return true
}

func (model *ProblemModel) Compare(other *ProblemModel) bool {
	if len(model.SubDomains) != len(other.SubDomains) {
		return false
	}

	for key := range model.SubDomains {
		ar := model.SubDomains[key]
		if !ar.Compare(other.SubDomains[key]) {
			return false
		}
	}

	return true
}

func (subDomain *SubDomain) Validate() bool {
	entityMap := make(map[string]int)
	for key := range subDomain.ARs {
		ar := subDomain.ARs[key]
		for _, entity := range ar.Entities {
			if _, ok := entityMap[entity.name]; !ok {
				entityMap[entity.name] = 0
			}
			entityMap[entity.name] = entityMap[entity.name] + 1
		}
	}
	for key := range entityMap {
		if entityMap[key] > 1 {
			return false
		}
	}
	return true
}

func (subDomain *SubDomain) Compare(other *SubDomain) bool {
	if len(subDomain.ARs) != len(other.ARs) {
		return false
	}
	if len(subDomain.Repos) != len(other.Repos) {
		return false
	}
	for key := range subDomain.ARs {
		ar := subDomain.ARs[key]
		if !ar.Compare(other.ARs[key]) {
			return false
		}
	}
	for key := range subDomain.Repos {
		repo := subDomain.Repos[key]
		if !repo.Compare(other.Repos[key]) {
			return false
		}
	}
	return true
}

func (subDomain *SubDomain) AddNode(cms *CommentMappingList, name, comment string) {
	for _, cm := range *cms {
		if cm.comment == comment {
			cm.mapping(subDomain, name)
			break
		}
	}
}

type SubDomainWhenThen struct {
	subDomain *SubDomain
	isMatch   bool
	current   interface{}
	src       string
	dsts      []string
}

func (subDomain *SubDomain) given(src string, dsts []string) *SubDomainWhenThen {
	return &SubDomainWhenThen{
		subDomain: subDomain,
		src:       src,
		dsts:      dsts,
	}
}

func (whenThen *SubDomainWhenThen) when(fined interface{}, ok bool) *SubDomainWhenThen {
	whenThen.isMatch = ok
	if ok {
		whenThen.current = fined
	}

	return whenThen
}

func (whenThen *SubDomainWhenThen) whenRepo() *SubDomainWhenThen {
	repo, ok := whenThen.subDomain.Repos[whenThen.src]
	return whenThen.when(repo, ok)
}

func (whenThen *SubDomainWhenThen) whenEntity() *SubDomainWhenThen {
	entity, ok := whenThen.subDomain.es[whenThen.src]
	return whenThen.when(entity, ok)
}
func (whenThen *SubDomainWhenThen) whenAggregateRoot() *SubDomainWhenThen {
	ar, ok := whenThen.subDomain.ARs[whenThen.src]
	return whenThen.when(ar, ok)
}

func (whenThen *SubDomainWhenThen) thenAdd(addRelations func(interface{}, []string)) *SubDomainWhenThen {
	if whenThen.isMatch {
		addRelations(whenThen.current, whenThen.dsts)
	}
	return whenThen
}

func (subDomain *SubDomain) AddRelations(src string, dsts []string) {
	subDomain.given(src, dsts).
		whenAggregateRoot().thenAdd(subDomain.addAggregateRootRelations).
		whenEntity().thenAdd(subDomain.addEntityRelations).
		whenRepo().thenAdd(subDomain.addRepoRelations)
}

func (subDomain *SubDomain) addRepoRelations(repo interface{}, dsts []string) {
	for _, dst := range dsts {
		repo.(*Repository).For = subDomain.ARs[dst]
	}
}
func (subDomain *SubDomain) addEntityRelations(entity interface{}, dsts []string) {
	for _, dst := range dsts {
		_entity := entity.(*Entity)
		if et, ok := subDomain.es[dst]; ok {
			_entity.Entities = append(_entity.Entities, et)
		}
		if vo, ok := subDomain.vos[dst]; ok {
			_entity.VOs = append(_entity.VOs, vo)
		}
	}
}
func (subDomain *SubDomain) addAggregateRootRelations(ar interface{}, dsts []string) {
	for _, dst := range dsts {
		_ar := ar.(*Entity)
		if ref, ok := subDomain.ARs[dst]; ok {
			_ar.Refs = append(_ar.Refs, ref)
		}
		if et, ok := subDomain.es[dst]; ok {
			_ar.Entities = append(_ar.Entities, et)
		}
		if vo, ok := subDomain.vos[dst]; ok {
			_ar.VOs = append(_ar.VOs, vo)
		}
	}
}

func NewSubDomain() *SubDomain {
	return &SubDomain{
		ARs:       make(map[string]*Entity),
		Repos:     make(map[string]*Repository),
		Providers: make(map[string]*Provider),
		es:        make(map[string]*Entity),
		vos:       make(map[string]*ValueObject),
	}
}

type CommentMapping struct {
	comment string
	mapping func(domain *SubDomain, name string)
}

type CommentMappingList []*CommentMapping

var addAggregateRootFunc = func(subDomain *SubDomain, name string) {
	subDomain.ARs[name] = NewEntity(name)
}
var addEntityFunc = func(subDomain *SubDomain, name string) {
	subDomain.es[name] = NewEntity(name)
}
var addValueObjectFunc = func(subDomain *SubDomain, name string) {
	subDomain.vos[name] = NewValueObject(name)
}

var addRepoFunc = func(subDomain *SubDomain, name string) {
	subDomain.Repos[name] = NewRepository(name)
}
var addProviderFunc = func(subDomain *SubDomain, name string) {
	subDomain.Providers[name] = NewProvider(name)
}

func InitCommentMapping() *CommentMappingList {
	return &CommentMappingList{
		{comment: "AR", mapping: addAggregateRootFunc},
		{comment: "E", mapping: addEntityFunc},
		{comment: "VO", mapping: addValueObjectFunc},
		{comment: "Repo", mapping: addRepoFunc},
		{comment: "Provider", mapping: addProviderFunc},
	}
}
