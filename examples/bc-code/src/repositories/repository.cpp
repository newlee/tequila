#include "repositories/repository.h"

using namespace repositories;

std::vector<Cargo*> cargo_list;

CargoRepository::CargoRepository()
{
}

void CargoRepository::Save(Cargo* cargo)
{
    cargo_list.push_back(cargo);
}

Cargo* CargoRepository::FindById(int id)
{
    for(Cargo* cargo : cargo_list){
        if(cargo->getId() == id) {
            return cargo;
        }
    }

    return NULL;
}
