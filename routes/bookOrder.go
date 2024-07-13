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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BookOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	database := db.Client.Database(os.Getenv("MONGODB_DATABASE"))
	orderQueue := database.Collection("orderQueue")
	orderProcess := database.Collection("orderProcess")

	session, err := db.Client.StartSession()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to start session", err.Error())
		return
	}

	defer session.EndSession(context.Background())

	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		var result Order
		err := orderQueue.FindOne(sc, bson.M{"order_id": order.OrderID, "rider_id": nil}).Decode(&result)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		result.RiderID = order.RiderID
		result.Timestamp = time.Now()
		_, err = orderProcess.InsertOne(sc, result)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		_, err = orderQueue.DeleteOne(sc, bson.M{"order_id": order.OrderID})
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error processing order", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Order booked successfully")
}
