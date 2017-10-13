#include <iostream>
#include "model.h"
#include "repository.h"
#include "service.h"

using namespace domain;
using namespace repositories;
using namespace services;

int main(int argc, char *argv[])
{
    CargoRepository* cargoRepo = new CargoRepository();
    CargoService* service = new CargoService(cargoRepo);
    service->Create(1, 10);
    service->Delay(1,2);

    std::cout<< cargoRepo->FindById(1)->getId()<<"\n";
    return 0;
}