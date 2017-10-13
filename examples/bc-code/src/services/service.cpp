#include "service.h"

using namespace services;

CargoService::CargoService(CargoRepository* cargoRepository)
    :cargoRepository_(cargoRepository)
{

}

void CargoService::Create(int id, int days)
{
    Delivery* delivery = new Delivery(10);
    Cargo* cargo = new Cargo(delivery, 1);
    this->cargoRepository_->Save(cargo);
}
void CargoService::Delay(int id, int days)
{
    Cargo* cargo = cargoRepository_->FindById(id);
    if(cargo != NULL) {
        cargo->Delay(days);
        cargoRepository_->Save(cargo);
    }
}