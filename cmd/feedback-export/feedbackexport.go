package main

import (
	"fmt"
	"log"

	"github.com/zekrotja/hermans/pkg/database"
)

func main() {
	db, err := database.New("db/orders.sqlite")
	if err != nil {
		log.Fatalf("DB konnte nicht geöffnet werden: %v", err)
	}

	feedbacks, err := db.GetAllFeedback()
	if err != nil {
		log.Fatalf("Feedbacks konnten nicht geladen werden: %v", err)
	}

	if len(feedbacks) == 0 {
		fmt.Println("Kein Feedback in der Datenbank gefunden.")
		return
	}

	fmt.Println("--- Gesammeltes Feedback ---")
	for _, fb := range feedbacks {
		fmt.Printf("[%s] [%s] on %s: %s\n",
			fb.Timestamp.Format("2006-01-02 15:04"),
			fb.Type,
			fb.Page,
			fb.Message)
	}
	fmt.Println("--------------------------")
}
