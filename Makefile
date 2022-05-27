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
	git commit -am "优化企业微信消息通知形式，文本，卡片"
	git push origin master
	git tag v1.7.45
	git push --tags
	@echo "\n tags 发布中..."
