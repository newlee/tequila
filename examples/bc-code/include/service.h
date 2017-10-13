#ifndef BC_DEMO_SERVICE_H
#define BC_DEMO_SERVICE_H

#include "model.h"
#include "repository.h"

using namespace domain;
using namespace repositories;

namespace services {
struct CargoService {
    CargoService(CargoRepository*);
    void Create(int id, int days);
    void Delay(int id, int days);
private:
    CargoRepository* cargoRepository_;
};
}
#endif //BC_DEMO_SERVICE_H
