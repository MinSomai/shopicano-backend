package queue

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/shopicano/shopicano-backend/machinery"
	tasks2 "github.com/shopicano/shopicano-backend/tasks"
	"time"
)

func SendOrderDetailsEmail(orderID string) error {
	now := time.Now().Add(time.Minute * 1)

	sig := &tasks.Signature{
		Name: tasks2.SendOrderDetailsEmailTaskName,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: orderID,
				Name:  "orderID",
			},
		},
		ETA: &now,
	}
	_, err := machinery.RabbitMQConnection().SendTask(sig)
	if err != nil {
		return err
	}
	return nil
}
