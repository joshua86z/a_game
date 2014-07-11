package models

type Config_Pay_Center struct {
	Id      int    `db:"pay_config_id"`
	Name    string `db:"pay_name"`
	Rmb     int    `db:"pay_rmb"`
	Diamond int    `db:"pay_diamond"`
}

var config_pay_center []*Config_Pay_Center

func init() {
	if _, err := DB().Select(&config_pay_center, "SELECT * FROM `config_pay_center` ORDER BY `pay_config_id` ASC "); err != nil {
		panic(err)
	}
}

// 充值商店
func ConfigPayCenterList() []*Config_Pay_Center {

	return config_pay_center
}

func GetPayCenterById(id int) *Config_Pay_Center {

	configlist := ConfigPayCenterList()

	for _, result := range configlist {
		if result.Id == id {
			return result
		}
	}

	return nil
}
