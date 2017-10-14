#ifndef BC_DEMO_GATEWAY_H
#define BC_DEMO_GATEWAY_H

#include "services/service.h"

using namespace services;

namespace gateways {
struct CargoProviderImpl: CargoProvider
{
    virtual void Confirm(Cargo *cargo) override;
};
}
#endif //BC_DEMO_GATEWAY_H
