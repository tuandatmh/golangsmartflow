package tsales_smart_flow

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const BaseUrl = `https://api.stage-smartflow.com/ext/v1`

type Auth struct {
	AuthToken string `json:"auth_token"`
}

type DataAuth struct {
	Message string `json:"message"`
	Data    Auth   `json:"data"`
}

func Login(params map[string]string, auths map[string]string) *Client {
	c := Client{}
	c.BaseUrl = BaseUrl
	resp := c.PostLogin("/login/en", params, auths)
	var dataresp DataAuth
	json.Unmarshal(resp, &dataresp)
	return &Client{AuthToken: dataresp.Data.AuthToken, BaseUrl: BaseUrl}
}

type DataResp struct {
	Message string `json:"message"`

	Data map[string]interface{} `json:"data"`
}

/*Department*/
const DepartmentUrl = "/en/departments/"

func (c *Client) GetDepartments() interface{} {
	resp := c.Get(DepartmentUrl, map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - name
//   - dept_code
//   - detail
//   - parent_id
func (c *Client) PostDepartment(params map[string]string) interface{} {
	resp := c.Post(DepartmentUrl, params)
	return ConvertToStruct(resp)
}

// params keys
//   - name
//   - dept_code
//   - detail
//   - parent_id
func (c *Client) PutDepartment(deartId string, params map[string]string) interface{} {
	resp := c.Put(DepartmentUrl+deartId, params)
	return ConvertToStruct(resp)
}

func (c *Client) GetDepartDetail(deartId string) interface{} {
	resp := c.Get(DepartmentUrl+deartId, map[string]string{})
	return ConvertToStruct(resp)
}

func (c *Client) DeleteDepart(deartId string) interface{} {
	resp := c.Delete(DepartmentUrl+deartId, map[string]string{})
	return ConvertToStruct(resp)
}

/*Position*/
const PositionUrl = "/en/positions/"

func (c *Client) GetPosition() interface{} {
	resp := c.Get(PositionUrl, map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - name
//   - p_code
//   - parent_id
func (c *Client) PostPosition(params map[string]string) interface{} {
	resp := c.Post(PositionUrl, params)
	return ConvertToStruct(resp)
}

func (c *Client) DeletePosition(posId string) interface{} {
	resp := c.Delete(PositionUrl+posId, map[string]string{})
	return ConvertToStruct(resp)
}

func (c *Client) GetPositionDetail(posId string) interface{} {
	resp := c.Get(PositionUrl+posId, map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - name
//   - p_code
//   - parent_id
func (c *Client) PutPosition(posId string, params map[string]string) interface{} {
	resp := c.Put(PositionUrl+posId, params)
	return ConvertToStruct(resp)
}

/*User*/
const RoleUrl = "/en/roles/"

type User struct {
	Username           string               `json:"username"`
	Email              string               `json:"email"`
	RoleID             string               `json:"role_id"`
	LastName           string               `json:"last_name"`
	FirstName          string               `json:"first_name"`
	EmployeeID         string               `json:"employee_id"`
	SlackMemberID      string               `json:"slack_member_id"`
	MUserDeptPositions []MUserDeptPositions `json:"m_user_dept_positions"`
}
type MUserDeptPositions struct {
	PosID  string `json:"pos_id"`
	DeptID string `json:"dept_id"`
}

func (c *Client) GetRoles() interface{} {
	resp := c.Get(RoleUrl, map[string]string{})
	return ConvertToStruct(resp)
}

const UserUrl = "/en/users/"

func (c *Client) GetUsers() interface{} {
	resp := c.Get(UserUrl, map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - username
//   - email
//   - role_id
//   - last_name
//   - first_name
//   - employee_id
//   - slack_member_id
func (c *Client) PostUser(params interface{}) interface{} {
	resp := c.ParamInterface("Post", UserUrl, params)
	return ConvertToStruct(resp)
}

// params keys
//   - role_id
//   - last_name
//   - first_name
//   - employee_id
//   - slack_member_id
func (c *Client) PutUser(userId string, params interface{}) interface{} {
	resp := c.ParamInterface("Put", UserUrl+userId, params)
	return ConvertToStruct(resp)
}

// params keys
//   - status
func (c *Client) PutUserStatus(userId string, params map[string]string) interface{} {
	resp := c.Put(UserUrl+userId+"/update_status", params)
	return ConvertToStruct(resp)
}

/*Application*/
const AppUrl = "/en/applications"

func (c *Client) GetApps() interface{} {
	resp := c.Get(AppUrl, map[string]string{})
	return ConvertToStruct(resp)
}

func (c *Client) CreateForm(appId string) interface{} {
	resp := c.Get("/"+appId+"/create_form/en", map[string]string{})
	return ConvertToStruct(resp)
}

func (c *Client) RequestForm(wfsId string) interface{} {
	resp := c.Get("/get_request_form/en", map[string]string{"wfs_id": wfsId})
	return ConvertToStruct(resp)
}

type TableRowData struct {
	ElementID string `json:"element_id"`
	Value     int    `json:"value"`
}
type Elements struct {
	ID             string         `json:"id"`
	ElementID      string         `json:"element_id"`
	Value          string         `json:"value"`
	IsTable        bool           `json:"is_table"`
	MElementTypeID int            `json:"m_element_type_id"`
	TableRowData   []TableRowData `json:"table_row_data"`
}
type RequestForm struct {
	Title    string     `json:"title"`
	Elements []Elements `json:"elements"`
}
type Body struct {
	Comment              string      `json:"comment"`
	SubmitStatus         string      `json:"submit_status"`
	TTodoTaskDetail      string      `json:"t_todo_task_detail"`
	WfsID                string      `json:"wfs_id"`
	SelectedDeptPosition string      `json:"selected_dept_position"`
	RequestForm          RequestForm `json:"request_form"`
}
type AppForm struct {
	Body Body `json:"body"`
}

// params keys
//   - comment
//   - submit_status
//   - t_todo_task_detail
//   - wfs_id
//   - request_form
func (c *Client) SubmitForm(params interface{}) interface{} {
	resp := c.ParamInterface("Post", "/form_submit/en", params)
	return ConvertToStruct(resp)
}

func (c *Client) DeleteApp(wfsId string) interface{} {
	resp := c.Delete("/"+wfsId+"/delete/en", map[string]string{})
	return ConvertToStruct(resp)
}

func (c *Client) FindStatus(wfsId string) interface{} {
	resp := c.Get("/find_status/en", map[string]string{"wfs_id": wfsId})
	return ConvertToStruct(resp)
}

// params keys
//   - id
//   - application_name
//   - User name
//   - assigned_user
//   - Original user
//   - Proxy user
//   - status
//   - all_tasks
//   - current_page_number
func (c *Client) GetTask(params map[string]string) interface{} {
	resp := c.RequestByForm("GET", "/en/tasks", params)
	return ConvertToStruct(resp)
}

// params keys
//   - file: string
//   - wfs_id
//   - id
//   - element_id
func (c *Client) UploadFile(params map[string]string) []byte {
	resp := c.RequestByForm("POST", "/upload_image/en", params)
	return resp
}

func (c *Client) DeleteImage(attachmentId string) interface{} {
	resp := c.Delete("/remove_image/"+attachmentId+"/en", map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - wfs_id
//   - type
func (c *Client) DownloadPdf(params map[string]string) interface{} {
	resp := c.Get("/download_pdf/en", params)
	return ConvertToStruct(resp)
}

func (c *Client) DownloadZip(wfs_id string) interface{} {
	resp := c.Get("/"+wfs_id+"/download_zip/en", map[string]string{})
	return ConvertToStruct(resp)
}

// params keys
//   - user_id
//   - type
func (c *Client) GetItemListing(params map[string]string) interface{} {
	resp := c.Get("/item_listing/en", params)
	return ConvertToStruct(resp)
}

// params keys
//   - comment
//   - wfs_id
func (c *Client) AddComment(params map[string]string) interface{} {
	resp := c.Post("/add_comment/en", params)
	return ConvertToStruct(resp)
}

type ErrorDetail struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

func ConvertToStruct(resp []byte) interface{} {
	var dataresp DataResp
	json.Unmarshal(resp, &dataresp)

	if dataresp.Data == nil && dataresp.Message != "Successful" {
		var errorresp ErrorDetail
		json.Unmarshal(resp, &errorresp)
		return errorresp
	}
	return dataresp
}

// Application for dayoff
type DataReturn struct {
	WfsID    int64      `json:"wfs_id"`
	Sections []Sections `json:"sections"`
}
type Sections struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	MFormDetailId int64     `json:"m_form_detail_id:"`
	IsTable       bool      `json:"is_table"`
	HelpText      string    `json:"help_text"`
	Element       []Element `json:"element"`
}

type Element struct {
	Id             int    `json:"id"`
	ElementID      string `json:"element_id"`
	Value          string `json:"value"`
	IsTable        bool   `json:"is_table"`
	MElementTypeID int    `json:"m_element_type_id"`
	Label          string `json:"label"`
}

func (c *Client) DayOff(wfsid string, params map[string]string) interface{} {

	// CreateForm and Get wfs_id
	createForm := c.CreateForm(wfsid)

	convertCreateForm := ConvertToByte(createForm)

	var dataCreateForm DataResp

	var err error

	err = json.Unmarshal(convertCreateForm, &dataCreateForm)
	if err != nil {
		fmt.Println("error:", err)
	}
	wfsId := dataCreateForm.Data["wfs_id"]

	// Get data return form RequesForm
	getRequestId := ConvertToByte(wfsId)

	getRequestForm := c.RequestForm(string(getRequestId))

	convertReForm := ConvertToByte(getRequestForm)

	var dataRequestForm DataResp

	err = json.Unmarshal(convertReForm, &dataRequestForm)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Get Param from element
	section := dataRequestForm.Data

	getElement := ConvertToByte(section)

	var element DataReturn

	err = json.Unmarshal(getElement, &element)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Post Formsubmit
	dataForm := AppForm{
		Body: Body{
			Comment:              params["Comment"],
			SubmitStatus:         params["SubmitStatus"],
			TTodoTaskDetail:      params["TTodoTaskDetail"],
			WfsID:                string(getRequestId),
			SelectedDeptPosition: params["SelectedDeptPosition"],
			RequestForm: RequestForm{
				Title: params["Title"],
				Elements: []Elements{
					{
						ID:             strconv.Itoa(element.Sections[0].Element[0].Id),
						ElementID:      element.Sections[0].Element[0].ElementID,
						Value:          params["Date"],
						IsTable:        element.Sections[0].Element[0].IsTable,
						MElementTypeID: element.Sections[0].Element[0].MElementTypeID,
					},
					{
						ID:             strconv.Itoa(element.Sections[0].Element[1].Id),
						ElementID:      element.Sections[0].Element[1].ElementID,
						Value:          params["Reason"],
						IsTable:        element.Sections[0].Element[1].IsTable,
						MElementTypeID: element.Sections[0].Element[1].MElementTypeID,
					},
				},
			},
		},
	}
	return c.SubmitForm(dataForm)
}

func ConvertToByte(resp interface{}) []byte {
	out, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return out
}
