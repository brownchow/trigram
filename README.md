# Go实现Trigram算法



Trigram算法用于 Google [CodeSearch](https://github.com/google/codesearch)， 是一个用于在大量源代码文件中进行索引和正则搜索的命令行工具集

## 如何工作?

trigram 算法

一些规则：

- 不会把大写转成小写
- 包含 空白符

自然语言处理中的 **N-Gram模型**，用于判断**两个字符串**之间的**相似度**

三元语言模型

三元语言模型的优缺点：

- 高阶 n-gram 对更多的上下文敏感
- 底阶 n-gram 考虑非常有限的上下文信息



## 参考：

https://en.wikipedia.org/wiki/Trigram

这个项目或许可以帮助理解类似ElasticSearch这类搜索引擎是怎么工作的

https://blog.csdn.net/CoderPai/article/details/80403897





