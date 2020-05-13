/**
 * @ClassName redisdb
 * @Description //TODO 
 * @Author liwei
 * @Date 2020/5/13 11:38
 * @Version go.redis.db V1.0
 **/

package redisdb

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.redis.db/redisdb/uuid"
	"strconv"
)

const Redis_Key_Prefix = "go.redis.db"

type DBModel struct {
	gorm.Model
	Uuid int64 `gorm:"column:user_id;unique_index;size:20;not null;comment:'服务器生成的有序的用户id'"`
}

type DBModelSet struct {
	Model DBModel
	Redis *redis.Client
	Db *gorm.DB
	TableName string
}

/**
* @Title AddModel
* @Description:添加模型
* @Description:添加模型到redis
* @Description:添加模型到db
* @Param:
* @return:
* @Author: liwei
* @Date: 2020/5/13
**/
func (DBModelSet) AddRedisDbModel(modelSet DBModelSet)bool{
	if modelSet.Model.Uuid==0{
		modelSet.Model.Uuid= uuid.NewUserId()
	}
	cacheKey:=Redis_Key_Prefix+ strconv.FormatInt(modelSet.Model.Uuid,10)
	modelJson, _ := json.Marshal(modelSet.Model)
	err:=modelSet.Redis.SetNX(cacheKey,modelJson,-1).Err()
	if err!=nil{
		fmt.Println("redis save fail")
		return false
	}
	go func() {
		modelSet.Db.Table(modelSet.TableName).Create(&modelSet.Model)
	}()
	return true
}

/**
* @Title DelRedisDbModel
* @Description:  删除缓存模型
* @Description:  删除缓存
* @Description:  删除数据库
* @Description:  警告 删除记录时，需要确保其主要字段具有值，GORM将使用主键删除记录，如果主要字段为空，GORM将删除模型的所有记录
* @Param:
* @return:
* @Author: liwei
* @Date: 2020/5/13
**/
func(DBModelSet) DelRedisDbModel(modelSet DBModelSet)bool{
	cacheKey:=Redis_Key_Prefix+ strconv.FormatInt(modelSet.Model.Uuid,10)
	modelSet.Redis.Del(cacheKey)
	go func() {
		modelSet.Db.Table(modelSet.TableName).Delete(&modelSet.Model)
	}()
	return true
}

/**
* @Title UpdateRedisDbModel
* @Description: 更新redis 和数据库
* @Param:
* @return:
* @Author: liwei
* @Date: 2020/5/13
**/
func(DBModelSet)UpdateRedisDbModel(modelSet DBModelSet) bool{
	cacheKey:=Redis_Key_Prefix+ strconv.FormatInt(modelSet.Model.Uuid,10)
	modelJson, _ := json.Marshal(modelSet.Model)
	err:=modelSet.Redis.SetNX(cacheKey,modelJson,-1).Err()
	if err!=nil{
		fmt.Println("redis save fail")
		return false
	}
	go func() {
		modelSet.Db.Table(modelSet.TableName).Model(&modelSet.Model).Updates(modelSet.Model)
	}()
	return true
}

/**
* @Title SelectRedisDbModel
* @Description:  查询属性
* @Param:
* @return:
* @Author: liwei
* @Date: 2020/5/13
**/
func (DBModelSet)SelectRedisDbModel(modelSet DBModelSet) DBModel{
	cacheKey:=Redis_Key_Prefix+ strconv.FormatInt(modelSet.Model.Uuid,10)
	status,_:=modelSet.Redis.Exists(cacheKey).Result()
	if status==0{
		modelSet.Db.Table(modelSet.TableName).Where("uuid=?",modelSet.Model.Uuid).First(&modelSet.Model)
		modelJson, _ := json.Marshal(modelSet.Model)
		modelSet.Redis.SetNX(cacheKey,modelJson,-1).Err()
	}
	modelJson,_:=modelSet.Redis.Get(cacheKey).Result()
	json.Unmarshal([]byte(modelJson),&modelSet.Model)
	return modelSet.Model
}

