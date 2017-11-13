package model

import (
	"errors"
)

type layer interface {
	Add(name, comment string)
	addRelations(src string, dsts []string)
	compare(other interface{}) bool
}

type Layer struct {
	name  string
	nodes map[string]string
	layer layer
}

type BCModel struct {
	Layers map[string]*Layer
}

type Service struct {
	name string
	Refs []string
}

func NewService(name string) *Service {
	return &Service{name: name}
}

type Api struct {
	name string
	Refs []string
}

type GateWay struct {
	name      string
	Implement string
}
type DomainLayer struct {
	ARs map[string]*Entity
	es  map[string]*Entity
	vos map[string]*ValueObject
}

type RepoLayer struct {
	Repos map[string]*Repository
}

type ServiceLayer struct {
	Providers map[string]*Provider
	Services  map[string]*Service
}

type ApiLayer struct {
	Apis map[string]*Api
}

type GatewayLayer struct {
	GateWays map[string]*GateWay
}

func NewBCModel() *BCModel {
	return &BCModel{Layers: make(map[string]*Layer)}
}

func (layer *DomainLayer) Add(name, comment string) {
	if comment == "AR" {
		layer.ARs[name] = NewEntity(name)
	}

	if comment == "E" {
		layer.es[name] = NewEntity(name)
	}

	if comment == "VO" {
		layer.vos[name] = NewValueObject(name)
	}
}

func (layer *RepoLayer) Add(name, comment string) {
	if comment == "Repo" {
		layer.Repos[name] = NewRepository(name)
	}
}

func (layer *GatewayLayer) Add(name, comment string) {
	if comment == "Provider" {
		layer.GateWays[name] = &GateWay{name: name}
	}
}

func (layer *ServiceLayer) Add(name, comment string) {
	if comment == "Service" {
		layer.Services[name] = &Service{name: name}
	}
	if comment == "Provider" {
		layer.Providers[name] = &Provider{name: name}
	}
}

func (layer *ApiLayer) Add(name, comment string) {
	if comment == "Api" {
		layer.Apis[name] = &Api{name: name}
	}
}

func newLayer(name string) layer {
	if name == "domain" {
		return &DomainLayer{ARs: make(map[string]*Entity),
			es:  make(map[string]*Entity),
			vos: make(map[string]*ValueObject)}
	}
	if name == "repositories" {
		return &RepoLayer{Repos: make(map[string]*Repository)}
	}

	if name == "gateways" {
		return &GatewayLayer{GateWays: make(map[string]*GateWay)}
	}

	if name == "services" {
		return &ServiceLayer{Services: make(map[string]*Service), Providers: make(map[string]*Provider)}
	}

	if name == "api" {
		return &ApiLayer{Apis: make(map[string]*Api)}
	}

	return nil
}

func (model *BCModel) AppendLayer(name string) {
	if _, ok := model.Layers[name]; !ok {
		model.Layers[name] = &Layer{name: name, nodes: make(map[string]string), layer: newLayer(name)}
	}
}

func (model *BCModel) AppendNode(layerName, nodeName string) {
	model.Layers[layerName].nodes[nodeName] = layerName
}

func (model *BCModel) findLayer(nodeName string) *Layer {
	for key := range model.Layers {
		layer := model.Layers[key]
		if _, ok := layer.nodes[nodeName]; ok {
			return layer
		}
	}
	return nil
}
func (model *BCModel) AddNode(name, comment string) {
	layer := model.findLayer(name)
	layer.layer.Add(name, comment)
}

func (model *BCModel) AddRepoToLayer(layerName string, repo *Repository) {
	layer := model.Layers[layerName]
	layer.layer.(*RepoLayer).Repos[repo.name] = repo
}

func (model *BCModel) AddARToLayer(layerName string, ar *Entity) {
	layer := model.Layers[layerName]
	layer.layer.(*DomainLayer).ARs[ar.name] = ar
}

func (model *BCModel) AddServiceToLayer(layerName string, service *Service) {
	layer := model.Layers[layerName]
	layer.layer.(*ServiceLayer).Services[service.name] = service
}

func (layer *DomainLayer) addEntityRelations(src string, dsts []string) {
	entity, ok := layer.es[src]
	if !ok {
		return
	}
	for _, dst := range dsts {

		if et, ok := layer.es[dst]; ok {
			entity.Entities = append(entity.Entities, et)
		}
		if vo, ok := layer.vos[dst]; ok {
			entity.VOs = append(entity.VOs, vo)
		}
	}
}
func (layer *DomainLayer) addAggregateRootRelations(src string, dsts []string) {
	ar, ok := layer.ARs[src]
	if !ok {
		return
	}
	for _, dst := range dsts {

		if ref, ok := layer.ARs[dst]; ok {
			ar.Refs = append(ar.Refs, ref)
		}
		if et, ok := layer.es[dst]; ok {
			ar.Entities = append(ar.Entities, et)
		}
		if vo, ok := layer.vos[dst]; ok {
			ar.VOs = append(ar.VOs, vo)
		}
	}
}

func (layer *DomainLayer) addRelations(src string, dsts []string) {
	layer.addAggregateRootRelations(src, dsts)
	layer.addEntityRelations(src, dsts)
}

func (layer *RepoLayer) addRelations(src string, dsts []string) {
	repo, ok := layer.Repos[src]
	if ok {
		for _, dst := range dsts {
			repo.For = dst
		}
	}
}

func (layer *GatewayLayer) addRelations(src string, dsts []string) {
	gateway, ok := layer.GateWays[src]
	if ok {
		for _, dst := range dsts {
			gateway.Implement = dst
		}
	}
}

func (layer *ServiceLayer) addRelations(src string, dsts []string) {
	service, ok := layer.Services[src]
	if ok {
		for _, dst := range dsts {
			service.Refs = append(service.Refs, dst)
		}
	}
}

func (layer *ApiLayer) addRelations(src string, dsts []string) {
	api, ok := layer.Apis[src]
	if ok {
		for _, dst := range dsts {
			api.Refs = append(api.Refs, dst)
		}
	}
}

func (layer *DomainLayer) compare(other interface{}) bool {
	o := (other).(*DomainLayer)
	if len(o.ARs) != len(layer.ARs) {
		return false
	}
	for key := range layer.ARs {
		ar := layer.ARs[key]
		if !ar.Compare(o.ARs[key]) {
			return false
		}
	}
	return true
}

func (layer *RepoLayer) compare(other interface{}) bool {
	o := (other).(*RepoLayer)
	if len(o.Repos) != len(layer.Repos) {
		return false
	}
	for key := range layer.Repos {
		repo := layer.Repos[key]
		if !repo.Compare(o.Repos[key]) {
			return false
		}
	}
	return true
}

func (layer *ServiceLayer) compare(other interface{}) bool {
	o := (other).(*ServiceLayer)
	if len(o.Services) != len(layer.Services) {
		return false
	}
	for key := range layer.Services {
		service := layer.Services[key]
		if !service.compare(o.Services[key]) {
			return false
		}
	}
	return true
}

func (layer *GatewayLayer) compare(other interface{}) bool {
	return true
}

func (layer *ApiLayer) compare(other interface{}) bool {
	return true
}

func (model *BCModel) AddRelations(src string, dsts []string) {
	layer := model.findLayer(src)
	layer.layer.addRelations(src, dsts)
}

func (service *Service) compare(other *Service) bool {
	if len(service.Refs) != len(other.Refs) {
		return false
	}
	return true
}
func (layer *Layer) Compare(other *Layer) bool {
	return layer.layer.compare(other.layer)
}
func (model *BCModel) Compare(other *BCModel) error {
	if len(model.Layers) != len(other.Layers) {
		return errors.New("diff layer number")
	}
	for key := range model.Layers {
		layer := model.Layers[key]
		if !layer.Compare(other.Layers[key]) {
			return errors.New("layer: " + key + " is diff")
		}
	}
	return nil
}
