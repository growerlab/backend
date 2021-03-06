## 依赖

### 数据库

PostgreSQL v12.x

KeyDB 5.x

#### 初始化数据库

创建数据库名称，数据库用户，用户密码均为 growerlab 的数据库

```
create database growerlab;
create user growerlab with encrypted password 'growerlab';
grant all privileges on database growerlab to growerlab;
```

#### 初始化数据库表结构

使用 `db/growerlab.sql` 文件初始化表结构

如果有种子数据，应该放入 `db/seed.sql` 文件中

```
初始账号 admin@growerlab.com
初始密码 growerlabadmin
```

#### GraphQL

基于 `gqlgen` 通过自动化生成GraphQL的基础代码

如果有修改 `*.graphql` 应该使用 `gqlgen` 工具生成代码
