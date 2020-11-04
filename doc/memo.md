

# MySQL バージョン確認
```
$ mysql --version

mysql  Ver 8.0.22 for osx10.15 on x86_64 (Homebrew)
```


# mysql起動と停止 
```
$ mysql.server start

$ mysql.server stop
```


# ユーザ情報の変更

https://qiita.com/azmint/items/6c7db6bc79f9789c5e6d  
https://dev.mysql.com/doc/refman/8.0/en/alter-user.html  
https://qiita.com/evansplus/items/66e06df882e9f40c092f  

## 認証形式を確認

それぞれのアカウントがどのように認証されるべきかを持っている 
確認方法： mysql.user テーブルには plugin というカラム  
  
MySQLの認証はチャレンジ・レスポンス認証なので、クライアントとサーバーで同じ認証形式をサポートしている必要がある  

MySQL 8.0.4とそれ以降は、認証プラグインを指定しないでユーザーを作成した場合のデフォルトが  
`caching_sha2_password`になっているので、  
MySQL 8.0.4とそれ以降にアップグレードしてから(プラグインを指定せずに)新しく作ったユーザーに対して、  
5.7とそれ以前と`caching_sha2_password`をサポートしていないライブラリーで接続しようとすると炸裂する。

```
mysql> use mysql
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> 
mysql> 
mysql> SELECT user, plugin FROM user;
+------------------+-----------------------+
| user             | plugin                |
+------------------+-----------------------+
| mysql.infoschema | caching_sha2_password |
| mysql.session    | caching_sha2_password |
| mysql.sys        | caching_sha2_password |
| root             | caching_sha2_password |
+------------------+-----------------------+
4 rows in set (0.00 sec)

```

## 接続ユーザの認証形式を変更
```
//パスワード値とともに、認証プラグインを指定する
ALTER USER "root"@"localhost" IDENTIFIED WITH 'mysql_native_password' BY "root";

```


# DB作成
```
CREATE DATABASE IF NOT EXISTS accountBook;
```


# 機能

## 購入品入力
金額は、以下のように管理する  
テーブル：税込金額のみ管理  
画面：税抜、税率、税込　それぞれ入力できるようにして、計算した結果をテーブルに格納する  

## 購入品一覧表示
購入品の名称、個数、使用個数などを表示する  

## 使用品入力
購入品から選ぶ形とする  
（条件：使用済、廃棄以外のもの）  
入力したら、購入品の個数と使用個数（過去分を含）比較し、購入品のステータスを更新する  

## 集計表示
月単位の購入金額（カテゴリー別）、予算、予算ー購入金額を表示する  



# 今後実装したいこと
全てGitHubのissueにて管理

- [ ] 更新機能の追加 
  - [ ] 購入品情報の更新機能追加
  - [ ] 使用品情報の更新機能追加
- [ ] 集計機能追加
  - [ ] 所属月、カテゴリー別購入金額の集計
  - [ ] Google Chart API を使って集計結果のグラフ化
- [ ] 使用品入力ダイアログの機能修正
  - [ ] 購入数をオーバーした数を入力できないように対応
- [ ] 一覧画面の機能追加
  - [ ] 表示条件を追加（所属月、カテゴリー）
  - [ ] 使用品情報についても、購入品のステータスが使用済、廃棄の場合には、グレーアウト表示とする
- [ ] 内部の処理で不要な項目までDBから取得している部分の修正（updateStatus -> getPurchasedItemSQL）
- [ ] 複数ユーザ対応
