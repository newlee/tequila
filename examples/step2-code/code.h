#include <iostream>

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

};

class AggregateRootA: public AggregateRoot
{
public:
	AggregateRootA(){};
	~AggregateRootA(){};
	EntityB* entity_b;
	ValueObjectC* vo_c; 
	void Init(){
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
