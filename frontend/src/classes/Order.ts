import {OrderItem} from '@/classes/OrderItem'


export default class Order {
    submitted_at: Date
	id:          string
	display_id: string
	items:    OrderItem[] 
	discount: number
	state: string
	started_at: Date
	comment: string

    constructor(){
        this.submitted_at = new Date()
        this.id = ""
        this.display_id = ""
        this.items = []
        this.discount = 0
        this.state = ""
        this.started_at = new Date()
        this.comment = ""
    }
}