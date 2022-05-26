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
	git commit -am "日志v3增加特性：错误界别的日志，单独记录新的文件"
	git push origin master
	git tag v1.7.43
	git push --tags
	@echo "\n tags 发布中..."
