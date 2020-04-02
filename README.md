# redmine-user-generator

[ダウンロードはこちらから](https://github.com/Kanatani28/redmine-user-generator/releases)

Windows: redmine-user-generator_windows_x86_64.zip   
Mac: redmine-user-generator_darwin_x86_64.tar.gz   
Linux: redmine-user-generator_linux_x86_64.tar.gz   


config.ymlとusers.csvに設定をしてから実行ファイルを叩くと  
Redmineのユーザーが一気に作成されます。

## Redmine側の設定
管理者ユーザーでログインし、管理 > 設定 > API   
「RESTによるWebサービスを有効にする」にチェックを入れる。

個人設定からAPIアクセスキーが確認できます。

## config.yml

|項目|説明|
|--|--|
|api_key|個人設定で確認できるAPIアクセスキー|
|host|ホスト名|
|auth_user|管理者ユーザーID（APIキーと対応したユーザー）|
|auth_pass|管理者ユーザーパスワード|

## users.csv

CSVデータなのでExcel等で編集するとちょっと楽です。

|列名|説明|
|--|--|
|login|ログイン時に使用するユーザーID|
|password|パスワード。8文字以上で設定すること。|
|firstname|名|
|lastname|姓|
|mail|メールアドレス。基本は〇〇@example.comで設定すること|
|admin|管理者権限を付与するかどうか。trueかfalseで設定すること（小文字で）。|
