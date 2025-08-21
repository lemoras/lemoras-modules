package note

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	NoteId     uuid.UUID `json:"noteId"`
	TitleName  string    `json:"titleName"`
	DetailText string    `json:"detailText"`
	Category   int       `json:"category"`

	Http CustomHttp `json:"http"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
	Method       string       `json:"method"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
}

func Invoke(in Request) (*u.Response, error) {

	var resp map[string]interface{}

	context := &u.Context{}
	// res, err := http.Get(os.Getenv("VALID_API_URL"))
	client := &http.Client{}

	req, _ := http.NewRequest("GET", os.Getenv("VALID_API_URL"), &bytes.Buffer{})
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", in.Http.CustomHeader.Authorization)

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		_, errRes := u.ResMessage(false, "0x11130:Missing auth token")
		return &errRes, nil
	}

	context.UserId = uuid.MustParse(res.Header.Get("userId"))
	roleId, _ := strconv.Atoi(res.Header.Get("roleId"))
	context.RoleId = roleId
	appId, _ := strconv.Atoi(res.Header.Get("appId"))
	context.AppId = appId
	context.MerchantId = uuid.MustParse(res.Header.Get("merchantId"))
	hasId, err := strconv.ParseBool(res.Header.Get("hasId"))
	context.HasId = hasId
	projectId, _ := strconv.Atoi(res.Header.Get("projectId"))
	context.ProjectId = projectId
	context.CustomData = res.Header.Get("customData")
	initCompleted, err := strconv.ParseBool(res.Header.Get("initCompleted"))
	context.InitCompleted = initCompleted

	tokenRoleLevel := u.GetRoleLevel(context.RoleId)

	if context.AppId > 0 && context.MerchantId != uuid.Nil && tokenRoleLevel == u.Member {

		note := &Note{}
		note.AppId = context.AppId
		note.MerchantId = context.MerchantId
		note.UserId = context.UserId
		note.RoleId = context.RoleId

		switch in.Http.Method {
		case "GET":
			resp = GetNotes(context.UserId, context.AppId, context.MerchantId, context.RoleId)
			break
		case "POST":
			note.TitleName = in.TitleName
			note.DetailText = in.DetailText
			note.Category = in.Category

			resp = note.Create() //Create note
			break
		case "PUT":
			resp = SetCategoryByNotId(in.NoteId, in.Category, context.UserId, context.AppId, context.MerchantId, context.RoleId)
			break
		case "DELETE":
			if in.NoteId != uuid.Nil {
				resp = Delete(in.NoteId, context.UserId, context.AppId, context.MerchantId, context.RoleId)
			}
			break
		default:
			resp = u.Message(false, "0x11028:Invalid request")
			break
		}
	}

	// resp := make(map[string]interface{})
	// resp["notes"] = members
	// return resp

	return u.Respond(resp)
}
