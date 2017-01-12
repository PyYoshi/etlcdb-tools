etlcdb-tools
============

現在ETL9Gをパーズする実装をただ書きなぐってる段階

破壊的な変更を逐次行うので利用しないほうがよい

# ゴール

ETL1~9までのファイルをPNG+JSON(メタデータ)なデータセットを生成できるようにする

またCLIの提供の他にライブラリとして利用できるようにする

# TODO

- ETL1パーズ
- ETL2パーズ
- ETL3パーズ
- ETL4パーズ
- ETL5パーズ
- ETL6パーズ
- ETL7パーズ
- ETL8Gパーズ
- ETL8Bパーズ
- ETL9Gパーズ
- ETL9Bパーズ
- 任意の解像度で出力可能にする
- 各フォーマットのレコード情報からメタデータを作成可能にする
- CLIの提供

# CLI

まだ実装していません

## 引数

### --format (-f)

以下の値を指定可能

- 1
- 2
- 3
- 4
- 5
- 6
- 7
- 8g
- 8b
- 9g
- 9b

指定しない場合はすべてのフォーマットをパーズする

### --etlcdb-dir (-e)

指定フォーマット(ETL1~9)が保存されているディレクトリパスを指定

階層構造は``etlcdbディレクトリの階層構造``を参照

### --datasets-dir (-o)

パーズしたデータを格納するためのディレクトリパスを指定

# メタデータの構造

まだ実装していません

## datasets/ETL1/etl1.json
## datasets/ETL2/etl2.json
## datasets/ETL3/etl3.json
## datasets/ETL4/etl4.json
## datasets/ETL5/etl5.json
## datasets/ETL6/etl6.json
## datasets/ETL7/etl7.json
## datasets/ETL8G/etl8g.json
## datasets/ETL8B/etl8b.json
## datasets/ETL9G/etl9g.json
## datasets/ETL9B/etl9b.json

# etlcdbディレクトリの階層構造

以下のようにETL1~9のファイルを展開すること

```
etlcdb
├─ETL1
│      ETL1C_01
│      ETL1C_02
│      ETL1C_03
│      ETL1C_04
│      ETL1C_05
│      ETL1C_06
│      ETL1C_07
│      ETL1C_08
│      ETL1C_09
│      ETL1C_10
│      ETL1C_11
│      ETL1C_12
│      ETL1C_13
│      ETL1INFO
│
├─ETL2
│      ETL2INFO
│      ETL2_1
│      ETL2_2
│      ETL2_3
│      ETL2_4
│      ETL2_5
│
├─ETL3
│      ETL3C_1
│      ETL3C_2
│      ETL3INFO
│
├─ETL4
│      ETL4C
│      ETL4INFO
│
├─ETL5
│      ETL5C
│      ETL5INFO
│
├─ETL6
│      ETL6C_01
│      ETL6C_02
│      ETL6C_03
│      ETL6C_04
│      ETL6C_05
│      ETL6C_06
│      ETL6C_07
│      ETL6C_08
│      ETL6C_09
│      ETL6C_10
│      ETL6C_11
│      ETL6C_12
│      ETL6INFO
│
├─ETL7
│      ETL7INFO
│      ETL7LC_1
│      ETL7LC_2
│      ETL7SC_1
│      ETL7SC_2
│
├─ETL8B
│      ETL8B2C1
│      ETL8B2C2
│      ETL8B2C3
│      ETL8INFO
│
├─ETL8G
│      ETL8G_01
│      ETL8G_02
│      ETL8G_03
│      ETL8G_04
│      ETL8G_05
│      ETL8G_06
│      ETL8G_07
│      ETL8G_08
│      ETL8G_09
│      ETL8G_10
│      ETL8G_11
│      ETL8G_12
│      ETL8G_13
│      ETL8G_14
│      ETL8G_15
│      ETL8G_16
│      ETL8G_17
│      ETL8G_18
│      ETL8G_19
│      ETL8G_20
│      ETL8G_21
│      ETL8G_22
│      ETL8G_23
│      ETL8G_24
│      ETL8G_25
│      ETL8G_26
│      ETL8G_27
│      ETL8G_28
│      ETL8G_29
│      ETL8G_30
│      ETL8G_31
│      ETL8G_32
│      ETL8G_33
│      ETL8INFO
│
├─ETL9B
│      ETL9B_1
│      ETL9B_2
│      ETL9B_3
│      ETL9B_4
│      ETL9B_5
│      ETL9INFO
│
├─ETL9G
│      ETL9G_01
│      ETL9G_02
│      ETL9G_03
│      ETL9G_04
│      ETL9G_05
│      ETL9G_06
│      ETL9G_07
│      ETL9G_08
│      ETL9G_09
│      ETL9G_10
│      ETL9G_11
│      ETL9G_12
│      ETL9G_13
│      ETL9G_14
│      ETL9G_15
│      ETL9G_16
│      ETL9G_17
│      ETL9G_18
│      ETL9G_19
│      ETL9G_20
│      ETL9G_21
│      ETL9G_22
│      ETL9G_23
│      ETL9G_24
│      ETL9G_25
│      ETL9G_26
│      ETL9G_27
│      ETL9G_28
│      ETL9G_29
│      ETL9G_30
│      ETL9G_31
│      ETL9G_32
│      ETL9G_33
│      ETL9G_34
│      ETL9G_35
│      ETL9G_36
│      ETL9G_37
│      ETL9G_38
│      ETL9G_39
│      ETL9G_40
│      ETL9G_41
│      ETL9G_42
│      ETL9G_43
│      ETL9G_44
│      ETL9G_45
│      ETL9G_46
│      ETL9G_47
│      ETL9G_48
│      ETL9G_49
│      ETL9G_50
│      ETL9INFO
```

# ライセンス

[The MIT License](https://pyyoshi.mit-license.org/)
