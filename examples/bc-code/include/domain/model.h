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
	int afterDays();
private:
	Delivery* delivery;
};

struct CargoDelayed {
	int CargoId;
};
}
#endif