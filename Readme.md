# sandwich

通过定义接口生成代码实现功能扩展(装饰器模式)

## Example 文件说明

- [example/origin.go](example/origin.go) 定义原本的对象
- [example/wrapper.go](example/wrapper.go) 定义包装的对象
- [example/wrapper_gen.go](example/wrapper_gen.go) 生成的代码
- [example/wrapper_test.go](example/wrapper_test.go) 测试代码

## 说明

### 假设需求

现有一个 Origin 对象(参考[example/origin.go](example/origin.go))，需要对 Origin 对象的方法进行扩展

1. 对 key 参数增加前缀
2. 统计耗时

### 实现

1. 定义包装代码([example/wrapper.go](example/wrapper.go))
2. `执行 go:generate` 生成代码([example/wrapper_gen.go](example/wrapper_gen.go))
3. 执行测试用例([example/wrapper_test.go](example/wrapper_test.go))
   最后执行测试代码后会输出

```
DoSomething1
start:2024-06-25 12:38:19.014668 +0800 CST m=+0.000469745
WrapperMethod                     # 包装方法输出  
DoSomething1, key->  prefix:key   # 改变后的入参，原来的方法输出
cost:0s                           # 包装方法输出，统计执行耗时  

DoSomething2
start:2024-06-25 12:38:19.014674 +0800 CST m=+0.000476388
WrapperMethod
DoSomething2, keys->  [prefix:key1 prefix:key2] 
cost:0s
```