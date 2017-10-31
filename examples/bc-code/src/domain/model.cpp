#include "domain/model.h"

using namespace domain;

int Entity::getId()
{
	return id;
}
Cargo::Cargo(Delivery* delivery, int id)
	:delivery(delivery)
{
	this->id = id;
}
Cargo::~Cargo()
{
	
}
void Cargo::Delay(int days)
{
	int after = this->delivery->AfterDays;
	this->delivery = new Delivery(after + days);
}

void Cargo::AddProduct(int productId)
{
	Product* product = new Product(productId);
	this->product_list.push_back(product);
}

int Cargo::afterDays()
{
	return this->delivery->AfterDays;
}

Delivery::Delivery(int afterDays)
	:AfterDays(afterDays)
{

}

Product::Product(int id)
    :productId(id)
{
}