#include "code.h"

int main(int argc, char const *argv[])
{
	// routerFactory.setRouter(new FakeRouter());
	router = new FakeRouter();
	AggregateRootARepo *repo = new AggregateRootARepo();
	repo->Save(new AggregateRootA());
	return 0;
}