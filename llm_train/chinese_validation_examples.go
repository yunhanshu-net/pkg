package llm_train

// ChineseValidationExamples 中国本土化验证规则示例
// 针对中国用户常用的验证场景，提供标准的验证规则
type ChineseValidationExamples struct {
	// ===== 个人信息验证 =====

	// 中国手机号验证（11位，1开头，第二位3-9）
	MobilePhone string `validate:"required,regexp=^1[3-9]\\d{9}$" runner:"code:mobile_phone;name:手机号;type:string;placeholder:请输入11位手机号" json:"mobile_phone"`

	// 中国身份证号验证（18位，包含校验位）
	IdCard string `validate:"required,regexp=^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$" runner:"code:id_card;name:身份证号;type:string;placeholder:请输入18位身份证号" json:"id_card"`

	// 中文姓名验证（2-20个中文字符）
	ChineseName string `validate:"required,regexp=^[\\p{Han}]{2,20}$" runner:"code:chinese_name;name:中文姓名;type:string;placeholder:请输入中文姓名" json:"chinese_name"`

	// 邮政编码验证（6位数字）
	PostalCode string `validate:"required,regexp=^\\d{6}$" runner:"code:postal_code;name:邮政编码;type:string;placeholder:请输入6位邮政编码" json:"postal_code"`

	// ===== 银行卡相关验证 =====

	// 银行卡号验证（16-19位数字）
	BankCard string `validate:"required,regexp=^\\d{16,19}$" runner:"code:bank_card;name:银行卡号;type:string;placeholder:请输入银行卡号" json:"bank_card"`

	// 银行预留手机号（同手机号格式）
	BankPhone string `validate:"required,regexp=^1[3-9]\\d{9}$" runner:"code:bank_phone;name:银行预留手机号;type:string;placeholder:请输入银行预留手机号" json:"bank_phone"`

	// ===== 企业信息验证 =====

	// 统一社会信用代码（18位，字母数字组合）
	UnifiedSocialCreditCode string `validate:"required,regexp=^[0-9A-HJ-NPQRTUWXY]{2}\\d{6}[0-9A-HJ-NPQRTUWXY]{10}$" runner:"code:unified_social_credit_code;name:统一社会信用代码;type:string;placeholder:请输入18位统一社会信用代码" json:"unified_social_credit_code"`

	// 营业执照注册号（15位数字）
	BusinessLicenseNumber string `validate:"required,regexp=^\\d{15}$" runner:"code:business_license_number;name:营业执照注册号;type:string;placeholder:请输入15位营业执照注册号" json:"business_license_number"`

	// 组织机构代码（9位，8位数字+1位校验码）
	OrganizationCode string `validate:"required,regexp=^[A-Z0-9]{8}-[A-Z0-9]$" runner:"code:organization_code;name:组织机构代码;type:string;placeholder:请输入组织机构代码" json:"organization_code"`

	// 税务登记号（15位或20位）
	TaxRegistrationNumber string `validate:"required,regexp=^\\d{15}(\\d{5})?$" runner:"code:tax_registration_number;name:税务登记号;type:string;placeholder:请输入税务登记号" json:"tax_registration_number"`

	// ===== 车辆相关验证 =====

	// 车牌号验证（新能源和传统车牌）
	LicensePlate string `validate:"required,regexp=^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{4,5}[A-HJ-NP-Z0-9挂学警港澳]$" runner:"code:license_plate;name:车牌号;type:string;placeholder:请输入车牌号" json:"license_plate"`

	// 车架号/VIN码验证（17位字母数字，不含I、O、Q）
	VinCode string `validate:"required,regexp=^[A-HJ-NPR-Z0-9]{17}$" runner:"code:vin_code;name:车架号;type:string;placeholder:请输入17位车架号" json:"vin_code"`

	// ===== 学历相关验证 =====

	// 学历证书编号验证（多种格式）
	EducationCertNumber string `validate:"required,regexp=^\\d{17,18}$" runner:"code:education_cert_number;name:学历证书编号;type:string;placeholder:请输入学历证书编号" json:"education_cert_number"`

	// 学位证书编号验证
	DegreeCertNumber string `validate:"required,regexp=^\\d{16,18}$" runner:"code:degree_cert_number;name:学位证书编号;type:string;placeholder:请输入学位证书编号" json:"degree_cert_number"`

	// ===== 医疗相关验证 =====

	// 医保卡号验证（各地格式不同，这里用通用格式）
	MedicalInsuranceNumber string `validate:"required,regexp=^\\d{9,18}$" runner:"code:medical_insurance_number;name:医保卡号;type:string;placeholder:请输入医保卡号" json:"medical_insurance_number"`

	// 社保卡号验证
	SocialSecurityNumber string `validate:"required,regexp=^\\d{9,18}$" runner:"code:social_security_number;name:社保卡号;type:string;placeholder:请输入社保卡号" json:"social_security_number"`

	// ===== 房产相关验证 =====

	// 房产证号验证（各地格式不同）
	PropertyCertNumber string `validate:"required,regexp=^[\\w\\d\\-]{10,30}$" runner:"code:property_cert_number;name:房产证号;type:string;placeholder:请输入房产证号" json:"property_cert_number"`

	// 土地证号验证
	LandCertNumber string `validate:"required,regexp=^[\\w\\d\\-]{10,30}$" runner:"code:land_cert_number;name:土地证号;type:string;placeholder:请输入土地证号" json:"land_cert_number"`

	// ===== 网络相关验证 =====

	// QQ号验证（5-11位数字，不能以0开头）
	QQNumber string `validate:"required,regexp=^[1-9]\\d{4,10}$" runner:"code:qq_number;name:QQ号;type:string;placeholder:请输入QQ号" json:"qq_number"`

	// 微信号验证（6-20位，字母数字下划线，字母开头）
	WeChatId string `validate:"required,regexp=^[a-zA-Z][a-zA-Z0-9_]{5,19}$" runner:"code:wechat_id;name:微信号;type:string;placeholder:请输入微信号" json:"wechat_id"`

	// ===== 地址相关验证 =====

	// 详细地址验证（包含中文字符、数字、常用符号）
	DetailedAddress string `validate:"required,regexp=^[\\p{Han}\\w\\d\\s\\-#()（）号栋单元室]{10,200}$" runner:"code:detailed_address;name:详细地址;type:string;placeholder:请输入详细地址" json:"detailed_address"`
}

// ChineseValidationRules 中国本土化常用验证规则映射
var ChineseValidationRules = map[string]string{
	// 个人信息
	"mobile_phone": "required,regexp=^1[3-9]\\d{9}$",
	"id_card":      "required,regexp=^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$",
	"chinese_name": "required,regexp=^[\\p{Han}]{2,20}$",
	"postal_code":  "required,regexp=^\\d{6}$",

	// 银行卡相关
	"bank_card":  "required,regexp=^\\d{16,19}$",
	"bank_phone": "required,regexp=^1[3-9]\\d{9}$",

	// 企业信息
	"unified_social_credit_code": "required,regexp=^[0-9A-HJ-NPQRTUWXY]{2}\\d{6}[0-9A-HJ-NPQRTUWXY]{10}$",
	"business_license":           "required,regexp=^\\d{15}$",
	"organization_code":          "required,regexp=^[A-Z0-9]{8}-[A-Z0-9]$",
	"tax_registration":           "required,regexp=^\\d{15}(\\d{5})?$",

	// 车辆相关
	"license_plate": "required,regexp=^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{4,5}[A-HJ-NP-Z0-9挂学警港澳]$",
	"vin_code":      "required,regexp=^[A-HJ-NPR-Z0-9]{17}$",

	// 学历相关
	"education_cert": "required,regexp=^\\d{17,18}$",
	"degree_cert":    "required,regexp=^\\d{16,18}$",

	// 医疗相关
	"medical_insurance": "required,regexp=^\\d{9,18}$",
	"social_security":   "required,regexp=^\\d{9,18}$",

	// 房产相关
	"property_cert": "required,regexp=^[\\w\\d\\-]{10,30}$",
	"land_cert":     "required,regexp=^[\\w\\d\\-]{10,30}$",

	// 网络相关
	"qq_number": "required,regexp=^[1-9]\\d{4,10}$",
	"wechat_id": "required,regexp=^[a-zA-Z][a-zA-Z0-9_]{5,19}$",

	// 地址相关
	"detailed_address": "required,regexp=^[\\p{Han}\\w\\d\\s\\-#()（）号栋单元室]{10,200}$",
}

// ChineseValidationExplanations 中国本土化验证规则说明
var ChineseValidationExplanations = map[string]string{
	"mobile_phone": "中国手机号：11位数字，1开头，第二位3-9",
	"id_card":      "中国身份证号：18位，包含地区码、出生日期、顺序码和校验位",
	"chinese_name": "中文姓名：2-20个中文字符",
	"postal_code":  "中国邮政编码：6位数字",

	"bank_card":  "银行卡号：16-19位数字",
	"bank_phone": "银行预留手机号：同手机号格式",

	"unified_social_credit_code": "统一社会信用代码：18位字母数字组合",
	"business_license":           "营业执照注册号：15位数字",
	"organization_code":          "组织机构代码：8位字母数字+1位校验码",
	"tax_registration":           "税务登记号：15位或20位数字",

	"license_plate": "车牌号：省份简称+字母+4-5位字母数字",
	"vin_code":      "车架号：17位字母数字，不含I、O、Q",

	"education_cert": "学历证书编号：17-18位数字",
	"degree_cert":    "学位证书编号：16-18位数字",

	"medical_insurance": "医保卡号：9-18位数字",
	"social_security":   "社保卡号：9-18位数字",

	"property_cert": "房产证号：10-30位字母数字组合",
	"land_cert":     "土地证号：10-30位字母数字组合",

	"qq_number": "QQ号：5-11位数字，不能以0开头",
	"wechat_id": "微信号：6-20位，字母开头，可含数字下划线",

	"detailed_address": "详细地址：10-200字符，包含中文、数字、常用符号",
}

// 使用示例说明
/*
中国本土化验证规则使用说明：

1. 手机号验证：
   validate:"required,regexp=^1[3-9]\\d{9}$"
   - 必须11位数字
   - 第一位必须是1
   - 第二位必须是3-9
   - 后面9位是任意数字

2. 身份证号验证：
   validate:"required,regexp=^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"
   - 18位字符
   - 前6位：地区码（不能以0开头）
   - 7-14位：出生日期（年月日）
   - 15-17位：顺序码
   - 18位：校验码（数字或X）

3. 中文姓名验证：
   validate:"required,regexp=^[\\p{Han}]{2,20}$"
   - 只能包含中文字符
   - 长度2-20个字符

4. 银行卡号验证：
   validate:"required,regexp=^\\d{16,19}$"
   - 16-19位纯数字
   - 覆盖主流银行卡格式

5. 车牌号验证：
   validate:"required,regexp=^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{4,5}[A-HJ-NP-Z0-9挂学警港澳]$"
   - 省份简称+字母+4-5位字母数字+特殊标识
   - 支持新能源车牌和传统车牌

注意事项：
- 正则表达式中的反斜杠需要双重转义
- 中文字符使用\\p{Han}匹配
- 某些证件号码格式可能因地区而异，需要根据实际情况调整
- 建议结合前端验证和后端验证，提供更好的用户体验
- 对于敏感信息，还需要考虑数据加密和隐私保护
*/
