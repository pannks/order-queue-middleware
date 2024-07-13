package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"lambda/db"
	"lambda/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOrderQueueHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	order.ID = primitive.NewObjectID()
	order.Timestamp = time.Now()

	database := db.Client.Database(os.Getenv("MONGODB_DATABASE"))
	orderQueue := database.Collection("orderQueue")

	_, err = orderQueue.InsertOne(context.Background(), order)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error inserting order", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Order inserted into orderQueue successfully")
}
