# Office 365 Web API

这是我第一次用 Go 写代码，只能保证可以运行，代码写的不好，请多见谅。

### 安装

直接从 release 页面下载对应平台 binary 文件即可。 下载完成后请复制一份 .app.ini.example 放到 binary 文件同级并命名为 app.ini

#### 配置

你需要在 portal.azure.com 创建一个 Azure Directory 的 Application 如何创建可以看这个 [How to: Use the portal to create an Azure AD application and service principal that can access resources](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal)

需要注意的是你需要在应用 -> Authencation -> Redirect URIs 配置好以你域名开始的 Redirect URI `https://domain.com/oauth/callback`

然后将 Overview 页面的相关信息和按照上文创建的 Secret Key 配置到 `app.ini` 配置项。

> 创建 Secret 的时候请设置一个较长的有效期

#### 如何授权

通常而言你配置好一个过后是可以授权多个账户的只需要访问 https://domain.com 按照顺序授权即可。目前我只测试了 Office 365 A1 和 Office 365 E3 其他的理论来说也是可以的。值得注意的是你只能使用管理员账户授权。

#### 运行

因为是 binary 文件您可以按照一下命令运行：

```bash
chmod +x binary
./binray
```

如果需要守护运行请自行采用 PM2、Supervisor 等.

### API 接口

访问接口采用 HTTP Basic 方式认证 其中 Username, Password 分别为您 app.ini 中配置的 AccessKey, AccessSecret。如果您不太熟悉 HTTP Basic 认证可以先自行 Google 查阅相关资料。

#### 已授权账户

`GET` https://domain.com/api/v1/accounts

##### 返回值

```json
{
  "data": {
    "user_id": "4aa04192-ef64-4f8b-82a2-fakefakefake",
    "email": "fake@fake.onmicrosoft.com",
    "created_at": "2019-08-27T17:01:23.652594211+08:00",
    "updated_at": "2019-08-28T09:56:44.042627983+08:00"
  },
  "error": ""
}
```

#### 账户 SKUS

`GET` https://domain.com/api/v1/skus/{USER_ID}

#### 返回值

```json
{
  "data": [
    {
      "sku_id": "314c4481-f395-4525-be8b-fakefakefake",
      "total": 5000,
      "used": 50,
      "friendly_name": "Microsoft Office 365 (Plan A2) for Students"
    }
  ],
  "error": ""
}
```

#### 创建用户

`POST` https://domain.com/api/v1/users

##### 请求参数

```json
{
  "account_id": "4aa04192-ef64-4f8b-82a2-fakefake",
  "enabled": true,
  "nickname": "用户名",
  "email": "youxiang",
  "password": "密码",
  "assign_license": true,
  "sku_id": "314c4481-f395-4525-be8b-2ec4bb1e9d91"
}
```

创建用户的时候有几点需要注意:

1. account_id 是 `获取已授权用户的中的 USER_ID` 这里有一点歧义没有设计好
2. nickname 随意设置，这里没有限制
3. email 是你最终登录的邮箱地址，请注意不要包含 @xx.com
4. sku_id 是 `获取账户 SKUS` 中返回的 sku_id

