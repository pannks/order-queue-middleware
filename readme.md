## Get Started

** Only For Local Development **

> This is service for recieve many order in time and process booking order
> for client who come first as queue timestamp
> and to protect bottle head to database

Stack: Go + MongoDB

## API Route

BASE_URL : http://localhost:8080

-   new Queue : _POST_ /api/orderQueue
-   regis Queue: _POST_ /api/bookOrder
-   kill Queue: _DELETE_ /api/orderProcess
-   reject Queue: _POST_ /api/rejectOrder
