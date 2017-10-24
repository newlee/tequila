#ifndef BC_DEMO_SERVICE_H
#define BC_DEMO_SERVICE_H

#include "repositories/repository.h"

using namespace repositories;

namespace services {
struct Provider
{

};

struct CargoProvider : Provider {
    virtual void Confirm(Cargo* cargo){};
};

struct CargoService {
    explicit CargoService(const std::shared_ptr<CargoRepository> cargoRepo, const std::shared_ptr<CargoProvider> cargoProvider)
        :cargoRepository_(cargoRepo)
        ,cargoProvider_(cargoProvider)
    {

    }
    void Create(int id, int days);
    void Delay(int id, int days);
private:
    std::shared_ptr<CargoRepository> cargoRepository_;
    std::shared_ptr<CargoProvider> cargoProvider_;
};

}
#endif //BC_DEMO_SERVICE_H
