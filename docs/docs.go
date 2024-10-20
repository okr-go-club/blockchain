// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/blockchain/mine/{id}": {
            "get": {
                "parameters": [
                    {
                        "type": "string",
                        "description": "Mining process ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.MineStatusResponse"
                            }
                        }
                    }
                }
            }
        },
        "/blocks/pool": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chain.Blockchain"
                        }
                    }
                }
            }
        },
        "/blocks/pool/": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/chain.Transaction"
                            }
                        }
                    }
                }
            }
        },
        "/transactions": {
            "post": {
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.AddTransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "api.AddTransactionRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string",
                    "example": "0"
                },
                "from": {
                    "type": "string"
                },
                "privateKey": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        },
        "api.MineStatusResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "chain.Block": {
            "type": "object",
            "properties": {
                "capacity": {
                    "type": "integer"
                },
                "hash": {
                    "type": "string"
                },
                "nonce": {
                    "type": "integer"
                },
                "previousHash": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/chain.Transaction"
                    }
                }
            }
        },
        "chain.Blockchain": {
            "type": "object",
            "properties": {
                "blocks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/chain.Block"
                    }
                },
                "difficulty": {
                    "type": "integer"
                },
                "maxBlockSize": {
                    "type": "integer"
                },
                "miningReward": {
                    "type": "number"
                },
                "pendingTransactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/chain.Transaction"
                    }
                },
                "storage": {}
            }
        },
        "chain.Transaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "fromAddress": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "toAddress": {
                    "type": "string"
                },
                "transactionId": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
