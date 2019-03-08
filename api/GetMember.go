package api

import (
	"fmt"
	"net/http"
	"tinkodnev/engine"
	"tinkodnev/utils"
)

var GetMemberMethod = func(w http.ResponseWriter, r *http.Request) {
	id, found := utils.RequireU64("id", r, w)
	resp := struct {
		Success bool        `json:"success"`
		Item    *engine.Member `json:"member,omitempty"`
		Error string `json:"error,omitempty"`
	}{}

	if !found {
		resp.Success = false;
		resp.Item = nil
		resp.Error = "Expected id parameter as integer"
	} else {
		member, err := engine.Database.Get(id)
		if err != nil {
			resp.Item = nil
			resp.Success = false
			resp.Error = "Not found"
			fmt.Println(err)
		} else {
			resp.Item = &member
			resp.Success = true
			resp.Error = ""
		}
	}

	utils.SendResponse(resp, w)
}
