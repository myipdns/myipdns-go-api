<div align="center">

# myipdns-go-api

[English](#english) | [简体中文](#简体中文) | [Русский](#русский) | [日本語](#日本語) | [Français](#français) | [Deutsch](#deutsch) | [Español](#español) | [Português](#português)

</div>

---

<a name="english"></a>
## 🇺🇸 English

**myipdns-go-api** is a high-performance IP address information query service built with Go and the Fiber v2 framework, providing IPv4/IPv6 geolocation and carrier information queries for [myipdns.com](https://myipdns.com).

### Technical Features

*   **High Performance Architecture**: Written in Go and based on the Fiber v2 Web framework, featuring extremely low memory footprint and high concurrency processing capabilities.
*   **Dual Stack Support**: Fully supports both IPv4 and IPv6 address queries.
*   **Multi-Source Data Fusion**: Combines **MaxMind GeoLite2** (primary geolocation) with **IP2Location/IP2Proxy** (network characteristics and risk control) to provide richer and more accurate IP profiling through intelligent fallback algorithms.
*   **Accurate Data**: Uses the MaxMind GeoLite2 database to provide city-level IP geolocation and ASN information.
*   **Smart ISP Translation**: Built-in multi-language ISP name translation engine automatically converts carrier names to the specified language based on request parameters.
*   **Flexible Response Modes**: Supports both plain text (IP only) and JSON (detailed info) response formats, intelligently switching based on the access path.

### Attribution
This project uses the following data sources:
*   **MaxMind GeoLite2**: Created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE**: IP2Location LITE data, available from [https://lite.ip2location.com](https://lite.ip2location.com).

### API Documentation
We provide free public API endpoints for testing and use.

#### 1. Standard Query Interface (JSON)
Returns detailed IP geolocation, ASN, and carrier information. Supports activating full data mode via the `full` parameter.

*   **Endpoint**: `GET https://api.myipdns.com/`
*   **Parameters**:

| Parameter | Type | Required | Description | Example |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | No | The IP address to query. If omitted, returns the current visitor's IP. | `?ip=1.1.1.1` |
| `lang` | String | No | Specifies the language for the result. Supports: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=en` |
| `full` | String | No | **[New]** Whether to enable full data mode. Activate by passing `true` to integrate IP2Proxy data and return more fields (e.g., VPN provider, fraud score, etc.). | `?full=true` |

*   **Response Example (Standard Mode)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Response Example (Full Mode)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Plain Text Interface (IPv4)
Returns only the visitor's IPv4 address, commonly used by scripts to obtain public IP.

*   **Endpoint**: `GET https://v4.api.myipdns.com/`
*   **Description**: Plain text format output.

#### 3. Plain Text Interface (IPv6)
Returns only the visitor's IPv6 address.

*   **Endpoint**: `GET https://v6.api.myipdns.com/`
*   **Description**: Plain text format output.

### Deployment
I don't think anyone will deploy this, so I won't write it. Let AI teach it.

---

<a name="简体中文"></a>
## 🇨🇳 简体中文

**myipdns-go-api** 是一个基于 Go 语言和 Fiber v2 框架构建的高性能 IP 地址信息查询服务，为 [myipdns.com](https://myipdns.com) 提供 IPv4/IPv6 地理位置及运营商信息查询服务。

### 技术特性

* **高性能架构**: 采用 Go 语言编写，基于 Fiber v2 Web 框架，具备极低的内存占用和极高的并发处理能力。
* **双栈支持**: 完美支持 IPv4 和 IPv6 地址查询。
* **多源数据融合**: 结合 **MaxMind GeoLite2** (主地理定位) 与 **IP2Location/IP2Proxy** (网络特征与风控)，通过智能回退算法提供更丰富、精准的 IP 画像。
* **智能 ISP 翻译**: 内置多语言 ISP 名称翻译引擎，可根据请求参数自动将运营商名称转换为指定语言。
* **灵活的响应模式**: 支持纯文本（仅 IP）、标准 JSON（基础信息）及全量 JSON（包含风控及详细 ISP）模式，根据访问路径及参数智能切换。

### 数据来源致谢
本项目使用了以下数据源：
* **MaxMind GeoLite2**: 由 MaxMind 创建，获取地址：[https://www.maxmind.com](https://www.maxmind.com)。
* **IP2Location LITE**: IP2Location LITE 数据，获取地址：[https://lite.ip2location.com](https://lite.ip2location.com)。

### API 接口文档
我们提供免费的公共 API 接口供测试和使用。

#### 1. 标准查询接口 (JSON)
返回详细的 IP 地理位置、ASN 及运营商信息。支持通过 `full` 参数激活全量数据模式。

* **接口地址**: `GET https://api.myipdns.com/`
* **参数说明**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | 否 | 指定查询的 IP 地址。若不传则返回当前访问者的 IP。 | `?ip=1.1.1.1` |
| `lang` | String | 否 | 指定返回结果的语言。支持：`en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`。 | `?lang=cn` |
| `full` | String | 否 | **[新]** 是否开启全量数据模式。传 `true` 时激活，将融合 IP2Proxy 数据并返回更多字段（如 VPN 提供商、欺诈评分等）。 | `?full=true` |

* **响应示例 (标准模式)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}

* **响应示例 (全量模式)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. 纯文本接口 (IPv4)
仅返回访问者的 IPv4 地址，常用于脚本获取公网 IP。

*   **接口地址**: `GET https://v4.api.myipdns.com/`
*   **说明**: 纯文本格式输出。

#### 3. 纯文本接口 (IPv6)
仅返回访问者的 IPv6 地址。

*   **接口地址**: `GET https://v6.api.myipdns.com/`
*   **说明**: 纯文本格式输出。

### 部署
我觉得不会有人部署这个，不写了。让ai教吧。

---

<a name="русский"></a>
## 🇷🇺 Русский

**myipdns-go-api** — это высокопроизводительный сервис запросов информации об IP-адресах, созданный на языке Go с использованием фреймворка Fiber v2. Он обеспечивает запросы геолокации IPv4/IPv6 и информации об операторе связи для [myipdns.com](https://myipdns.com).

### Технические характеристики

*   **Высокопроизводительная архитектура**: Написан на Go и основан на веб-фреймворке Fiber v2, отличается крайне низким потреблением памяти и высокой способностью обработки одновременных запросов.
*   **Поддержка двойного стека**: Полная поддержка запросов адресов IPv4 и IPv6.
*   **Слияние данных из нескольких источников**: Объединяет **MaxMind GeoLite2** (основная геолокация) с **IP2Location/IP2Proxy** (сетевые характеристики и контроль рисков) для предоставления более полного и точного профилирования IP с помощью интеллектуальных алгоритмов резервирования.
*   **Точные данные**: Использует базу данных MaxMind GeoLite2 для предоставления геолокации IP на уровне города и информации об ASN.
*   **Умный перевод ISP**: Встроенный многоязычный движок перевода названий ISP автоматически преобразует названия операторов на указанный язык на основе параметров запроса.
*   **Гибкие режимы ответа**: Поддерживает как простой текст (только IP), так и формат JSON (подробная информация), интеллектуально переключаясь в зависимости от пути доступа.

### Атрибуция
Этот проект использует следующие источники данных:
*   **MaxMind GeoLite2**: Создано MaxMind, доступно по адресу [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE**: Данные IP2Location LITE, доступны по адресу [https://lite.ip2location.com](https://lite.ip2location.com).

### Документация API
Мы предоставляем бесплатные публичные API для тестирования и использования.

#### 1. Стандартный интерфейс запроса (JSON)
Возвращает подробную геолокацию IP, ASN и информацию об операторе. Поддерживает активацию режима полных данных через параметр `full`.

*   **URL**: `GET https://api.myipdns.com/`
*   **Параметры**:

| Параметр | Тип | Обяз. | Описание | Пример |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | Нет | IP-адрес для запроса. Если не указан, возвращает IP текущего посетителя. | `?ip=1.1.1.1` |
| `lang` | String | Нет | Указывает язык результата. Поддерживаются: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=ru` |
| `full` | String | Нет | **[Новое]** Включить ли режим полных данных. Активируйте передачей `true` для интеграции данных IP2Proxy и возврата дополнительных полей (например, провайдер VPN, оценка мошенничества и т.д.). | `?full=true` |

*   **Пример ответа (Стандартный режим)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Пример ответа (Режим полных данных)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Текстовый интерфейс (IPv4)
Возвращает только IPv4-адрес посетителя, часто используется скриптами для получения публичного IP.

*   **URL**: `GET https://v4.api.myipdns.com/`
*   **Описание**: Вывод в текстовом формате.

#### 3. Текстовый интерфейс (IPv6)
Возвращает только IPv6-адрес посетителя.

*   **URL**: `GET https://v6.api.myipdns.com/`
*   **Описание**: Вывод в текстовом формате.

### Развертывание
Я думаю, никто не будет это развертывать, так что не буду писать. Пусть ИИ научит.

---

<a name="日本語"></a>
## 🇯🇵 日本語

**myipdns-go-api** は、Go言語とFiber v2フレームワークで構築された高性能なIPアドレス情報照会サービスで、[myipdns.com](https://myipdns.com) にIPv4/IPv6の地理位置情報および通信事業者情報の照会サービスを提供します。

### 技術的特徴

*   **高性能アーキテクチャ**: Go言語で記述され、Fiber v2 Webフレームワークに基づいており、極めて低いメモリ消費と高い並行処理能力を備えています。
*   **デュアルスタック対応**: IPv4およびIPv6アドレス照会を完全にサポートしています。
*   **マルチソースデータ融合**: **MaxMind GeoLite2**（主要な位置情報）と **IP2Location/IP2Proxy**（ネットワーク特性とリスク管理）を組み合わせ、インテリジェントなフォールバックアルゴリズムを通じて、より豊富で正確なIPプロファイリングを提供します。
*   **正確なデータ**: MaxMind GeoLite2データベースを採用し、都市レベルのIP位置情報とASN情報を提供します。
*   **スマートISP翻訳**: 内蔵の多言語ISP名翻訳エンジンにより、リクエストパラメータに基づいてキャリア名を指定された言語に自動的に変換します。
*   **柔軟なレスポンスモード**: プレーンテキスト（IPのみ）とJSON（詳細情報）の2つのレスポンス形式をサポートし、アクセスパスに基づいてインテリジェントに切り替えます。

### 帰属表示
本プロジェクトは以下のデータソースを使用しています：
*   **MaxMind GeoLite2**: MaxMind によって作成され、[https://www.maxmind.com](https://www.maxmind.com) から入手可能です。
*   **IP2Location LITE**: IP2Location LITE データ、[https://lite.ip2location.com](https://lite.ip2location.com) から入手可能です。

### API ドキュメント
テストや利用のために、無料の公開APIを提供しています。

#### 1. 標準照会インターフェース (JSON)
詳細なIPジオロケーション、ASN、およびキャリア情報を返します。`full` パラメータによる全データモードの有効化をサポートします。

*   **エンドポイント**: `GET https://api.myipdns.com/`
*   **パラメータ**:

| パラメータ名 | 型 | 必須 | 説明 | 例 |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | いいえ | 照会するIPアドレス。省略した場合、現在の訪問者のIPを返します。 | `?ip=1.1.1.1` |
| `lang` | String | いいえ | 結果の言語を指定します。対応言語: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`。 | `?lang=ja` |
| `full` | String | いいえ | **[新規]** 全データモードを有効にするかどうか。`true` を渡して有効にすると、IP2Proxy データが統合され、より多くのフィールド（VPNプロバイダー、詐欺スコアなど）が返されます。 | `?full=true` |

*   **レスポンス例 (標準モード)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **レスポンス例 (全データモード)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. プレーンテキストインターフェース (IPv4)
訪問者のIPv4アドレスのみを返します。スクリプトでパブリックIPを取得するためによく使用されます。

*   **エンドポイント**: `GET https://v4.api.myipdns.com/`
*   **説明**: プレーンテキスト形式の出力。

#### 3. プレーンテキストインターフェース (IPv6)
訪問者のIPv6アドレスのみを返します。

*   **エンドポイント**: `GET https://v6.api.myipdns.com/`
*   **説明**: プレーンテキスト形式の出力。

### デプロイ
これをデプロイする人はいないと思うので、書きません。AIに教えてもらってください。

---

<a name="français"></a>
## 🇫🇷 Français

**myipdns-go-api** est un service de requête d'informations sur les adresses IP haute performance construit avec Go et le framework Fiber v2, fournissant des services de requête de géolocalisation IPv4/IPv6 et d'informations sur les opérateurs pour [myipdns.com](https://myipdns.com).

### Caractéristiques Techniques

*   **Architecture Haute Performance**: Écrit en Go et basé sur le framework Web Fiber v2, offrant une empreinte mémoire extrêmement faible et des capacités de traitement de haute simultanéité.
*   **Support Double Pile**: Prend entièrement en charge les requêtes d'adresses IPv4 et IPv6.
*   **Fusion de Données Multi-Sources**: Combine **MaxMind GeoLite2** (géolocalisation principale) avec **IP2Location/IP2Proxy** (caractéristiques réseau et contrôle des risques) pour fournir un profilage IP plus riche et plus précis grâce à des algorithmes de repli intelligents.
*   **Données Précises**: Utilise la base de données MaxMind GeoLite2 pour fournir une géolocalisation IP au niveau de la ville et des informations ASN.
*   **Traduction Intelligente des FAI**: Le moteur de traduction de noms de FAI multilingue intégré convertit automatiquement les noms des opérateurs dans la langue spécifiée en fonction des paramètres de la requête.
*   **Modes de Réponse Flexibles**: Prend en charge les formats de réponse texte brut (IP uniquement) et JSON (informations détaillées), basculant intelligemment en fonction du chemin d'accès.

### Attribution
Ce projet utilise les sources de données suivantes :
*   **MaxMind GeoLite2** : Créé par MaxMind, disponible sur [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE** : Données IP2Location LITE, disponibles sur [https://lite.ip2location.com](https://lite.ip2location.com).

### Documentation API
Nous fournissons des API publiques gratuites pour les tests et l'utilisation.

#### 1. Interface de Requête Standard (JSON)
Renvoie la géolocalisation IP détaillée, l'ASN et les informations sur l'opérateur. Prend en charge l'activation du mode données complètes via le paramètre `full`.

*   **Endpoint**: `GET https://api.myipdns.com/`
*   **Paramètres**:

| Paramètre | Type | Requis | Description | Exemple |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | Non | L'adresse IP à interroger. Si omis, renvoie l'IP du visiteur actuel. | `?ip=1.1.1.1` |
| `lang` | String | Non | Spécifie la langue du résultat. Supporté : `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=fr` |
| `full` | String | Non | **[Nouveau]** Activer ou non le mode données complètes. Activez en passant `true` pour intégrer les données IP2Proxy et renvoyer plus de champs (par ex., fournisseur VPN, score de fraude, etc.). | `?full=true` |

*   **Exemple de Réponse (Mode Standard)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Exemple de Réponse (Mode Données Complètes)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Interface Texte Brut (IPv4)
Renvoie uniquement l'adresse IPv4 du visiteur, couramment utilisé par les scripts pour obtenir l'IP publique.

*   **Endpoint**: `GET https://v4.api.myipdns.com/`
*   **Description**: Sortie au format texte brut.

#### 3. Interface Texte Brut (IPv6)
Renvoie uniquement l'adresse IPv6 du visiteur.

*   **Endpoint**: `GET https://v6.api.myipdns.com/`
*   **Description**: Sortie au format texte brut.

### Déploiement
Je pense que personne ne déploiera ceci, donc je ne l'écris pas. Laissez l'IA l'enseigner.

---

<a name="deutsch"></a>
## 🇩🇪 Deutsch

**myipdns-go-api** ist ein hochleistungsfähiger Dienst zur Abfrage von IP-Adressinformationen, der mit Go und dem Fiber v2-Framework erstellt wurde und IPv4/IPv6-Geolokalisierungs- und Betreiberinformationsabfragedienste für [myipdns.com](https://myipdns.com) bereitstellt.

### Technische Merkmale

*   **Hochleistungsarchitektur**: Geschrieben in Go und basierend auf dem Fiber v2 Web-Framework, zeichnet es sich durch extrem geringen Speicherbedarf und hohe Verarbeitungsfähigkeiten bei Gleichzeitigkeit aus.
*   **Dual-Stack-Unterstützung**: Unterstützt vollständig sowohl IPv4- als auch IPv6-Adressabfragen.
*   **Multi-Source-Datenfusion**: Kombiniert **MaxMind GeoLite2** (primäre Geolokalisierung) mit **IP2Location/IP2Proxy** (Netzwerkmerkmale und Risikokontrolle), um durch intelligente Fallback-Algorithmen umfassendere und genauere IP-Profilierung zu ermöglichen.
*   **Präzise Daten**: Verwendet die MaxMind GeoLite2-Datenbank, um IP-Geolokalisierung auf Stadtebene und ASN-Informationen bereitzustellen.
*   **Intelligente ISP-Übersetzung**: Die integrierte mehrsprachige ISP-Namensübersetzungs-Engine konvertiert Betreibernamen automatisch basierend auf Anfrageparametern in die angegebene Sprache.
*   **Flexible Antwortmodi**: Unterstützt sowohl Klartext- (nur IP) als auch JSON-Antwortformate (detaillierte Infos) und schaltet basierend auf dem Zugriffspfad intelligent um.

### Danksagung
Dieses Projekt verwendet die folgenden Datenquellen:
*   **MaxMind GeoLite2**: Erstellt von MaxMind, verfügbar unter [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE**: IP2Location LITE-Daten, verfügbar unter [https://lite.ip2location.com](https://lite.ip2location.com).

### API-Dokumentation
Wir bieten kostenlose öffentliche APIs zum Testen und Verwenden an.

#### 1. Standard-Abfrageschnittstelle (JSON)
Gibt detaillierte IP-Geolokalisierung, ASN und Betreiberinformationen zurück. Unterstützt die Aktivierung des vollständigen Datenmodus über den Parameter `full`.

*   **Endpunkt**: `GET https://api.myipdns.com/`
*   **Parameter**:

| Parameter | Typ | Erfor. | Beschreibung | Beispiel |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | Nein | Die abzufragende IP-Adresse. Wenn weggelassen, wird die IP des aktuellen Besuchers zurückgegeben. | `?ip=1.1.1.1` |
| `lang` | String | Nein | Gibt die Sprache für das Ergebnis an. Unterstützt: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=de` |
| `full` | String | Nein | **[Neu]** Ob der vollständige Datenmodus aktiviert werden soll. Aktivieren Sie ihn durch Übergabe von `true`, um IP2Proxy-Daten zu integrieren und weitere Felder (z. B. VPN-Anbieter, Betrugswert usw.) zurückzugeben. | `?full=true` |

*   **Antwortbeispiel (Standardmodus)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Antwortbeispiel (Vollständiger Datenmodus)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Klartext-Schnittstelle (IPv4)
Gibt nur die IPv4-Adresse des Besuchers zurück, häufig von Skripten verwendet, um die öffentliche IP zu erhalten.

*   **Endpunkt**: `GET https://v4.api.myipdns.com/`
*   **Beschreibung**: Ausgabe im Klartextformat.

#### 3. Klartext-Schnittstelle (IPv6)
Gibt nur die IPv6-Adresse des Besuchers zurück.

*   **Endpunkt**: `GET https://v6.api.myipdns.com/`
*   **Beschreibung**: Ausgabe im Klartextformat.

### Bereitstellung
Ich glaube nicht, dass das jemand deployen wird, also schreibe ich es nicht. Lass es dir von der KI beibringen.

---

<a name="español"></a>
## 🇪🇸 Español

**myipdns-go-api** es un servicio de consulta de información de direcciones IP de alto rendimiento construido con Go y el framework Fiber v2, que proporciona servicios de consulta de geolocalización IPv4/IPv6 e información del operador para [myipdns.com](https://myipdns.com).

### Características Técnicas

*   **Arquitectura de Alto Rendimiento**: Escrito en Go y basado en el framework web Fiber v2, presenta una huella de memoria extremadamente baja y altas capacidades de procesamiento concurrente.
*   **Soporte de Doble Pila**: Soporta completamente consultas de direcciones tanto IPv4 como IPv6.
*   **Fusión de Datos de Múltiples Fuentes**: Combina **MaxMind GeoLite2** (geolocalización principal) con **IP2Location/IP2Proxy** (características de red y control de riesgos) para proporcionar un perfilado de IP más rico y preciso a través de algoritmos de respaldo inteligentes.
*   **Datos Precisos**: Utiliza la base de datos MaxMind GeoLite2 para proporcionar geolocalización IP a nivel de ciudad e información ASN.
*   **Traducción Inteligente de ISP**: El motor integrado de traducción de nombres de ISP multilingüe convierte automáticamente los nombres de los operadores al idioma especificado basándose en los parámetros de la solicitud.
*   **Modos de Respuesta Flexibles**: Soporta formatos de respuesta de texto plano (solo IP) y JSON (información detallada), cambiando inteligentemente según la ruta de acceso.

### Atribución
Este proyecto utiliza las siguientes fuentes de datos:
*   **MaxMind GeoLite2**: Creado por MaxMind, disponible en [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE**: Datos de IP2Location LITE, disponibles en [https://lite.ip2location.com](https://lite.ip2location.com).

### Documentación de la API
Ofrecemos API públicas gratuitas para pruebas y uso.

#### 1. Interfaz de Consulta Estándar (JSON)
Devuelve geolocalización IP detallada, ASN e información del operador. Soporta la activación del modo de datos completos a través del parámetro `full`.

*   **Endpoint**: `GET https://api.myipdns.com/`
*   **Parámetros**:

| Parámetro | Tipo | Req. | Descripción | Ejemplo |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | No | La dirección IP a consultar. Si se omite, devuelve la IP del visitante actual. | `?ip=1.1.1.1` |
| `lang` | String | No | Especifica el idioma para el resultado. Soporta: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=es` |
| `full` | String | No | **[Nuevo]** Si se habilita el modo de datos completos. Actívelo pasando `true` para integrar datos de IP2Proxy y devolver más campos (por ejemplo, proveedor de VPN, puntuación de fraude, etc.). | `?full=true` |

*   **Ejemplo de Respuesta (Modo Estándar)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Ejemplo de Respuesta (Modo de Datos Completos)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Interfaz de Texto Plano (IPv4)
Devuelve solo la dirección IPv4 del visitante, comúnmente usado por scripts para obtener la IP pública.

*   **Endpoint**: `GET https://v4.api.myipdns.com/`
*   **Descripción**: Salida en formato de texto plano.

#### 3. Interfaz de Texto Plano (IPv6)
Devuelve solo la dirección IPv6 del visitante.

*   **Endpoint**: `GET https://v6.api.myipdns.com/`
*   **Descripción**: Salida en formato de texto plano.

### Despliegue
No creo que nadie despliegue esto, así que no lo escribiré. Que la IA te enseñe.

---

<a name="português"></a>
## 🇵🇹 Português

**myipdns-go-api** é um serviço de consulta de informações de endereço IP de alto desempenho construído com Go e o framework Fiber v2, fornecendo serviços de consulta de geolocalização IPv4/IPv6 e informações da operadora para [myipdns.com](https://myipdns.com).

### Características Técnicas

*   **Arquitetura de Alto Desempenho**: Escrito em Go e baseado no framework web Fiber v2, apresentando uma pegada de memória extremamente baixa e altas capacidades de processamento simultâneo.
*   **Suporte Dual Stack**: Suporta totalmente consultas de endereços IPv4 e IPv6.
*   **Fusão de Dados de Múltiplas Fontes**: Combina **MaxMind GeoLite2** (geolocalização primária) com **IP2Location/IP2Proxy** (características de rede e controle de risco) para fornecer perfis de IP mais ricos e precisos por meio de algoritmos de fallback inteligentes.
*   **Dados Precisos**: Usa o banco de dados MaxMind GeoLite2 para fornecer geolocalização IP em nível de cidade e informações ASN.
*   **Tradução Inteligente de ISP**: O mecanismo integrado de tradução de nomes de ISP multilíngue converte automaticamente os nomes das operadoras para o idioma especificado com base nos parâmetros da solicitação.
*   **Modos de Resposta Flexíveis**: Suporta formatos de resposta de texto simples (apenas IP) e JSON (informações detalhadas), alternando de forma inteligente com base no caminho de acesso.

### Atribuição
Este projeto usa as seguintes fontes de dados:
*   **MaxMind GeoLite2**: Criado pela MaxMind, disponível em [https://www.maxmind.com](https://www.maxmind.com).
*   **IP2Location LITE**: Dados IP2Location LITE, disponíveis em [https://lite.ip2location.com](https://lite.ip2location.com).

### Documentação da API
Oferecemos APIs públicas gratuitas para testes e uso.

#### 1. Interface de Consulta Padrão (JSON)
Retorna geolocalização IP detalhada, ASN e informações da operadora. Suporta a ativação do modo de dados completos através do parâmetro `full`.

*   **Endpoint**: `GET https://api.myipdns.com/`
*   **Parâmetros**:

| Parámetro | Tipo | Obrig. | Descrição | Exemplo |
| :--- | :--- | :--- | :--- | :--- |
| `ip` | String | Não | O endereço IP a ser consultado. Se omitido, retorna o IP do visitante atual. | `?ip=1.1.1.1` |
| `lang` | String | Não | Especifica o idioma para o resultado. Suporta: `en`, `cn`, `ru`, `ja`, `fr`, `de`, `es`, `pt`. | `?lang=pt` |
| `full` | String | Não | **[Novo]** Se deve habilitar o modo de dados completos. Ative passando `true` para integrar dados IP2Proxy e retornar mais campos (por exemplo, provedor VPN, pontuação de fraude, etc.). | `?full=true` |

*   **Exemplo de Resposta (Modo Padrão)**:
`https://api.myipdns.com/?ip=8.8.4.4&lang=en`
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

*   **Exemplo de Resposta (Modo de Dados Completos)**:
https://api.myipdns.com/?ip=8.8.4.4&lang=en&full=true
```json
{
  "ip": "8.8.4.4",
  "continent": "North America",
  "continent_code": "NA",
  "country": "United States",
  "country_code": "US",
  "is_eu": false,
  "region": "Massachusetts",
  "region_code": "MA",
  "city": "Westfield",
  "time_zone": "America/New_York",
  "latitude": 42.1293,
  "longitude": -72.7522,
  "asn": 15169,
  "as_org": "GOOGLE",
  "isp": "Google LLC",
  "provider": "",
  "usage_type": "DCH",
  "domain": "google.com",
  "fraud_score": 0,
  "threat": "-",
  "is_proxy": false,
  "is_anycast": true,
  "is_satellite": false
}
```

#### 2. Interface de Texto Simples (IPv4)
Retorna apenas o endereço IPv4 do visitante, comumente usado por scripts para obter IP público.

*   **Endpoint**: `GET https://v4.api.myipdns.com/`
*   **Descrição**: Saída em formato de texto simples.

#### 3. Interface de Texto Simples (IPv6)
Retorna apenas o endereço IPv6 do visitante.

*   **Endpoint**: `GET https://v6.api.myipdns.com/`
*   **Descrição**: Saída em formato de texto simples.

### Implantação
Acho que ninguém vai implantar isso, então não vou escrever. Deixe a IA ensinar.
