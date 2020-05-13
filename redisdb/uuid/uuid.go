/**
 * @ClassName uuid
 * @Description //TODO 创建唯一到id
 * @Author liwei
 * @Date 2020/5/13 11:53
 * @Version go.redis.db V1.0
 **/

package uuid

import "github.com/beinan/fastid"

// @Title NewUserId
// @Description 唯一id
// @Description https://zhuanlan.zhihu.com/p/38308576
// @Description 参考文档
// @Return snowflake.ID
func NewUserId()(int64)  {
	id:=fastid.CommonConfig.GenInt64ID()
	return id
}
