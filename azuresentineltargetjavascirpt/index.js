"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
var uuid_1 = require("uuid");
var axios_1 = require("axios");
var log4js = require("log4js");
var logger = log4js.getLogger('azure-sentinel-dispatcher');
logger.level = process.env.LOG_LEVEL || 'info';
var express = require("express");
var HTTP = require("cloudevents").HTTP;
var app = express();
var msal = require('@azure/msal-node');
var AZURE_SUBSCRIPTION_ID = process.env.AZURE_SUBSCRIPTION_ID;
var AZURE_RESOURCE_GROUP = process.env.AZURE_RESOURCE_GROUP;
var AZURE_WORKSPACE = process.env.AZURE_WORKSPACE;
var AZURE_CLIENT_SECRET = process.env.AZURE_CLIENT_SECRET;
var AZURE_CLIENT_ID = process.env.AZURE_CLIENT_ID;
var AZURE_TENANT_ID = process.env.AZURE_TENANT_ID;
var AZURE_API_VERSION = '2020-01-01';
var incidentId = (0, uuid_1.v4)();
var requestUrl = "https://management.azure.com/subscriptions/".concat(AZURE_SUBSCRIPTION_ID, "/")
    + "resourceGroups/".concat(AZURE_RESOURCE_GROUP, "/")
    + "providers/Microsoft.OperationalInsights/workspaces/".concat(AZURE_WORKSPACE, "/")
    + "providers/Microsoft.SecurityInsights/incidents/".concat(incidentId, "?api-version=").concat(AZURE_API_VERSION);
app.use(function (req, res, next) {
    var data = "";
    req.setEncoding("utf8");
    req.on("data", function (chunk) {
        data += chunk;
    });
    req.on("end", function () {
        req.body = data;
        next();
    });
});
app.post("/", function (req, res) { return __awaiter(void 0, void 0, void 0, function () {
    var msalConfig, cca, tokenRequest, getTokenResponse, event, incident, result, e_1;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                msalConfig = {
                    auth: {
                        clientId: AZURE_CLIENT_ID,
                        authority: "https://login.microsoftonline.com/" + AZURE_TENANT_ID,
                        clientSecret: AZURE_CLIENT_SECRET
                    }
                };
                cca = new msal.ConfidentialClientApplication(msalConfig);
                tokenRequest = {
                    scopes: ['https://management.azure.com/.default']
                };
                return [4 /*yield*/, cca.acquireTokenByClientCredential(tokenRequest)];
            case 1:
                getTokenResponse = _a.sent();
                try {
                    event = HTTP.toEvent({ headers: req.headers, body: req.body });
                }
                catch (err) {
                    console.error(err);
                    res.status(415).header("Content-Type", "application/json").send(JSON.stringify(err));
                }
                console.log(event.data.event);
                incident = {
                    "properties": {
                        "severity": "High",
                        "status": "Active",
                        "title": "test",
                        "description": "test",
                        "additionalData": {
                        // "alertProductNames": [
                        //   // event.data.resource,
                        //   // event.data.event.resources[0].platform,
                        //   // event.data.event.resources[0].accountId,
                        //   // event.data.event.resources[0].region,
                        //   // event.data.event.resources[0].service,
                        //   // event.data.event.resources[0].type + ':' + event.data.event.resources[0].name + ':' + event.data.event.resources[0].guid
                        // ]
                        },
                        "labels": [
                            {
                                "labelName": "test",
                                "labelType": "User"
                            }
                        ]
                    }
                };
                console.log(incident);
                _a.label = 2;
            case 2:
                _a.trys.push([2, 4, , 5]);
                return [4 /*yield*/, axios_1["default"].put(requestUrl, incident, {
                        headers: {
                            'Authorization': "Bearer ".concat(getTokenResponse.accessToken),
                            'content-type': 'application/json'
                        }
                    })];
            case 3:
                result = _a.sent();
                // console.log(result);
                if (result.status === 200) {
                    logger.debug('incident successfully dispatched');
                    console.log('incident successfully dispatched');
                }
                else {
                    logger.error('failed to dispatch incident', result.data);
                    res.status(500).send(result);
                }
                return [3 /*break*/, 5];
            case 4:
                e_1 = _a.sent();
                logger.error('failed to dispatch incident', e_1);
                res.status(500).send(e_1);
                return [3 /*break*/, 5];
            case 5:
                res.status(200).send("ok");
                return [2 /*return*/];
        }
    });
}); });
app.listen(8080, function () {
    console.log("Example app listening on port 8080!");
});
