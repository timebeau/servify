# Handler Testdata

本目录存放 `apps/server/internal/handlers` 包测试需要长期保留的静态样本。

约定：

- `testdata/` 用于可提交、可复用、稳定的小型样本
- 运行时上传目录不作为测试样本目录使用
- 临时文件仍优先使用 `t.TempDir()`
- 如果样本可读性强，优先使用文本格式
- 大文件、敏感文件、可运行文件不应放入本目录
