#ifndef BC_DEMO__MODEL_H__
#define BC_DEMO__MODEL_H__

#include <vector>

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

struct Product: ValueObject
{
	Product(int id);
private:
	int productId;
};

struct Cargo: AggregateRoot
{
	Cargo(Delivery*, int);
	~Cargo();
	void Delay(int);
	void AddProduct(int);
	int afterDays();
private:
	Delivery* delivery;
	std::vector<Product*> product_list;
};

struct CargoDelayed {
	int CargoId;
};
}
#endif