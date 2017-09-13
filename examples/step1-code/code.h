public class Entity
{
public:
	Entity();
	~Entity();
	int Id;
};

public class AggregateRoot: public Entity
{
public:
	AggregateRoot();
	~AggregateRoot();
	
};


public class ValueObject
{
public:
	ValueObject();
	~ValueObject();
	
};



public class ValueObjectC: public ValueObject
{
public:
	ValueObjectC();
	~ValueObjectC();
	
};

public class ValueObjectD: public ValueObject
{
public:
	ValueObjectD();
	~ValueObjectD();
	
};

public class EntityB: public Entity
{
public:
	EntityB();
	~EntityB();
	ValueObjectD* vo_d; 
};

public class AggregateRootA: public AggregateRoot
{
public:
	AggregateRootA();
	~AggregateRootA();
	EntityB* entity_b;
	ValueObjectC* vo_c; 
};

