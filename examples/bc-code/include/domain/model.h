#ifndef BC_DEMO__MODEL_H__
#define BC_DEMO__MODEL_H__

namespace domain{
struct Entity
{
	int getId();
protected:
    int id;
};

struct AggregateRoot: Entity
{

};

struct ValueObject
{

};

struct Provider
{

};

struct Delivery: ValueObject
{
	Delivery(int);

	int AfterDays;
};

struct Cargo: AggregateRoot
{
	Cargo(Delivery*, int);
	~Cargo();
	void Delay(int);
private:
	Delivery* delivery;
};

}
#endif