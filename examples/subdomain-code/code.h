#include <iostream>

namespace subdomain1{
namespace Domain {
class Entity
{
public:
	int Id;
};

class AggregateRoot: public Entity
{
};


class ValueObject
{
};

class Provider
{

};

class Router: public Provider
{
public:
	virtual int Selete() = 0;
};

static Router* router;

class ValueObjectC: public ValueObject
{
public:
	ValueObjectC(){};
	~ValueObjectC(){};
	
};

class ValueObjectD: public ValueObject
{
public:
	ValueObjectD(){};
	~ValueObjectD(){};
	
};

class EntityB: public Entity
{
public:
	EntityB(){};
	~EntityB(){};
	ValueObjectD* vo_d; 
	void init(){
		vo_d = new ValueObjectD();
		std::cout << "entity b init" << "\n";
	};
};

class AggregateRootA: public AggregateRoot
{
public:
	AggregateRootA(){};
	~AggregateRootA(){};
	EntityB* entity_b;
	ValueObjectC* vo_c; 
	void Init(){
		entity_b = new EntityB();
		entity_b->init();
		router->Selete();
	};
};

class AggregateRootB: public AggregateRoot
{
public:
	AggregateRootB(){};
	~AggregateRootB(){};
	AggregateRootA* a;
};


}

namespace Repositories {
using namespace Domain;

class Repository{
};

class AggregateRootARepo: public Repository
{
public:
	AggregateRootARepo(){};
	~AggregateRootARepo(){};
	void Save(AggregateRootA *a){
		a->Init();
		std::cout << "saved" << "\n";
	};
};

}

namespace Gateways {
using namespace Domain;

class FakeRouter: public Router
{
public:
	FakeRouter(){};
	~FakeRouter(){};
	int Selete(){
		std::cout << "routed" << "\n";
		return 1;
	}
};

}
}

namespace subdomain2 {
namespace Domain {	
class Entity
{
public:
	int Id;
};

class AggregateRoot: public Entity
{
};


class ValueObject
{
};
class EntityC: public Entity
{
public:
	EntityC(){};
	~EntityC(){};
};

class AggregateRootC: public AggregateRoot
{
public:
	AggregateRootC(){};
	~AggregateRootC(){};
	EntityC* entity_c;
};
}
}