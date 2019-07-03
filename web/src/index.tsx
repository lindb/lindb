import * as React from "react";
import * as ReactDOM from "react-dom";
import {HashRouter as Router, Route, Switch} from "react-router-dom";
import AppPage from "./App";
import "./index.css";

ReactDOM.render(
    <Router>
        <Switch>
            <Route path="/" component={AppPage}/>
        </Switch>
    </Router>,
    document.getElementById("root") as HTMLElement
);
