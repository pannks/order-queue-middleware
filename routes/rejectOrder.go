package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"lambda/db"
	"lambda/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RejectOrderHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OrderID string `json:"order_id"`
		RiderID string `json:"rider_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input", err.Error())
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
		err := orderProcess.FindOne(sc, bson.M{"order_id": request.OrderID, "rider_id": request.RiderID}).Decode(&result)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		result.RiderID = ""
		_, err = orderQueue.InsertOne(sc, bson.M{
			"_id":       result.ID,
			"order_id":  result.OrderID,
			"rider_id":  nil,
			"timestamp": result.Timestamp,
		})
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		_, err = orderProcess.DeleteOne(sc, bson.M{"order_id": request.OrderID, "rider_id": request.RiderID})
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error processing order rejection", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Order rejected successfully")
}
