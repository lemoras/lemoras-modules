package drive

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
	BucketId     uuid.UUID `json:"bucketId"`
	BucketItemId uuid.UUID `json:"bucketItemId"`
	ItemType     int       `json:"itemId"`
	ItemName     string    `json:"itemName"`
	ItemSize     float64   `json:"itemSize"`
	ItemURL      string    `json:"itemURL"`

	Http CustomHttp `json:"http"`
}

type RequestDto struct {
	BucketId uuid.UUID `json:"bucketId"`

	BucketItemId uuid.UUID `json:"bucketItemId"`
	ItemType     int       `json:"itemId"`
	ItemName     string    `json:"itemName"`
	ItemSize     float64   `json:"itemSize"`
	ItemURL      string    `json:"itemURL"`

	ItemCount int     `json:"itemCount"`
	TotalSize float64 `json:"totalSize"`

	UserId     uuid.UUID `json:"userId"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
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

		req := &RequestDto{}
		req.AppId = context.AppId
		req.MerchantId = context.MerchantId
		req.UserId = context.UserId
		req.RoleId = context.RoleId

		switch in.Http.Method {
		case "GET":
			resp = GetBucketsWithItems(context.UserId, context.MerchantId, context.AppId, context.RoleId, in.BucketId, in.BucketItemId, 3, initCompleted)
			break
		case "POST":
			req.BucketItemId = in.BucketItemId
			req.ItemType = in.ItemType
			req.ItemName = in.ItemName
			req.ItemSize = in.ItemSize
			req.ItemURL = in.ItemURL

			resp = req.Create() //Create
			break
		case "PUT":
			// resp = SetCategoryByNotId(in.BucketItemId, in.Category, context.UserId, context.AppId, context.MerchantId, context.RoleId)
			break
		case "DELETE":
			if in.BucketItemId != uuid.Nil {
				resp = SoftDeleteBucketItemRecursive(in.BucketItemId, context.UserId, context.MerchantId, context.AppId, context.RoleId)
			}
			break
		default:
			resp = u.Message(false, "0x11028:Invalid request")
			break
		}
	}

	return u.Respond(resp)
}
