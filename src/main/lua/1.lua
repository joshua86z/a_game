require "luacurl"

-- is tongbu login
function isLogin(uid, session, sign, other)

	url = "http://tgi.tongbu.com/checkv2.aspx?k=" .. session

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

	html = table.concat(result)

	if uid == html and uid > '0' then
		return uid, true
	end

	return '0', false
end
