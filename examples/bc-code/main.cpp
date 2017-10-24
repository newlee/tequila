#include <iostream>
#include "repositories/repository.h"
#include "services/service.h"
#include "gateways/gateway.h"
#include "interface/api.h"
#include "Hypodermic/ContainerBuilder.h"

using namespace repositories;
using namespace services;
using namespace gateways;
using namespace Hypodermic;

int main(int argc, char *argv[])
{
    ContainerBuilder builder;
    builder.registerType< CargoRepository >().singleInstance();
    builder.registerType< CargoProviderImpl >().singleInstance().as< CargoProvider >();
    builder.registerType< CargoService >().singleInstance();
    builder.registerType<api::Api>().singleInstance();

    auto container = builder.build();

    std::shared_ptr<api::Api> api = container->resolve<api::Api>();
    std::shared_ptr<CargoRepository> cargoRepo = container->resolve<CargoRepository>();
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = 1;
    msg->AfterDays = 10;
    api->CreateCargo(msg);
    std::cout<< "hello" << "\n";
    std::cout<< cargoRepo->FindById(1)->getId()<<"\n";
    return 0;
}