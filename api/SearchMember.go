package api

import (
	"fmt"
	"net/http"
	"tinkodnev/engine"
	"tinkodnev/utils"
)

var SearchMemberMethod = func(w http.ResponseWriter, r *http.Request) {
	data, found := utils.RequireString("query", r, w)
	resp := struct {
		Success bool            `json:"success"`
		Items   []engine.Member `json:"members,omitempty"`
		Error   string          `json:"error,omitempty"`
	}{}

	if !found {
		resp.Success = false
		resp.Error = "Expected query parameter as a string"
	} else {
		res, err := engine.Database.Search(data, 10)
		if err != nil {
			fmt.Println(err)
			resp.Success = false
			resp.Error = "Not found"
			resp.Items = nil
		} else {
			if len(res) == 0 {
				resp.Success = false
				resp.Error = "Not found"
			} else {
				fmt.Print(res)
				resp.Success = true
				resp.Items = res
			}
		}
	}

	utils.SendResponse(resp, w)
}
