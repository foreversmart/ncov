## ncov

全国新型冠状病毒，肺炎疫情增长情况打分和全国市级打分结果展示。china novel coronavirus (nCoV) city level score and visualization.

* 由于数据和算法都存在很大缺陷打分结果只供一定的参考，切勿造谣传谣信谣 

#### 初衷

这些天一直在看网上公开的数据，很多数据都是静态的数据，特别是到市这一级，只能查到当前确诊人数多少，
历史人数和趋势无法查看。我所在的城市据我观察已经很多天都是 7 人没有变化了。所以我在想能否根据一个城市
最近一段时间新增确诊人数的变化和总量来一定程度上去衡量一个城市疫情的严重程度。

* 当前想了一个简单的线性算法来进行打分

    score = 前一天的 score + 今天新增确诊人数 x 10 - （累计前10之内新增确诊人数的总和）
    
##### 运行

```
// 可视化
// 打开前确保 index.html、data.txt、merge.json 文件存在
在浏览器中 打开 index.html

```

```
// 更新结果数据
// 运行前保证 main.go DXYArea.csv 文件的存在
cd ncov/
go run main.go
// 运行后会生成 结果数据存到 data.txt 中
```

```cassandraql
// 自定义 Calc 函数来自定义城市的打分计算函数
// city 数据结构中可以方便的获取该城市历史确诊人数的数据点，这些数据点以数组的形式按更新时间降序排列
func (c *City) Calc() int {
  return 0
}

```

##### 数据来源

* 疫情历史数据：

> DXYArea.csv 到 DXY-2019-nCoV-Data 中获取最新数据
把 DXY-2019-nCoV-Data/csv 下的 DXYArea.csv 替换当前文件 
 
* 数据源：https://github.com/BlankerL/DXY-2019-nCoV-Data

##### 可视化

* 使用 echart 快速的进行可视化 https://www.echartsjs.com/en/index.html

# 中国加油，武汉加油！

