package ifile

import (
	"log"
	"path"
	"testing"
)

func TestExt(t *testing.T) {
	img := `https://p6-pc-sign.douyinpic.com/tos-cn-i-0813c001/oAIeFDyAKlaaqdHpzAEiRtfBLAACJgNIhQAAvr~tplv-dy-aweme-images:q75.webp?biz_tag=aweme_images&from=3213915784&s=PackSourceEnum_AWEME_DETAIL&sc=image&se=false&x-expires=1706104800&x-signature=5acd8zq%2BuPIUm0rz1y/xWX/yiHQ=`
	log.Println(FileExt(img))
	log.Println(path.Ext(img))

	log.Println(NewFileName(img).GetExt())

}
