package cache_redis

/*

	注意事项
	因为any类型，所以存储的时候，都转为json保存，加一个字段表示当前数据的类型
	GET之后再按支持的类型，做相应的数据转换，然后再返回

*/


import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/depend/redigo/redis"
	"errors"
	"encoding/json"
)

type (
	RedisCacheBase struct {
		name    string
		conn    *RedisCacheConnect
	}
)



//关闭连接
func (base *RedisCacheBase) Close() (error) {
	//redis的缓存库，不需要关闭任何东西
	//如果是mongodb，sql类的，就可能要了
	return nil
}


//获取数据
func (base *RedisCacheBase) Get(key string) (Any,error) {
	if base.conn.connect == nil {
		return nil,errors.New("连接失败")
	} else {
		conn := base.conn.connect.Get()
		defer conn.Close()

		realkey := base.conn.config.Prefix + key

		val,err := redis.String(conn.Do("GET", realkey))
		if err != nil {
			return nil,err
		} else {

			m := Map{}
			err := json.Unmarshal([]byte(val), &m)
			if err != nil {
				return nil,err
			} else {
				//解包
				return noggo.Cache.UnPack(m)
			}
		}
	}
}
//设置数据
func (base *RedisCacheBase) Set(key string, val Any, exps ...int64) (error) {

	//先打包
	value,err := noggo.Cache.Pack(val)
	if err != nil {
		return err
	} else {

		//再解析JSON
		jsonVal,err := json.Marshal(value)
		if err != nil {
			return err
		} else {

			//过期时间
			expiry := base.conn.config.Expiry
			if len(exps) > 0 {
				expiry = exps[0]
			}



			//检查连接
			if base.conn.connect == nil {
				return errors.New("连接失败")
			} else {
				//获取一个连接
				conn := base.conn.connect.Get()
				defer conn.Close()

				//带前缀
				realkey := base.conn.config.Prefix + key
				//传值表数据，如果有过期，才加上过期EX
				args := []interface{}{ realkey, jsonVal }
				if expiry > int64(0) {
					args = append(args, "EX", int(expiry))
				}

				_,err := conn.Do("SET", args...)
				if err != nil {
					return err
				} else {
					//成功
					return nil
				}
			}

		}
	}
}
//删除数据
func (base *RedisCacheBase) Del(key string) (error) {
	if base.conn.connect == nil {
		return errors.New("连接失败")
	} else {
		conn := base.conn.connect.Get()
		defer conn.Close()

		//key要加上前缀
		realkey := base.conn.config.Prefix + key

		_,err := conn.Do("DEL", realkey)
		if err != nil {
			return err
		} else {
			//成功
			return nil
		}
	}
}
//清理数据
func (base *RedisCacheBase) Clear(prefixs ...string) (error) {
	return nil
}

//获取keys
//暂时不支持前缀查询
func (base *RedisCacheBase) Keys(prefixs ...string) ([]string,error) {
	keys := []string{}



	return keys,nil
}
