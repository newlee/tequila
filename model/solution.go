package model

import (
	"errors"
)

type layer interface {
	Add(name, comment string)
	getNodes() []string
	addRelations(src string, dsts []string)
	getRelations() map[string][]string
	compare(other interface{}) bool
}

type Layer struct {
	Name  string
	nodes map[string]string
	layer layer
}

func (layer *Layer) GetNodes() []string {
	return layer.layer.getNodes()
}

func (layer *Layer) GetRelations() map[string][]string {
	return layer.layer.getRelations()
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

func (layer *DomainLayer) getNodes() []string {
	result := make([]string, 0)
	for key := range layer.ARs {
		result = append(result, layer.ARs[key].name)
	}
	for key := range layer.es {
		result = append(result, layer.es[key].name)
	}

	for key := range layer.vos {
		result = append(result, layer.vos[key].name)
	}
	return result
}

func (layer *ServiceLayer) getNodes() []string {
	result := make([]string, 0)
	for key := range layer.Services {
		result = append(result, layer.Services[key].name)
	}
	for key := range layer.Providers {
		result = append(result, layer.Providers[key].name)
	}
	return result
}

func (layer *RepoLayer) getNodes() []string {
	result := make([]string, 0)
	for key := range layer.Repos {
		result = append(result, layer.Repos[key].name)
	}
	return result
}

func (layer *GatewayLayer) getNodes() []string {
	result := make([]string, 0)
	return result
}

func (layer *ApiLayer) getNodes() []string {
	result := make([]string, 0)
	return result
}

func (layer *DomainLayer) getRelations() map[string][]string {
	result := make(map[string][]string)
	for key := range layer.ARs {
		ar := layer.ARs[key]
		result[ar.name] = make([]string, 0)
		for _, entity := range ar.Entities {
			result[ar.name] = append(result[ar.name], entity.name)
		}
		for _, vo := range ar.VOs {
			result[ar.name] = append(result[ar.name], vo.name)
		}
	}
	return result
}
func (layer *RepoLayer) getRelations() map[string][]string {
	result := make(map[string][]string)
	for key := range layer.Repos {
		repo := layer.Repos[key]
		result[repo.name] = make([]string, 0)
		result[repo.name] = append(result[repo.name], repo.For)
	}
	return result
}
func (layer *ServiceLayer) getRelations() map[string][]string {
	result := make(map[string][]string)
	for key := range layer.Services {
		service := layer.Services[key]
		result[service.name] = service.Refs
	}

	return result
}
func (layer *ApiLayer) getRelations() map[string][]string {
	result := make(map[string][]string)

	return result
}
func (layer *GatewayLayer) getRelations() map[string][]string {
	result := make(map[string][]string)

	return result
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
		model.Layers[name] = &Layer{Name: name, nodes: make(map[string]string), layer: newLayer(name)}
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
	domainLayer := layer.layer.(*DomainLayer)
	domainLayer.ARs[ar.name] = ar
	//TODO: recursive entitys
	for _, entity := range ar.Entities {
		domainLayer.es[entity.name] = entity
	}
	for _, vo := range ar.VOs {
		domainLayer.vos[vo.name] = vo
	}
}

func (model *BCModel) AddServiceToLayer(layerName string, service *Service) {
	layer := model.Layers[layerName]
	layer.layer.(*ServiceLayer).Services[service.name] = service
}

func (model *BCModel) AddProviderToLayer(layerName string, provider *Provider) {
	layer := model.Layers[layerName]
	layer.layer.(*ServiceLayer).Providers[provider.name] = provider
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
