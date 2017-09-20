package domain;


public class AggregateRootA extends AggregateRoot {
 
    private EntityB entity_b;
    private ValueObjectC vo_c; 
    private Router router;
    public AggregateRootA(Router router) {
        this.router = router;
    }
    public void init() {
        router.select();
    }
	 
}