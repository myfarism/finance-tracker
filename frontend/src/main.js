"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var react_1 = require("react");
var client_1 = require("react-dom/client");
var react_router_dom_1 = require("react-router-dom");
var LoginPage_1 = require("./pages/LoginPage");
var RegisterPage_1 = require("./pages/RegisterPage");
var DashboardPage_1 = require("./pages/DashboardPage");
var ProtectedRoute_1 = require("./components/ProtectedRoute");
require("./index.css");
client_1.default.createRoot(document.getElementById("root")).render(<react_1.default.StrictMode>
    <react_router_dom_1.BrowserRouter>
      <react_router_dom_1.Routes>
        <react_router_dom_1.Route path="/login" element={<LoginPage_1.default />}/>
        <react_router_dom_1.Route path="/register" element={<RegisterPage_1.default />}/>
        <react_router_dom_1.Route path="/dashboard" element={<ProtectedRoute_1.default>
              <DashboardPage_1.default />
            </ProtectedRoute_1.default>}/>
        <react_router_dom_1.Route path="*" element={<react_router_dom_1.Navigate to="/login" replace/>}/>
      </react_router_dom_1.Routes>
    </react_router_dom_1.BrowserRouter>
  </react_1.default.StrictMode>);
