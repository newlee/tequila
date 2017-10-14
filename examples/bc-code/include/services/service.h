#ifndef BC_DEMO_SERVICE_H
#define BC_DEMO_SERVICE_H

#include "domain/model.h"
#include "repositories/repository.h"

using namespace domain;
using namespace repositories;

namespace services {
struct CargoProvider : Provider {
    virtual void Confirm(Cargo* cargo){};
};

struct CargoService {
    CargoService(CargoRepository* cargoRepo, CargoProvider* cargoProvider);
    void Create(int id, int days);
    void Delay(int id, int days);
private:
    CargoRepository* cargoRepository_;
    CargoProvider* cargoProvider_;
};


}
#endif //BC_DEMO_SERVICE_H
