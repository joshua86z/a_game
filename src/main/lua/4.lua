require "json"
require "luacurl"
require "md5"

-- is 91 login
function isLogin(uid, session, sign, other)

	sign = md5.sumhexa('1145944' .. uid .. session .. 'efe1a6db9ae7d27937c9eb9ae5ba43378a2168b10ddc6452')
	url = string.format("http://service.sj.91.com/usercenter/AP.aspx?AppId=114594&Act=4&Uin=%s&SessionId=%s&Sign=%s", uid, session, sign)

	local result = {}

	c = curl.new()
	c:setopt(curl.OPT_URL, url)
	c:setopt(curl.OPT_WRITEDATA, result)
	c:setopt(curl.OPT_WRITEFUNCTION,
		function(tab, buffer)
			table.insert(tab, buffer)
			return #buffer
		end)

	local ok = c:perform()

	c:close()

	local html = table.concat(result)
	local data = json.decode(html)

	if data["ErrorCode"] == "1" then
		return uid, true
	end

	return '0', false
end
