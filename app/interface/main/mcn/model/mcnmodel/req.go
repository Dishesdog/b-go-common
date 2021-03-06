package mcnmodel

import (
	mcnadminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/tool/validate"
	"go-common/library/time"
)

//CookieMidInterface cookie set interface, set mid from cookie to arg
type CookieMidInterface interface {
	SetMid(midFromCookie int64)
}

//McnCommonReq common mcn
type McnCommonReq struct {
	McnCheatReq
	SignID int64 `form:"sign_id"`
	McnMid int64
}

//CheatInterface cheat interface
type CheatInterface interface {
	// Cheat return true if cheated, false if not cheated
	Cheat() bool
}

//McnCheatReq cheat
type McnCheatReq struct {
	TMcnMid int64 `form:"t_mcn_mid"`
}

//Cheat .
func (m *McnCommonReq) Cheat() bool {
	if m.TMcnMid == 0 {
		return false
	}
	m.SetMid(m.TMcnMid)
	return true
}

//SetMid set mid
func (m *McnCommonReq) SetMid(midFromCookie int64) {
	m.McnMid = midFromCookie
}

//UpCommonReq common up
type UpCommonReq struct {
	UpMid int64
}

//SetMid set mid
func (m *UpCommonReq) SetMid(midFromCookie int64) {
	m.UpMid = midFromCookie
}

//GetStateReq get state
type GetStateReq struct {
	McnCommonReq
}

//McnApplyReq apply req
type McnApplyReq struct {
	McnCommonReq
	CompanyName        string `form:"company_name"`
	CompanyLicenseID   string `form:"company_license_id"`
	ContactName        string `form:"contact_name"`
	ContactTitle       string `form:"contact_title"`
	ContactIdcard      string `form:"contact_idcard" validate:"idcheck"`
	ContactPhone       string `form:"contact_phone" validate:"phonecheck"`
	CompanyLicenseLink string `form:"company_license_link" validate:"httpcheck"`
	ContractLink       string `form:"contract_link" validate:"httpcheck"`
}

//CopyTo .
func (m *McnApplyReq) CopyTo(v *McnSign) {
	if v == nil {
		return
	}
	v.McnMid = m.McnMid
	v.CompanyName = m.CompanyName
	v.CompanyLicenseID = m.CompanyLicenseID
	v.ContactName = m.ContactName
	v.ContactTitle = m.ContactTitle
	v.ContactIdcard = m.ContactIdcard
	v.ContactPhone = m.ContactPhone
	v.CompanyLicenseLink = m.CompanyLicenseLink
	v.ContractLink = m.ContractLink
}

//McnBindUpApplyReq .
type McnBindUpApplyReq struct {
	McnCommonReq
	UpMid        int64     `form:"up_mid"`
	BeginDate    time.Time `form:"begin_date"`
	EndDate      time.Time `form:"end_date"`
	ContractLink string    `form:"contract_link"` // ????????????http??????
	UpAuthLink   string    `form:"up_auth_link"`  // ????????????http??????
	UpType       int8      `form:"up_type"`       // ???????????????0????????????1?????????
	SiteLink     string    `form:"site_link"`     //up?????????????????????, ??????up type???1???????????????
	mcnadminmodel.Permits
	PublicationPrice int64 `form:"publication_price"` // ?????????1/1000 ???
}

//IsSiteInfoOk ????????????up???????????????OK?????????????????????Up??????????????????ok
func (m *McnBindUpApplyReq) IsSiteInfoOk() bool {
	if m.UpType == 0 {
		return true
	}
	return validate.RegHTTPCheck.MatchString(m.SiteLink)
}

//CopyTo .
func (m *McnBindUpApplyReq) CopyTo(v *McnUp) {
	v.UpMid = m.UpMid
	v.McnMid = m.McnMid
	v.BeginDate = m.BeginDate
	v.EndDate = m.EndDate
	v.ContractLink = m.ContractLink
	v.UpAuthLink = m.UpAuthLink
	v.UpType = m.UpType
	v.SiteLink = m.SiteLink
	v.Permission = uint32(m.GetAttrPermitVal())
	v.PublicationPrice = m.PublicationPrice
}

//McnUpConfirmReq .
type McnUpConfirmReq struct {
	UpCommonReq
	BindID int64 `form:"bind_id"`
	Choice bool  `form:"choice"`
}

//McnUpGetBindReq .
type McnUpGetBindReq struct {
	UpCommonReq
	BindID int64 `form:"bind_id"`
}

// McnGetDataSummaryReq req
type McnGetDataSummaryReq = McnCommonReq

//McnGetUpListReq req
type McnGetUpListReq struct {
	McnCommonReq
	UpMid int64 `form:"up_mid"`
	model.PageArg
}

//McnGetAccountReq req
type McnGetAccountReq struct {
	Mid int64 `form:"mid"`
}

// McnGetMcnOldInfoReq req
type McnGetMcnOldInfoReq struct {
	McnCommonReq
}

// McnGetRankReq req to ????????????
type McnGetRankReq struct {
	McnCommonReq
	Tid      int16    `form:"tid"` // ??????
	DataType DataType `form:"data_type"`
}

// McnGetRecommendPoolReq get recommend pool
type McnGetRecommendPoolReq struct {
	McnCommonReq
	model.PageArg
	Tid        int16  `form:"tid"`
	OrderField string `form:"order_field"`
	Sort       string `form:"sort"`
}

// McnGetRecommendPoolTidListReq common req
type McnGetRecommendPoolTidListReq = McnCommonReq

// ------inner request

// McnGetRankAPIReq req to ????????????
type McnGetRankAPIReq struct {
	SignID   int64    `form:"sign_id"`
	Tid      int16    `form:"tid"` // ??????
	DataType DataType `form:"data_type"`
}

// ??????/??????/??????/??????/??????/??????/?????????
const (
	ActionTypePlay  = "play"  //??????
	ActionTypeDanmu = "danmu" //??????
	ActionTypeReply = "reply" //??????
	ActionTypeShare = "share" //??????
	ActionTypeCoin  = "coin"  //??????
	ActionTypeFav   = "fav"   //??????
	ActionTypeLike  = "like"  //?????????
)

const (
	// UserTypeGuest .
	UserTypeGuest = "guest" // ??????
	// UserTypeFans .
	UserTypeFans = "fans" // ??????
)

//McnGetIndexIncReq ????????????
type McnGetIndexIncReq struct {
	McnCommonReq
	Type string `form:"type"`
}

//McnGetIndexSourceReq ????????????
type McnGetIndexSourceReq = McnGetIndexIncReq

//McnGetPlaySourceReq ????????????????????????
type McnGetPlaySourceReq struct {
	McnCommonReq
}

//McnGetMcnFansReq mcn
type McnGetMcnFansReq = McnCommonReq

//McnGetMcnFansIncReq mcn?????????????????????
type McnGetMcnFansIncReq = McnCommonReq

//McnGetMcnFansDecReq mcn?????????????????????
type McnGetMcnFansDecReq = McnCommonReq

//McnGetMcnFansAttentionWayReq mcn??????????????????
type McnGetMcnFansAttentionWayReq = McnCommonReq

// McnGetBaseFansAttrReq  mcn ?????????????????????????????????
type McnGetBaseFansAttrReq struct {
	McnCommonReq
	UserType string `form:"user_type"`
}

// McnGetFansAreaReq mcn ??????????????????
type McnGetFansAreaReq = McnGetBaseFansAttrReq

// McnGetFansTypeReq  mcn  ??????/????????????????????????
type McnGetFansTypeReq = McnGetBaseFansAttrReq

// McnGetFansTagReq  mcn  ??????/??????????????????????????????
type McnGetFansTagReq = McnGetBaseFansAttrReq

//McnChangePermitReq change permit
type McnChangePermitReq struct {
	McnCommonReq
	UpMid int64 `form:"up_mid"`
	mcnadminmodel.Permits
	UpAuthLink string `form:"up_auth_link" validate:"httpcheck"`
}

//McnPublicationPriceChangeReq change publication price
type McnPublicationPriceChangeReq struct {
	McnCommonReq
	Price int64 `form:"price"`
	UpMid int64 `form:"up_mid"`
}
