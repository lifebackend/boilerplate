// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Stats operator API",
    "version": "0.1.0"
  },
  "paths": {
    "/_livenessProbe": {
      "get": {
        "description": "Liveness Probe",
        "tags": [
          "health"
        ],
        "summary": "Liveness Probe",
        "operationId": "getLivenessProbe",
        "responses": {
          "200": {
            "description": "Successful Response",
            "schema": {
              "$ref": "#/definitions/LivenessProbe"
            }
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/_readinessProbe": {
      "get": {
        "description": "Readiness Probe",
        "tags": [
          "health"
        ],
        "summary": "Readiness Probe",
        "operationId": "getReadinessProbe",
        "responses": {
          "200": {
            "description": "Successful Response"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/app-codes": {
      "get": {
        "description": "Get list of application codes",
        "tags": [
          "general"
        ],
        "summary": "Get List of Application Codes",
        "operationId": "getAppCodes",
        "responses": {
          "200": {
            "description": "Successful Response",
            "schema": {
              "$ref": "#/definitions/GetAppMessagesResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorMessage"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "AppMessage": {
      "type": "object",
      "title": "AppMessage",
      "required": [
        "message",
        "code",
        "attributes"
      ],
      "properties": {
        "attributes": {
          "description": "attributes",
          "type": "array",
          "items": {
            "type": "string"
          },
          "x-nullable": false
        },
        "code": {
          "type": "string",
          "title": "Code",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "title": "Message",
          "x-nullable": false
        }
      }
    },
    "ErrorMessage": {
      "type": "object",
      "title": "ErrorMessage",
      "required": [
        "message"
      ],
      "properties": {
        "attributes": {
          "description": "attributes",
          "type": "object",
          "additionalProperties": {
            "type": "string"
          },
          "x-nullable": false
        },
        "code": {
          "type": "string",
          "title": "Code",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "title": "Message",
          "x-nullable": false
        }
      }
    },
    "GetAppMessagesResponse": {
      "type": "array",
      "title": "GetAppCodesResponse",
      "items": {
        "$ref": "#/definitions/AppMessage"
      }
    },
    "LivenessProbe": {
      "type": "object",
      "title": "LivenessProbe",
      "required": [
        "tag"
      ],
      "properties": {
        "components": {
          "type": "array",
          "title": "LivenessProbeComponents",
          "items": {
            "$ref": "#/definitions/LivenessProbeComponent"
          }
        },
        "tag": {
          "type": "string",
          "title": "Tag",
          "x-nullable": false
        }
      }
    },
    "LivenessProbeComponent": {
      "type": "object",
      "title": "LivenessProbeComponent",
      "required": [
        "name",
        "status"
      ],
      "properties": {
        "name": {
          "type": "string",
          "title": "Name",
          "x-nullable": false
        },
        "status": {
          "type": "boolean",
          "title": "Status",
          "x-nullable": false
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Stats operator API",
    "version": "0.1.0"
  },
  "paths": {
    "/_livenessProbe": {
      "get": {
        "description": "Liveness Probe",
        "tags": [
          "health"
        ],
        "summary": "Liveness Probe",
        "operationId": "getLivenessProbe",
        "responses": {
          "200": {
            "description": "Successful Response",
            "schema": {
              "$ref": "#/definitions/LivenessProbe"
            }
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/_readinessProbe": {
      "get": {
        "description": "Readiness Probe",
        "tags": [
          "health"
        ],
        "summary": "Readiness Probe",
        "operationId": "getReadinessProbe",
        "responses": {
          "200": {
            "description": "Successful Response"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/app-codes": {
      "get": {
        "description": "Get list of application codes",
        "tags": [
          "general"
        ],
        "summary": "Get List of Application Codes",
        "operationId": "getAppCodes",
        "responses": {
          "200": {
            "description": "Successful Response",
            "schema": {
              "$ref": "#/definitions/GetAppMessagesResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorMessage"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "AppMessage": {
      "type": "object",
      "title": "AppMessage",
      "required": [
        "message",
        "code",
        "attributes"
      ],
      "properties": {
        "attributes": {
          "description": "attributes",
          "type": "array",
          "items": {
            "type": "string"
          },
          "x-nullable": false
        },
        "code": {
          "type": "string",
          "title": "Code",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "title": "Message",
          "x-nullable": false
        }
      }
    },
    "ErrorMessage": {
      "type": "object",
      "title": "ErrorMessage",
      "required": [
        "message"
      ],
      "properties": {
        "attributes": {
          "description": "attributes",
          "type": "object",
          "additionalProperties": {
            "type": "string"
          },
          "x-nullable": false
        },
        "code": {
          "type": "string",
          "title": "Code",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "title": "Message",
          "x-nullable": false
        }
      }
    },
    "GetAppMessagesResponse": {
      "type": "array",
      "title": "GetAppCodesResponse",
      "items": {
        "$ref": "#/definitions/AppMessage"
      }
    },
    "LivenessProbe": {
      "type": "object",
      "title": "LivenessProbe",
      "required": [
        "tag"
      ],
      "properties": {
        "components": {
          "type": "array",
          "title": "LivenessProbeComponents",
          "items": {
            "$ref": "#/definitions/LivenessProbeComponent"
          }
        },
        "tag": {
          "type": "string",
          "title": "Tag",
          "x-nullable": false
        }
      }
    },
    "LivenessProbeComponent": {
      "type": "object",
      "title": "LivenessProbeComponent",
      "required": [
        "name",
        "status"
      ],
      "properties": {
        "name": {
          "type": "string",
          "title": "Name",
          "x-nullable": false
        },
        "status": {
          "type": "boolean",
          "title": "Status",
          "x-nullable": false
        }
      }
    }
  }
}`))
}
