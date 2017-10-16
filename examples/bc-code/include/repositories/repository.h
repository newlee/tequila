#ifndef BC_DEMO__REPO_H__
#define BC_DEMO__REPO_H__

#include "domain/model.h"
#include <vector>

using namespace domain;

namespace repositories {
struct Repository
{

};

struct CargoRepository: Repository
{
    CargoRepository();
    void Save(Cargo* cargo);
    Cargo* FindById(int id);
private:
    std::vector<Cargo*> cargo_list;
};

}

#endif