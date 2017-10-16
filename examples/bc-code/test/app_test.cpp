#include <gtest/gtest.h>
#include "../include/interface/api.h"
#include "../include/services/service.h"
#include "../include/repositories/repository.h"

struct StubCargoProvider : services::CargoProvider{
    int cargo_id;
    int after_days;
    virtual void Confirm(Cargo *cargo) override;
};
StubCargoProvider* provider = new StubCargoProvider();

api::Api* createApi()  {
    repositories::CargoRepository* cargoRepo = new repositories::CargoRepository();

    services::CargoService* service = new services::CargoService(cargoRepo, provider);
    api::Api* api = new api::Api(service);
    return api;
}

TEST(bc_demo_test, create_cargo)
{
    api::Api* api = createApi();
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = 1;
    msg->AfterDays = 10;
    api->CreateCargo(msg);
    EXPECT_EQ(msg->Id, provider->cargo_id);
    EXPECT_EQ(10, provider->after_days);
}

TEST(bc_demo_test, delay_cargo)
{
    api::Api* api = createApi();
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = 1;
    msg->AfterDays = 10;
    api->CreateCargo(msg);
    api->Delay(1,2);
    EXPECT_EQ(msg->Id, provider->cargo_id);
    EXPECT_EQ(12, provider->after_days);
}


void StubCargoProvider::Confirm(Cargo *cargo) {
    this->cargo_id = cargo->getId();
    this->after_days = cargo->afterDays();
}
