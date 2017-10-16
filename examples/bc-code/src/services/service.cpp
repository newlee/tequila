#include "services/service.h"

using namespace services;

CargoService::CargoService(CargoRepository* cargoRepo, CargoProvider* cargoProvider)
    :cargoRepository_(cargoRepo)
    ,cargoProvider_(cargoProvider)
{

}

void CargoService::Create(int id, int days)
{
    Delivery* delivery = new Delivery(days);
    Cargo* cargo = new Cargo(delivery, id);
    this->cargoRepository_->Save(cargo);
    this->cargoProvider_->Confirm(cargo);
}
void CargoService::Delay(int id, int days)
{
    Cargo* cargo = cargoRepository_->FindById(id);
    if(cargo != NULL) {
        cargo->Delay(days);
        cargoRepository_->Save(cargo);
        this->cargoProvider_->Confirm(cargo);
    }
}