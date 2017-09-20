package repositories;

import domain.*;


public class AggregateRootARepo extends Repository {
    public void save(AggregateRootA a){
        a.init();
        System.out.println("saved\n");		
	};
}