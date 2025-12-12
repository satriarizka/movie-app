package response

type DailyRevenueResponse struct {
	Date        string  `json:"date"`
	TotalAmount float64 `json:"total_amount"`
	Count       int64   `json:"transaction_count"`
}

type TopMovieResponse struct {
	MovieID    string  `json:"movie_id"`
	Title      string  `json:"title"`
	TotalSold  int64   `json:"total_tickets_sold"`
	TotalSales float64 `json:"total_sales_revenue"`
}
