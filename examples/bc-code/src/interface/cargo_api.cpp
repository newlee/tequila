#include "interface/api.h"

using namespace api;
using namespace services;

Api::Api(CargoService*cargoService)
    :cargoService_(cargoService)
{
}

void Api::CreateCargo(CreateCargoMsg * msg)
{
    this->cargoService_->Create(msg->Id,msg->AfterDays);
}