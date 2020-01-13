# Microsoft Online RESTful API

### 安装

直接从 release 页面下载对应平台 binary 文件即可。 下载完成过后同时复制一份仓库的 `configs/config.toml` 到您的工作目录。

#### 配置

你需要在 portal.azure.com 创建一个 Azure Directory 的 Application 如何创建可以看这个 [How to: Use the portal to create an Azure AD application and service principal that can access resources](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal)

需要注意的是你需要在应用 -> Authencation -> Redirect URIs 配置好以你域名开始的 Redirect URI `https://domain.com/oauth/callback`

然后将 Overview 页面的相关信息和按照上文创建的 Secret Key 配置到 `config.toml` 配置项。

如果您使用的是世纪互联版本请按照配置文件注释修改对应的请求地址。

> 创建 Secret 的时候请设置一个较长的有效期

#### 如何授权

通常而言你配置好一个过后是可以授权多个账户的只需要访问 https://domain.com/oauth/authorize 按照顺序授权即可。目前我只测试了 Office 365 A1 和 Office 365 E3 其他的理论来说也是可以的。另外请注意您需要使用管理员来授权，不然会没有权限操作相关 API.

#### 运行

因为是 binary 文件您可以按照一下命令运行：

```bash
chmod +x msonline
GIN_MODE=release ./msonline serve --conf configs/config.toml
```

如果需要守护运行请自行采用 PM2、Supervisor 等. 同时我在在 `init` 目录中提供了 `systemd` 的配置文件可以复制到 `/etc/systemd/system/msonline.service` 使用。

### API 接口

访问接口采用 HTTP Basic 方式认证 其中 Username, Password 分别为您 `config.toml` 中配置的 AccessKey, AccessSecret。如果您不太熟悉 HTTP Basic 认证可以先自行 Google 查阅相关资料。

#### 已授权账户

`GET` https://domain.com/api/v1/accounts

##### 返回值

```json
{
    "message": "",
    "status_code": 200,
    "data": [
        {
        "email": "admin@msstu.onmicrosoft.com",
        "id": "a9046d12-cba0-4c11-9bbc-xxxx"
        }
    ]
}
```

#### 删除已授权账号

这里删除只是从本地数据库中删除，如果要取消后台请前往 Microsoft Azure 后台操作

`GET` https://domain.com/api/v1/accounts/{id}

##### 返回值

```json
{
    "message": "",
    "status_code": 200,
    "data": []
}
```

#### 账户 SKUS

`GET` https://domain.com/api/v1/accounts/{id}/skus

#### 返回值

```json
{
    "message": "",
    "status_code": 200,
    "data": [
        {
        "id": "314c4481-f395-4525-be8b-xxxxx",
        "name": "Microsoft Office 365 (Plan A2) for Students",
        "consumed_units": 250,
        "prepaid_units": 5000
        }
    ]
}
```

#### 创建用户

`POST` https://domain.com/api/v1/accounts/{id}/users

##### 请求参数

```json
{
  "name": "用户名",
  "principal_name": "youxiang", // 邮箱前缀
  "domain": "ns.onmicrosoft.com", // 如果不传 默认按照授权的管理员账户后缀分配
  "password": "密码",
  "sku_id": "314c4481-f395-4525-be8b-2ec4bb1e9d91" // 如果不传 默认不分配授权
}
```
#### 返回值

```json
{
    "message": "",
    "status_code": 200,
    "data": []
}
```

#### 删除用户

`POST` https://domain.com/api/v1/accounts/{id}/users/{uid/principal_name}

#### 返回值

```json
{
    "message": "",
    "status_code": 200,
    "data": []
}
```
