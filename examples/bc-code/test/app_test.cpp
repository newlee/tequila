#include <gtest/gtest.h>
#include "../include/interface/api.h"
#include "../include/services/service.h"
#include "../include/repositories/repository.h"

struct StubCargoProvider : services::CargoProvider{
    int cargo_id;
    virtual void Confirm(Cargo *cargo) override;
};
TEST(bc_demo_test, create_cargo)
{
    repositories::CargoRepository* cargoRepo = new repositories::CargoRepository();
    StubCargoProvider* provider = new StubCargoProvider();
    services::CargoService* service = new services::CargoService(cargoRepo, provider);
    api::Api* api = new api::Api(service);
    api::CreateCargoMsg* msg = new api::CreateCargoMsg();
    msg->Id = 1;
    msg->AfterDays = 10;
    api->CreateCargo(msg);
    EXPECT_EQ(msg->Id, provider->cargo_id);
}

void StubCargoProvider::Confirm(Cargo *cargo) {
    cargo_id = cargo->getId();
}
