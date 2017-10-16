#include "repositories/repository.h"

using namespace repositories;

CargoRepository::CargoRepository()
{
}

void CargoRepository::Save(Cargo* cargo)
{
    this->cargo_list.push_back(cargo);
}

Cargo* CargoRepository::FindById(int id)
{
    for(Cargo* cargo : this->cargo_list){
        if(cargo->getId() == id) {
            return cargo;
        }
    }
    return NULL;
}
