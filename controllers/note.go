package controllers

import (
	"log"
	"net/http"
	"pratice/config"
	"pratice/structs"

	"github.com/gin-gonic/gin"
)

func GetNote(c *gin.Context) {
	var notes structs.Note
	var arr_note []structs.Note

	db := config.Dbconn()
	defer db.Close()

	rows, err := db.Query("SELECT id_note, title, note FROM tb_note")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&notes.Id, &notes.Title, &notes.Note); err != nil {
			log.Fatal(err.Error())

		} else {
			arr_note = append(arr_note, notes)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   arr_note[0],
	})
}
