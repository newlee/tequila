#include "interface/api.h"

using namespace api;
using namespace services;

Api::Api(std::shared_ptr<CargoService> cargoService)
    :cargoService_(cargoService)
{
}

void Api::CreateCargo(CreateCargoMsg * msg)
{
    this->cargoService_->Create(msg->Id,msg->AfterDays);
}

void Api::Delay(int cargoId, int days) {
    this->cargoService_->Delay(cargoId,days);
}
