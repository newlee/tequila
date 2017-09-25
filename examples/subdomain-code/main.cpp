#include "code.h"

using namespace subdomain1::Gateways;
using namespace subdomain1::Repositories;

int main(int argc, char const *argv[])
{
	router = new FakeRouter();
	AggregateRootA* a = new AggregateRootA();
	a->Init();
	return 0;
}