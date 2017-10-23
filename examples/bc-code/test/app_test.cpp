#include <gtest/gtest.h>
#include "../include/interface/api.h"
#include "../include/services/service.h"
#include "../include/repositories/repository.h"
#include "Hypodermic/ContainerBuilder.h"

using namespace Hypodermic;

struct StubCargoProvider : services::CargoProvider{
    int cargo_id;
    int after_days;
    virtual void Confirm(Cargo *cargo) override;
};

static const int ID = 1;
static const int AFTER_DAYS = 10;

//StubCargoProvider* provider = new StubCargoProvider();
auto provider = std::make_shared< StubCargoProvider >();

api::Api* createApi()  {
    ContainerBuilder builder;
    builder.registerType< CargoRepository >().singleInstance();
    builder.registerInstance(provider).as<CargoProvider>();
    builder.registerType< CargoService >().singleInstance();
    builder.registerType<api::Api>().singleInstance();

    auto container = builder.build();

    std::shared_ptr<api::Api> api = container->resolve<api::Api>();

    return api.get();
}

void createCargo(api::CreateCargoMsg* msg) {
    api::Api* api = createApi();
    api->CreateCargo(msg);
}

TEST(bc_demo_test, create_cargo)
{
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = ID;
    msg->AfterDays = AFTER_DAYS;
    createCargo(msg);
    EXPECT_EQ(msg->Id, provider->cargo_id);
    EXPECT_EQ(msg->AfterDays, provider->after_days);
}

TEST(bc_demo_test, delay_cargo)
{
    api::Api* api = createApi();
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = ID;
    msg->AfterDays = AFTER_DAYS;
    api->CreateCargo(msg);
    api->Delay(ID,2);
    EXPECT_EQ(ID, provider->cargo_id);
    EXPECT_EQ(12, provider->after_days);
}


void StubCargoProvider::Confirm(Cargo *cargo) {
    this->cargo_id = cargo->getId();
    this->after_days = cargo->afterDays();
}
