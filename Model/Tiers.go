package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb" //Va con guion bajo porque obvio
	"github.com/gorilla/mux"

	//https://pkg.go.dev/github.com/denisenkom/go-mssqldb#section-readme
	DbHandle "rankapi/DBHandle"
)

func GetAllTiersDB() ([]tier, error) {
	db, err := DbHandle.SetDBConnection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Tiers") // WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//var returnTiers []tier
	returnTiers := make([]tier, 0)
	for rows.Next() {
		var nwTier tier
		if err := rows.Scan(&nwTier.Id, &nwTier.RowName, &nwTier.RowNumber, &nwTier.NumCells); err != nil {
			return returnTiers, err
		}
		returnTiers = append(returnTiers, nwTier)
	}
	if err = rows.Err(); err != nil {
		return returnTiers, err
	}
	if err != nil {
		return returnTiers, err
		//json.NewEncoder(w).Encode(err)
	}
	return returnTiers, err
}

func GetAllTiers(w http.ResponseWriter, r *http.Request) {
	returnTiers, err := GetAllTiersDB()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	json.NewEncoder(w).Encode(returnTiers)
}

type tier struct {
	Id        int32  `json:"id"`
	RowName   string `json:"rowName"`
	RowNumber int32  `json:"rowNumber"`
	NumCells  int32  `json:"numCells"`
}

func maxId(tierPool []tier) int32 {
	max := int32(1)
	for _, v := range tierPool { //i es el indice, v es el valor
		if v.Id > max {
			max = v.Id
		}
	}
	return int32(max)
}

func GetTier(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idTier, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	tiers, _ := GetAllTiersDB()
	for _, tier := range tiers {
		if tier.Id == int32(idTier) {
			json.NewEncoder(w).Encode(tier)
		}
	}
}

func UpdateTier(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	var updatedTier tier
	idTier, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	if idTier < 1 {
		fmt.Fprintf(w, "no existe")
		return
	}

	json.Unmarshal(reqBody, &updatedTier)
	updCMD := fmt.Sprintf("EXEC dbo.EditTier %d, '%s', %d", idTier, updatedTier.RowName, updatedTier.NumCells)
	err = DbHandle.RunCommand(updCMD)
	if err != nil {
		fmt.Fprintf(w, "error al actualizar")
		return
	}
	fmt.Fprintf(w, "Se actualizo el tier %v", idTier)
}

func DeleteTier(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idTier, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	err = EmptyTier(idTier)
	//emptyTierCMD := fmt.Sprintf("EXEC dbo.EmptyTier %d", idTier)
	//err = DbHandle.RunCommand(emptyTierCMD)
	if err != nil {
		fmt.Fprintf(w, "error al actualizar")
		return
	}
	delCMD := fmt.Sprintf("DELETE FROM Tiers WHERE Id =%d", idTier)
	err = DbHandle.RunCommand(delCMD)
	if err != nil {
		fmt.Fprintf(w, "error al actualizar")
		return
	}
	fmt.Fprintf(w, "Se elimino el tier %v", idTier)
}

func EmptyTier(idTier int) error {
	emptyTierCMD := fmt.Sprintf("EXEC dbo.EmptyTier %d", idTier)
	err := DbHandle.RunCommand(emptyTierCMD)
	return err
}

func CreateTier(w http.ResponseWriter, r *http.Request) {
	var newTier tier
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Ingrese datos validos")
	}
	json.Unmarshal(requestBody, &newTier)

	//tiers, _ := GetAllTiersDB()
	maxId, scalarError := DbHandle.GetScalarVal("SELECT dbo.GetMaxTierIdNumber()")
	if scalarError != nil {
		fmt.Fprintf(w, scalarError.Error())
		return
	}
	newTier.Id = int32(maxId) + 1

	rowNum, scalarError2 := DbHandle.GetScalarVal("SELECT dbo.GetMaxTierRowNumber()")
	if scalarError2 != nil {
		fmt.Fprintf(w, scalarError2.Error())
		return
	}
	newTier.RowNumber = int32(rowNum) + 1

	//tiers = append(tiers, newTier)

	insCMD := fmt.Sprintf("INSERT INTO Tiers (Id, rowName, rowNumber, numCells) VALUES(%d, '%s', %d, %d)", newTier.Id, newTier.RowName, newTier.RowNumber, newTier.NumCells)
	err = DbHandle.RunCommand(insCMD)
	if err != nil {
		fmt.Fprintf(w, "error al insertar")
		return
	}

	//fmt.Fprintf(w, "Se agrego el tier %v", newTier.Id)

	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTier)
}
