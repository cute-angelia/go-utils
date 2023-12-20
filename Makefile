.PHONY: up tag

up:
	git pull origin master
	git add .
	git commit -am "update"
	git push origin master
	@echo "\n 发布中..."

tag:
	git pull origin master
	git add .
	git commit -am "增加位权限操作，方便压缩数据"
	git push origin master
	git tag v1.8.33
	git push --tags
	@echo "\n tags 发布中..."
