// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "kevin"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/case/create": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "case"
                ],
                "summary": "發案/案件建立",
                "parameters": [
                    {
                        "description": "發案",
                        "name": "case",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateCase"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "/case/getAll": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "case"
                ],
                "summary": "接案查詢",
                "parameters": [
                    {
                        "type": "string",
                        "description": "接案查詢",
                        "name": "case",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        },
        "/case/getDetail/{caseid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "case"
                ],
                "summary": "案件詳細資料",
                "parameters": [
                    {
                        "type": "string",
                        "description": "案件詳細資料",
                        "name": "case",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        },
        "/download/{filename}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "下載檔案",
                "parameters": [
                    {
                        "type": "string",
                        "description": "下載檔案",
                        "name": "files",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    }
                }
            }
        },
        "/member/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "登入",
                "parameters": [
                    {
                        "description": "登入成功回傳 token",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/member/profile/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "GetProfile",
                "parameters": [
                    {
                        "type": "string",
                        "default": "”",
                        "description": "帳號",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/member/registerStep1": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "註冊 Step1",
                "parameters": [
                    {
                        "description": "檢查帳號是否已存在",
                        "name": "register",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RegisterStep1"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "/member/registerStep3": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "註冊 Step3",
                "parameters": [
                    {
                        "description": "註冊帳號",
                        "name": "register",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RegisterStep3"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "/profile/save": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "個人資料儲存",
                "parameters": [
                    {
                        "description": "修改成功回傳 boolean",
                        "name": "MyAccount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ProfileSave"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "上傳檔案",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "上傳檔案",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.CreateCase": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "cityTalk": {
                    "type": "string",
                    "example": "04"
                },
                "cityTalk2": {
                    "type": "string",
                    "example": "26232356"
                },
                "contactTime": {
                    "type": "string",
                    "example": "m,a"
                },
                "email": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "expectDate": {
                    "type": "string",
                    "example": "2022/02/02"
                },
                "expectDateChk": {
                    "type": "string",
                    "example": "1"
                },
                "expectMoney": {
                    "type": "string",
                    "example": "5000"
                },
                "extension": {
                    "type": "string",
                    "example": "0000"
                },
                "filesName": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "[a.jpg]"
                    ]
                },
                "kind": {
                    "type": "string",
                    "example": "o:一般案件 i:急件"
                },
                "line": {
                    "type": "string",
                    "example": "imlineid"
                },
                "name": {
                    "type": "string",
                    "example": "kevin"
                },
                "phone": {
                    "type": "string",
                    "example": "0999999999"
                },
                "title": {
                    "type": "string",
                    "example": "電傷平台架設"
                },
                "type": {
                    "type": "string",
                    "example": "程式開發"
                },
                "workArea": {
                    "type": "string",
                    "example": "台北市 信義區"
                },
                "workAreaChk": {
                    "type": "string",
                    "example": "1"
                },
                "workContent": {
                    "type": "string",
                    "example": "電傷平台架設，伺服器維護..."
                }
            }
        },
        "model.Login": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "model.ProfileSave": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "introduction": {
                    "type": "string",
                    "example": "我有 8000 名部下"
                },
                "name": {
                    "type": "string",
                    "example": "桐谷和人"
                },
                "oldPassword": {
                    "type": "string",
                    "example": "123456"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "phone": {
                    "type": "string",
                    "example": "0999999999"
                },
                "pwdSwitch": {
                    "type": "boolean",
                    "example": true
                },
                "zipcode": {
                    "type": "string",
                    "example": "200"
                }
            }
        },
        "model.RegisterStep1": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                }
            }
        },
        "model.RegisterStep3": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "email": {
                    "type": "string",
                    "example": "kevin@gmail.com"
                },
                "introduction": {
                    "type": "string",
                    "example": "我有 8000 名部下"
                },
                "name": {
                    "type": "string",
                    "example": "桐谷和人"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "phone": {
                    "type": "string",
                    "example": "0999999999"
                },
                "zipcode": {
                    "type": "string",
                    "example": "200"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.8.10",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Gin swagger",
	Description:      "Gin swagger",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
