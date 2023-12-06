package entity

type Order_Queue struct {
	Order []*Order
}

func (order_queue *Order_Queue) Less(i, j int) bool {
	return order_queue.Order[i].Price > order_queue.Order[j].Price
}

func (order_queue *Order_Queue) Swap(i, j int) {
	order_queue.Order[i], order_queue.Order[j] = order_queue.Order[j], order_queue.Order[i]
}

func (order_queue *Order_Queue) Len() int {
	return len(order_queue.Order)
}

func (order_queue *Order_Queue) Push(x any) {
	order_queue.Order = append(order_queue.Order, x.(*Order))
}

func (order_queue *Order_Queue) Pop() any {
	old := order_queue.Order
	n := len(old)
	x := old[n-1]
	order_queue.Order = old[0 : n-1]
	return x
}

func NewOrderQueue() *Order_Queue {
	return &Order_Queue{}
}
