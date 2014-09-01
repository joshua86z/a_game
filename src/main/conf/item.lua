
index = "500,501,502,503,504,505,506,507,508,509"

-- ID , 名字 , 描述 , 能力数值 , 成长 , 掉落概率
function item(itemId)
	if itemId == 500 then
		return "磁铁", "吸引全图的金币", 5, 1, 10

	elseif itemId == 501 then
		return "金币", "这个东西很多就是高富帅啦", 1, 1, 30

	elseif itemId == 502 then
		return "能量瓶", "药不能停,吃多了会暴走", 10, 2, 15

	elseif itemId == 503 then
		return "保护罩", "罩子破了,会出人命的", 5, 1, 10

	elseif itemId == 504 then
		return "鞋子", "增加移动速度一段时间", 5, 1, 5

	elseif itemId == 505 then
		return "清醒剂", "一段时间内,对debuff免疫", 5, 1, 5

	elseif itemId == 506 then
		return "狂暴药水", "增加攻击速度一段时间", 5, 1, 5

	elseif itemId == 507 then
		return "极速药水", "这个东西很多就是高富帅啦", 1, 1, 5

	elseif itemId == 508 then
		return "血瓶", "回复一定血量", 1, 1, 10

	elseif itemId == 509 then
		return "钱箱", "会爆出满屏钱币的神奇箱子", 5, 1, 5

--	elseif itemId == 600 then
--		return 600, 生命, 死亡之后可以复活一次, 1, 0, 0
--	elseif itemId == 601 then
--		return 601, 加速, 增加移动速度很长一段时间, 20, 0, 0
--	elseif itemId == 602 then
--		return 602, 护盾, 增加一个护盾很长一段时间, 20, 0, 0
--	elseif itemId == 603 then
--		return 603, 攻速, 增加攻击速度很长一段时间, 20, 0, 0

	else
		return  "", "", 0, 0, 0
	end
end

-- 临时道具消耗的钻石
function tempItemDiamond()
	return 5, 6, 7, 8
end

-- 升级需要的金币
function levelUpCoin(level)
	if level == 0 then
		return 100
	elseif level == 1 then
		return 1000
	elseif level == 2 then
		return 5000
	elseif level == 3 then
		return 20000
	elseif level == 4 then
		return 50000
	else
		return 0
	end
end

max_level = 5
