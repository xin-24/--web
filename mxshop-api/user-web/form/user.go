package form

type PassWordLoginForm struct{
	Mobile string `form:"mobile"json:"mobile"binding:"required,mobile"`//手机号码有规范可寻，自定义validator 
	PassWord string `form:"password" json:"password"binding:"required,min=3,max=10"`//约束条件之间不能加空格
}