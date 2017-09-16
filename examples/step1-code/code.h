class Entity
{
public:
	Entity(){};
	~Entity(){};
	int Id;
};

class AggregateRoot: public Entity
{
public:
	AggregateRoot(){};
	~AggregateRoot(){};
	void Init(){
		
	}
};


class ValueObject
{
public:
	ValueObject(){};
	~ValueObject(){};
	
};

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
};