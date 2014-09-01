require "json"
require "luacurl"

-- is anysdk login
function isLogin(uid, session, sign, other)

	url = 'http://oauth.anysdk.com/api/User/LoginOauth/'

	local result = {}

	c = curl.new()
	c:setopt(curl.OPT_URL, url)
	c:setopt(curl.OPT_WRITEDATA, result)
	c:setopt(curl.OPT_POSTFIELDS, other) --POST
	c:setopt(curl.OPT_WRITEFUNCTION,
		function(tab, buffer)
			table.insert(tab, buffer)
			return #buffer
		end)

	local ok = c:perform()

	c:close()

	local html = table.concat(result)
	local data = json.decode(html)

	if data["status"] == "ok" then
		return data["data"]["id"], true
	end

	return '0', false
end
