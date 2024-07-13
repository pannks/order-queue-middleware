package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"lambda/db"
	"lambda/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteOrderProcessHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	riderID := r.URL.Query().Get("rider_id")

	if orderID == "" || riderID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Order ID and Rider ID are required", "No Order ID or Rider ID provided")
		return
	}

	database := db.Client.Database(os.Getenv("MONGODB_DATABASE"))
	orderProcess := database.Collection("orderProcess")

	res, err := orderProcess.DeleteOne(context.Background(), bson.M{"order_id": orderID, "rider_id": riderID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting order", err.Error())
		return
	}

	if res.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found", "No order matched the given criteria")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Order deleted from orderProcess successfully")
}
