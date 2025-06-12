package dev

import "time"

// 这个文件包含了所有组件的完整使用示例
// 展示了如何在实际业务场景中使用各种组件和标签

// ===== 用户管理示例 =====

// UserExample 用户管理完整示例
type UserExample struct {
	// 基础信息字段
	ID       uint   `runner:"code:id;name:ID;show:list,detail" json:"id"`
	Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名" validate:"required,min=3,max=20" json:"username"`
	Password string `runner:"code:password;name:密码;widget:input;mode:password;hidden:list,detail" validate:"required,min=6" json:"password"`
	Email    string `runner:"code:email;name:邮箱;widget:input;mode:line_text;placeholder:请输入邮箱地址" validate:"required,email" json:"email"`

	// 个人信息
	RealName string `runner:"code:real_name;name:真实姓名;widget:input;mode:line_text;placeholder:请输入真实姓名" validate:"required,max=50" json:"real_name"`
	Age      int    `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;unit:岁;placeholder:请输入年龄" validate:"min=1,max=150" json:"age"`
	Gender   string `runner:"code:gender;name:性别;widget:radio;options:male(男),female(女),other(其他);direction:horizontal" validate:"required,oneof=male female other" json:"gender"`
	Birthday string `runner:"code:birthday;name:生日;widget:date;format:YYYY-MM-DD;placeholder:请选择生日" validate:"required" json:"birthday"`

	// 状态控制
	Status    string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用),banned(封禁);default_value:active" validate:"required,oneof=active inactive banned" json:"status"`
	IsVIP     bool   `runner:"code:is_vip;name:VIP用户;widget:switch;on_text:是;off_text:否" json:"is_vip"`
	IsEnabled bool   `runner:"code:is_enabled;name:是否启用;widget:switch;on_text:启用;off_text:禁用" json:"is_enabled"`

	// 权限管理
	Roles       []string `runner:"code:roles;name:角色;widget:checkbox;options:admin(管理员),user(普通用户),guest(访客);min_select:1" validate:"required,min=1" json:"roles"`
	Permissions string   `runner:"code:permissions;name:权限;widget:checkbox;options:read(读取),write(写入),delete(删除),admin(管理)" json:"permissions"`

	// 联系方式
	Phone   string `runner:"code:phone;name:手机号;widget:input;mode:line_text;placeholder:请输入手机号" validate:"required,len=11,numeric" json:"phone"`
	Address string `runner:"code:address;name:地址;widget:input;mode:text_area;placeholder:请输入详细地址" validate:"max=200" json:"address"`

	// 头像上传
	Avatar string `runner:"code:avatar;name:头像;widget:file;accept:.jpg,.png,.gif;max_size:2097152;upload_text:选择头像" json:"avatar"`

	// 系统字段
	CreatedAt time.Time `runner:"code:created_at;name:创建时间;show:list,detail" json:"created_at"`
	UpdatedAt time.Time `runner:"code:updated_at;name:更新时间;show:detail" json:"updated_at"`
}

// ===== 商品管理示例 =====

// ProductExample 商品管理完整示例
type ProductExample struct {
	// 基础信息
	ID          uint   `runner:"code:id;name:ID;show:list,detail" json:"id"`
	Name        string `runner:"code:name;name:商品名称;widget:input;mode:line_text;placeholder:请输入商品名称" validate:"required,max=100" json:"name"`
	Description string `runner:"code:description;name:商品描述;widget:input;mode:text_area;placeholder:请输入商品描述" validate:"max=1000" json:"description"`
	SKU         string `runner:"code:sku;name:商品编码;widget:input;mode:line_text;placeholder:请输入SKU编码" validate:"required,max=50" json:"sku"`

	// 价格信息
	Price       float64 `runner:"code:price;name:售价;widget:number;min:0;step:0.01;precision:2;unit:元;placeholder:请输入售价" validate:"required,min=0" json:"price"`
	CostPrice   float64 `runner:"code:cost_price;name:成本价;widget:number;min:0;step:0.01;precision:2;unit:元;placeholder:请输入成本价" validate:"min=0" json:"cost_price"`
	MarketPrice float64 `runner:"code:market_price;name:市场价;widget:number;min:0;step:0.01;precision:2;unit:元;placeholder:请输入市场价" validate:"min=0" json:"market_price"`

	// 库存管理
	Stock    int `runner:"code:stock;name:库存数量;widget:number;min:0;step:1;unit:件;placeholder:请输入库存数量" validate:"min=0" json:"stock"`
	MinStock int `runner:"code:min_stock;name:最低库存;widget:number;min:0;step:1;unit:件;placeholder:请输入最低库存" validate:"min=0" json:"min_stock"`
	MaxStock int `runner:"code:max_stock;name:最高库存;widget:number;min:0;step:1;unit:件;placeholder:请输入最高库存" validate:"min=0" json:"max_stock"`

	// 分类和标签
	Category string   `runner:"code:category;name:商品分类;widget:select;options:electronics(电子产品),clothing(服装),books(图书),food(食品)" validate:"required" json:"category"`
	Tags     []string `runner:"code:tags;name:商品标签;widget:checkbox;options:hot(热销),new(新品),discount(折扣),limited(限量)" json:"tags"`
	Brand    string   `runner:"code:brand;name:品牌;widget:select;options:apple(苹果),samsung(三星),huawei(华为),xiaomi(小米)" json:"brand"`

	// 状态控制
	Status string `runner:"code:status;name:状态;widget:select;options:draft(草稿),published(已发布),offline(下架);default_value:draft" validate:"required,oneof=draft published offline" json:"status"`
	IsHot  bool   `runner:"code:is_hot;name:热销商品;widget:switch;on_text:是;off_text:否" json:"is_hot"`
	IsNew  bool   `runner:"code:is_new;name:新品;widget:switch;on_text:是;off_text:否" json:"is_new"`

	// 规格参数
	Weight float64 `runner:"code:weight;name:重量;widget:number;min:0;step:0.01;precision:2;unit:kg;placeholder:请输入重量" validate:"min=0" json:"weight"`
	Size   string  `runner:"code:size;name:尺寸;widget:input;mode:line_text;placeholder:长x宽x高(cm)" json:"size"`
	Color  string  `runner:"code:color;name:颜色;widget:radio;options:red(红色),blue(蓝色),green(绿色),black(黑色),white(白色);direction:horizontal" json:"color"`

	// 图片上传
	MainImage string   `runner:"code:main_image;name:主图;widget:file;accept:.jpg,.png,.gif;max_size:5242880;upload_text:选择主图" validate:"required" json:"main_image"`
	Images    []string `runner:"code:images;name:商品图片;widget:file;accept:.jpg,.png,.gif;max_size:5242880;multiple:true;upload_text:选择图片" json:"images"`

	// 销售信息
	SalesCount int     `runner:"code:sales_count;name:销量;show:list,detail" json:"sales_count"`
	Rating     float64 `runner:"code:rating;name:评分;show:list,detail" json:"rating"`

	// 时间字段
	PublishTime time.Time `runner:"code:publish_time;name:发布时间;widget:date;format:YYYY-MM-DD HH:mm;show_time:true" json:"publish_time"`
	CreatedAt   time.Time `runner:"code:created_at;name:创建时间;show:list,detail" json:"created_at"`
	UpdatedAt   time.Time `runner:"code:updated_at;name:更新时间;show:detail" json:"updated_at"`
}

// ===== 订单管理示例 =====

// OrderExample 订单管理完整示例
type OrderExample struct {
	// 基础信息
	ID       uint   `runner:"code:id;name:订单ID;show:list,detail" json:"id"`
	OrderNo  string `runner:"code:order_no;name:订单号;show:list,detail" json:"order_no"`
	UserID   uint   `runner:"code:user_id;name:用户ID;widget:number;min:1;step:1" validate:"required,min=1" json:"user_id"`
	Username string `runner:"code:username;name:用户名;show:list,detail" json:"username"`

	// 订单状态
	Status    string `runner:"code:status;name:订单状态;widget:select;options:pending(待付款),paid(已付款),shipped(已发货),delivered(已送达),cancelled(已取消);default_value:pending" validate:"required,oneof=pending paid shipped delivered cancelled" json:"status"`
	PayStatus string `runner:"code:pay_status;name:支付状态;widget:select;options:unpaid(未支付),paid(已支付),refunded(已退款)" validate:"required,oneof=unpaid paid refunded" json:"pay_status"`

	// 金额信息
	TotalAmount    float64 `runner:"code:total_amount;name:订单总额;widget:number;min:0;step:0.01;precision:2;unit:元" validate:"required,min=0" json:"total_amount"`
	PayAmount      float64 `runner:"code:pay_amount;name:实付金额;widget:number;min:0;step:0.01;precision:2;unit:元" validate:"required,min=0" json:"pay_amount"`
	DiscountAmount float64 `runner:"code:discount_amount;name:优惠金额;widget:number;min:0;step:0.01;precision:2;unit:元" validate:"min=0" json:"discount_amount"`
	ShippingFee    float64 `runner:"code:shipping_fee;name:运费;widget:number;min:0;step:0.01;precision:2;unit:元" validate:"min=0" json:"shipping_fee"`

	// 支付信息
	PayMethod string    `runner:"code:pay_method;name:支付方式;widget:radio;options:alipay(支付宝),wechat(微信支付),bank(银行卡);direction:horizontal" json:"pay_method"`
	PayTime   time.Time `runner:"code:pay_time;name:支付时间;widget:date;format:YYYY-MM-DD HH:mm;show_time:true" json:"pay_time"`

	// 收货信息
	ReceiverName    string `runner:"code:receiver_name;name:收货人;widget:input;mode:line_text;placeholder:请输入收货人姓名" validate:"required,max=50" json:"receiver_name"`
	ReceiverPhone   string `runner:"code:receiver_phone;name:收货电话;widget:input;mode:line_text;placeholder:请输入收货电话" validate:"required,len=11,numeric" json:"receiver_phone"`
	ReceiverAddress string `runner:"code:receiver_address;name:收货地址;widget:input;mode:text_area;placeholder:请输入收货地址" validate:"required,max=200" json:"receiver_address"`

	// 物流信息
	ShippingCompany string    `runner:"code:shipping_company;name:物流公司;widget:select;options:sf(顺丰),ems(邮政EMS),yt(圆通),sto(申通),zto(中通)" json:"shipping_company"`
	TrackingNumber  string    `runner:"code:tracking_number;name:快递单号;widget:input;mode:line_text;placeholder:请输入快递单号" json:"tracking_number"`
	ShippedTime     time.Time `runner:"code:shipped_time;name:发货时间;widget:date;format:YYYY-MM-DD HH:mm;show_time:true" json:"shipped_time"`
	DeliveredTime   time.Time `runner:"code:delivered_time;name:送达时间;widget:date;format:YYYY-MM-DD HH:mm;show_time:true" json:"delivered_time"`

	// 备注信息
	Remark    string `runner:"code:remark;name:订单备注;widget:input;mode:text_area;placeholder:请输入订单备注" validate:"max=500" json:"remark"`
	AdminNote string `runner:"code:admin_note;name:管理员备注;widget:input;mode:text_area;placeholder:管理员备注" validate:"max=500" json:"admin_note"`

	// 标记字段
	IsUrgent bool `runner:"code:is_urgent;name:紧急订单;widget:switch;on_text:是;off_text:否" json:"is_urgent"`
	IsGift   bool `runner:"code:is_gift;name:礼品订单;widget:switch;on_text:是;off_text:否" json:"is_gift"`

	// 时间字段
	CreatedAt time.Time `runner:"code:created_at;name:下单时间;show:list,detail" json:"created_at"`
	UpdatedAt time.Time `runner:"code:updated_at;name:更新时间;show:detail" json:"updated_at"`
}

// ===== 系统配置示例 =====

// SystemConfigExample 系统配置完整示例
type SystemConfigExample struct {
	// 基础配置
	ID          uint   `runner:"code:id;name:配置ID;show:list,detail" json:"id"`
	ConfigKey   string `runner:"code:config_key;name:配置键;widget:input;mode:line_text;placeholder:请输入配置键" validate:"required,max=100" json:"config_key"`
	ConfigValue string `runner:"code:config_value;name:配置值;widget:input;mode:text_area;placeholder:请输入配置值" validate:"required" json:"config_value"`
	Description string `runner:"code:description;name:配置描述;widget:input;mode:text_area;placeholder:请输入配置描述" validate:"max=500" json:"description"`

	// 配置分组
	Group    string `runner:"code:group;name:配置分组;widget:select;options:system(系统配置),email(邮件配置),sms(短信配置),payment(支付配置),storage(存储配置)" validate:"required" json:"group"`
	Category string `runner:"code:category;name:配置分类;widget:radio;options:basic(基础),advanced(高级),security(安全);direction:horizontal" validate:"required" json:"category"`

	// 数据类型
	DataType string `runner:"code:data_type;name:数据类型;widget:select;options:string(字符串),number(数字),boolean(布尔值),json(JSON对象)" validate:"required,oneof=string number boolean json" json:"data_type"`

	// 状态控制
	IsEnabled   bool `runner:"code:is_enabled;name:是否启用;widget:switch;on_text:启用;off_text:禁用" json:"is_enabled"`
	IsReadonly  bool `runner:"code:is_readonly;name:只读配置;widget:switch;on_text:是;off_text:否" json:"is_readonly"`
	IsEncrypted bool `runner:"code:is_encrypted;name:是否加密;widget:switch;on_text:是;off_text:否" json:"is_encrypted"`

	// 验证规则
	ValidationRule string `runner:"code:validation_rule;name:验证规则;widget:input;mode:text_area;placeholder:请输入验证规则(正则表达式)" json:"validation_rule"`
	DefaultValue   string `runner:"code:default_value;name:默认值;widget:input;mode:line_text;placeholder:请输入默认值" json:"default_value"`

	// 权限控制
	RequiredRoles []string `runner:"code:required_roles;name:所需角色;widget:checkbox;options:admin(管理员),operator(操作员),viewer(查看者)" json:"required_roles"`

	// 排序和显示
	SortOrder int  `runner:"code:sort_order;name:排序;widget:number;min:0;step:1;placeholder:请输入排序值" validate:"min=0" json:"sort_order"`
	IsVisible bool `runner:"code:is_visible;name:是否显示;widget:switch;on_text:显示;off_text:隐藏" json:"is_visible"`

	// 时间字段
	CreatedAt time.Time `runner:"code:created_at;name:创建时间;show:list,detail" json:"created_at"`
	UpdatedAt time.Time `runner:"code:updated_at;name:更新时间;show:detail" json:"updated_at"`
}

// ===== 使用场景说明 =====

/*
使用场景说明：

1. 用户管理 (UserExample):
   - 展示了完整的用户信息管理
   - 包含基础信息、状态控制、权限管理等
   - 演示了各种输入组件的使用

2. 商品管理 (ProductExample):
   - 展示了电商商品管理的复杂场景
   - 包含价格、库存、分类、规格等信息
   - 演示了数值输入、文件上传等组件

3. 订单管理 (OrderExample):
   - 展示了订单流程管理
   - 包含状态流转、金额计算、物流跟踪等
   - 演示了日期时间、状态选择等组件

4. 系统配置 (SystemConfigExample):
   - 展示了系统配置管理
   - 包含配置分组、数据类型、权限控制等
   - 演示了复杂的业务逻辑配置

这些示例涵盖了：
- 所有已实现的组件类型
- 常见的业务场景
- 复杂的字段关系
- 完整的CRUD操作需求
- 权限和状态控制
- 数据验证规则

可以作为实际开发的参考模板。
*/
