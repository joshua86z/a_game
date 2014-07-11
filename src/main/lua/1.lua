
--  get tongu sign url
function getUrl(uid, session, sign, other)
	return "http://tgi.tongbu.com/checkv2.aspx?k=" .. session
end


-- get tongu sign url
function isPost()
	return false
end

-- get http post data
function getPost(uid, session, sign, other)
	return ""
end

-- is tongbu login
function isLogin(uid, session, sign, other, str)

    if uid == str and uid > 0 then
        return true
    end

    return false

end
