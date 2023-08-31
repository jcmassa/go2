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

type rankitem struct { //Los campos todos arrancan con mayuscula porque obvio
	Id       int32  `json:"id"`
	Titulo   string `json:"titulo"`
	ImageId  int32  `json:"imageId"`
	Ranking  int32  `json:"ranking"`
	ItemType int32  `json:"itemType"`
	TierId   int32  `json:"tierId"`
}

func GetAllItemsDB() ([]rankitem, error) {
	db, err := DbHandle.SetDBConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id, titulo, imageId, ranking, itemType, tierId FROM Item") // WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//var returnItems []rankitem
	returnItems := make([]rankitem, 0)
	for rows.Next() {
		var nwItem rankitem
		if err := rows.Scan(&nwItem.Id, &nwItem.Titulo, &nwItem.ImageId, &nwItem.Ranking, &nwItem.ItemType, &nwItem.TierId); err != nil {
			return returnItems, err
		}
		returnItems = append(returnItems, nwItem)
	}
	if err = rows.Err(); err != nil {
		return returnItems, err
	}
	if err != nil {
		return returnItems, err
		//json.NewEncoder(w).Encode(err)
	}
	return returnItems, err
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	returnItems, err := GetAllItemsDB()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	json.NewEncoder(w).Encode(returnItems)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idItem, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	items, _ := GetAllItemsDB()
	for _, itm := range items {
		if itm.Id == int32(idItem) {
			json.NewEncoder(w).Encode(itm)
			//fmt.Fprintf(w, itm.titulo)
		}
	}
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	var updatedItem rankitem
	iditm, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	if iditm < 1 {
		fmt.Fprintf(w, "no existe")
		return
	}

	json.Unmarshal(reqBody, &updatedItem)
	updCMD := fmt.Sprintf("UPDATE Item SET ranking = %d, tierId = %d WHERE Id = %d", updatedItem.Ranking, updatedItem.TierId, updatedItem.Id)
	err = DbHandle.RunCommand(updCMD)
	if err != nil {
		fmt.Fprintf(w, "error al actualizar")
		return
	}
	fmt.Fprintf(w, "Se actualizo el item %v", iditm)

}

func ResetType(w http.ResponseWriter, r *http.Request) {
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idType, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	resetTypeCMD := fmt.Sprintf("EXEC dbo.ResetType %d", idType)
	err = DbHandle.RunCommand(resetTypeCMD)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	return
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idItem, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	delCMD := fmt.Sprintf("DELETE FROM Item WHERE Id =%d", idItem)
	err = DbHandle.RunCommand(delCMD)
	if err != nil {
		fmt.Fprintf(w, "error al actualizar")
		return
	}
	fmt.Fprintf(w, "Se elimino el item %v", idItem)
}

func AssignTier(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	AddCORSHeader(&w)
	paramVars := mux.Vars(r)
	idItem, err := strconv.Atoi(paramVars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	idTier, errTier := strconv.Atoi(paramVars["tierId"])
	if errTier != nil {
		fmt.Fprintf(w, errTier.Error())
		return
	}
	ranking, errRank := strconv.Atoi(paramVars["ranking"])
	if errRank != nil {
		fmt.Fprintf(w, errRank.Error())
		return
	}

	updCMD := fmt.Sprintf("UPDATE Item SET ranking = %d , tierId = %d WHERE Id =  %d", ranking, idTier, idItem)
	err = DbHandle.RunCommand(updCMD)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	//fmt.Fprintf(w, "Se consulta el item %d tier %v ranking %d", idItem, idTier, ranking)
	fmt.Fprintf(w, "Se actualizo el item %v", idItem)
}
