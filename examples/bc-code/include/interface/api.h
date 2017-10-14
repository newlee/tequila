#ifndef BC_DEMO_API_H
#define BC_DEMO_API_H

#include "services/service.h"
#include "msg.h"

using namespace services;

namespace api {
struct Api {
    Api(CargoService*);
    void CreateCargo(CreateCargoMsg* msg);

private:
    CargoService* cargoService_;
};
}
#endif //BC_DEMO_API_H
