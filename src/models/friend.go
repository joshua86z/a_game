package models



func init() {

}

var Friend FriendModel

type FriendModel struct {
}

func (this FriendModel) List(uid int64) []int64 {

	var friendList []struct {
		Uid int64 `db:"uid"`
		Fid int64 `db:"friend_uid"`
	}

	sql := "SELECT * FROM friends WHERE `uid` = ? OR `friend_uid` = ?"
	DB().Select(&friendList, sql, uid, uid)

	var result []int64
	for _, val := range friendList {
		if val.Uid != uid {
			result = append(result, val.Uid)
		}
		if val.Fid != uid {
			result = append(result, val.Fid)
		}
	}

	return result
}
